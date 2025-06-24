package frontend

import (
	"fmt"
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
	categoryModel   *model.CategoryModel
	articleModel    *model.ArticleModel
	templateService *service.TemplateService
}

// NewCategoryController 创建前台栏目控制器
func NewCategoryController(db *database.DB, cache cache.Cache, config *config.Config) *CategoryController {
	return &CategoryController{
		db:              db,
		cache:           cache,
		config:          config,
		categoryModel:   model.NewCategoryModel(db),
		articleModel:    model.NewArticleModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 栏目列表页
func (c *CategoryController) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeIDStr := vars["typeid"]
	typeID, err := strconv.ParseInt(typeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 检查是否有静态列表页
	if c.config.Site.StaticList && r.URL.Query().Get("upcache") == "" {
		staticPath := fmt.Sprintf("list/%d_%d.html", typeID, page)
		if page == 1 {
			staticPath = fmt.Sprintf("list/%d.html", typeID)
		}
		if _, err := http.Dir(".").Open(staticPath); err == nil {
			http.ServeFile(w, r, staticPath)
			return
		}
	}

	// 获取栏目信息
	category, err := c.categoryModel.GetByID(typeID)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", typeID, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 获取文章列表
	pageSize := 20
	articles, total, err := c.articleModel.GetByTypeID(typeID, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "typeid", typeID, "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (int64(total) + int64(pageSize) - 1) / int64(pageSize)
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     int64(page) < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 获取子栏目
	subCategories, err := c.categoryModel.GetChildCategories(typeID)
	if err != nil {
		logger.Error("获取子栏目失败", "typeid", typeID, "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":       globals,
		"Category":      category,
		"Articles":      articles,
		"Pagination":    pagination,
		"SubCategories": subCategories,
		"PageTitle":     category.TypeName + " - " + c.config.Site.Name,
		"Keywords":      category.Keywords,
		"Description":   category.Description,
		"SiteName":      c.config.Site.Name,
		"Title":         category.TypeName + " - " + c.config.Site.Name,
		"CurrentTime":   "2025-05-31 12:00:00",
		"Version":       "1.0.0",
		"SiteURL":       c.config.Site.URL,
		"CopyRight":     "© 2025 " + c.config.Site.Name + ". All rights reserved.",
	}

	// 确定模板文件
	var tplFile string
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	} else {
		tplFile = c.config.Template.DefaultTpl + "/list.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模板失败", "template", tplFile, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ShowByDir 通过目录显示栏目
func (c *CategoryController) ShowByDir(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dir := vars["dir"]

	// 获取栏目
	category, err := c.categoryModel.GetByDir(dir)
	if err != nil {
		logger.Error("获取栏目失败", "dir", dir, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目是否存在
	if category == nil {
		logger.Error("栏目不存在", "dir", dir)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目状态
	if category.Status != 1 {
		logger.Error("栏目已禁用", "dir", dir, "status", category.Status)
		http.Error(w, "Category is disabled", http.StatusForbidden)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取文章列表
	pageSize := 20
	articles, total, err := c.articleModel.GetByTypeID(category.ID, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "typeid", category.ID, "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (int64(total) + int64(pageSize) - 1) / int64(pageSize)
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     int64(page) < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 获取子栏目
	subCategories, err := c.categoryModel.GetChildCategories(category.ID)
	if err != nil {
		logger.Error("获取子栏目失败", "typeid", category.ID, "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":       globals,
		"Category":      category,
		"Articles":      articles,
		"Pagination":    pagination,
		"SubCategories": subCategories,
		"PageTitle":     category.TypeName + " - " + c.config.Site.Name,
		"Keywords":      category.Keywords,
		"Description":   category.Description,
		"SiteName":      c.config.Site.Name,
		"Title":         category.TypeName + " - " + c.config.Site.Name,
		"CurrentTime":   "2025-05-31 12:00:00",
		"Version":       "1.0.0",
		"SiteURL":       c.config.Site.URL,
		"CopyRight":     "© 2025 " + c.config.Site.Name + ". All rights reserved.",
	}

	// 确定模板文件
	var tplFile string
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	} else {
		tplFile = c.config.Template.DefaultTpl + "/list.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模板失败", "template", tplFile, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ShowByPath 通过多级路径显示栏目
func (c *CategoryController) ShowByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dir := vars["dir"]
	subdir := vars["subdir"]

	// 获取栏目
	category, err := c.categoryModel.GetByPath(dir, subdir)
	if err != nil {
		logger.Error("获取栏目失败", "dir", dir, "subdir", subdir, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目是否存在
	if category == nil {
		logger.Error("栏目不存在", "dir", dir, "subdir", subdir)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 检查栏目状态
	if category.Status != 1 {
		logger.Error("栏目已禁用", "dir", dir, "subdir", subdir, "status", category.Status)
		http.Error(w, "Category is disabled", http.StatusForbidden)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取文章列表
	pageSize := 20
	articles, total, err := c.articleModel.GetByTypeID(category.ID, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "typeid", category.ID, "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (int64(total) + int64(pageSize) - 1) / int64(pageSize)
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     int64(page) < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 获取子栏目
	subCategories, err := c.categoryModel.GetChildCategories(category.ID)
	if err != nil {
		logger.Error("获取子栏目失败", "typeid", category.ID, "error", err)
	}

	// 获取父栏目信息
	var parentCategory *model.Category
	if category.ParentID > 0 {
		parentCategory, err = c.categoryModel.GetByID(category.ParentID)
		if err != nil {
			logger.Error("获取父栏目失败", "parentID", category.ParentID, "error", err)
		}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":        globals,
		"Category":       category,
		"ParentCategory": parentCategory,
		"Articles":       articles,
		"Pagination":     pagination,
		"SubCategories":  subCategories,
		"PageTitle":      category.TypeName + " - " + c.config.Site.Name,
		"Keywords":       category.Keywords,
		"Description":    category.Description,
		"SiteName":       c.config.Site.Name,
		"Title":          category.TypeName + " - " + c.config.Site.Name,
		"CurrentTime":    "2025-05-31 12:00:00",
		"Version":        "1.0.0",
		"SiteURL":        c.config.Site.URL,
		"CopyRight":      "© 2025 " + c.config.Site.Name + ". All rights reserved.",
		"Path":           "/" + dir + "/" + subdir,
	}

	// 确定模板文件
	var tplFile string
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	} else {
		tplFile = c.config.Template.DefaultTpl + "/list.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模板失败", "template", tplFile, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
