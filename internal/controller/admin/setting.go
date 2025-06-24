package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SettingController 设置控制器
type SettingController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	templateService *service.TemplateService
}

// NewSettingController 创建设置控制器
func NewSettingController(db *database.DB, cache cache.Cache, config *config.Config) *SettingController {
	return &SettingController{
		db:              db,
		cache:           cache,
		config:          config,
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 设置首页
func (c *SettingController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "setting",
		"PageTitle":   "系统设置",
	}

	// 渲染模板
	tplFile := "admin/setting_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染设置首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Basic 基本设置
func (c *SettingController) Basic(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "setting",
		"PageTitle":   "基本设置",
	}

	// 检查是否有成功消息
	if r.URL.Query().Get("success") == "1" {
		data["Message"] = "基本设置保存成功！"
		data["MessageType"] = "success"
	}

	// 渲染模板
	tplFile := "admin/setting_basic.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染基本设置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoBasic 处理基本设置
func (c *SettingController) DoBasic(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	siteName := r.FormValue("site_name")
	siteURL := r.FormValue("site_url")
	siteKeywords := r.FormValue("site_keywords")
	siteDescription := r.FormValue("site_description")
	siteICP := r.FormValue("site_icp")
	siteCopyright := r.FormValue("site_copyright")
	siteStatCode := r.FormValue("site_statcode")
	defaultTpl := r.FormValue("default_template")
	timezone := r.FormValue("timezone")

	// 更新配置
	c.config.Site.Name = siteName
	c.config.Site.URL = siteURL
	c.config.Site.Keywords = siteKeywords
	c.config.Site.Description = siteDescription
	c.config.Site.ICP = siteICP
	c.config.Site.CopyRight = siteCopyright
	c.config.Site.StatCode = siteStatCode
	c.config.Template.DefaultTpl = defaultTpl

	// 时区设置可以在这里处理，暂时忽略
	_ = timezone

	// 保存配置
	err := c.config.Save()
	if err != nil {
		logger.Error("保存配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Clear()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "基本设置保存成功",
		})
	} else {
		// 普通表单提交，重定向到设置页面并显示成功消息
		http.Redirect(w, r, "/aq3cms/setting/basic?success=1", http.StatusFound)
	}
}

// Upload 上传设置
func (c *SettingController) Upload(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "setting",
		"PageTitle":   "上传设置",
	}

	// 检查是否有成功消息
	if r.URL.Query().Get("success") == "1" {
		data["Message"] = "上传设置保存成功！"
		data["MessageType"] = "success"
	}

	// 渲染模板
	tplFile := "admin/setting_upload.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染上传设置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoUpload 处理上传设置
func (c *SettingController) DoUpload(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	uploadDir := r.FormValue("upload_dir")
	uploadMaxSizeStr := r.FormValue("max_file_size")
	uploadAllowedExts := r.FormValue("allowed_types")
	autoRename := r.FormValue("auto_rename")
	createThumbnail := r.FormValue("create_thumbnail")
	imageQualityStr := r.FormValue("image_quality")
	thumbnailSize := r.FormValue("thumbnail_size")
	watermarkText := r.FormValue("watermark_text")
	watermarkPosition := r.FormValue("watermark_position")
	watermarkOpacityStr := r.FormValue("watermark_opacity")

	// 更新配置
	c.config.Upload.Dir = uploadDir
	c.config.Upload.AllowedExts = uploadAllowedExts
	c.config.Upload.WatermarkText = watermarkText
	c.config.Upload.WatermarkPos = watermarkPosition

	// 解析整数值
	uploadMaxSize, _ := strconv.Atoi(uploadMaxSizeStr)
	c.config.Upload.MaxSize = uploadMaxSize

	imageQuality, _ := strconv.Atoi(imageQualityStr)
	_ = imageQuality // 暂时忽略，配置结构中没有这个字段

	watermarkOpacity, _ := strconv.Atoi(watermarkOpacityStr)
	c.config.Upload.WatermarkAlpha = float64(watermarkOpacity) / 100.0

	// 解析布尔值
	if autoRename == "1" {
		// 自动重命名功能，暂时忽略
		_ = autoRename
	}

	if createThumbnail == "1" {
		// 自动生成缩略图功能，暂时忽略
		_ = createThumbnail
	}

	// 缩略图尺寸处理，暂时忽略
	_ = thumbnailSize

	// 保存配置
	err := c.config.Save()
	if err != nil {
		logger.Error("保存配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "上传设置保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/setting/upload?success=1", http.StatusFound)
	}
}

// Static 静态化设置
func (c *SettingController) Static(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "setting",
		"PageTitle":   "静态化设置",
	}

	// 检查是否有成功消息
	if r.URL.Query().Get("success") == "1" {
		data["Message"] = "静态化设置保存成功！"
		data["MessageType"] = "success"
	}

	// 渲染模板
	tplFile := "admin/setting_static.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染静态化设置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoStatic 处理静态化设置
func (c *SettingController) DoStatic(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	staticIndexStr := r.FormValue("static_index")
	staticListStr := r.FormValue("static_list")
	staticArticleStr := r.FormValue("static_article")
	staticSpecialStr := r.FormValue("static_special")
	staticTagStr := r.FormValue("static_tag")
	staticMobileStr := r.FormValue("static_mobile")
	staticDir := r.FormValue("static_dir")
	staticSuffix := r.FormValue("static_suffix")
	staticMobileDir := r.FormValue("static_mobile_dir")
	staticMobileSuffix := r.FormValue("static_mobile_suffix")

	// 更新配置
	c.config.Site.StaticDir = staticDir
	c.config.Site.StaticSuffix = staticSuffix
	c.config.Site.StaticMobileDir = staticMobileDir
	c.config.Site.StaticMobileSuffix = staticMobileSuffix

	// 解析布尔值
	if staticIndexStr == "1" {
		c.config.Site.StaticIndex = true
	} else {
		c.config.Site.StaticIndex = false
	}

	if staticListStr == "1" {
		c.config.Site.StaticList = true
	} else {
		c.config.Site.StaticList = false
	}

	if staticArticleStr == "1" {
		c.config.Site.StaticArticle = true
	} else {
		c.config.Site.StaticArticle = false
	}

	if staticSpecialStr == "1" {
		c.config.Site.StaticSpecial = true
	} else {
		c.config.Site.StaticSpecial = false
	}

	if staticTagStr == "1" {
		c.config.Site.StaticTag = true
	} else {
		c.config.Site.StaticTag = false
	}

	if staticMobileStr == "1" {
		c.config.Site.StaticMobile = true
	} else {
		c.config.Site.StaticMobile = false
	}

	// 保存配置
	err := c.config.Save()
	if err != nil {
		logger.Error("保存配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "静态化设置保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/setting/static?success=1", http.StatusFound)
	}
}

// Cache 缓存设置
func (c *SettingController) Cache(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "setting",
		"PageTitle":   "缓存设置",
	}

	// 检查是否有成功消息
	successType := r.URL.Query().Get("success")
	if successType != "" {
		switch successType {
		case "1":
			data["Message"] = "缓存设置保存成功！"
			data["MessageType"] = "success"
		case "clear":
			data["Message"] = "缓存清除成功！"
			data["MessageType"] = "success"
		}
	}

	// 渲染模板
	tplFile := "admin/setting_cache.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染缓存设置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoCache 处理缓存设置
func (c *SettingController) DoCache(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	cacheTypeStr := r.FormValue("cache_type")
	cacheTimeStr := r.FormValue("cache_time")
	cacheIndexStr := r.FormValue("cache_index")
	cacheListStr := r.FormValue("cache_list")
	cacheArticleStr := r.FormValue("cache_article")
	cacheSearchStr := r.FormValue("cache_search")
	cacheSearchTimeStr := r.FormValue("cache_search_time")
	redisHost := r.FormValue("redis_host")
	redisPortStr := r.FormValue("redis_port")
	redisPassword := r.FormValue("redis_password")

	// 更新配置
	c.config.Cache.Type = cacheTypeStr
	c.config.Cache.Host = redisHost
	c.config.Cache.Password = redisPassword

	// 解析整数值
	cacheTime, _ := strconv.Atoi(cacheTimeStr)
	c.config.Cache.Expire = cacheTime

	cacheSearchTime, _ := strconv.Atoi(cacheSearchTimeStr)
	c.config.Cache.SearchCacheTime = cacheSearchTime

	redisPort, _ := strconv.Atoi(redisPortStr)
	c.config.Cache.Port = redisPort

	// 解析布尔值
	if cacheIndexStr == "1" {
		c.config.Cache.EnableArticleCache = true
	} else {
		c.config.Cache.EnableArticleCache = false
	}

	if cacheListStr == "1" {
		c.config.Cache.EnableListCache = true
	} else {
		c.config.Cache.EnableListCache = false
	}

	if cacheArticleStr == "1" {
		c.config.Cache.EnableArticleCache = true
	} else {
		c.config.Cache.EnableArticleCache = false
	}

	if cacheSearchStr == "1" {
		c.config.Cache.EnableSearchCache = true
	} else {
		c.config.Cache.EnableSearchCache = false
	}

	// 保存配置
	err := c.config.Save()
	if err != nil {
		logger.Error("保存配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Clear()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "缓存设置保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/setting/cache?success=1", http.StatusFound)
	}
}

// ClearCache 清除缓存
func (c *SettingController) ClearCache(w http.ResponseWriter, r *http.Request) {
	// 清除缓存
	c.cache.Clear()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "缓存清除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/setting/cache?success=clear", http.StatusFound)
	}
}
