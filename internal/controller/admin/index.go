package admin

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
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
	memberModel     *model.MemberModel
	commentModel    *model.CommentModel
	templateService *service.TemplateService
}

// NewIndexController 创建首页控制器
func NewIndexController(db *database.DB, cache cache.Cache, config *config.Config) *IndexController {
	return &IndexController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		memberModel:     model.NewMemberModel(db),
		commentModel:    model.NewCommentModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 首页
func (c *IndexController) Index(w http.ResponseWriter, r *http.Request) {
	logger.Info("访问后台首页", "path", r.URL.Path)

	// 获取管理员信息
	adminName := middleware.GetAdminName(r)
	if adminName == "" {
		adminName = "管理员"
	}

	// 获取统计信息
	articleCount, _ := c.articleModel.GetCount()
	categoryCount, _ := c.categoryModel.GetCount()
	memberCount, _ := c.memberModel.GetCount()
	commentCount, _ := c.commentModel.GetCount()

	// 获取系统信息
	systemInfo := map[string]interface{}{
		"GoVersion":    runtime.Version(),
		"NumCPU":       runtime.NumCPU(),
		"GOOS":         runtime.GOOS,
		"GOARCH":       runtime.GOARCH,
		"ServerTime":   time.Now().Format("2006-01-02 15:04:05"),
		"ServerUptime": getUptime(),
	}

	// 简单的HTML管理界面
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>管理中心 - aq3cms</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background: #f5f5f5; }
        .header { background: #2c3e50; color: white; padding: 15px 20px; display: flex; justify-content: space-between; align-items: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .header .user-info { display: flex; align-items: center; gap: 15px; }
        .header .user-info a { color: white; text-decoration: none; }
        .header .user-info a:hover { text-decoration: underline; }
        .container { max-width: 1200px; margin: 20px auto; padding: 0 20px; }
        .dashboard { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .card h3 { margin: 0 0 15px 0; color: #333; }
        .stat-card { text-align: center; }
        .stat-number { font-size: 36px; font-weight: bold; color: #3498db; margin: 10px 0; }
        .stat-label { color: #666; font-size: 14px; }
        .menu-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; }
        .menu-item { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: center; transition: transform 0.2s; }
        .menu-item:hover { transform: translateY(-2px); }
        .menu-item a { text-decoration: none; color: #333; }
        .menu-item .icon { font-size: 32px; margin-bottom: 10px; }
        .menu-item .title { font-weight: bold; margin-bottom: 5px; }
        .menu-item .desc { font-size: 12px; color: #666; }
        .system-info { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-top: 20px; }
        .system-info h3 { margin: 0 0 15px 0; color: #333; }
        .info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; }
        .info-item { padding: 10px; background: #f8f9fa; border-radius: 5px; }
        .info-label { font-weight: bold; color: #555; }
        .info-value { color: #333; margin-top: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 aq3cms 管理中心</h1>
        <div class="user-info">
            <span>欢迎，` + adminName + `</span>
            <a href="/">查看网站</a>
            <a href="/aq3cms/logout">退出登录</a>
        </div>
    </div>

    <div class="container">
        <div class="dashboard">
            <div class="card stat-card">
                <h3>📄 文章统计</h3>
                <div class="stat-number">` + fmt.Sprintf("%d", articleCount) + `</div>
                <div class="stat-label">篇文章</div>
            </div>

            <div class="card stat-card">
                <h3>📁 栏目统计</h3>
                <div class="stat-number">` + fmt.Sprintf("%d", categoryCount) + `</div>
                <div class="stat-label">个栏目</div>
            </div>

            <div class="card stat-card">
                <h3>👥 会员统计</h3>
                <div class="stat-number">` + fmt.Sprintf("%d", memberCount) + `</div>
                <div class="stat-label">位会员</div>
            </div>

            <div class="card stat-card">
                <h3>💬 评论统计</h3>
                <div class="stat-number">` + fmt.Sprintf("%d", commentCount) + `</div>
                <div class="stat-label">条评论</div>
            </div>
        </div>

        <div class="card">
            <h3>🛠️ 快速操作</h3>
            <div class="menu-grid">
                <div class="menu-item">
                    <a href="/aq3cms/article">
                        <div class="icon">📝</div>
                        <div class="title">文章管理</div>
                        <div class="desc">管理网站文章内容</div>
                    </a>
                </div>

                <div class="menu-item">
                    <a href="/aq3cms/category">
                        <div class="icon">📂</div>
                        <div class="title">栏目管理</div>
                        <div class="desc">管理网站栏目分类</div>
                    </a>
                </div>

                <div class="menu-item">
                    <a href="/aq3cms/member">
                        <div class="icon">👤</div>
                        <div class="title">会员管理</div>
                        <div class="desc">管理注册会员</div>
                    </a>
                </div>

                <div class="menu-item">
                    <a href="/aq3cms/comment">
                        <div class="icon">💭</div>
                        <div class="title">评论管理</div>
                        <div class="desc">管理用户评论</div>
                    </a>
                </div>

                <div class="menu-item">
                    <a href="/aq3cms/template">
                        <div class="icon">🎨</div>
                        <div class="title">模板管理</div>
                        <div class="desc">管理网站模板</div>
                    </a>
                </div>

                <div class="menu-item">
                    <a href="/aq3cms/setting">
                        <div class="icon">⚙️</div>
                        <div class="title">系统设置</div>
                        <div class="desc">配置系统参数</div>
                    </a>
                </div>
            </div>
        </div>

        <div class="system-info">
            <h3>💻 系统信息</h3>
            <div class="info-grid">
                <div class="info-item">
                    <div class="info-label">Go 版本</div>
                    <div class="info-value">` + systemInfo["GoVersion"].(string) + `</div>
                </div>

                <div class="info-item">
                    <div class="info-label">CPU 核心</div>
                    <div class="info-value">` + fmt.Sprintf("%d", systemInfo["NumCPU"].(int)) + ` 核</div>
                </div>

                <div class="info-item">
                    <div class="info-label">操作系统</div>
                    <div class="info-value">` + systemInfo["GOOS"].(string) + ` / ` + systemInfo["GOARCH"].(string) + `</div>
                </div>

                <div class="info-item">
                    <div class="info-label">服务器时间</div>
                    <div class="info-value">` + systemInfo["ServerTime"].(string) + `</div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// Menu 菜单
func (c *IndexController) Menu(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)
	adminRank := middleware.GetAdminRank(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"AdminRank":   adminRank,
		"CurrentMenu": "menu",
		"PageTitle":   "系统菜单",
	}

	// 渲染模板
	tplFile := "admin/menu.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染菜单模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Main 主框架
func (c *IndexController) Main(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "main",
		"PageTitle":   "管理中心",
	}

	// 渲染模板
	tplFile := "admin/main.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染主框架模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Top 顶部框架
func (c *IndexController) Top(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "top",
		"PageTitle":   "顶部导航",
		"SiteName":    c.config.Site.Name,
	}

	// 渲染模板
	tplFile := "admin/top.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染顶部框架模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Frame 框架页
func (c *IndexController) Frame(w http.ResponseWriter, r *http.Request) {
	logger.Info("访问后台框架页", "path", r.URL.Path)

	// 直接调用Index方法显示后台首页
	c.Index(w, r)
}

// 获取系统运行时间
func getUptime() string {
	// 这里简化处理，实际应该记录程序启动时间
	return "未知"
}
