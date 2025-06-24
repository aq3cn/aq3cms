package api

import (
	"net/http"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CategoryController 栏目API控制器
type CategoryController struct {
	*BaseController
}

// NewCategoryController 创建栏目API控制器
func NewCategoryController(db *database.DB, cache cache.Cache, config *config.Config) *CategoryController {
	return &CategoryController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// List 栏目列表
func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取查询参数
	parentID := c.GetQueryInt64(r, "parent_id", 0)
	_ = c.GetQueryInt(r, "is_nav", -1) // 暂时不使用

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		c.Error(w, 500, "Failed to get categories")
		return
	}

	// 过滤栏目
	var filteredCategories []*model.Category
	for _, category := range categories {
		// 过滤父栏目
		if parentID > 0 && category.ParentID != parentID {
			continue
		}

		filteredCategories = append(filteredCategories, category)
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"categories": filteredCategories,
	})
}

// Detail 栏目详情
func (c *CategoryController) Detail(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取栏目ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid category ID")
		return
	}

	// 获取栏目
	category, err := c.categoryModel.GetByID(id)
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		c.Error(w, 404, "Category not found")
		return
	}

	// 获取子栏目
	children, err := c.categoryModel.GetChildCategories(id)
	if err != nil {
		logger.Error("获取子栏目失败", "id", id, "error", err)
	}

	// 获取栏目文章
	articles, total, err := c.articleModel.GetByTypeID(id, 1, 10)
	if err != nil {
		logger.Error("获取栏目文章失败", "id", id, "error", err)
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"category":  category,
		"children":  children,
		"articles":  articles,
		"total":     total,
	})
}
