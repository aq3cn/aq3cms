/*
 * @Author: dnsmap 5553557@gmail.com
 * @Date: 2025-05-30 10:33:06
 * @LastEditors: dnsmap 5553557@gmail.com
 * @LastEditTime: 2025-05-31 13:49:47
 * @FilePath: /aq3cms/internal/controller/frontend/index.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package frontend

import (
	"net/http"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// IndexController 首页控制器
type IndexController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	tagModel        *model.TagModel
	linkModel       *model.LinkModel
	templateService *service.TemplateService
	articleService  *service.ArticleService
	categoryService *service.CategoryService
	linkService     *service.LinkService
}

// NewIndexController 创建首页控制器
func NewIndexController(db *database.DB, cache cache.Cache, config *config.Config) *IndexController {
	return &IndexController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		tagModel:        model.NewTagModel(db),
		linkModel:       model.NewLinkModel(db),
		templateService: service.NewTemplateService(db, cache, config),
		articleService:  service.NewArticleService(db, cache, config),
		categoryService: service.NewCategoryService(db, cache, config),
		linkService:     service.NewLinkService(db, cache, config),
	}
}

// Index 首页
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	logger.Info("访问首页", "path", r.URL.Path)

	// 准备模板数据
	data := map[string]interface{}{
		"Title":       c.config.Site.Name + " - 首页",
		"Keywords":    c.config.Site.Keywords,
		"Description": c.config.Site.Description,
		"SiteName":    c.config.Site.Name,
		"SiteURL":     c.config.Site.URL,
		"CopyRight":   c.config.Site.CopyRight,
		"CurrentTime": time.Now().Format("2006-01-02 15:04:05"),
		"Version":     "1.0.0",

		// 全局变量
		"Globals": map[string]interface{}{
			"cfg_webname":     c.config.Site.Name,
			"cfg_weburl":      c.config.Site.URL,
			"cfg_keywords":    c.config.Site.Keywords,
			"cfg_description": c.config.Site.Description,
			"cfg_copyright":   c.config.Site.CopyRight,
		},

		// 字段数据（用于模板标签）
		"Fields": map[string]interface{}{
			"title":       c.config.Site.Name,
			"keywords":    c.config.Site.Keywords,
			"description": c.config.Site.Description,
		},
	}

	// 渲染模板
	tplFile := "default/index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染首页模板失败", "error", err)
		// 如果模板渲染失败，返回简单的错误页面
		http.Error(w, "页面加载失败", http.StatusInternalServerError)
		return
	}
}
