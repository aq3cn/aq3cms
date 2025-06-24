package front

import (
	"net/http"
	"strconv"

	"aq3cms/config"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// IndexController 前台首页控制器
type IndexController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleService  *service.ArticleService
	categoryService *service.CategoryService
	adService       *service.AdService
	templateService *service.TemplateService
}

// NewIndexController 创建前台首页控制器
func NewIndexController(db *database.DB, cache cache.Cache, config *config.Config) *IndexController {
	return &IndexController{
		db:              db,
		cache:           cache,
		config:          config,
		articleService:  service.NewArticleService(db, cache, config),
		categoryService: service.NewCategoryService(db, cache, config),
		adService:       service.NewAdService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 首页
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取站点配置
	siteConfig := c.config.Site

	// 获取首页推荐文章
	recommendArticles, err := c.articleService.GetRecommendArticles(10)
	if err != nil {
		logger.Error("获取首页推荐文章失败", "error", err)
	}

	// 获取首页最新文章
	latestArticles, err := c.articleService.GetLatestArticles(10)
	if err != nil {
		logger.Error("获取首页最新文章失败", "error", err)
	}

	// 获取首页热门文章
	hotArticles, err := c.articleService.GetHotArticles(10)
	if err != nil {
		logger.Error("获取首页热门文章失败", "error", err)
	}

	// 获取顶级栏目
	topCategories, err := c.categoryService.GetTopCategories()
	if err != nil {
		logger.Error("获取顶级栏目失败", "error", err)
	}

	// 获取首页广告
	var topBanner string
	var sideBanner string
	var footerBanner string

	// 获取顶部广告
	topBannerHTML, err := c.adService.GetPositionHTMLByCode("index_top")
	if err == nil {
		topBanner = topBannerHTML
	}

	// 获取侧边栏广告
	sideBannerHTML, err := c.adService.GetPositionHTMLByCode("index_side")
	if err == nil {
		sideBanner = sideBannerHTML
	}

	// 获取底部广告
	footerBannerHTML, err := c.adService.GetPositionHTMLByCode("index_footer")
	if err == nil {
		footerBanner = footerBannerHTML
	}

	// 准备模板数据
	data := map[string]interface{}{
		"SiteConfig":        siteConfig,
		"RecommendArticles": recommendArticles,
		"LatestArticles":    latestArticles,
		"HotArticles":       hotArticles,
		"TopCategories":     topCategories,
		"TopBanner":         topBanner,
		"SideBanner":        sideBanner,
		"FooterBanner":      footerBanner,
		"PageTitle":         "首页",
	}

	// 渲染模板
	tplFile := "index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Search 搜索
func (c *IndexController) Search(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	query := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 搜索文章
	articles, total, err := c.articleService.SearchArticles(query, page, 20)
	if err != nil {
		logger.Error("搜索文章失败", "error", err)
		http.Error(w, "Failed to search articles", http.StatusInternalServerError)
		return
	}

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
		"Query":      query,
		"Articles":   articles,
		"Pagination": pagination,
		"PageTitle":  "搜索结果: " + query,
	}

	// 渲染模板
	tplFile := "search.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染搜索结果模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Tag 标签页
func (c *IndexController) Tag(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	tag := r.URL.Query().Get("tag")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 获取标签文章
	articles, total, err := c.articleService.GetArticlesByTag(tag, page, 20)
	if err != nil {
		logger.Error("获取标签文章失败", "error", err)
		http.Error(w, "Failed to get tag articles", http.StatusInternalServerError)
		return
	}

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
		"Tag":        tag,
		"Articles":   articles,
		"Pagination": pagination,
		"PageTitle":  "标签: " + tag,
	}

	// 渲染模板
	tplFile := "tag.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Sitemap 站点地图
func (c *IndexController) Sitemap(w http.ResponseWriter, r *http.Request) {
	// 获取所有栏目
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 获取最新文章
	articles, err := c.articleService.GetLatestArticles(100)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		http.Error(w, "Failed to get latest articles", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Categories": categories,
		"Articles":   articles,
		"PageTitle":  "站点地图",
	}

	// 渲染模板
	tplFile := "sitemap.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染站点地图模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// SitemapXML 站点地图XML
func (c *IndexController) SitemapXML(w http.ResponseWriter, r *http.Request) {
	// 获取所有栏目
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 获取最新文章
	articles, err := c.articleService.GetLatestArticles(1000)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		http.Error(w, "Failed to get latest articles", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Categories": categories,
		"Articles":   articles,
		"SiteURL":    c.config.Site.URL,
	}

	// 设置Content-Type
	w.Header().Set("Content-Type", "application/xml")

	// 渲染模板
	tplFile := "sitemap.xml"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染站点地图XML模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RSS RSS订阅
func (c *IndexController) RSS(w http.ResponseWriter, r *http.Request) {
	// 获取最新文章
	articles, err := c.articleService.GetLatestArticles(50)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		http.Error(w, "Failed to get latest articles", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Articles": articles,
		"SiteURL":  c.config.Site.URL,
		"SiteTitle": "网站标题",
		"SiteDescription": "网站描述",
	}

	// 设置Content-Type
	w.Header().Set("Content-Type", "application/xml")

	// 渲染模板
	tplFile := "rss.xml"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染RSS模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
