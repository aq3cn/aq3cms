package service

import (
	"fmt"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SEOService SEO服务
type SEOService struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	articleModel  *model.ArticleModel
	categoryModel *model.CategoryModel
	tagModel      *model.TagModel
}

// NewSEOService 创建SEO服务
func NewSEOService(db *database.DB, cache cache.Cache, config *config.Config) *SEOService {
	return &SEOService{
		db:            db,
		cache:         cache,
		config:        config,
		articleModel:  model.NewArticleModel(db),
		categoryModel: model.NewCategoryModel(db),
		tagModel:      model.NewTagModel(db),
	}
}

// GenerateSitemap 生成站点地图
func (s *SEOService) GenerateSitemap() (string, error) {
	// 缓存键
	cacheKey := "sitemap"

	// 检查缓存
	if s.config.Cache.EnableSitemapCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if sitemap, ok := cached.(string); ok {
				return sitemap, nil
			}
		}
	}

	// 构建站点地图
	var sitemap strings.Builder

	// 添加XML头
	sitemap.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`)

	// 添加首页
	sitemap.WriteString(fmt.Sprintf(`  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
`, s.config.Site.URL, time.Now().Format("2006-01-02")))

	// 添加栏目页
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
	} else {
		for _, category := range categories {
			sitemap.WriteString(fmt.Sprintf(`  <url>
    <loc>%s/list/%d.html</loc>
    <lastmod>%s</lastmod>
    <changefreq>daily</changefreq>
    <priority>0.8</priority>
  </url>
`, s.config.Site.URL, category.ID, time.Now().Format("2006-01-02")))
		}
	}

	// 添加文章页
	articles, _, err := s.articleModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
	} else {
		for _, article := range articles {
			sitemap.WriteString(fmt.Sprintf(`  <url>
    <loc>%s/article/%d.html</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.6</priority>
  </url>
`, s.config.Site.URL, article.ID, article.PubDate.Format("2006-01-02")))
		}
	}

	// 添加标签页
	tags, err := s.tagModel.GetHotTags(100)
	if err != nil {
		logger.Error("获取标签列表失败", "error", err)
	} else {
		for _, tag := range tags {
			sitemap.WriteString(fmt.Sprintf(`  <url>
    <loc>%s/tag/%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.4</priority>
  </url>
`, s.config.Site.URL, tag.Tag, time.Now().Format("2006-01-02")))
		}
	}

	// 添加XML尾
	sitemap.WriteString(`</urlset>`)

	// 缓存站点地图
	if s.config.Cache.EnableSitemapCache {
		cache.SafeSet(s.cache, cacheKey, sitemap.String(), time.Duration(s.config.Cache.SitemapCacheTime)*time.Second)
	}

	return sitemap.String(), nil
}

// GenerateRSS 生成RSS
func (s *SEOService) GenerateRSS() (string, error) {
	// 缓存键
	cacheKey := "rss"

	// 检查缓存
	if s.config.Cache.EnableRSSCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if rss, ok := cached.(string); ok {
				return rss, nil
			}
		}
	}

	// 构建RSS
	var rss strings.Builder

	// 添加XML头
	rss.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>`)
	rss.WriteString(s.config.Site.Name)
	rss.WriteString(`</title>
    <link>`)
	rss.WriteString(s.config.Site.URL)
	rss.WriteString(`</link>
    <description>`)
	rss.WriteString(s.config.Site.Description)
	rss.WriteString(`</description>
    <language>zh-cn</language>
    <pubDate>`)
	rss.WriteString(time.Now().Format(time.RFC1123Z))
	rss.WriteString(`</pubDate>
    <lastBuildDate>`)
	rss.WriteString(time.Now().Format(time.RFC1123Z))
	rss.WriteString(`</lastBuildDate>
    <docs>http://blogs.law.harvard.edu/tech/rss</docs>
    <generator>aq3cms</generator>
    <managingEditor>`)
	rss.WriteString(s.config.Site.Email)
	rss.WriteString(`</managingEditor>
    <webMaster>`)
	rss.WriteString(s.config.Site.Email)
	rss.WriteString(`</webMaster>
`)

	// 添加文章
	articles, _, err := s.articleModel.GetList(0, 1, 20)
	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
	} else {
		for _, article := range articles {
			rss.WriteString(`    <item>
      <title>`)
			rss.WriteString(article.Title)
			rss.WriteString(`</title>
      <link>`)
			rss.WriteString(fmt.Sprintf("%s/article/%d.html", s.config.Site.URL, article.ID))
			rss.WriteString(`</link>
      <description>`)
			rss.WriteString(article.Description)
			rss.WriteString(`</description>
      <pubDate>`)
			rss.WriteString(article.PubDate.Format(time.RFC1123Z))
			rss.WriteString(`</pubDate>
      <guid>`)
			rss.WriteString(fmt.Sprintf("%s/article/%d.html", s.config.Site.URL, article.ID))
			rss.WriteString(`</guid>
    </item>
`)
		}
	}

	// 添加XML尾
	rss.WriteString(`  </channel>
</rss>`)

	// 缓存RSS
	if s.config.Cache.EnableRSSCache {
		cache.SafeSet(s.cache, cacheKey, rss.String(), time.Duration(s.config.Cache.RSSCacheTime)*time.Second)
	}

	return rss.String(), nil
}

// GenerateCategoryRSS 生成栏目RSS
func (s *SEOService) GenerateCategoryRSS(categoryID int64) (string, error) {
	// 缓存键
	cacheKey := fmt.Sprintf("rss_category_%d", categoryID)

	// 检查缓存
	if s.config.Cache.EnableRSSCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if rss, ok := cached.(string); ok {
				return rss, nil
			}
		}
	}

	// 获取栏目
	category, err := s.categoryModel.GetByID(categoryID)
	if err != nil {
		logger.Error("获取栏目失败", "id", categoryID, "error", err)
		return "", err
	}

	// 构建RSS
	var rss strings.Builder

	// 添加XML头
	rss.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>`)
	rss.WriteString(category.TypeName)
	rss.WriteString(` - `)
	rss.WriteString(s.config.Site.Name)
	rss.WriteString(`</title>
    <link>`)
	rss.WriteString(fmt.Sprintf("%s/list/%d.html", s.config.Site.URL, category.ID))
	rss.WriteString(`</link>
    <description>`)
	rss.WriteString(category.Description)
	rss.WriteString(`</description>
    <language>zh-cn</language>
    <pubDate>`)
	rss.WriteString(time.Now().Format(time.RFC1123Z))
	rss.WriteString(`</pubDate>
    <lastBuildDate>`)
	rss.WriteString(time.Now().Format(time.RFC1123Z))
	rss.WriteString(`</lastBuildDate>
    <docs>http://blogs.law.harvard.edu/tech/rss</docs>
    <generator>aq3cms</generator>
    <managingEditor>`)
	rss.WriteString(s.config.Site.Email)
	rss.WriteString(`</managingEditor>
    <webMaster>`)
	rss.WriteString(s.config.Site.Email)
	rss.WriteString(`</webMaster>
`)

	// 添加文章
	articles, _, err := s.articleModel.GetList(categoryID, 1, 20)
	if err != nil {
		logger.Error("获取文章列表失败", "typeid", categoryID, "error", err)
	} else {
		for _, article := range articles {
			rss.WriteString(`    <item>
      <title>`)
			rss.WriteString(article.Title)
			rss.WriteString(`</title>
      <link>`)
			rss.WriteString(fmt.Sprintf("%s/article/%d.html", s.config.Site.URL, article.ID))
			rss.WriteString(`</link>
      <description>`)
			rss.WriteString(article.Description)
			rss.WriteString(`</description>
      <pubDate>`)
			rss.WriteString(article.PubDate.Format(time.RFC1123Z))
			rss.WriteString(`</pubDate>
      <guid>`)
			rss.WriteString(fmt.Sprintf("%s/article/%d.html", s.config.Site.URL, article.ID))
			rss.WriteString(`</guid>
    </item>
`)
		}
	}

	// 添加XML尾
	rss.WriteString(`  </channel>
</rss>`)

	// 缓存RSS
	if s.config.Cache.EnableRSSCache {
		cache.SafeSet(s.cache, cacheKey, rss.String(), time.Duration(s.config.Cache.RSSCacheTime)*time.Second)
	}

	return rss.String(), nil
}

// GenerateRobotsTxt 生成robots.txt
func (s *SEOService) GenerateRobotsTxt() string {
	var robotsTxt strings.Builder

	robotsTxt.WriteString(`User-agent: *
Disallow: /aq3cms/
Disallow: /member/
Disallow: /search
Allow: /

Sitemap: `)
	robotsTxt.WriteString(s.config.Site.URL)
	robotsTxt.WriteString(`/sitemap.xml
`)

	return robotsTxt.String()
}

// GetMetaTags 获取Meta标签
func (s *SEOService) GetMetaTags(title, keywords, description string) map[string]string {
	metaTags := make(map[string]string)

	// 标题
	if title == "" {
		title = s.config.Site.Name
	} else {
		title = title + " - " + s.config.Site.Name
	}
	metaTags["title"] = title

	// 关键词
	if keywords == "" {
		keywords = s.config.Site.Keywords
	}
	metaTags["keywords"] = keywords

	// 描述
	if description == "" {
		description = s.config.Site.Description
	}
	metaTags["description"] = description

	// 其他Meta标签
	metaTags["author"] = s.config.Site.Name
	metaTags["generator"] = "aq3cms"
	metaTags["robots"] = "index,follow"
	metaTags["viewport"] = "width=device-width, initial-scale=1.0"

	return metaTags
}

// GetOpenGraphTags 获取Open Graph标签
func (s *SEOService) GetOpenGraphTags(title, description, url, image string) map[string]string {
	ogTags := make(map[string]string)

	// 标题
	if title == "" {
		title = s.config.Site.Name
	}
	ogTags["og:title"] = title

	// 描述
	if description == "" {
		description = s.config.Site.Description
	}
	ogTags["og:description"] = description

	// URL
	if url == "" {
		url = s.config.Site.URL
	}
	ogTags["og:url"] = url

	// 图片
	if image == "" {
		image = s.config.Site.URL + "/static/images/logo.png"
	}
	ogTags["og:image"] = image

	// 其他Open Graph标签
	ogTags["og:type"] = "website"
	ogTags["og:site_name"] = s.config.Site.Name

	return ogTags
}

// GetTwitterCardTags 获取Twitter Card标签
func (s *SEOService) GetTwitterCardTags(title, description, image string) map[string]string {
	twitterTags := make(map[string]string)

	// 标题
	if title == "" {
		title = s.config.Site.Name
	}
	twitterTags["twitter:title"] = title

	// 描述
	if description == "" {
		description = s.config.Site.Description
	}
	twitterTags["twitter:description"] = description

	// 图片
	if image == "" {
		image = s.config.Site.URL + "/static/images/logo.png"
	}
	twitterTags["twitter:image"] = image

	// 其他Twitter Card标签
	twitterTags["twitter:card"] = "summary_large_image"
	twitterTags["twitter:site"] = "@" + s.config.Site.Name

	return twitterTags
}

// GetCanonicalURL 获取规范URL
func (s *SEOService) GetCanonicalURL(path string) string {
	if path == "" {
		return s.config.Site.URL
	}
	return s.config.Site.URL + path
}

// GetAlternateURLs 获取备用URL
func (s *SEOService) GetAlternateURLs(path string) map[string]string {
	alternateURLs := make(map[string]string)

	// 获取可用语言
	langs := s.GetAvailableLangs()
	for _, lang := range langs {
		if path == "" {
			alternateURLs[lang["Code"]] = s.config.Site.URL + "?lang=" + lang["Code"]
		} else {
			alternateURLs[lang["Code"]] = s.config.Site.URL + path + "?lang=" + lang["Code"]
		}
	}

	return alternateURLs
}

// GetAvailableLangs 获取可用语言
func (s *SEOService) GetAvailableLangs() []map[string]string {
	// 创建I18n服务
	i18nService := NewI18nService(s.db, s.cache, s.config)
	return i18nService.GetAvailableLangs()
}
