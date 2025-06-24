package frontend

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// RSSController RSS控制器
type RSSController struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	articleModel  *model.ArticleModel
	categoryModel *model.CategoryModel
}

// NewRSSController 创建RSS控制器
func NewRSSController(db *database.DB, cache cache.Cache, config *config.Config) *RSSController {
	return &RSSController{
		db:            db,
		cache:         cache,
		config:        config,
		articleModel:  model.NewArticleModel(db),
		categoryModel: model.NewCategoryModel(db),
	}
}

// RSS 结构体
type RSS struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	Channel     Channel  `xml:"channel"`
}

// Channel 结构体
type Channel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []Item    `xml:"item"`
}

// Item 结构体
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
}

// Index 首页RSS
func (c *RSSController) Index(w http.ResponseWriter, r *http.Request) {
	// 检查缓存
	if cached, ok := c.cache.Get("rss:index"); ok {
		if rssXML, ok := cached.(string); ok {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(rssXML))
			return
		}
	}

	// 获取最新文章
	articles, err := c.articleModel.GetLatestArticles(20)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 创建RSS
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         c.config.Site.Name,
			Link:          c.config.Site.URL,
			Description:   c.config.Site.Description,
			Language:      "zh-cn",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         make([]Item, 0, len(articles)),
		},
	}

	// 添加文章
	for _, article := range articles {
		// 清理内容，防止XSS攻击
		description := security.StripTags(article.Description)
		if description == "" {
			// 如果描述为空，使用内容的前200个字符
			description = security.StripTags(article.Body)
			if len(description) > 200 {
				description = description[:200] + "..."
			}
		}

		item := Item{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/article/%d.html", c.config.Site.URL, article.ID),
			Description: description,
			PubDate:     article.PubDate.Format(time.RFC1123Z),
			GUID:        fmt.Sprintf("%s/article/%d.html", c.config.Site.URL, article.ID),
			Author:      article.Writer,
		}

		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	// 生成XML
	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		logger.Error("生成RSS XML失败", "error", err)
		http.Error(w, "Failed to generate RSS", http.StatusInternalServerError)
		return
	}

	// 添加XML头
	xmlOutput := []byte(xml.Header + string(output))

	// 缓存RSS
	c.cache.Set("rss:index", string(xmlOutput), time.Hour)

	// 输出RSS
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlOutput)
}

// Category 栏目RSS
func (c *RSSController) Category(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeidStr := vars["typeid"]
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 检查缓存
	cacheKey := fmt.Sprintf("rss:category:%d", typeid)
	if cached, ok := c.cache.Get(cacheKey); ok {
		if rssXML, ok := cached.(string); ok {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(rssXML))
			return
		}
	}

	// 获取栏目信息
	category, err := c.categoryModel.GetByID(typeid)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", typeid, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 获取栏目文章
	articles, _, err := c.articleModel.GetList(typeid, 1, 20)
	if err != nil {
		logger.Error("获取栏目文章失败", "typeid", typeid, "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 创建RSS
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         category.TypeName + " - " + c.config.Site.Name,
			Link:          fmt.Sprintf("%s/list/%d.html", c.config.Site.URL, typeid),
			Description:   category.Description,
			Language:      "zh-cn",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         make([]Item, 0, len(articles)),
		},
	}

	// 添加文章
	for _, article := range articles {
		// 清理内容，防止XSS攻击
		description := security.StripTags(article.Description)
		if description == "" {
			// 如果描述为空，使用内容的前200个字符
			description = security.StripTags(article.Body)
			if len(description) > 200 {
				description = description[:200] + "..."
			}
		}

		item := Item{
			Title:       article.Title,
			Link:        fmt.Sprintf("%s/article/%d.html", c.config.Site.URL, article.ID),
			Description: description,
			PubDate:     article.PubDate.Format(time.RFC1123Z),
			GUID:        fmt.Sprintf("%s/article/%d.html", c.config.Site.URL, article.ID),
			Author:      article.Writer,
			Category:    category.TypeName,
		}

		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	// 生成XML
	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		logger.Error("生成RSS XML失败", "error", err)
		http.Error(w, "Failed to generate RSS", http.StatusInternalServerError)
		return
	}

	// 添加XML头
	xmlOutput := []byte(xml.Header + string(output))

	// 缓存RSS
	c.cache.Set(cacheKey, string(xmlOutput), time.Hour)

	// 输出RSS
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlOutput)
}
