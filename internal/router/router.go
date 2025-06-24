package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/controller"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
)

// New 创建路由
func New(db *database.DB, cache cache.Cache, config *config.Config) http.Handler {
	// 创建路由
	r := mux.NewRouter()

	// 注册路由
	controller.RegisterRoutes(r, db, cache, config)

	return r
}
