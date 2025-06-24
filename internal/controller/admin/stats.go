package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// StatsController 统计控制器
type StatsController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	statsService    *service.StatsService
	templateService *service.TemplateService
}

// NewStatsController 创建统计控制器
func NewStatsController(db *database.DB, cache cache.Cache, config *config.Config) *StatsController {
	return &StatsController{
		db:              db,
		cache:           cache,
		config:          config,
		statsService:    service.NewStatsService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 统计首页
func (c *StatsController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取站点统计
	now := time.Now().Unix()
	startTime := now - 30*24*60*60 // 30天前
	siteStats, err := c.statsService.GetSiteStats(startTime, now)
	if err != nil {
		logger.Error("获取站点统计失败", "error", err)
		http.Error(w, "Failed to get site stats", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Stats":       siteStats,
		"CurrentMenu": "stats",
		"PageTitle":   "站点统计",
	}

	// 渲染模板
	tplFile := "admin/stats_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染统计首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Category 栏目统计
func (c *StatsController) Category(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目统计
	now := time.Now().Unix()
	startTime := now - 30*24*60*60 // 30天前
	categoryStats, err := c.statsService.GetCategoryStats(startTime, now)
	if err != nil {
		logger.Error("获取栏目统计失败", "error", err)
		http.Error(w, "Failed to get category stats", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Stats":       categoryStats,
		"CurrentMenu": "stats",
		"PageTitle":   "栏目统计",
	}

	// 渲染模板
	tplFile := "admin/stats_category.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染栏目统计模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Member 会员统计
func (c *StatsController) Member(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取会员统计
	now := time.Now().Unix()
	startTime := now - 30*24*60*60 // 30天前
	memberStats, err := c.statsService.GetMemberStats(startTime, now)
	if err != nil {
		logger.Error("获取会员统计失败", "error", err)
		http.Error(w, "Failed to get member stats", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Stats":       memberStats,
		"CurrentMenu": "stats",
		"PageTitle":   "会员统计",
	}

	// 渲染模板
	tplFile := "admin/stats_member.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染会员统计模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Visit 访问统计
func (c *StatsController) Visit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取天数
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// 获取访问统计
	now := time.Now().Unix()
	startTime := now - int64(days)*24*60*60
	visitStats, err := c.statsService.GetVisitStats(startTime, now)
	if err != nil {
		logger.Error("获取访问统计失败", "error", err)
		http.Error(w, "Failed to get visit stats", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Stats":       visitStats,
		"Days":        days,
		"CurrentMenu": "stats",
		"PageTitle":   "访问统计",
	}

	// 渲染模板
	tplFile := "admin/stats_visit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染访问统计模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Search 搜索统计
func (c *StatsController) Search(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取天数
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// 获取搜索统计
	now := time.Now().Unix()
	startTime := now - int64(days)*24*60*60
	searchStats, err := c.statsService.GetSearchStats(startTime, now)
	if err != nil {
		logger.Error("获取搜索统计失败", "error", err)
		http.Error(w, "Failed to get search stats", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Stats":       searchStats,
		"Days":        days,
		"CurrentMenu": "stats",
		"PageTitle":   "搜索统计",
	}

	// 渲染模板
	tplFile := "admin/stats_search.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染搜索统计模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Export 导出统计数据
func (c *StatsController) Export(w http.ResponseWriter, r *http.Request) {
	// 导出统计数据
	data, err := c.statsService.ExportStats()
	if err != nil {
		logger.Error("导出统计数据失败", "error", err)
		http.Error(w, "Failed to export stats", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=stats_"+time.Now().Format("20060102")+".json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	// 输出数据
	w.Write([]byte(data))
}

// GetData 获取统计数据
func (c *StatsController) GetData(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取数据类型
	dataType := r.FormValue("type")
	if dataType == "" {
		http.Error(w, "Missing data type", http.StatusBadRequest)
		return
	}

	// 获取天数
	daysStr := r.FormValue("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// 获取数据
	var data interface{}
	var err error

	switch dataType {
	case "site":
		now := time.Now().Unix()
		startTime := now - 30*24*60*60 // 30天前
		data, err = c.statsService.GetSiteStats(startTime, now)
	case "category":
		now := time.Now().Unix()
		startTime := now - 30*24*60*60 // 30天前
		data, err = c.statsService.GetCategoryStats(startTime, now)
	case "member":
		now := time.Now().Unix()
		startTime := now - 30*24*60*60 // 30天前
		data, err = c.statsService.GetMemberStats(startTime, now)
	case "visit":
		now := time.Now().Unix()
		startTime := now - int64(days)*24*60*60
		data, err = c.statsService.GetVisitStats(startTime, now)
	case "search":
		now := time.Now().Unix()
		startTime := now - int64(days)*24*60*60
		data, err = c.statsService.GetSearchStats(startTime, now)
	default:
		http.Error(w, "Invalid data type", http.StatusBadRequest)
		return
	}

	if err != nil {
		logger.Error("获取统计数据失败", "type", dataType, "error", err)
		http.Error(w, "Failed to get stats data", http.StatusInternalServerError)
		return
	}

	// 返回数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}
