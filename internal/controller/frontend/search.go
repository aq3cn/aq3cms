package frontend

import (
	"net/http"
	"strconv"
	"strings"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SearchController 搜索控制器
type SearchController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	searchService   *service.SearchService
	categoryModel   *model.CategoryModel
	templateService *service.TemplateService
}

// NewSearchController 创建搜索控制器
func NewSearchController(db *database.DB, cache cache.Cache, config *config.Config) *SearchController {
	return &SearchController{
		db:              db,
		cache:           cache,
		config:          config,
		searchService:   service.NewSearchService(db, cache, config),
		categoryModel:   model.NewCategoryModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Search 搜索
func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	// 获取搜索参数
	keyword := r.URL.Query().Get("keyword")
	channelTypeStr := r.URL.Query().Get("channeltype")
	typeIDStr := r.URL.Query().Get("typeid")
	orderBy := r.URL.Query().Get("orderby")
	timeRange := r.URL.Query().Get("timerange")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pagesize")
	fieldsStr := r.URL.Query().Get("fields")

	// 解析参数
	channelType := 0
	if channelTypeStr != "" {
		channelType, _ = strconv.Atoi(channelTypeStr)
	}

	typeID := int64(0)
	if typeIDStr != "" {
		typeID, _ = strconv.ParseInt(typeIDStr, 10, 64)
	}

	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
		if pageSize < 1 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
	}

	// 解析搜索字段
	fields := []string{"a.title", "a.keywords", "a.description"}
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, ",")
	}

	// 创建搜索选项
	options := &service.SearchOptions{
		Keyword:     keyword,
		ChannelType: channelType,
		TypeID:      typeID,
		OrderBy:     orderBy,
		TimeRange:   timeRange,
		Page:        page,
		PageSize:    pageSize,
		Fields:      fields,
	}

	// 执行搜索
	results, total, err := c.searchService.Search(options)
	if err != nil {
		logger.Error("搜索失败", "keyword", keyword, "error", err)
		http.Error(w, "Failed to search", http.StatusInternalServerError)
		return
	}

	// 获取栏目信息
	var category *model.Category
	if typeID > 0 {
		category, err = c.categoryModel.GetByID(typeID)
		if err != nil {
			logger.Error("获取栏目信息失败", "typeid", typeID, "error", err)
		}
	}

	// 获取所有栏目
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
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
		"Results":     results,
		"Pagination":  pagination,
		"Keyword":     keyword,
		"ChannelType": channelType,
		"TypeID":      typeID,
		"OrderBy":     orderBy,
		"TimeRange":   timeRange,
		"Category":    category,
		"Categories":  categories,
		"PageTitle":   "搜索: " + keyword + " - " + c.config.Site.Name,
		"Keywords":    keyword + "," + c.config.Site.Keywords,
		"Description": "搜索: " + keyword + " - " + c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/search.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染搜索模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AdvancedSearch 高级搜索
func (c *SearchController) AdvancedSearch(w http.ResponseWriter, r *http.Request) {
	// 获取所有栏目
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":    globals,
		"Categories": categories,
		"PageTitle":  "高级搜索 - " + c.config.Site.Name,
		"Keywords":   "高级搜索," + c.config.Site.Keywords,
		"Description": "高级搜索 - " + c.config.Site.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/search_advanced.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染高级搜索模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
