package frontend

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

// SpecialController 专题控制器
type SpecialController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	specialModel    *model.SpecialModel
	articleModel    *model.ArticleModel
	templateService *service.TemplateService
}

// NewSpecialController 创建专题控制器
func NewSpecialController(db *database.DB, cache cache.Cache, config *config.Config) *SpecialController {
	return &SpecialController{
		db:              db,
		cache:           cache,
		config:          config,
		specialModel:    model.NewSpecialModel(db),
		articleModel:    model.NewArticleModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 专题列表
func (c *SpecialController) List(w http.ResponseWriter, r *http.Request) {
	// 检查是否有静态专题列表页
	if c.config.Site.StaticSpecial && r.URL.Query().Get("upcache") == "" {
		http.ServeFile(w, r, "special/index.html")
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取专题列表
	pageSize := 10
	specials, err := c.specialModel.GetAll()

	// 计算总数
	total := len(specials)

	// 分页处理
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		start = 0
		end = 0
		specials = []*model.Special{}
	} else if end > total {
		end = total
		specials = specials[start:end]
	} else {
		specials = specials[start:end]
	}
	if err != nil {
		logger.Error("获取专题列表失败", "error", err)
		http.Error(w, "Failed to get specials", http.StatusInternalServerError)
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

	// 获取热门专题
	hotSpecials, err := c.specialModel.GetHotSpecials(5)
	if err != nil {
		logger.Error("获取热门专题失败", "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Specials":    specials,
		"HotSpecials": hotSpecials,
		"Pagination":  pagination,
		"PageTitle":   "专题列表 - " + c.config.Site.Name,
		"Keywords":    "专题,专题列表," + c.config.Site.Keywords,
		"Description": "专题列表 - " + c.config.Site.Description,
	}

	// 确定模板文件
	tplFile := c.config.Template.DefaultTpl + "/specials.htm"

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染专题列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticSpecial {
		go func() {
			staticPath := "special/index.html"
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticPath); err != nil {
				logger.Error("生成静态专题列表页失败", "error", err)
			}
		}()
	}
}

// Detail 专题详情
func (c *SpecialController) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	if filename == "" {
		http.Error(w, "Invalid special filename", http.StatusBadRequest)
		return
	}

	// 检查是否有静态专题页
	if c.config.Site.StaticSpecial && r.URL.Query().Get("upcache") == "" {
		staticPath := "special/" + filename + ".html"
		if _, err := http.Dir(".").Open(staticPath); err == nil {
			http.ServeFile(w, r, staticPath)
			return
		}
	}

	// 获取专题信息
	special, err := c.specialModel.GetByFilename(filename)
	if err != nil {
		logger.Error("获取专题信息失败", "filename", filename, "error", err)
		http.Error(w, "Special not found", http.StatusNotFound)
		return
	}

	// 增加点击量
	go c.specialModel.IncrementClick(special.ID)

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取专题文章
	pageSize := 10
	articles, total, err := c.specialModel.GetArticles(special.ID, page, pageSize)
	if err != nil {
		logger.Error("获取专题文章失败", "id", special.ID, "error", err)
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

	// 获取热门专题
	hotSpecials, err := c.specialModel.GetHotSpecials(5)
	if err != nil {
		logger.Error("获取热门专题失败", "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Special":     special,
		"Articles":    articles,
		"HotSpecials": hotSpecials,
		"Pagination":  pagination,
		"PageTitle":   special.Title + " - " + c.config.Site.Name,
		"Keywords":    special.Keywords,
		"Description": special.Description,
	}

	// 确定模板文件
	var tplFile string
	if special.Template != "" {
		tplFile = special.Template
	} else {
		tplFile = c.config.Template.DefaultTpl + "/special.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染专题详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 如果需要生成静态页面
	if c.config.Site.StaticSpecial {
		go func() {
			staticPath := "special/" + filename + ".html"
			if err := c.templateService.GenerateStaticPage(tplFile, data, staticPath); err != nil {
				logger.Error("生成静态专题页失败", "filename", filename, "error", err)
			}
		}()
	}
}
