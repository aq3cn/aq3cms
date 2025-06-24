package service

import (
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// LinkService 友情链接服务
type LinkService struct {
	db        *database.DB
	cache     cache.Cache
	config    *config.Config
	linkModel *model.LinkModel
}

// NewLinkService 创建友情链接服务
func NewLinkService(db *database.DB, cache cache.Cache, config *config.Config) *LinkService {
	return &LinkService{
		db:        db,
		cache:     cache,
		config:    config,
		linkModel: model.NewLinkModel(db),
	}
}

// GetLinks 获取友情链接
func (s *LinkService) GetLinks(limit int) ([]*model.Link, error) {
	// 缓存键
	cacheKey := "links_" + string(limit)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if links, ok := cached.([]*model.Link); ok {
			return links, nil
		}
	}

	// 获取友情链接
	links, err := s.linkModel.GetLinks(limit)
	if err != nil {
		logger.Error("获取友情链接失败", "error", err)
		return nil, err
	}

	// 缓存友情链接
	cache.SafeSet(s.cache, cacheKey, links, time.Duration(3600)*time.Second)

	return links, nil
}

// GetLinkByID 根据ID获取友情链接
func (s *LinkService) GetLinkByID(id int64) (*model.Link, error) {
	// 缓存键
	cacheKey := "link_" + string(id)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if link, ok := cached.(*model.Link); ok {
			return link, nil
		}
	}

	// 获取友情链接
	link, err := s.linkModel.GetByID(id)
	if err != nil {
		logger.Error("获取友情链接失败", "id", id, "error", err)
		return nil, err
	}

	// 缓存友情链接
	cache.SafeSet(s.cache, cacheKey, link, time.Duration(3600)*time.Second)

	return link, nil
}

// GetLinksByType 根据类型获取友情链接
func (s *LinkService) GetLinksByType(typeID int64, limit int) ([]*model.Link, error) {
	// 缓存键
	cacheKey := "links_type_" + string(typeID) + "_" + string(limit)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if links, ok := cached.([]*model.Link); ok {
			return links, nil
		}
	}

	// 获取友情链接
	links, err := s.linkModel.GetByType(typeID, limit)
	if err != nil {
		logger.Error("获取友情链接失败", "typeID", typeID, "error", err)
		return nil, err
	}

	// 缓存友情链接
	cache.SafeSet(s.cache, cacheKey, links, time.Duration(3600)*time.Second)

	return links, nil
}
