package service

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// AdService 广告服务
type AdService struct {
	db             *database.DB
	cache          cache.Cache
	config         *config.Config
	adModel        *model.AdModel
	adPositionModel *model.AdPositionModel
}

// NewAdService 创建广告服务
func NewAdService(db *database.DB, cache cache.Cache, config *config.Config) *AdService {
	return &AdService{
		db:             db,
		cache:          cache,
		config:         config,
		adModel:        model.NewAdModel(db),
		adPositionModel: model.NewAdPositionModel(db),
	}
}

// GetAdsByPosition 根据广告位获取广告
func (s *AdService) GetAdsByPosition(positionID int64) ([]*model.Ad, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("ads:position:%d", positionID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if ads, ok := cached.([]*model.Ad); ok {
			return ads, nil
		}
	}

	// 获取广告
	ads, err := s.adModel.GetByPositionID(positionID)
	if err != nil {
		return nil, err
	}

	// 缓存广告
	cache.SafeSet(s.cache, cacheKey, ads, time.Hour)

	return ads, nil
}

// GetAdsByCode 根据广告位代码获取广告
func (s *AdService) GetAdsByCode(code string) ([]*model.Ad, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("ads:code:%s", code)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if ads, ok := cached.([]*model.Ad); ok {
			return ads, nil
		}
	}

	// 获取广告位
	position, err := s.adPositionModel.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// 获取广告
	ads, err := s.adModel.GetByPositionID(position.ID)
	if err != nil {
		return nil, err
	}

	// 缓存广告
	cache.SafeSet(s.cache, cacheKey, ads, time.Hour)

	return ads, nil
}

// GetAdHTML 获取广告HTML
func (s *AdService) GetAdHTML(ad *model.Ad) (string, error) {
	// 根据广告类型生成HTML
	switch ad.Type {
	case 0: // 图片
		return s.getImageAdHTML(ad)
	case 1: // Flash
		return s.getFlashAdHTML(ad)
	case 2: // 代码
		return ad.Code, nil
	case 3: // 文字
		return s.getTextAdHTML(ad)
	default:
		return "", fmt.Errorf("unknown ad type: %d", ad.Type)
	}
}

// getImageAdHTML 获取图片广告HTML
func (s *AdService) getImageAdHTML(ad *model.Ad) (string, error) {
	// 构建HTML
	html := fmt.Sprintf(`<a href="%s" target="%s" title="%s"><img src="%s" width="%d" height="%d" alt="%s" /></a>`,
		ad.URL, ad.Target, ad.Title, ad.Image, ad.Width, ad.Height, ad.Title)
	return html, nil
}

// getFlashAdHTML 获取Flash广告HTML
func (s *AdService) getFlashAdHTML(ad *model.Ad) (string, error) {
	// 构建HTML
	html := fmt.Sprintf(`<object classid="clsid:D27CDB6E-AE6D-11cf-96B8-444553540000" codebase="http://download.macromedia.com/pub/shockwave/cabs/flash/swflash.cab#version=7,0,19,0" width="%d" height="%d">
		<param name="movie" value="%s" />
		<param name="quality" value="high" />
		<embed src="%s" quality="high" pluginspage="http://www.macromedia.com/go/getflashplayer" type="application/x-shockwave-flash" width="%d" height="%d"></embed>
	</object>`,
		ad.Width, ad.Height, ad.Flash, ad.Flash, ad.Width, ad.Height)
	return html, nil
}

// getTextAdHTML 获取文字广告HTML
func (s *AdService) getTextAdHTML(ad *model.Ad) (string, error) {
	// 构建HTML
	html := fmt.Sprintf(`<a href="%s" target="%s" title="%s">%s</a>`,
		ad.URL, ad.Target, ad.Title, ad.Text)
	return html, nil
}

// GetPositionHTML 获取广告位HTML
func (s *AdService) GetPositionHTML(positionID int64) (string, error) {
	// 获取广告位
	position, err := s.adPositionModel.GetByID(positionID)
	if err != nil {
		return "", err
	}

	// 获取广告
	ads, err := s.GetAdsByPosition(positionID)
	if err != nil {
		return "", err
	}

	// 如果没有广告，返回空
	if len(ads) == 0 {
		return "", nil
	}

	// 如果没有模板，使用默认模板
	if position.Template == "" {
		// 获取第一个广告的HTML
		return s.GetAdHTML(ads[0])
	}

	// 解析模板
	tmpl, err := template.New("ad").Parse(position.Template)
	if err != nil {
		logger.Error("解析广告位模板失败", "error", err)
		return "", err
	}

	// 准备广告HTML
	adHTMLs := make([]string, 0, len(ads))
	for _, ad := range ads {
		html, err := s.GetAdHTML(ad)
		if err != nil {
			logger.Error("获取广告HTML失败", "error", err)
			continue
		}
		adHTMLs = append(adHTMLs, html)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Position": position,
		"Ads":      ads,
		"AdHTMLs":  adHTMLs,
	}

	// 渲染模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		logger.Error("渲染广告位模板失败", "error", err)
		return "", err
	}

	return buf.String(), nil
}

// GetPositionHTMLByCode 根据广告位代码获取广告位HTML
func (s *AdService) GetPositionHTMLByCode(code string) (string, error) {
	// 获取广告位
	position, err := s.adPositionModel.GetByCode(code)
	if err != nil {
		return "", err
	}

	// 获取广告位HTML
	return s.GetPositionHTML(position.ID)
}

// ClickAd 点击广告
func (s *AdService) ClickAd(id int64) error {
	// 增加点击次数
	err := s.adModel.IncrementClick(id)
	if err != nil {
		return err
	}

	// 清除缓存
	s.ClearCache(id)

	return nil
}

// ClearCache 清除缓存
func (s *AdService) ClearCache(id int64) {
	// 获取广告
	ad, err := s.adModel.GetByID(id)
	if err != nil {
		logger.Error("获取广告失败", "id", id, "error", err)
		return
	}

	// 清除广告位缓存
	cacheKey := fmt.Sprintf("ads:position:%d", ad.PositionID)
	s.cache.Delete(cacheKey)

	// 获取广告位
	position, err := s.adPositionModel.GetByID(ad.PositionID)
	if err != nil {
		logger.Error("获取广告位失败", "id", ad.PositionID, "error", err)
		return
	}

	// 清除广告位代码缓存
	cacheKey = fmt.Sprintf("ads:code:%s", position.Code)
	s.cache.Delete(cacheKey)
}

// ClearAllCache 清除所有缓存
func (s *AdService) ClearAllCache() {
	// 获取所有广告位
	positions, err := s.adPositionModel.GetAll(-1)
	if err != nil {
		logger.Error("获取所有广告位失败", "error", err)
		return
	}

	// 清除所有广告位缓存
	for _, position := range positions {
		cacheKey := fmt.Sprintf("ads:position:%d", position.ID)
		s.cache.Delete(cacheKey)

		cacheKey = fmt.Sprintf("ads:code:%s", position.Code)
		s.cache.Delete(cacheKey)
	}
}
