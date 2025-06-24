package api

import (
	"net/http"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SpecialController 专题API控制器
type SpecialController struct {
	*BaseController
}

// NewSpecialController 创建专题API控制器
func NewSpecialController(db *database.DB, cache cache.Cache, config *config.Config) *SpecialController {
	return &SpecialController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// List 专题列表
func (c *SpecialController) List(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

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

	// 获取专题列表
	specials, err := c.specialModel.GetAll()
	if err != nil {
		logger.Error("获取专题列表失败", "error", err)
		c.Error(w, 500, "Failed to get specials")
		return
	}

	// 计算分页
	total := len(specials)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		start = 0
		end = 0
	}
	if end > total {
		end = total
	}

	// 分页专题
	var pagedSpecials []*model.Special
	if start < end {
		pagedSpecials = specials[start:end]
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
		"specials":   pagedSpecials,
		"pagination": pagination,
	})
}

// Detail 专题详情
func (c *SpecialController) Detail(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取专题ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid special ID")
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

	// 获取专题
	special, err := c.specialModel.GetByID(id)
	if err != nil {
		logger.Error("获取专题失败", "id", id, "error", err)
		c.Error(w, 404, "Special not found")
		return
	}

	// 获取专题文章
	articles, total, err := c.specialModel.GetArticles(id, page, pageSize)
	if err != nil {
		logger.Error("获取专题文章失败", "id", id, "error", err)
		c.Error(w, 500, "Failed to get special articles")
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
		"special":    special,
		"articles":   articles,
		"pagination": pagination,
	})
}
