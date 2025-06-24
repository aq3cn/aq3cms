package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CollectService 采集服务
type CollectService struct {
	db               *database.DB
	cache            cache.Cache
	config           *config.Config
	collectRuleModel *model.CollectRuleModel
	collectItemModel *model.CollectItemModel
	articleModel     *model.ArticleModel
	categoryModel    *model.CategoryModel
	contentModel     *model.ContentModelModel
}

// NewCollectService 创建采集服务
func NewCollectService(db *database.DB, cache cache.Cache, config *config.Config) *CollectService {
	return &CollectService{
		db:               db,
		cache:            cache,
		config:           config,
		collectRuleModel: model.NewCollectRuleModel(db),
		collectItemModel: model.NewCollectItemModel(db),
		articleModel:     model.NewArticleModel(db),
		categoryModel:    model.NewCategoryModel(db),
		contentModel:     model.NewContentModelModel(db),
	}
}

// CollectRule 采集规则
func (s *CollectService) CollectRule(ruleID int64) (int, error) {
	// 获取采集规则
	rule, err := s.collectRuleModel.GetByID(ruleID)
	if err != nil {
		return 0, err
	}

	// 检查规则状态
	if rule.Status != 1 {
		return 0, fmt.Errorf("collect rule is disabled")
	}

	// 采集计数
	count := 0

	// 根据来源类型采集
	if rule.SourceType == 0 {
		// 列表采集
		for page := rule.StartPage; page <= rule.EndPage; page++ {
			// 构建URL
			url := rule.SourceURL
			if rule.PageRule != "" {
				url = strings.Replace(rule.PageRule, "{page}", fmt.Sprintf("%d", page), -1)
			}

			// 采集列表
			n, err := s.collectList(rule, url)
			if err != nil {
				logger.Error("采集列表失败", "url", url, "error", err)
				continue
			}
			count += n
		}
	} else if rule.SourceType == 1 {
		// RSS采集
		n, err := s.collectRSS(rule)
		if err != nil {
			logger.Error("采集RSS失败", "url", rule.SourceURL, "error", err)
			return 0, err
		}
		count += n
	}

	// 更新最后采集时间
	s.collectRuleModel.UpdateLastTime(ruleID)

	return count, nil
}

// 采集列表
func (s *CollectService) collectList(rule *model.CollectRule, url string) (int, error) {
	// 获取页面内容
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 读取页面内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// 解析页面内容
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return 0, err
	}

	// 采集计数
	count := 0

	// 查找列表项
	doc.Find(rule.ListRule).Each(func(i int, selection *goquery.Selection) {
		// 获取链接
		link, exists := selection.Find("a").Attr("href")
		if !exists {
			return
		}

		// 处理相对链接
		if !strings.HasPrefix(link, "http") {
			baseURL := getBaseURL(url)
			link = baseURL + link
		}

		// 检查链接是否已采集
		exists, err := s.checkLinkExists(link)
		if err != nil {
			logger.Error("检查链接是否已采集失败", "link", link, "error", err)
			return
		}
		if exists {
			return
		}

		// 采集内容
		err = s.collectContent(rule, link)
		if err != nil {
			logger.Error("采集内容失败", "link", link, "error", err)
			return
		}

		count++
	})

	return count, nil
}

// 采集RSS
func (s *CollectService) collectRSS(rule *model.CollectRule) (int, error) {
	// 获取RSS内容
	resp, err := http.Get(rule.SourceURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 读取RSS内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// 解析RSS内容
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return 0, err
	}

	// 采集计数
	count := 0

	// 查找RSS项
	doc.Find("item").Each(func(i int, selection *goquery.Selection) {
		// 获取链接
		link := selection.Find("link").Text()
		if link == "" {
			return
		}

		// 检查链接是否已采集
		exists, err := s.checkLinkExists(link)
		if err != nil {
			logger.Error("检查链接是否已采集失败", "link", link, "error", err)
			return
		}
		if exists {
			return
		}

		// 获取标题
		title := selection.Find("title").Text()

		// 获取内容
		content := selection.Find("description").Text()

		// 创建采集项目
		item := &model.CollectItem{
			RuleID:  rule.ID,
			Title:   title,
			URL:     link,
			Content: content,
			Status:  0,
		}

		// 保存采集项目
		_, err = s.collectItemModel.Create(item)
		if err != nil {
			logger.Error("保存采集项目失败", "error", err)
			return
		}

		count++
	})

	return count, nil
}

// 检查链接是否已采集
func (s *CollectService) checkLinkExists(link string) (bool, error) {
	// 构建查询
	qb := database.NewQueryBuilder(s.db, "collect_item")
	qb.Where("url = ?", link)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 采集内容
func (s *CollectService) collectContent(rule *model.CollectRule, link string) error {
	// 获取页面内容
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取页面内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析页面内容
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	// 获取标题
	title := ""
	if rule.TitleRule != "" {
		title = doc.Find(rule.TitleRule).Text()
	}
	if title == "" {
		title = doc.Find("title").Text()
	}

	// 获取内容
	content := ""
	if rule.ContentRule != "" {
		content, _ = doc.Find(rule.ContentRule).Html()
	}

	// 获取字段数据
	fieldData := make(map[string]interface{})
	if rule.FieldRules != "" {
		var fieldRules []model.FieldRule
		err = json.Unmarshal([]byte(rule.FieldRules), &fieldRules)
		if err != nil {
			logger.Error("解析字段规则失败", "error", err)
		} else {
			for _, fieldRule := range fieldRules {
				// 获取字段值
				value := ""
				if fieldRule.Rule != "" {
					// 使用选择器
					if strings.HasPrefix(fieldRule.Rule, "#") || strings.HasPrefix(fieldRule.Rule, ".") {
						value = doc.Find(fieldRule.Rule).Text()
					} else {
						// 使用正则表达式
						re := regexp.MustCompile(fieldRule.Rule)
						matches := re.FindStringSubmatch(string(body))
						if len(matches) > 1 {
							value = matches[1]
						}
					}
				}

				// 使用默认值
				if value == "" && fieldRule.Default != "" {
					value = fieldRule.Default
				}

				// 设置字段值
				fieldData[fieldRule.Field] = value
			}
		}
	}

	// 序列化字段数据
	fieldDataJSON, err := json.Marshal(fieldData)
	if err != nil {
		logger.Error("序列化字段数据失败", "error", err)
		return err
	}

	// 创建采集项目
	item := &model.CollectItem{
		RuleID:    rule.ID,
		Title:     title,
		URL:       link,
		Content:   content,
		Status:    0,
		FieldData: string(fieldDataJSON),
	}

	// 保存采集项目
	_, err = s.collectItemModel.Create(item)
	if err != nil {
		logger.Error("保存采集项目失败", "error", err)
		return err
	}

	return nil
}

// PublishItem 发布采集项目
func (s *CollectService) PublishItem(itemID int64) error {
	// 获取采集项目
	item, err := s.collectItemModel.GetByID(itemID)
	if err != nil {
		return err
	}

	// 获取采集规则
	rule, err := s.collectRuleModel.GetByID(item.RuleID)
	if err != nil {
		return err
	}

	// 获取字段数据
	fieldData, err := s.collectItemModel.GetFieldData(itemID)
	if err != nil {
		return err
	}

	// 创建文章
	article := &model.Article{
		TypeID:      rule.TypeID,
		Title:       item.Title,
		ShortTitle:  "",
		Color:       "",
		Source:      "采集",
		Writer:      "采集",
		LitPic:      "",
		Keywords:    "",
		Description: "",
		Body:        item.Content,
		PubDate:     time.Now(),
		SendDate:    time.Now(),
		Click:       0,
		IsTop:       0,
		IsRecommend: 0,
		IsHot:       0,
		ArcRank:     0,
	}

	// 保存文章
	articleID, err := s.articleModel.Create(article)
	if err != nil {
		logger.Error("保存文章失败", "error", err)
		return err
	}

	// 保存扩展模型内容
	if rule.ModelID > 0 {
		err = s.contentModel.SaveContent(rule.ModelID, articleID, fieldData)
		if err != nil {
			logger.Error("保存扩展模型内容失败", "error", err)
			return err
		}
	}

	// 更新采集项目状态
	err = s.collectItemModel.UpdateStatus(itemID, 2)
	if err != nil {
		logger.Error("更新采集项目状态失败", "error", err)
		return err
	}

	return nil
}

// BatchPublish 批量发布采集项目
func (s *CollectService) BatchPublish(ruleID int64) (int, error) {
	// 获取采集规则
	_, err := s.collectRuleModel.GetByID(ruleID)
	if err != nil {
		return 0, err
	}

	// 获取未处理的采集项目
	items, _, err := s.collectItemModel.GetByRuleID(ruleID, 0, 1, 100)
	if err != nil {
		return 0, err
	}

	// 发布计数
	count := 0

	// 批量发布
	for _, item := range items {
		err = s.PublishItem(item.ID)
		if err != nil {
			logger.Error("发布采集项目失败", "id", item.ID, "error", err)
			continue
		}
		count++
	}

	return count, nil
}

// DeleteRule 删除采集规则
func (s *CollectService) DeleteRule(ruleID int64) error {
	// 删除采集项目
	err := s.collectItemModel.DeleteByRuleID(ruleID)
	if err != nil {
		return err
	}

	// 删除采集规则
	return s.collectRuleModel.Delete(ruleID)
}

// 获取基础URL
func getBaseURL(url string) string {
	// 查找最后一个斜杠
	index := strings.LastIndex(url, "/")
	if index == -1 {
		return url
	}

	// 获取基础URL
	baseURL := url[:index+1]
	return baseURL
}
