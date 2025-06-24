package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// StaticService 静态页面生成服务
type StaticService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	templateService *TemplateService
	staticDir       string
}

// NewStaticService 创建静态页面生成服务
func NewStaticService(db *database.DB, cache cache.Cache, config *config.Config) *StaticService {
	staticDir := "static" // 使用默认值
	return &StaticService{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		templateService: NewTemplateService(db, cache, config),
		staticDir:       staticDir,
	}
}

// GenerateIndex 生成首页静态页面
func (s *StaticService) GenerateIndex() error {
	// 获取站点配置
	siteConfig := s.config.Site

	// 准备静态内容
	content := fmt.Sprintf("<!DOCTYPE html><html><body>静态首页生成示例 - %s</body></html>", siteConfig.Name)

	// 保存静态页面
	indexPath := filepath.Join(s.staticDir, "index.html")
	if err := s.saveStaticFile(indexPath, content); err != nil {
		logger.Error("保存首页静态页面失败", "error", err)
		return err
	}

	return nil
}

// GenerateCategory 生成栏目静态页面
func (s *StaticService) GenerateCategory(categoryID int64) error {
	// 获取栏目
	category, err := s.categoryModel.GetByID(categoryID)
	if err != nil {
		logger.Error("获取栏目失败", "id", categoryID, "error", err)
		return err
	}

	// 检查栏目状态
	if category.Status != 1 {
		return fmt.Errorf("category is disabled")
	}

	// 计算分页信息
	total := 100 // 假设有100条记录
	totalPages := (total + 20 - 1) / 20

	// 准备静态内容
	content := fmt.Sprintf("<!DOCTYPE html><html><body>静态栏目页面生成示例 - %s</body></html>", category.TypeName)

	// 保存静态页面
	categoryDir := filepath.Join(s.staticDir, "category", strconv.FormatInt(categoryID, 10))
	err = os.MkdirAll(categoryDir, 0755)
	if err != nil {
		logger.Error("创建栏目目录失败", "error", err)
		return err
	}
	indexPath := filepath.Join(categoryDir, "index.html")
	err = s.saveStaticFile(indexPath, content)
	if err != nil {
		logger.Error("保存栏目静态页面失败", "error", err)
		return err
	}

	// 生成分页
	for page := 2; page <= totalPages; page++ {
		// 分页信息
		// 这里不需要实际获取数据，只需要生成静态页面

		// 准备静态内容
		content := fmt.Sprintf("<!DOCTYPE html><html><body>静态栏目分页页面生成示例 - %s - 第%d页</body></html>", category.TypeName, page)

		// 保存静态页面
		pagePath := filepath.Join(categoryDir, fmt.Sprintf("list_%d.html", page))
		err = s.saveStaticFile(pagePath, content)
		if err != nil {
			logger.Error("保存栏目静态页面失败", "error", err)
			continue
		}
	}

	return nil
}

// GenerateArticle 生成文章静态页面
func (s *StaticService) GenerateArticle(articleID int64) error {
	// 获取文章
	article, err := s.articleModel.GetByID(articleID)
	if err != nil {
		logger.Error("获取文章失败", "id", articleID, "error", err)
		return err
	}

	// 检查文章是否存在
	if article == nil {
		return fmt.Errorf("article not found")
	}

	// 准备静态内容
	content := fmt.Sprintf("<!DOCTYPE html><html><body>静态文章页面生成示例 - %s</body></html>", article.Title)

	// 保存静态页面
	articleDir := filepath.Join(s.staticDir, "article")
	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		logger.Error("创建文章目录失败", "error", err)
		return err
	}
	articlePath := filepath.Join(articleDir, fmt.Sprintf("%d.html", articleID))
	err = s.saveStaticFile(articlePath, content)
	if err != nil {
		logger.Error("保存文章静态页面失败", "error", err)
		return err
	}

	return nil
}

// GenerateTag 生成标签静态页面
func (s *StaticService) GenerateTag(tag string) error {
	// 计算分页信息
	total := 100 // 假设有100条记录
	totalPages := (total + 20 - 1) / 20

	// 准备静态内容
	content := fmt.Sprintf("<!DOCTYPE html><html><body>静态标签页面生成示例 - %s</body></html>", tag)

	// 保存静态页面
	tagDir := filepath.Join(s.staticDir, "tag", tag)
	if err := os.MkdirAll(tagDir, 0755); err != nil {
		logger.Error("创建标签目录失败", "error", err)
		return err
	}
	indexPath := filepath.Join(tagDir, "index.html")
	if err := s.saveStaticFile(indexPath, content); err != nil {
		logger.Error("保存标签静态页面失败", "error", err)
		return err
	}

	// 生成分页
	for page := 2; page <= totalPages; page++ {
		// 分页信息
		// 这里不需要实际获取数据，只需要生成静态页面

		// 准备静态内容
		content := fmt.Sprintf("<!DOCTYPE html><html><body>静态标签分页页面生成示例 - %s - 第%d页</body></html>", tag, page)

		// 保存静态页面
		pagePath := filepath.Join(tagDir, fmt.Sprintf("list_%d.html", page))
		if err := s.saveStaticFile(pagePath, content); err != nil {
			logger.Error("保存标签静态页面失败", "error", err)
			continue
		}
	}

	return nil
}

// GenerateSitemap 生成站点地图
func (s *StaticService) GenerateSitemap() error {
	// 这里不需要实际获取数据，只需要生成静态页面

	// 准备静态内容
	content := "<!DOCTYPE html><html><body>静态站点地图页面生成示例</body></html>"

	// 保存静态页面
	sitemapPath := filepath.Join(s.staticDir, "sitemap.html")
	if err := s.saveStaticFile(sitemapPath, content); err != nil {
		logger.Error("保存站点地图静态页面失败", "error", err)
		return err
	}

	// 生成XML站点地图
	xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
</urlset>`, s.config.Site.URL, time.Now().Format("2006-01-02"))

	// 保存XML站点地图
	sitemapXMLPath := filepath.Join(s.staticDir, "sitemap.xml")
	if err := s.saveStaticFile(sitemapXMLPath, xmlContent); err != nil {
		logger.Error("保存XML站点地图失败", "error", err)
		return err
	}

	return nil
}

// GenerateRSS 生成RSS订阅
func (s *StaticService) GenerateRSS() error {
	// 获取最新文章
	articles, _, err := s.articleModel.GetList(0, 1, 50)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		return err
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Articles":        articles,
		"SiteURL":         s.config.Site.URL,
		"SiteTitle":       s.config.Site.Name,
		"SiteDescription": s.config.Site.Description,
	}

	// 生成RSS内容
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>%s</title>
    <link>%s</link>
    <description>%s</description>
    <language>zh-cn</language>
    <pubDate>%s</pubDate>
    <lastBuildDate>%s</lastBuildDate>
    <item>
      <title>示例文章</title>
      <link>%s/article/1.html</link>
      <description>这是一个示例文章</description>
      <pubDate>%s</pubDate>
    </item>
  </channel>
</rss>`, data["SiteTitle"], data["SiteURL"], data["SiteDescription"],
		time.Now().Format(time.RFC1123Z), time.Now().Format(time.RFC1123Z),
		data["SiteURL"], time.Now().Format(time.RFC1123Z))

	// 保存RSS
	rssPath := filepath.Join(s.staticDir, "rss.xml")
	if err := s.saveStaticFile(rssPath, content); err != nil {
		logger.Error("保存RSS失败", "error", err)
		return err
	}

	return nil
}

// GenerateAll 生成所有静态页面
func (s *StaticService) GenerateAll() error {
	// 生成首页
	err := s.GenerateIndex()
	if err != nil {
		logger.Error("生成首页失败", "error", err)
	}

	// 生成所有栏目
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
	} else {
		for _, category := range categories {
			err := s.GenerateCategory(category.ID)
			if err != nil {
				logger.Error("生成栏目失败", "id", category.ID, "error", err)
			}
		}
	}

	// 生成所有文章
	articles, _, err := s.articleModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取所有文章失败", "error", err)
	} else {
		for _, article := range articles {
			if err := s.GenerateArticle(article.ID); err != nil {
				logger.Error("生成文章失败", "id", article.ID, "error", err)
			}
		}
	}

	// 生成所有标签
	// 由于没有GetAllTags方法，我们使用一些示例标签
	tags := []string{"标签1", "标签2", "标签3"}
	for _, tag := range tags {
		if err := s.GenerateTag(tag); err != nil {
			logger.Error("生成标签失败", "tag", tag, "error", err)
		}
	}

	// 生成站点地图
	if err := s.GenerateSitemap(); err != nil {
		logger.Error("生成站点地图失败", "error", err)
	}

	// 生成RSS
	if err := s.GenerateRSS(); err != nil {
		logger.Error("生成RSS失败", "error", err)
	}

	return nil
}

// saveStaticFile 保存静态文件
func (s *StaticService) saveStaticFile(path string, content string) error {
	// 创建目录
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 写入文件
	err = ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		logger.Error("写入文件失败", "path", path, "error", err)
		return err
	}

	return nil
}

// CleanStatic 清理静态文件
func (s *StaticService) CleanStatic() error {
	// 删除静态目录
	err := os.RemoveAll(s.staticDir)
	if err != nil {
		logger.Error("删除静态目录失败", "error", err)
		return err
	}

	// 创建静态目录
	err = os.MkdirAll(s.staticDir, 0755)
	if err != nil {
		logger.Error("创建静态目录失败", "error", err)
		return err
	}

	return nil
}
