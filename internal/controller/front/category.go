package front

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CategoryController 前台栏目控制器
type CategoryController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryService *service.CategoryService
	adService       *service.AdService
	templateService *service.TemplateService
}

// NewCategoryController 创建前台栏目控制器
func NewCategoryController(db *database.DB, cache cache.Cache, config *config.Config) *CategoryController {
	return &CategoryController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryService: service.NewCategoryService(db, cache, config),
		adService:       service.NewAdService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Show 显示栏目
func (c *CategoryController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取栏目ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 获取栏目
	category, err := c.categoryService.GetCategoryByID(id)
	if err != nil {
		logger.Error("获取栏目失败", "id", id, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目状态
	if category.Status != 1 {
		http.Error(w, "Category is disabled", http.StatusForbidden)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取栏目文章
	articles, total, err := c.articleModel.GetByTypeID(id, page, 20)
	if err != nil {
		logger.Error("获取栏目文章失败", "error", err)
		http.Error(w, "Failed to get category articles", http.StatusInternalServerError)
		return
	}

	// 获取子栏目
	childCategories, err := c.categoryService.GetChildCategories(id)
	if err != nil {
		logger.Error("获取子栏目失败", "error", err)
	}

	// 获取栏目路径
	categoryPath, err := c.categoryService.GetCategoryPath(id)
	if err != nil {
		logger.Error("获取栏目路径失败", "error", err)
	}

	// 广告相关
	var topBanner string
	var sideBanner string

	// 简化处理，实际应该从广告服务获取
	topBanner = ""
	sideBanner = ""

	// 计算分页信息
	totalPages := (total + 20 - 1) / 20
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Category":        category,
		"Articles":        articles,
		"ChildCategories": childCategories,
		"CategoryPath":    categoryPath,
		"TopBanner":       topBanner,
		"SideBanner":      sideBanner,
		"Pagination":      pagination,
		"PageTitle":       category.TypeName,
	}

	// 渲染模板
	tplFile := "category.htm"
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	}
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ShowByDir 通过目录显示栏目
func (c *CategoryController) ShowByDir(w http.ResponseWriter, r *http.Request) {
	// 获取栏目目录
	vars := mux.Vars(r)
	dir := vars["dir"]

	// 获取栏目
	category, err := c.categoryService.GetCategoryByDir(dir)
	if err != nil {
		logger.Error("获取栏目失败", "dir", dir, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目状态
	if category.Status != 1 {
		http.Error(w, "Category is disabled", http.StatusForbidden)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取栏目文章
	articles, total, err := c.articleModel.GetByTypeID(category.ID, page, 20)
	if err != nil {
		logger.Error("获取栏目文章失败", "error", err)
		http.Error(w, "Failed to get category articles", http.StatusInternalServerError)
		return
	}

	// 获取子栏目
	childCategories, err := c.categoryService.GetChildCategories(category.ID)
	if err != nil {
		logger.Error("获取子栏目失败", "error", err)
	}

	// 获取栏目路径
	categoryPath, err := c.categoryService.GetCategoryPath(category.ID)
	if err != nil {
		logger.Error("获取栏目路径失败", "error", err)
	}

	// 广告相关
	var topBanner string
	var sideBanner string

	// 简化处理，实际应该从广告服务获取
	topBanner = ""
	sideBanner = ""

	// 计算分页信息
	totalPages := (total + 20 - 1) / 20
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Category":        category,
		"Articles":        articles,
		"ChildCategories": childCategories,
		"CategoryPath":    categoryPath,
		"TopBanner":       topBanner,
		"SideBanner":      sideBanner,
		"Pagination":      pagination,
		"PageTitle":       category.TypeName,
	}

	// 渲染模板
	tplFile := "category.htm"
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	}
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
