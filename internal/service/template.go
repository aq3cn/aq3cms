package service

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"aq3cms/config"
	"aq3cms/internal/interfaces"
	tmpl "aq3cms/internal/template"
	"aq3cms/internal/template/tags"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// TemplateService 模板服务
type TemplateService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	engine          *tmpl.Engine
	globals         map[string]interface{}
	globalsMtx      sync.RWMutex
	templateCache   map[string]*template.Template
	templateCacheMtx sync.RWMutex
	i18nService     interfaces.I18nServiceInterface
	seoService      interfaces.SEOServiceInterface
	statsService    interfaces.StatsServiceInterface
}

// NewTemplateService 创建模板服务
func NewTemplateService(db *database.DB, cache cache.Cache, config *config.Config) *TemplateService {
	// 创建模板引擎
	engine := tmpl.New(&config.Template, cache)

	// 创建模板服务
	service := &TemplateService{
		db:              db,
		cache:           cache,
		config:          config,
		engine:          engine,
		globals:         make(map[string]interface{}),
		templateCache:   make(map[string]*template.Template),
	}

	// 创建国际化服务
	service.i18nService = NewI18nService(db, cache, config)

	// 创建SEO服务
	service.seoService = NewSEOService(db, cache, config)

	// 创建统计服务
	service.statsService = NewStatsService(db, cache, config)

	// 注册标签处理器
	service.registerTagHandlers()

	// 加载全局变量
	service.loadGlobals()

	return service
}

// Render 渲染模板
func (s *TemplateService) Render(w io.Writer, name string, data interface{}) error {
	// 添加全局变量
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		dataMap = make(map[string]interface{})
	}

	// 添加全局变量
	s.globalsMtx.RLock()
	dataMap["Globals"] = s.globals
	s.globalsMtx.RUnlock()

	// 添加辅助函数
	dataMap["FormatDate"] = s.formatDate
	dataMap["FormatTime"] = s.formatTime
	dataMap["FormatDateTime"] = s.formatDateTime
	dataMap["Truncate"] = s.truncate
	dataMap["StripTags"] = s.stripTags
	dataMap["URLEncode"] = s.urlEncode
	dataMap["HTMLEncode"] = s.htmlEncode
	dataMap["Translate"] = s.translate
	dataMap["GetMetaTags"] = s.getMetaTags
	dataMap["GetOpenGraphTags"] = s.getOpenGraphTags
	dataMap["GetCanonicalURL"] = s.getCanonicalURL
	dataMap["GetAlternateURLs"] = s.getAlternateURLs
	dataMap["GetAvailableLangs"] = s.getAvailableLangs

	// 渲染模板
	return s.engine.Render(w, name, dataMap)
}

// GenerateStaticPage 生成静态页面
func (s *TemplateService) GenerateStaticPage(tplFile string, data interface{}, outputFile string) error {
	// 创建缓冲区
	var buf bytes.Buffer

	// 渲染模板
	err := s.Render(&buf, tplFile, data)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
}

// GetGlobals 获取全局变量
func (s *TemplateService) GetGlobals() map[string]interface{} {
	s.globalsMtx.RLock()
	defer s.globalsMtx.RUnlock()

	// 复制全局变量
	globals := make(map[string]interface{})
	for k, v := range s.globals {
		globals[k] = v
	}

	return globals
}

// SetGlobal 设置全局变量
func (s *TemplateService) SetGlobal(name string, value interface{}) {
	s.globalsMtx.Lock()
	defer s.globalsMtx.Unlock()

	s.globals[name] = value
}

// 加载全局变量
func (s *TemplateService) loadGlobals() {
	s.globalsMtx.Lock()
	defer s.globalsMtx.Unlock()

	// 从缓存中获取全局变量
	if cached, ok := s.cache.Get("globals"); ok {
		if globals, ok := cached.(map[string]interface{}); ok {
			s.globals = globals
			return
		}
	}

	// 从数据库加载全局变量
	qb := database.NewQueryBuilder(s.db, "sysconfig")
	qb.Select("*")

	results, err := qb.Get()
	if err != nil {
		logger.Error("加载全局变量失败", "error", err)
		return
	}

	// 处理结果
	for _, result := range results {
		if varname, ok := result["varname"].(string); ok {
			s.globals[varname] = result["value"]
		}
	}

	// 添加站点配置
	s.globals["cfg_webname"] = s.config.Site.Name
	s.globals["cfg_weburl"] = s.config.Site.URL
	s.globals["cfg_keywords"] = s.config.Site.Keywords
	s.globals["cfg_description"] = s.config.Site.Description
	s.globals["cfg_icp"] = s.config.Site.ICP
	s.globals["cfg_statcode"] = s.config.Site.StatCode
	s.globals["cfg_copyright"] = s.config.Site.CopyRight

	// 缓存全局变量
	cache.SafeSet(s.cache, "globals", s.globals, 0)
}

// ClearCache 清除缓存
func (s *TemplateService) ClearCache() {
	// 清除模板缓存
	s.templateCacheMtx.Lock()
	s.templateCache = make(map[string]*template.Template)
	s.templateCacheMtx.Unlock()

	// 清除引擎缓存
	// s.engine.ClearCache() // 引擎没有ClearCache方法，暂时注释掉

	// 清除全局变量缓存
	s.cache.Delete("globals")

	// 重新加载全局变量
	s.loadGlobals()
}

// 格式化日期
func (s *TemplateService) formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 格式化时间
func (s *TemplateService) formatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// 格式化日期时间
func (s *TemplateService) formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 截断字符串
func (s *TemplateService) truncate(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// 去除HTML标签
func (s *TemplateService) stripTags(html string) string {
	return security.StripTags(html)
}

// URL编码
func (s *TemplateService) urlEncode(str string) string {
	return url.QueryEscape(str)
}

// HTML编码
func (s *TemplateService) htmlEncode(str string) string {
	return template.HTMLEscapeString(str)
}

// 翻译
func (s *TemplateService) translate(lang, key string, args ...interface{}) string {
	return s.i18nService.T(lang, key, args...)
}

// 获取Meta标签
func (s *TemplateService) getMetaTags(title, keywords, description string) map[string]string {
	return s.seoService.GetMetaTags(title, keywords, description)
}

// 获取Open Graph标签
func (s *TemplateService) getOpenGraphTags(title, description, url, image string) map[string]string {
	return s.seoService.GetOpenGraphTags(title, description, url, image)
}

// 获取规范URL
func (s *TemplateService) getCanonicalURL(path string) string {
	return s.seoService.GetCanonicalURL(path)
}

// 获取备用URL
func (s *TemplateService) getAlternateURLs(path string) map[string]string {
	return s.seoService.GetAlternateURLs(path)
}

// 获取可用语言
func (s *TemplateService) getAvailableLangs() []map[string]string {
	return s.seoService.GetAvailableLangs()
}

// 注册标签处理器
func (s *TemplateService) registerTagHandlers() {
	// 文章列表标签
	s.engine.RegisterTag("arclist", &tags.ArcListTag{
		DB: s.db,
	})

	// 栏目标签
	s.engine.RegisterTag("channel", &tags.ChannelTag{
		DB: s.db,
	})

	// 字段标签
	s.engine.RegisterTag("field", &tags.FieldTag{})

	// 全局变量标签
	s.engine.RegisterTag("global", &tags.GlobalTag{
		Globals: s.globals,
	})

	// 包含标签
	s.engine.RegisterTag("include", &tags.IncludeTag{
		Config: &s.config.Template,
	})

	// 分页标签
	s.engine.RegisterTag("pagelist", &tags.PageListTag{})

	// 友情链接标签
	s.engine.RegisterTag("flink", &tags.FLinkTag{
		DB: s.db,
	})

	// 投票标签
	s.engine.RegisterTag("vote", &tags.VoteTag{
		DB: s.db,
	})

	// 广告标签
	s.engine.RegisterTag("myad", &tags.MyAdTag{
		DB: s.db,
	})

	// 标签标签
	s.engine.RegisterTag("tag", &tags.TagTag{
		DB: s.db,
	})

	// 评论标签
	// s.engine.RegisterTag("comment", &tags.CommentTag{
	// 	DB: s.db,
	// })

	// 专题标签
	// s.engine.RegisterTag("special", &tags.SpecialTag{
	// 	DB: s.db,
	// })

	// 会员标签
	// s.engine.RegisterTag("member", &tags.MemberTag{
	// 	DB: s.db,
	// })

	// 搜索标签
	// s.engine.RegisterTag("search", &tags.SearchTag{
	// 	DB: s.db,
	// })

	// 国际化标签
	// s.engine.RegisterTag("i18n", &tags.I18nTag{
	// 	I18nService: s.i18nService,
	// })

	// SEO标签
	// s.engine.RegisterTag("seo", &tags.SEOTag{
	// 	SEOService: s.seoService,
	// })

	// 统计标签
	// s.engine.RegisterTag("stats", &tags.StatsTag{
	// 	StatsService: s.statsService,
	// })
}
