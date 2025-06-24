package frontend

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ArticleController 文章控制器
type ArticleController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	templateService *service.TemplateService
}

// NewArticleController 创建文章控制器
func NewArticleController(db *database.DB, cache cache.Cache, config *config.Config) *ArticleController {
	return &ArticleController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Detail 文章详情
func (c *ArticleController) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 检查是否有静态文章页
	if c.config.Site.StaticArticle && r.URL.Query().Get("upcache") == "" {
		staticPath := fmt.Sprintf("a/%d.html", id)
		if _, err := http.Dir(".").Open(staticPath); err == nil {
			http.ServeFile(w, r, staticPath)
			return
		}
	}

	// 获取文章详情
	article, err := c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章详情失败", "id", id, "error", err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// 增加点击量
	go func() {
		if err := c.articleModel.IncrementClick(id); err != nil {
			logger.Error("增加点击量失败", "id", id, "error", err)
		}
	}()

	// 获取栏目信息
	category, err := c.categoryModel.GetByID(article.TypeID)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", article.TypeID, "error", err)
	}

	// 获取上一篇文章
	prevArticle, err := c.articleModel.GetPrevArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取上一篇文章失败", "id", id, "error", err)
	}

	// 获取下一篇文章
	nextArticle, err := c.articleModel.GetNextArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取下一篇文章失败", "id", id, "error", err)
	}

	// 获取相关文章
	relatedArticles, err := c.articleModel.GetRelatedArticles(article.Keywords, article.ID, 10)
	if err != nil {
		logger.Error("获取相关文章失败", "id", id, "error", err)
	}

	// 获取热门文章（按点击量排序）
	hotArticles, _, err := c.articleModel.GetList(0, 1, 5)
	if err != nil {
		logger.Error("获取热门文章失败", "error", err)
	}

	// 获取本栏目热门文章
	categoryHotArticles, _, err := c.articleModel.GetList(article.TypeID, 1, 5)
	if err != nil {
		logger.Error("获取本栏目热门文章失败", "typeid", article.TypeID, "error", err)
	}

	// 获取最新文章
	latestArticles, err := c.articleModel.GetLatestArticles(5)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
	}

	// 获取所有栏目（使用简单的方法）
	var categories []*model.Category
	// 这里可以后续添加获取栏目的逻辑

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":             globals,
		"Article":             article,
		"Category":            category,
		"PrevArticle":         prevArticle,
		"NextArticle":         nextArticle,
		"RelatedArticles":     relatedArticles,
		"HotArticles":         hotArticles,
		"CategoryHotArticles": categoryHotArticles,
		"LatestArticles":      latestArticles,
		"Categories":          categories,
		"PageTitle":           article.Title + " - " + c.config.Site.Name,
		"Keywords":            article.Keywords,
		"Description":         article.Description,
		"Title":               article.Title + " - " + c.config.Site.Name,
		"SiteName":            c.config.Site.Name,
		"CurrentTime":         "2025-05-31 12:00:00",
		"Version":             "1.0.0",
		"SiteURL":             c.config.Site.URL,
		"CopyRight":           "© 2025 " + c.config.Site.Name + ". All rights reserved.",
	}

	// 确定模板文件
	var tplFile string
	if article.TemplateFile != "" {
		tplFile = article.TemplateFile
	} else if category != nil && category.ArticleTpl != "" {
		tplFile = category.ArticleTpl
	} else {
		tplFile = c.config.Template.DefaultTpl + "/article.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticArticle {
		go func() {
			staticPath := fmt.Sprintf("a/%d.html", id)
			dir := filepath.Dir(staticPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				logger.Error("创建静态文章目录失败", "dir", dir, "error", err)
				return
			}
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticPath); err != nil {
				logger.Error("生成静态文章页失败", "id", id, "error", err)
			}
		}()
	}
}

// List 文章列表
func (c *ArticleController) List(w http.ResponseWriter, r *http.Request) {
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
	pageSize := 10
	articles, total, err := c.articleModel.GetList(typeID, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "typeid", typeID, "error", err)
		http.Error(w, "Failed to get article list", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
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
		logger.Error("渲染列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticList {
		go func() {
			staticPath := fmt.Sprintf("list/%d_%d.html", typeID, page)
			if page == 1 {
				staticPath = fmt.Sprintf("list/%d.html", typeID)
			}
			dir := filepath.Dir(staticPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				logger.Error("创建静态列表目录失败", "dir", dir, "error", err)
				return
			}
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticPath); err != nil {
				logger.Error("生成静态列表页失败", "typeid", typeID, "page", page, "error", err)
			}
		}()
	}
}

// Search 搜索文章
func (c *ArticleController) Search(w http.ResponseWriter, r *http.Request) {
	// 获取搜索关键词
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取文章列表
	pageSize := 10
	articles, total, err := c.articleModel.Search(keyword, page, pageSize)
	if err != nil {
		logger.Error("搜索文章失败", "keyword", keyword, "error", err)
		http.Error(w, "Failed to search articles", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Keyword":     keyword,
		"Articles":    articles,
		"Pagination":  pagination,
		"PageTitle":   "搜索: " + keyword + " - " + c.config.Site.Name,
		"Keywords":    keyword,
		"Description": "搜索结果: " + keyword,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/search.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染搜索模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AllArticles 显示所有文章列表
func (c *ArticleController) AllArticles(w http.ResponseWriter, r *http.Request) {
	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 检查是否有静态列表页
	if c.config.Site.StaticList && r.URL.Query().Get("upcache") == "" {
		staticPath := fmt.Sprintf("articles_%d.html", page)
		if page == 1 {
			staticPath = "articles.html"
		}
		if _, err := http.Dir(".").Open(staticPath); err == nil {
			http.ServeFile(w, r, staticPath)
			return
		}
	}

	// 获取文章列表（所有栏目，typeID=0表示所有）
	pageSize := 15
	articles, total, err := c.articleModel.GetList(0, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
		http.Error(w, "Failed to get article list", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 获取所有栏目
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
	}

	// 获取热门文章
	hotArticles, _, err := c.articleModel.GetList(0, 1, 10)
	if err != nil {
		logger.Error("获取热门文章失败", "error", err)
	}

	// 获取最新文章
	latestArticles, err := c.articleModel.GetLatestArticles(10)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":        globals,
		"Articles":       articles,
		"Pagination":     pagination,
		"Categories":     categories,
		"HotArticles":    hotArticles,
		"LatestArticles": latestArticles,
		"PageTitle":      "文章列表 - " + c.config.Site.Name,
		"Keywords":       "文章,列表," + c.config.Site.Name,
		"Description":    "浏览所有文章内容",
		"Title":          "文章列表 - " + c.config.Site.Name,
		"SiteName":       c.config.Site.Name,
		"CurrentTime":    "2025-05-31 12:00:00",
		"Version":        "1.0.0",
		"SiteURL":        c.config.Site.URL,
		"CopyRight":      "© 2025 " + c.config.Site.Name + ". All rights reserved.",
	}

	// 使用文章列表模板
	tplFile := c.config.Template.DefaultTpl + "/articles.htm"

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticList {
		go func() {
			staticPath := fmt.Sprintf("articles_%d.html", page)
			if page == 1 {
				staticPath = "articles.html"
			}
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticPath); err != nil {
				logger.Error("生成静态文章列表页失败", "page", page, "error", err)
			}
		}()
	}
}
