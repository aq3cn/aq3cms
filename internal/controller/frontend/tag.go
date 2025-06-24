package frontend

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// TagController 标签控制器
type TagController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	tagModel        *model.TagModel
	articleModel    *model.ArticleModel
	templateService *service.TemplateService
	articleService  *service.ArticleService
}

// NewTagController 创建标签控制器
func NewTagController(db *database.DB, cache cache.Cache, config *config.Config) *TagController {
	return &TagController{
		db:              db,
		cache:           cache,
		config:          config,
		tagModel:        model.NewTagModel(db),
		articleModel:    model.NewArticleModel(db),
		templateService: service.NewTemplateService(db, cache, config),
		articleService:  service.NewArticleService(db, cache, config),
	}
}

// List 标签列表
func (c *TagController) List(w http.ResponseWriter, r *http.Request) {
	// 检查是否有静态页面
	if c.config.Site.StaticTag && r.URL.Query().Get("upcache") == "" {
		http.ServeFile(w, r, "tags.html")
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取标签列表
	pageSize := 200
	tags, total, err := c.tagModel.GetList(page, pageSize)
	if err != nil {
		logger.Error("获取标签列表失败", "error", err)
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
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

	// 获取热门标签
	hotTags, err := c.tagModel.GetHotTags(50)
	if err != nil {
		logger.Error("获取热门标签失败", "error", err)
	}

	// 获取最新标签
	newTags, err := c.tagModel.GetNewTags(30)
	if err != nil {
		logger.Error("获取最新标签失败", "error", err)
	}

	// 获取热门文章
	hotArticles, err := c.articleService.GetHotArticles(10)
	if err != nil {
		logger.Error("获取热门文章失败", "error", err)
	}

	// 获取最新文章
	latestArticles, err := c.articleService.GetLatestArticles(10)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":       globals,
		"Tags":          tags,
		"HotTags":       hotTags,
		"NewTags":       newTags,
		"HotArticles":   hotArticles,
		"LatestArticles": latestArticles,
		"Pagination":    pagination,
		"PageTitle":     "标签云 - " + c.config.Site.Name,
		"Keywords":      "标签,tag,标签云," + c.config.Site.Keywords,
		"Description":   "标签云 - " + c.config.Site.Description,
		"CurrentTime":   time.Now().Format("2006-01-02 15:04:05"),
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/taglist.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticTag {
		go func() {
			if err := c.templateService.GenerateStaticPage(tplFile, data, "tags.html"); err != nil {
				logger.Error("生成静态标签列表页失败", "error", err)
			}
		}()
	}
}

// Detail 标签详情
func (c *TagController) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName := vars["tag"]
	if tagName == "" {
		http.Error(w, "Invalid tag name", http.StatusBadRequest)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 检查是否有静态页面
	if c.config.Site.StaticTag && page == 1 && r.URL.Query().Get("upcache") == "" {
		staticFile := "tag/" + tagName + ".html"
		http.ServeFile(w, r, staticFile)
		return
	}

	// 获取标签信息
	tag, err := c.tagModel.GetByName(tagName)
	if err != nil {
		logger.Error("获取标签信息失败", "tag", tagName, "error", err)
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	// 获取标签相关文章
	pageSize := 10
	articles, total, err := c.articleService.GetArticlesByTag(tagName, page, pageSize)
	if err != nil {
		logger.Error("获取标签相关文章失败", "tag", tagName, "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
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

	// 获取热门标签
	hotTags, err := c.tagModel.GetHotTags(30)
	if err != nil {
		logger.Error("获取热门标签失败", "error", err)
	}

	// 获取最新标签
	newTags, err := c.tagModel.GetNewTags(30)
	if err != nil {
		logger.Error("获取最新标签失败", "error", err)
	}

	// 获取热门文章
	hotArticles, err := c.articleService.GetHotArticles(10)
	if err != nil {
		logger.Error("获取热门文章失败", "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Tag":         tag,
		"Articles":    articles,
		"HotTags":     hotTags,
		"NewTags":     newTags,
		"HotArticles": hotArticles,
		"Pagination":  pagination,
		"PageTitle":   "标签: " + tagName + " - " + c.config.Site.Name,
		"Keywords":    tagName + "," + c.config.Site.Keywords,
		"Description": "标签: " + tagName + " - " + c.config.Site.Description,
		"CurrentTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/tag.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticTag && page == 1 {
		go func() {
			staticFile := "tag/" + tagName + ".html"
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticFile); err != nil {
				logger.Error("生成静态标签详情页失败", "error", err)
			}
		}()
	}
}
