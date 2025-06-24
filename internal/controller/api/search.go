package api

import (
	"net/http"

	"aq3cms/config"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SearchController 搜索API控制器
type SearchController struct {
	*BaseController
}

// NewSearchController 创建搜索API控制器
func NewSearchController(db *database.DB, cache cache.Cache, config *config.Config) *SearchController {
	return &SearchController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// Search 搜索
func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取查询参数
	keyword := c.GetQueryString(r, "keyword", "")
	_ = c.GetQueryInt64(r, "typeid", 0) // 未使用的变量
	_ = c.GetQueryInt(r, "flag", 0) // 未使用的变量
	page := c.GetQueryInt(r, "page", 1)
	pageSize := c.GetQueryInt(r, "pagesize", 20)

	// 验证参数
	if keyword == "" {
		c.Error(w, 400, "Missing keyword")
		return
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 搜索文章
	articles, total, err := c.articleModel.Search(keyword, page, pageSize)
	if err != nil {
		logger.Error("搜索文章失败", "error", err)
		c.Error(w, 500, "Failed to search articles")
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"current_page": page,
		"total_pages":  totalPages,
		"total_items":  total,
		"page_size":    pageSize,
	}

	// 记录搜索
	c.statsService.RecordSearch(r, keyword, total)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"keyword":    keyword,
		"articles":   articles,
		"pagination": pagination,
	})
}

// Hot 热门搜索
func (c *SearchController) Hot(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取查询参数
	limit := c.GetQueryInt(r, "limit", 10)

	// 验证参数
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 获取热门搜索
	hotSearches, err := c.statsService.GetHotSearches(limit)
	if err != nil {
		logger.Error("获取热门搜索失败", "error", err)
		c.Error(w, 500, "Failed to get hot searches")
		return
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"hot_searches": hotSearches,
	})
}
