package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"aq3cms/config"
	"aq3cms/internal/controller/admin"
	"aq3cms/internal/controller/api"
	"aq3cms/internal/controller/frontend"
	"aq3cms/internal/middleware"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/i18n"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/plugin"

	"github.com/gorilla/mux"
)

// HealthController 健康检查控制器
type HealthController struct {
	db     *database.DB
	cache  cache.Cache
	config *config.Config
}

var startTime = time.Now()

// Health 健康检查端点
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	// 检查数据库连接
	dbStatus := "healthy"
	dbError := ""
	if err := c.db.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbError = err.Error()
	}

	// 检查缓存连接
	cacheStatus := "healthy"
	cacheError := ""
	if c.cache != nil {
		testKey := "health_check_" + time.Now().Format("20060102150405")
		if err := c.cache.Set(testKey, "test", 10); err != nil {
			cacheStatus = "unhealthy"
			cacheError = err.Error()
		} else {
			c.cache.Delete(testKey)
		}
	}

	// 系统信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	systemInfo := map[string]interface{}{
		"go_version":    runtime.Version(),
		"num_cpu":       runtime.NumCPU(),
		"num_goroutine": runtime.NumGoroutine(),
		"memory_alloc":  m.Alloc,
		"memory_total":  m.TotalAlloc,
		"memory_sys":    m.Sys,
		"gc_runs":       m.NumGC,
	}

	// 服务状态
	services := map[string]interface{}{
		"database": map[string]interface{}{
			"status": dbStatus,
			"error":  dbError,
		},
		"cache": map[string]interface{}{
			"status": cacheStatus,
			"error":  cacheError,
		},
	}

	// 确定整体状态
	overallStatus := "healthy"
	if dbStatus != "healthy" || cacheStatus != "healthy" {
		overallStatus = "unhealthy"
	}

	// 构建响应
	response := map[string]interface{}{
		"status":    overallStatus,
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"uptime":    time.Since(startTime).String(),
		"system":    systemInfo,
		"services":  services,
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	if overallStatus != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// 返回JSON响应
	json.NewEncoder(w).Encode(response)
}

// Ready 就绪检查端点
func (c *HealthController) Ready(w http.ResponseWriter, r *http.Request) {
	ready := true
	errors := make([]string, 0)

	// 检查数据库
	if err := c.db.Ping(); err != nil {
		ready = false
		errors = append(errors, "database: "+err.Error())
	}

	// 检查缓存
	if c.cache != nil {
		testKey := "ready_check_" + time.Now().Format("20060102150405")
		if err := c.cache.Set(testKey, "test", 10); err != nil {
			ready = false
			errors = append(errors, "cache: "+err.Error())
		} else {
			c.cache.Delete(testKey)
		}
	}

	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now().Format(time.RFC3339),
		"errors":    errors,
	}

	w.Header().Set("Content-Type", "application/json")
	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

// Live 存活检查端点
func (c *HealthController) Live(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes 注册路由
func RegisterRoutes(router *mux.Router, db *database.DB, cache cache.Cache, cfg *config.Config) {

	// 初始化插件系统
	// 确保插件目录存在
	if cfg.Plugin.Dir != "" {
		if err := os.MkdirAll(cfg.Plugin.Dir, 0755); err != nil {
			logger.Error("创建插件目录失败", "dir", cfg.Plugin.Dir, "error", err)
		}
	}

	// 确保插件配置文件的父目录存在
	if cfg.Plugin.ConfigFile != "" {
		configDir := filepath.Dir(cfg.Plugin.ConfigFile)
		if configDir != "" && configDir != "." {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				logger.Error("创建插件配置目录失败", "dir", configDir, "error", err)
			}
		}
	}

	pluginManager := plugin.NewManager(cfg.Plugin.Dir, cfg.Plugin.ConfigFile)
	if err := pluginManager.LoadPlugins(); err != nil {
		// 记录错误但继续运行
		logger.Error("加载插件失败", "error", err)
	}

	// 初始化国际化
	i18nInstance, err := i18n.New(cfg.Site.DefaultLang, filepath.Join(cfg.Template.Dir, "lang"))
	if err != nil {
		// 创建一个默认的国际化实例
		logger.Error("初始化国际化失败", "error", err)
		// 创建一个内存中的国际化实例
		i18nInstance = i18n.NewMemory("zh-cn")
		// 添加一些基本的翻译
		i18nInstance.AddLang("zh-cn", map[string]string{
			"site_name": "aq3cms",
			"home":      "首页",
			"login":     "登录",
			"register":  "注册",
			"logout":    "退出",
			"search":    "搜索",
		})
	}

	// 初始化认证中间件
	middleware.InitAuth(cfg)

	// 添加中间件
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)
	router.Use(middleware.Security)
	router.Use(middleware.I18nMiddleware(i18nInstance))

	// 添加速率限制中间件
	if cfg.Server.EnableRateLimit {
		router.Use(middleware.RateLimit(cfg.Server.RateLimitRequests, time.Duration(cfg.Server.RateLimitWindow)*time.Second))
	}

	// 健康检查端点
	healthController := &HealthController{db: db, cache: cache, config: cfg}
	router.HandleFunc("/health", healthController.Health).Methods("GET")
	router.HandleFunc("/ready", healthController.Ready).Methods("GET")
	router.HandleFunc("/live", healthController.Live).Methods("GET")

	// 静态文件
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// 注册API路由
	api.RegisterRoutes(router, db, cache, cfg, pluginManager)

	// 前台控制器
	indexController := frontend.NewIndexController(db, cache, cfg)
	articleController := frontend.NewArticleController(db, cache, cfg)
	categoryController := frontend.NewCategoryController(db, cache, cfg)
	memberController := frontend.NewMemberController(db, cache, cfg)
	commentController := frontend.NewCommentController(db, cache, cfg)
	tagController := frontend.NewTagController(db, cache, cfg)
	specialController := frontend.NewSpecialController(db, cache, cfg)
	rssController := frontend.NewRSSController(db, cache, cfg)
	messageController := frontend.NewMessageController(db, cache, cfg)
	formController := frontend.NewFormController(db, cache, cfg)
	searchController := frontend.NewSearchController(db, cache, cfg)

	// 前台路由
	router.HandleFunc("/", indexController.Index).Methods("GET")

	// 文章路由
	router.HandleFunc("/article/{id:[0-9]+}.html", articleController.Detail).Methods("GET")
	router.HandleFunc("/list/{typeid:[0-9]+}.html", articleController.List).Methods("GET")
	router.HandleFunc("/articles", articleController.AllArticles).Methods("GET")

	// 栏目路由
	router.HandleFunc("/category/{typeid:[0-9]+}.html", categoryController.List).Methods("GET")

	// 多级目录路由 - 支持 /dir/subdir 格式，排除保留路径
	router.HandleFunc("/{dir:(?!aq3cms|member|admin|api)[a-zA-Z0-9_-]+}/{subdir:[a-zA-Z0-9_-]+}", categoryController.ShowByPath).Methods("GET")

	// 通用目录路由（必须放在最后，因为它会匹配所有路径），排除保留路径
	router.HandleFunc("/{dir:(?!aq3cms|member|admin|api)[a-zA-Z0-9_-]+}", categoryController.ShowByDir).Methods("GET")

	// 产品路由
	router.HandleFunc("/product/{id:[0-9]+}.html", articleController.Detail).Methods("GET")

	// 下载路由
	router.HandleFunc("/download/{id:[0-9]+}.html", articleController.Detail).Methods("GET")

	// 评论路由
	router.HandleFunc("/comment/list/{aid:[0-9]+}", commentController.List).Methods("GET")
	router.HandleFunc("/comment/post", commentController.Post).Methods("POST")
	router.HandleFunc("/comment/vote", commentController.Vote).Methods("POST")
	router.HandleFunc("/comment/reply", commentController.Reply).Methods("POST")

	// 标签路由
	router.HandleFunc("/tags", tagController.List).Methods("GET")
	router.HandleFunc("/tag/{tag}", tagController.Detail).Methods("GET")

	// 专题路由
	router.HandleFunc("/special", specialController.List).Methods("GET")
	router.HandleFunc("/special/{filename}", specialController.Detail).Methods("GET")

	// RSS路由
	router.HandleFunc("/rss", rssController.Index).Methods("GET")
	router.HandleFunc("/rss/{typeid:[0-9]+}", rssController.Category).Methods("GET")

	// 表单路由
	router.HandleFunc("/form/{id:[0-9]+}", formController.Show).Methods("GET")
	router.HandleFunc("/form/submit", formController.Submit).Methods("POST")
	router.HandleFunc("/form/success/{id:[0-9]+}", formController.Success).Methods("GET")

	// 搜索路由
	router.HandleFunc("/search", searchController.Search).Methods("GET")
	router.HandleFunc("/search/advanced", searchController.AdvancedSearch).Methods("GET")

	// 会员路由
	memberRouter := router.PathPrefix("/member").Subrouter()
	memberRouter.HandleFunc("/login", memberController.Login).Methods("GET")
	memberRouter.HandleFunc("/login", memberController.DoLogin).Methods("POST")
	memberRouter.HandleFunc("/logout", memberController.Logout).Methods("GET")
	memberRouter.HandleFunc("/register", memberController.Register).Methods("GET")
	memberRouter.HandleFunc("/register", memberController.DoRegister).Methods("POST")

	// 需要登录的会员路由
	memberAuthRouter := memberRouter.NewRoute().Subrouter()
	memberAuthRouter.Use(middleware.MemberAuth)
	memberAuthRouter.HandleFunc("/", memberController.Index).Methods("GET")
	memberAuthRouter.HandleFunc("/profile", memberController.Profile).Methods("GET")
	memberAuthRouter.HandleFunc("/profile", memberController.UpdateProfile).Methods("POST")
	memberAuthRouter.HandleFunc("/password", memberController.ChangePassword).Methods("GET")
	memberAuthRouter.HandleFunc("/password", memberController.DoChangePassword).Methods("POST")
	memberAuthRouter.HandleFunc("/articles", memberController.Articles).Methods("GET")

	// 会员消息路由
	memberAuthRouter.HandleFunc("/inbox", messageController.Inbox).Methods("GET")
	memberAuthRouter.HandleFunc("/outbox", messageController.Outbox).Methods("GET")
	memberAuthRouter.HandleFunc("/message/read", messageController.Read).Methods("GET")
	memberAuthRouter.HandleFunc("/message/send", messageController.Send).Methods("GET")
	memberAuthRouter.HandleFunc("/message/send", messageController.DoSend).Methods("POST")
	memberAuthRouter.HandleFunc("/message/delete", messageController.Delete).Methods("POST")

	// 后台控制器
	// 创建会话存储（使用已初始化的会话存储）

	// 初始化后台控制器
	adminLoginController := admin.NewLoginController(db, cache, cfg)
	adminIndexController := admin.NewIndexController(db, cache, cfg)
	adminArticleController := admin.NewArticleController(db, cache, cfg)
	adminCategoryController := admin.NewCategoryController(db, cache, cfg)
	adminTagController := admin.NewTagController(db, cache, cfg)
	adminMemberController := admin.NewMemberController(db, cache, cfg)
	adminCommentController := admin.NewCommentController(db, cache, cfg)
	adminTemplateController := admin.NewTemplateController(db, cache, cfg)
	adminSettingController := admin.NewSettingController(db, cache, cfg)
	adminSystemController := admin.NewSystemController(db, cache, cfg)
	adminHtmlController := admin.NewHtmlController(db, cache, cfg)
	adminI18nController := admin.NewI18nController(db, cache, cfg)
	adminStatsController := admin.NewStatsController(db, cache, cfg)
	adminModelController := admin.NewModelController(db, cache, cfg)
	adminPluginController := admin.NewPluginController(db, cache, cfg, pluginManager)
	adminCollectController := admin.NewCollectController(db, cache, cfg)
	adminVoteController := admin.NewVoteController(db, cache, cfg)
	adminLinkController := admin.NewLinkController(db, cache, cfg)
	adminDBFixController := admin.NewDBFixController(db)

	// 后台路由
	adminRouter := router.PathPrefix("/aq3cms").Subrouter()
	adminRouter.HandleFunc("/login", adminLoginController.Login).Methods("GET")
	adminRouter.HandleFunc("/login", adminLoginController.DoLogin).Methods("POST")
	adminRouter.HandleFunc("/logout", adminLoginController.Logout).Methods("GET")
	adminRouter.HandleFunc("/captcha", adminLoginController.Captcha).Methods("GET")

	// 需要登录的后台路由
	adminAuthRouter := adminRouter.NewRoute().Subrouter()
	adminAuthRouter.Use(middleware.AdminAuth)

	// 首页相关
	adminAuthRouter.HandleFunc("/", adminIndexController.Frame).Methods("GET")
	adminAuthRouter.HandleFunc("/index", adminIndexController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/menu", adminIndexController.Menu).Methods("GET")
	adminAuthRouter.HandleFunc("/main", adminIndexController.Main).Methods("GET")
	adminAuthRouter.HandleFunc("/top", adminIndexController.Top).Methods("GET")
	adminAuthRouter.HandleFunc("/frame", adminIndexController.Frame).Methods("GET")

	// 文章管理
	adminAuthRouter.HandleFunc("/article", adminArticleController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/article_list", adminArticleController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/article_add", adminArticleController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/article_add", adminArticleController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/article_edit/{id:[0-9]+}", adminArticleController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/article_edit/{id:[0-9]+}", adminArticleController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/article_delete/{id:[0-9]+}", adminArticleController.Delete).Methods("GET")

	// 栏目管理
	adminAuthRouter.HandleFunc("/category", adminCategoryController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/category_list", adminCategoryController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/category_add", adminCategoryController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/category_add", adminCategoryController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/category_edit/{id:[0-9]+}", adminCategoryController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/category_edit/{id:[0-9]+}", adminCategoryController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/category_delete/{id:[0-9]+}", adminCategoryController.Delete).Methods("GET")

	// 标签管理
	adminAuthRouter.HandleFunc("/tag", adminTagController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/tag_list", adminTagController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/tag_add", adminTagController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/tag_add", adminTagController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/tag_edit/{id:[0-9]+}", adminTagController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/tag_edit/{id:[0-9]+}", adminTagController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/tag_delete/{id:[0-9]+}", adminTagController.Delete).Methods("GET")

	// 会员管理
	adminAuthRouter.HandleFunc("/member", adminMemberController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/member_list", adminMemberController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/member_add", adminMemberController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/member_add", adminMemberController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/member_edit/{id:[0-9]+}", adminMemberController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/member_edit/{id:[0-9]+}", adminMemberController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/member_delete/{id:[0-9]+}", adminMemberController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/member_batch_delete", adminMemberController.BatchDelete).Methods("POST")

	// 数据库修复
	adminAuthRouter.HandleFunc("/dbfix_member_table", adminDBFixController.FixMemberTable).Methods("POST")
	adminAuthRouter.HandleFunc("/dbfix_create_feedback_table", adminDBFixController.CreateFeedbackTable).Methods("POST")

	// 系统设置
	adminAuthRouter.HandleFunc("/system_config", adminSystemController.Config).Methods("GET")
	adminAuthRouter.HandleFunc("/system_config", adminSystemController.SaveConfig).Methods("POST")
	adminAuthRouter.HandleFunc("/system_database", adminSystemController.Database).Methods("GET")
	adminAuthRouter.HandleFunc("/system_backup", adminSystemController.Backup).Methods("POST")
	adminAuthRouter.HandleFunc("/system_restore", adminSystemController.Restore).Methods("POST")
	adminAuthRouter.HandleFunc("/system_cache", adminSystemController.Cache).Methods("GET")
	adminAuthRouter.HandleFunc("/system_clear_cache", adminSystemController.ClearCache).Methods("POST")
	adminAuthRouter.HandleFunc("/system_log", adminSystemController.Log).Methods("GET")
	adminAuthRouter.HandleFunc("/system_clear_log", adminSystemController.ClearLog).Methods("POST")

	// 网站设置
	adminAuthRouter.HandleFunc("/setting", adminSettingController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/setting/basic", adminSettingController.Basic).Methods("GET")
	adminAuthRouter.HandleFunc("/setting/basic", adminSettingController.DoBasic).Methods("POST")
	adminAuthRouter.HandleFunc("/setting/upload", adminSettingController.Upload).Methods("GET")
	adminAuthRouter.HandleFunc("/setting/upload", adminSettingController.DoUpload).Methods("POST")
	adminAuthRouter.HandleFunc("/setting/static", adminSettingController.Static).Methods("GET")
	adminAuthRouter.HandleFunc("/setting/static", adminSettingController.DoStatic).Methods("POST")
	adminAuthRouter.HandleFunc("/setting/cache", adminSettingController.Cache).Methods("GET")
	adminAuthRouter.HandleFunc("/setting/cache", adminSettingController.DoCache).Methods("POST")
	adminAuthRouter.HandleFunc("/setting/clear_cache", adminSettingController.ClearCache).Methods("POST")

	// 评论管理
	adminAuthRouter.HandleFunc("/comment", adminCommentController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_list", adminCommentController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_detail/{id:[0-9]+}", adminCommentController.Detail).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_approve/{id:[0-9]+}", adminCommentController.Approve).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_reject/{id:[0-9]+}", adminCommentController.Reject).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_delete/{id:[0-9]+}", adminCommentController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/comment_batch_approve", adminCommentController.BatchApprove).Methods("POST")
	adminAuthRouter.HandleFunc("/comment_batch_delete", adminCommentController.BatchDelete).Methods("POST")

	// 模板管理
	adminAuthRouter.HandleFunc("/template", adminTemplateController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/template_list", adminTemplateController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/template_edit", adminTemplateController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/template_edit", adminTemplateController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/template_create", adminTemplateController.Create).Methods("GET")
	adminAuthRouter.HandleFunc("/template_create", adminTemplateController.DoCreate).Methods("POST")
	adminAuthRouter.HandleFunc("/template_delete/{path:.*}", adminTemplateController.Delete).Methods("GET")

	// 已移动到 /setting/ 路径下

	// 静态页面生成
	adminAuthRouter.HandleFunc("/html_index", adminHtmlController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/html_index", adminHtmlController.DoIndex).Methods("POST")
	adminAuthRouter.HandleFunc("/html_list", adminHtmlController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/html_list", adminHtmlController.DoList).Methods("POST")
	adminAuthRouter.HandleFunc("/html_category", adminHtmlController.List).Methods("GET")    // 别名
	adminAuthRouter.HandleFunc("/html_category", adminHtmlController.DoList).Methods("POST") // 别名
	adminAuthRouter.HandleFunc("/html_article", adminHtmlController.Article).Methods("GET")
	adminAuthRouter.HandleFunc("/html_article", adminHtmlController.DoArticle).Methods("POST")
	adminAuthRouter.HandleFunc("/html_special", adminHtmlController.Special).Methods("GET")
	adminAuthRouter.HandleFunc("/html_special", adminHtmlController.DoSpecial).Methods("POST")
	adminAuthRouter.HandleFunc("/html_tag", adminHtmlController.Tag).Methods("GET")
	adminAuthRouter.HandleFunc("/html_tag", adminHtmlController.DoTag).Methods("POST")
	adminAuthRouter.HandleFunc("/html_all", adminHtmlController.All).Methods("GET")
	adminAuthRouter.HandleFunc("/html_all", adminHtmlController.DoAll).Methods("POST")

	// 多语言管理
	adminAuthRouter.HandleFunc("/i18n_list", adminI18nController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_add", adminI18nController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_add", adminI18nController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/i18n_edit/{lang}", adminI18nController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_edit", adminI18nController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/i18n_delete/{lang}", adminI18nController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_set_default/{lang}", adminI18nController.SetDefault).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_import", adminI18nController.Import).Methods("GET")
	adminAuthRouter.HandleFunc("/i18n_import", adminI18nController.DoImport).Methods("POST")
	adminAuthRouter.HandleFunc("/i18n_export/{lang}", adminI18nController.Export).Methods("GET")

	// 统计管理
	adminAuthRouter.HandleFunc("/stats_index", adminStatsController.Index).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_category", adminStatsController.Category).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_member", adminStatsController.Member).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_visit", adminStatsController.Visit).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_search", adminStatsController.Search).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_export", adminStatsController.Export).Methods("GET")
	adminAuthRouter.HandleFunc("/stats_data", adminStatsController.GetData).Methods("POST")

	// 内容模型管理
	adminAuthRouter.HandleFunc("/model_list", adminModelController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/model_add", adminModelController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/model_add", adminModelController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/model_edit/{id:[0-9]+}", adminModelController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/model_edit", adminModelController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/model_delete/{id:[0-9]+}", adminModelController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/model_fields/{id:[0-9]+}", adminModelController.Fields).Methods("GET")
	adminAuthRouter.HandleFunc("/model_content/{id:[0-9]+}", adminModelController.Content).Methods("GET")
	adminAuthRouter.HandleFunc("/model_content_add/{id:[0-9]+}", adminModelController.ContentAdd).Methods("GET")
	adminAuthRouter.HandleFunc("/model_content_add", adminModelController.ContentDoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/model_content_edit/{id:[0-9]+}/{aid:[0-9]+}", adminModelController.ContentEdit).Methods("GET")
	adminAuthRouter.HandleFunc("/model_content_edit", adminModelController.ContentDoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/model_content_delete/{id:[0-9]+}/{aid:[0-9]+}", adminModelController.ContentDelete).Methods("GET")

	// 插件管理
	adminAuthRouter.HandleFunc("/plugin_list", adminPluginController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/plugin_enable/{name}", adminPluginController.Enable).Methods("GET")
	adminAuthRouter.HandleFunc("/plugin_disable/{name}", adminPluginController.Disable).Methods("GET")
	adminAuthRouter.HandleFunc("/plugin_config/{name}", adminPluginController.Config).Methods("GET")
	adminAuthRouter.HandleFunc("/plugin_config", adminPluginController.DoConfig).Methods("POST")
	adminAuthRouter.HandleFunc("/plugin_upload", adminPluginController.Upload).Methods("GET")
	adminAuthRouter.HandleFunc("/plugin_upload", adminPluginController.DoUpload).Methods("POST")
	adminAuthRouter.HandleFunc("/plugin_delete/{name}", adminPluginController.Delete).Methods("GET")

	// 采集管理
	adminAuthRouter.HandleFunc("/collect_rule_list", adminCollectController.RuleList).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_rule_add", adminCollectController.RuleAdd).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_rule_add", adminCollectController.RuleDoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/collect_rule_edit/{id:[0-9]+}", adminCollectController.RuleEdit).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_rule_edit", adminCollectController.RuleDoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/collect_rule_delete/{id:[0-9]+}", adminCollectController.RuleDelete).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_rule_collect/{id:[0-9]+}", adminCollectController.RuleCollect).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_item_list/{id:[0-9]+}", adminCollectController.ItemList).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_item_detail/{id:[0-9]+}", adminCollectController.ItemDetail).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_item_publish/{id:[0-9]+}", adminCollectController.ItemPublish).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_item_delete/{id:[0-9]+}", adminCollectController.ItemDelete).Methods("GET")
	adminAuthRouter.HandleFunc("/collect_batch_publish/{id:[0-9]+}", adminCollectController.BatchPublish).Methods("GET")

	// 投票管理
	adminAuthRouter.HandleFunc("/vote_list", adminVoteController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/vote_add", adminVoteController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/vote_add", adminVoteController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/vote_edit/{id:[0-9]+}", adminVoteController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/vote_edit", adminVoteController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/vote_delete/{id:[0-9]+}", adminVoteController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/vote_result/{id:[0-9]+}", adminVoteController.Result).Methods("GET")
	adminAuthRouter.HandleFunc("/vote_logs/{id:[0-9]+}", adminVoteController.Logs).Methods("GET")

	// 友情链接管理
	adminAuthRouter.HandleFunc("/link_list", adminLinkController.List).Methods("GET")
	adminAuthRouter.HandleFunc("/link_add", adminLinkController.Add).Methods("GET")
	adminAuthRouter.HandleFunc("/link_add", adminLinkController.DoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/link_edit/{id:[0-9]+}", adminLinkController.Edit).Methods("GET")
	adminAuthRouter.HandleFunc("/link_edit", adminLinkController.DoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/link_delete/{id:[0-9]+}", adminLinkController.Delete).Methods("GET")
	adminAuthRouter.HandleFunc("/link_type_list", adminLinkController.TypeList).Methods("GET")
	adminAuthRouter.HandleFunc("/link_type_add", adminLinkController.TypeAdd).Methods("GET")
	adminAuthRouter.HandleFunc("/link_type_add", adminLinkController.TypeDoAdd).Methods("POST")
	adminAuthRouter.HandleFunc("/link_type_edit/{id:[0-9]+}", adminLinkController.TypeEdit).Methods("GET")
	adminAuthRouter.HandleFunc("/link_type_edit", adminLinkController.TypeDoEdit).Methods("POST")
	adminAuthRouter.HandleFunc("/link_type_delete/{id:[0-9]+}", adminLinkController.TypeDelete).Methods("GET")
}
