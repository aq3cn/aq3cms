package service

import (
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ArticleService 文章服务
type ArticleService struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	articleModel  *model.ArticleModel
	categoryModel *model.CategoryModel
	tagModel      *model.TagModel
}

// NewArticleService 创建文章服务
func NewArticleService(db *database.DB, cache cache.Cache, config *config.Config) *ArticleService {
	return &ArticleService{
		db:            db,
		cache:         cache,
		config:        config,
		articleModel:  model.NewArticleModel(db),
		categoryModel: model.NewCategoryModel(db),
		tagModel:      model.NewTagModel(db),
	}
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(id int64) (*model.Article, error) {
	// 缓存键
	cacheKey := "article:" + string(id)

	// 检查缓存
	if s.config.Cache.EnableArticleCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if article, ok := cached.(*model.Article); ok {
				return article, nil
			}
		}
	}

	// 获取文章
	article, err := s.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		return nil, err
	}

	// 缓存文章
	if s.config.Cache.EnableArticleCache {
		cache.SafeSet(s.cache, cacheKey, article, time.Duration(s.config.Cache.ArticleCacheTime)*time.Second)
	}

	return article, nil
}

// GetArticles 获取文章列表
func (s *ArticleService) GetArticles(page, pageSize int, orderBy string) ([]*model.Article, int, error) {
	// 缓存键
	cacheKey := "articles:" + string(page) + ":" + string(pageSize) + ":" + orderBy

	// 检查缓存
	if s.config.Cache.EnableListCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if result, ok := cached.(map[string]interface{}); ok {
				articles, _ := result["articles"].([]*model.Article)
				total, _ := result["total"].(int)
				return articles, total, nil
			}
		}
	}

	// 获取文章列表
	articles, total, err := s.articleModel.GetList(0, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
		return nil, 0, err
	}

	// 缓存文章列表
	if s.config.Cache.EnableListCache {
		result := map[string]interface{}{
			"articles": articles,
			"total":    total,
		}
		cache.SafeSet(s.cache, cacheKey, result, time.Duration(s.config.Cache.ListCacheTime)*time.Second)
	}

	return articles, total, nil
}

// GetArticlesByCategoryID 根据栏目ID获取文章列表
func (s *ArticleService) GetArticlesByCategoryID(categoryID int64, page, pageSize int) ([]*model.Article, int, error) {
	// 缓存键
	cacheKey := "category_articles:" + string(categoryID) + ":" + string(page) + ":" + string(pageSize)

	// 检查缓存
	if s.config.Cache.EnableListCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if result, ok := cached.(map[string]interface{}); ok {
				articles, _ := result["articles"].([]*model.Article)
				total, _ := result["total"].(int)
				return articles, total, nil
			}
		}
	}

	// 获取文章列表
	articles, total, err := s.articleModel.GetList(categoryID, page, pageSize)
	if err != nil {
		logger.Error("获取栏目文章列表失败", "categoryID", categoryID, "error", err)
		return nil, 0, err
	}

	// 缓存文章列表
	if s.config.Cache.EnableListCache {
		result := map[string]interface{}{
			"articles": articles,
			"total":    total,
		}
		cache.SafeSet(s.cache, cacheKey, result, time.Duration(s.config.Cache.ListCacheTime)*time.Second)
	}

	return articles, total, nil
}

// GetLatestArticles 获取最新文章
func (s *ArticleService) GetLatestArticles(limit int) ([]*model.Article, error) {
	// 缓存键
	cacheKey := "latest_articles:" + string(limit)

	// 检查缓存
	if s.config.Cache.EnableListCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if articles, ok := cached.([]*model.Article); ok {
				return articles, nil
			}
		}
	}

	// 获取文章列表
	articles, _, err := s.articleModel.GetList(0, 1, limit)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		return nil, err
	}

	// 缓存文章列表
	if s.config.Cache.EnableListCache {
		cache.SafeSet(s.cache, cacheKey, articles, time.Duration(s.config.Cache.ListCacheTime)*time.Second)
	}

	return articles, nil
}

// GetRecommendArticles 获取推荐文章
func (s *ArticleService) GetRecommendArticles(limit int) ([]*model.Article, error) {
	// 实现获取推荐文章的逻辑
	return nil, nil
}

// GetHotArticles 获取热门文章
func (s *ArticleService) GetHotArticles(limit int) ([]*model.Article, error) {
	// 实现获取热门文章的逻辑
	return nil, nil
}

// GetFocusArticles 获取焦点图文章
func (s *ArticleService) GetFocusArticles(limit int) ([]*model.Article, error) {
	// 实现获取焦点图文章的逻辑
	return nil, nil
}

// SearchArticles 搜索文章
func (s *ArticleService) SearchArticles(keyword string, page, pageSize int) ([]*model.Article, int, error) {
	// 实现搜索文章的逻辑
	return nil, 0, nil
}

// GetArticlesByTag 获取标签文章
func (s *ArticleService) GetArticlesByTag(tag string, page, pageSize int) ([]*model.Article, int, error) {
	// 实现获取标签文章的逻辑
	return nil, 0, nil
}

// GetAllTags 获取所有标签
func (s *ArticleService) GetAllTags() ([]string, error) {
	// 实现获取所有标签的逻辑
	return nil, nil
}
