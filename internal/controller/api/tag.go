package api

import (
	"net/http"

	"aq3cms/config"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// TagController 标签API控制器
type TagController struct {
	*BaseController
}

// NewTagController 创建标签API控制器
func NewTagController(db *database.DB, cache cache.Cache, config *config.Config) *TagController {
	return &TagController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// List 标签列表
func (c *TagController) List(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取查询参数
	limit := c.GetQueryInt(r, "limit", 100)

	// 验证参数
	if limit < 1 || limit > 1000 {
		limit = 100
	}

	// 获取热门标签
	tags, err := c.tagModel.GetHotTags(limit)
	if err != nil {
		logger.Error("获取标签列表失败", "error", err)
		c.Error(w, 500, "Failed to get tags")
		return
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"tags": tags,
	})
}

// Detail 标签详情
func (c *TagController) Detail(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取标签名称
	name := c.GetQueryString(r, "name", "")
	if name == "" {
		c.Error(w, 400, "Missing tag name")
		return
	}

	// 获取查询参数
	page := c.GetQueryInt(r, "page", 1)
	pageSize := c.GetQueryInt(r, "pagesize", 20)

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取标签
	tag, err := c.tagModel.GetByName(name)
	if err != nil {
		logger.Error("获取标签失败", "name", name, "error", err)
		c.Error(w, 404, "Tag not found")
		return
	}

	// 获取标签文章
	articles, total, err := c.tagModel.GetArticles(name, page, pageSize)
	if err != nil {
		logger.Error("获取标签文章失败", "name", name, "error", err)
		c.Error(w, 500, "Failed to get tag articles")
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

	// 返回数据
	c.Success(w, map[string]interface{}{
		"tag":        tag,
		"articles":   articles,
		"pagination": pagination,
	})
}
