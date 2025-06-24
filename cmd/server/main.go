package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"aq3cms/config"
	"aq3cms/internal/controller"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 启动服务器
func run() error {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Path); err != nil {
		return fmt.Errorf("初始化日志失败: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	defer db.Close()

	logger.Info("数据库连接成功", "type", cfg.Database.Type, "host", cfg.Database.Host)

	// 初始化默认管理员
	adminModel := model.NewAdminModel(db)
	if err := adminModel.InitDefaultAdmin(); err != nil {
		logger.Error("初始化默认管理员失败", "error", err)
	}

	// 初始化缓存
	var cacheProvider cache.Cache
	logger.Info("初始化缓存", "type", cfg.Cache.Type)
	switch cfg.Cache.Type {
	case "memory":
		logger.Info("使用内存缓存")
		cacheProvider = cache.NewMemoryCache()
	case "redis":
		logger.Info("使用Redis缓存")
		cacheProvider, err = cache.NewRedisCache(cfg.Cache)
		if err != nil {
			return fmt.Errorf("Redis缓存初始化失败: %v", err)
		}
	default:
		logger.Info("使用文件缓存")
		cacheProvider = cache.NewFileCache(cfg.Cache.Path)
	}

	// 确保必要的目录存在
	ensureDirectories()

	// 初始化路由
	router := mux.NewRouter()
	controller.RegisterRoutes(router, db, cacheProvider, cfg)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("正在关闭服务器...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("服务器关闭出错", "error", err)
		}
	}()

	// 启动服务器
	logger.Info("服务器启动", "地址", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %v", err)
	}

	return nil
}

// 确保必要的目录存在
func ensureDirectories() {
	dirs := []string{
		"uploads",
		"uploads/images",
		"uploads/media",
		"uploads/files",
		"data/cache",
		"data/backup",
		"data/sessions",
		"data/logs",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error("创建目录失败", "目录", dir, "错误", err)
		}
	}

	// 检查模板目录是否存在
	if _, err := os.Stat("templets"); os.IsNotExist(err) {
		logger.Warn("模板目录不存在，将创建默认模板目录", nil)
		os.MkdirAll("templets/default", 0755)
	}
}
