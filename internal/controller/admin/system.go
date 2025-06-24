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

// SystemController 系统设置控制器
type SystemController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	templateService *service.TemplateService
	adminService    *service.AdminService
}

// NewSystemController 创建系统设置控制器
func NewSystemController(db *database.DB, cache cache.Cache, config *config.Config) *SystemController {
	return &SystemController{
		db:              db,
		cache:           cache,
		config:          config,
		templateService: service.NewTemplateService(db, cache, config),
		adminService:    service.NewAdminService(db, cache, config),
	}
}

// Config 系统配置
func (c *SystemController) Config(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "system",
		"SubMenu":     "config",
		"PageTitle":   "系统配置",
	}

	// 检查是否有成功消息
	if r.URL.Query().Get("success") == "1" {
		data["Message"] = "系统配置保存成功！"
		data["MessageType"] = "success"
	}

	// 渲染模板
	tplFile := "admin/system_config.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染系统配置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// SaveConfig 保存系统配置
func (c *SystemController) SaveConfig(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	siteTitle := r.FormValue("site_title")
	siteKeywords := r.FormValue("site_keywords")
	siteDescription := r.FormValue("site_description")
	siteURL := r.FormValue("site_url")
	siteICP := r.FormValue("site_icp")
	siteCopyright := r.FormValue("site_copyright")
	siteStatistics := r.FormValue("site_statistics")

	// 更新配置
	c.config.Site.Name = siteTitle
	c.config.Site.Keywords = siteKeywords
	c.config.Site.Description = siteDescription
	c.config.Site.URL = siteURL
	c.config.Site.ICP = siteICP
	c.config.Site.CopyRight = siteCopyright
	c.config.Site.StatCode = siteStatistics

	// 保存配置
	if err := c.config.Save(); err != nil {
		logger.Error("保存系统配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "系统配置保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/system_config?success=1", http.StatusFound)
	}
}

// Database 数据库管理
func (c *SystemController) Database(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "system",
		"SubMenu":     "database",
		"PageTitle":   "数据库管理",
	}

	// 检查是否有成功消息
	successType := r.URL.Query().Get("success")
	if successType != "" {
		switch successType {
		case "backup":
			data["Message"] = "数据库备份成功！"
			data["MessageType"] = "success"
		case "restore":
			data["Message"] = "数据库恢复成功！"
			data["MessageType"] = "success"
		}
	}

	// 渲染模板
	tplFile := "admin/system_database.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染数据库管理模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Backup 数据库备份
func (c *SystemController) Backup(w http.ResponseWriter, r *http.Request) {
	// 执行数据库备份
	backupFile, err := c.db.Backup()
	if err != nil {
		logger.Error("数据库备份失败", "error", err)
		http.Error(w, "Failed to backup database", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "数据库备份成功",
			"file":    backupFile,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/system_database?success=backup", http.StatusFound)
	}
}

// Restore 数据库恢复
func (c *SystemController) Restore(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取上传文件
	file, _, err := r.FormFile("backup_file")
	if err != nil {
		http.Error(w, "Invalid backup file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 执行数据库恢复
	err = c.db.Restore(file)
	if err != nil {
		logger.Error("数据库恢复失败", "error", err)
		http.Error(w, "Failed to restore database", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "数据库恢复成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/system_database?success=restore", http.StatusFound)
	}
}

// Cache 缓存管理
func (c *SystemController) Cache(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Config":      c.config,
		"CurrentMenu": "system",
		"SubMenu":     "cache",
		"PageTitle":   "缓存管理",
	}

	// 渲染模板
	tplFile := "admin/system_cache.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染缓存管理模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ClearCache 清除缓存
func (c *SystemController) ClearCache(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取缓存类型
	cacheType := r.FormValue("type")

	// 清除缓存
	switch cacheType {
	case "all":
		c.cache.Clear()
	case "article":
		c.cache.DeleteByPrefix("article:")
	case "category":
		c.cache.DeleteByPrefix("category:")
	case "member":
		c.cache.DeleteByPrefix("member:")
	case "template":
		c.cache.DeleteByPrefix("template:")
	case "static":
		c.cache.DeleteByPrefix("static:")
	default:
		http.Error(w, "Invalid cache type", http.StatusBadRequest)
		return
	}

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
		http.Redirect(w, r, "/admin/system_cache", http.StatusFound)
	}
}

// Log 日志管理
func (c *SystemController) Log(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	level := r.URL.Query().Get("level")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取日志列表
	logs, total, err := logger.GetLogs(level, page, 20)
	if err != nil {
		logger.Error("获取日志列表失败", "error", err)
		http.Error(w, "Failed to get logs", http.StatusInternalServerError)
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
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Logs":        logs,
		"Level":       level,
		"Pagination":  pagination,
		"CurrentMenu": "system",
		"SubMenu":     "log",
		"PageTitle":   "日志管理",
	}

	// 渲染模板
	tplFile := "admin/system_log.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染日志管理模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ClearLog 清除日志
func (c *SystemController) ClearLog(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取日志级别
	level := r.FormValue("level")

	// 清除日志
	err := logger.ClearLogs(level)
	if err != nil {
		logger.Error("清除日志失败", "error", err)
		http.Error(w, "Failed to clear logs", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "日志清除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/system_log", http.StatusFound)
	}
}
