package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// HtmlController 静态页面生成控制器
type HtmlController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	tagModel        *model.TagModel
	specialModel    *model.SpecialModel
	htmlService     *service.HtmlService
	templateService *service.TemplateService
}

// NewHtmlController 创建静态页面生成控制器
func NewHtmlController(db *database.DB, cache cache.Cache, config *config.Config) *HtmlController {
	return &HtmlController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		tagModel:        model.NewTagModel(db),
		specialModel:    model.NewSpecialModel(db),
		htmlService:     service.NewHtmlService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 首页静态化
func (c *HtmlController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "html",
		"PageTitle":   "首页静态化",
	}

	// 渲染模板
	tplFile := "admin/html_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染首页静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoIndex 处理首页静态化
func (c *HtmlController) DoIndex(w http.ResponseWriter, r *http.Request) {
	// 检查是否启用静态化
	if !c.config.Site.StaticIndex {
		http.Error(w, "Static index is disabled", http.StatusBadRequest)
		return
	}

	// 生成首页静态文件
	err := c.htmlService.GenerateIndex()
	if err != nil {
		logger.Error("生成首页静态文件失败", "error", err)
		http.Error(w, "Failed to generate index", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "首页静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/index", http.StatusFound)
	}
}

// List 栏目页静态化
func (c *HtmlController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"Categories":  categories,
		"CurrentMenu": "html",
		"PageTitle":   "栏目页静态化",
	}

	// 渲染模板
	tplFile := "admin/html_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目页静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoList 处理栏目页静态化
func (c *HtmlController) DoList(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	typeIDStr := r.FormValue("typeid")
	startPageStr := r.FormValue("startpage")
	endPageStr := r.FormValue("endpage")

	// 检查是否启用静态化
	if !c.config.Site.StaticList {
		http.Error(w, "Static list is disabled", http.StatusBadRequest)
		return
	}

	// 解析参数
	typeID := int64(0)
	if typeIDStr != "" {
		typeID, _ = strconv.ParseInt(typeIDStr, 10, 64)
	}

	startPage := 1
	if startPageStr != "" {
		startPage, _ = strconv.Atoi(startPageStr)
		if startPage < 1 {
			startPage = 1
		}
	}

	endPage := 1
	if endPageStr != "" {
		endPage, _ = strconv.Atoi(endPageStr)
		if endPage < startPage {
			endPage = startPage
		}
	}

	// 生成栏目页静态文件
	if typeID > 0 {
		// 生成指定栏目
		err := c.htmlService.GenerateList(typeID)
		if err != nil {
			logger.Error("生成栏目页静态文件失败", "typeid", typeID, "error", err)
			http.Error(w, "Failed to generate list", http.StatusInternalServerError)
			return
		}
	} else {
		// 生成所有栏目
		categories, err := c.categoryModel.GetAll()
		if err != nil {
			logger.Error("获取栏目列表失败", "error", err)
			http.Error(w, "Failed to get categories", http.StatusInternalServerError)
			return
		}

		for _, category := range categories {
			err := c.htmlService.GenerateList(category.ID)
			if err != nil {
				logger.Error("生成栏目页静态文件失败", "typeid", category.ID, "error", err)
			}
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "栏目页静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/list", http.StatusFound)
	}
}

// Article 文章页静态化
func (c *HtmlController) Article(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"Categories":  categories,
		"CurrentMenu": "html",
		"PageTitle":   "文章页静态化",
	}

	// 渲染模板
	tplFile := "admin/html_article.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章页静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoArticle 处理文章页静态化
func (c *HtmlController) DoArticle(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	typeIDStr := r.FormValue("typeid")
	startIDStr := r.FormValue("startid")
	endIDStr := r.FormValue("endid")
	startDateStr := r.FormValue("startdate")
	endDateStr := r.FormValue("enddate")

	// 检查是否启用静态化
	if !c.config.Site.StaticArticle {
		http.Error(w, "Static article is disabled", http.StatusBadRequest)
		return
	}

	// 解析参数
	typeID := int64(0)
	if typeIDStr != "" {
		typeID, _ = strconv.ParseInt(typeIDStr, 10, 64)
	}

	startID := int64(0)
	if startIDStr != "" {
		startID, _ = strconv.ParseInt(startIDStr, 10, 64)
	}

	endID := int64(0)
	if endIDStr != "" {
		endID, _ = strconv.ParseInt(endIDStr, 10, 64)
	}

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	// 生成文章页静态文件
	if startID > 0 && endID > 0 {
		// 按ID范围生成
		successCount := 0
		for id := startID; id <= endID; id++ {
			err := c.htmlService.GenerateArticle(id)
			if err != nil {
				logger.Error("生成文章页静态文件失败", "id", id, "error", err)
			} else {
				successCount++
			}
		}
		logger.Info("按ID范围生成文章页完成", "startID", startID, "endID", endID, "successCount", successCount)
	} else if !startDate.IsZero() && !endDate.IsZero() {
		// 按日期范围生成
		articles, err := c.articleModel.GetByDateRange(typeID, startDate, endDate)
		if err != nil {
			logger.Error("获取文章列表失败", "error", err)
			http.Error(w, "Failed to get articles", http.StatusInternalServerError)
			return
		}

		for _, article := range articles {
			err := c.htmlService.GenerateArticle(article.ID)
			if err != nil {
				logger.Error("生成文章页静态文件失败", "id", article.ID, "error", err)
			}
		}
	} else if typeID > 0 {
		// 按栏目生成
		articles, _, err := c.articleModel.GetByTypeID(typeID, 1, 1000)
		if err != nil {
			logger.Error("获取文章列表失败", "error", err)
			http.Error(w, "Failed to get articles", http.StatusInternalServerError)
			return
		}

		for _, article := range articles {
			err := c.htmlService.GenerateArticle(article.ID)
			if err != nil {
				logger.Error("生成文章页静态文件失败", "id", article.ID, "error", err)
			}
		}
	} else {
		// 生成所有文章
		articles, _, err := c.articleModel.GetAll(1, 1000)
		if err != nil {
			logger.Error("获取文章列表失败", "error", err)
			http.Error(w, "Failed to get articles", http.StatusInternalServerError)
			return
		}

		for _, article := range articles {
			err := c.htmlService.GenerateArticle(article.ID)
			if err != nil {
				logger.Error("生成文章页静态文件失败", "id", article.ID, "error", err)
			}
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "文章页静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/article", http.StatusFound)
	}
}

// Special 专题页静态化
func (c *HtmlController) Special(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取专题列表
	specials, err := c.specialModel.GetAll()
	if err != nil {
		logger.Error("获取专题列表失败", "error", err)
		http.Error(w, "Failed to get specials", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"Specials":    specials,
		"CurrentMenu": "html",
		"PageTitle":   "专题页静态化",
	}

	// 渲染模板
	tplFile := "admin/html_special.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染专题页静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoSpecial 处理专题页静态化
func (c *HtmlController) DoSpecial(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	specialIDStr := r.FormValue("specialid")

	// 检查是否启用静态化
	if !c.config.Site.StaticSpecial {
		http.Error(w, "Static special is disabled", http.StatusBadRequest)
		return
	}

	// 解析参数
	specialID := int64(0)
	if specialIDStr != "" {
		specialID, _ = strconv.ParseInt(specialIDStr, 10, 64)
	}

	// 生成专题页静态文件
	if specialID > 0 {
		// 生成指定专题
		err := c.htmlService.GenerateSpecial(specialID)
		if err != nil {
			logger.Error("生成专题页静态文件失败", "specialid", specialID, "error", err)
			http.Error(w, "Failed to generate special", http.StatusInternalServerError)
			return
		}
	} else {
		// 生成所有专题
		specials, err := c.specialModel.GetAll()
		if err != nil {
			logger.Error("获取专题列表失败", "error", err)
			http.Error(w, "Failed to get specials", http.StatusInternalServerError)
			return
		}

		for _, special := range specials {
			err := c.htmlService.GenerateSpecial(special.ID)
			if err != nil {
				logger.Error("生成专题页静态文件失败", "specialid", special.ID, "error", err)
			}
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "专题页静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/special", http.StatusFound)
	}
}

// Tag 标签页静态化
func (c *HtmlController) Tag(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取标签列表
	tags, err := c.tagModel.GetHotTags(100)
	if err != nil {
		logger.Error("获取标签列表失败", "error", err)
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"Tags":        tags,
		"CurrentMenu": "html",
		"PageTitle":   "标签页静态化",
	}

	// 渲染模板
	tplFile := "admin/html_tag.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签页静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoTag 处理标签页静态化
func (c *HtmlController) DoTag(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	tagName := r.FormValue("tagname")

	// 检查是否启用静态化
	if !c.config.Site.StaticTag {
		http.Error(w, "Static tag is disabled", http.StatusBadRequest)
		return
	}

	// 生成标签页静态文件
	if tagName != "" {
		// 生成指定标签
		err := c.htmlService.GenerateTag(tagName)
		if err != nil {
			logger.Error("生成标签页静态文件失败", "tagname", tagName, "error", err)
			http.Error(w, "Failed to generate tag", http.StatusInternalServerError)
			return
		}
	} else {
		// 生成所有标签
		tags, err := c.tagModel.GetHotTags(100)
		if err != nil {
			logger.Error("获取标签列表失败", "error", err)
			http.Error(w, "Failed to get tags", http.StatusInternalServerError)
			return
		}

		for _, tag := range tags {
			err := c.htmlService.GenerateTag(tag.Tag)
			if err != nil {
				logger.Error("生成标签页静态文件失败", "tagname", tag.Tag, "error", err)
			}
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "标签页静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/tag", http.StatusFound)
	}
}

// All 全站静态化
func (c *HtmlController) All(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "html",
		"PageTitle":   "全站静态化",
	}

	// 渲染模板
	tplFile := "admin/html_all.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染全站静态化模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAll 处理全站静态化
func (c *HtmlController) DoAll(w http.ResponseWriter, r *http.Request) {
	// 生成首页静态文件
	if c.config.Site.StaticIndex {
		err := c.htmlService.GenerateIndex()
		if err != nil {
			logger.Error("生成首页静态文件失败", "error", err)
		}
	}

	// 生成栏目页静态文件
	if c.config.Site.StaticList {
		categories, err := c.categoryModel.GetAll()
		if err != nil {
			logger.Error("获取栏目列表失败", "error", err)
		} else {
			for _, category := range categories {
				err := c.htmlService.GenerateList(category.ID)
				if err != nil {
					logger.Error("生成栏目页静态文件失败", "typeid", category.ID, "error", err)
				}
			}
		}
	}

	// 生成文章页静态文件
	if c.config.Site.StaticArticle {
		articles, _, err := c.articleModel.GetAll(1, 1000)
		if err != nil {
			logger.Error("获取文章列表失败", "error", err)
		} else {
			for _, article := range articles {
				err := c.htmlService.GenerateArticle(article.ID)
				if err != nil {
					logger.Error("生成文章页静态文件失败", "id", article.ID, "error", err)
				}
			}
		}
	}

	// 生成专题页静态文件
	if c.config.Site.StaticSpecial {
		specials, err := c.specialModel.GetAll()
		if err != nil {
			logger.Error("获取专题列表失败", "error", err)
		} else {
			for _, special := range specials {
				err := c.htmlService.GenerateSpecial(special.ID)
				if err != nil {
					logger.Error("生成专题页静态文件失败", "specialid", special.ID, "error", err)
				}
			}
		}
	}

	// 生成标签页静态文件
	if c.config.Site.StaticTag {
		tags, err := c.tagModel.GetHotTags(100)
		if err != nil {
			logger.Error("获取标签列表失败", "error", err)
		} else {
			for _, tag := range tags {
				err := c.htmlService.GenerateTag(tag.Tag)
				if err != nil {
					logger.Error("生成标签页静态文件失败", "tagname", tag.Tag, "error", err)
				}
			}
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "全站静态化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/html/all", http.StatusFound)
	}
}
