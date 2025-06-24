package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/plugin"
)

// RegisterRoutes 注册API路由
func RegisterRoutes(router *mux.Router, db *database.DB, cache cache.Cache, config *config.Config, pluginManager *plugin.Manager) {
	// 创建API路由
	apiRouter := router.PathPrefix("/api").Subrouter()

	// 创建API控制器
	articleController := NewArticleController(db, cache, config)
	memberController := NewMemberController(db, cache, config)
	commentController := NewCommentController(db, cache, config)
	categoryController := NewCategoryController(db, cache, config)
	tagController := NewTagController(db, cache, config)
	specialController := NewSpecialController(db, cache, config)
	searchController := NewSearchController(db, cache, config)
	uploadController := NewUploadController(db, cache, config)

	// 文章API
	apiRouter.HandleFunc("/articles", articleController.List).Methods("GET")
	apiRouter.HandleFunc("/articles/{id:[0-9]+}", articleController.Detail).Methods("GET")
	apiRouter.HandleFunc("/articles", articleController.Create).Methods("POST")
	apiRouter.HandleFunc("/articles/{id:[0-9]+}", articleController.Update).Methods("PUT")
	apiRouter.HandleFunc("/articles/{id:[0-9]+}", articleController.Delete).Methods("DELETE")

	// 会员API
	apiRouter.HandleFunc("/members/login", memberController.Login).Methods("POST")
	apiRouter.HandleFunc("/members/register", memberController.Register).Methods("POST")
	apiRouter.HandleFunc("/members/profile", memberController.Profile).Methods("GET")
	apiRouter.HandleFunc("/members/profile", memberController.UpdateProfile).Methods("PUT")
	apiRouter.HandleFunc("/members/password", memberController.ChangePassword).Methods("PUT")
	apiRouter.HandleFunc("/members/logout", memberController.Logout).Methods("POST")

	// 评论API
	apiRouter.HandleFunc("/comments/article/{aid:[0-9]+}", commentController.List).Methods("GET")
	apiRouter.HandleFunc("/comments", commentController.Create).Methods("POST")
	apiRouter.HandleFunc("/comments/{id:[0-9]+}", commentController.Delete).Methods("DELETE")
	apiRouter.HandleFunc("/comments/{id:[0-9]+}/vote", commentController.Vote).Methods("POST")

	// 栏目API
	apiRouter.HandleFunc("/categories", categoryController.List).Methods("GET")
	apiRouter.HandleFunc("/categories/{id:[0-9]+}", categoryController.Detail).Methods("GET")

	// 标签API
	apiRouter.HandleFunc("/tags", tagController.List).Methods("GET")
	apiRouter.HandleFunc("/tags/{name}", tagController.Detail).Methods("GET")

	// 专题API
	apiRouter.HandleFunc("/specials", specialController.List).Methods("GET")
	apiRouter.HandleFunc("/specials/{id:[0-9]+}", specialController.Detail).Methods("GET")

	// 搜索API
	apiRouter.HandleFunc("/search", searchController.Search).Methods("GET")
	apiRouter.HandleFunc("/search/hot", searchController.Hot).Methods("GET")

	// 上传API
	apiRouter.HandleFunc("/upload/image", uploadController.Image).Methods("POST")
	apiRouter.HandleFunc("/upload/file", uploadController.File).Methods("POST")

	// 添加CORS中间件
	apiRouter.Use(corsMiddleware)

	// 应用插件钩子
	if pluginManager != nil {
		pluginManager.ApplyHooks("api_register_routes", apiRouter, db, cache, config)
	}
}

// corsMiddleware CORS中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
