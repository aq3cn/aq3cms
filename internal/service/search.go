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

// SearchService 搜索服务
type SearchService struct {
	db            *database.DB
	cache         cache.Cache
	config        *config.Config
	articleModel  *model.ArticleModel
	productModel  *model.ProductModel
	downloadModel *model.DownloadModel
}

// SearchResult 搜索结果
type SearchResult struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pubdate"`
	URL         string    `json:"url"`
	ChannelType int       `json:"channeltype"`
	TypeID      int64     `json:"typeid"`
	TypeName    string    `json:"typename"`
	Highlight   string    `json:"highlight"`
}

// SearchOptions 搜索选项
type SearchOptions struct {
	Keyword     string   `json:"keyword"`     // 关键词
	ChannelType int      `json:"channeltype"` // 频道类型
	TypeID      int64    `json:"typeid"`      // 栏目ID
	OrderBy     string   `json:"orderby"`     // 排序
	TimeRange   string   `json:"timerange"`   // 时间范围
	Page        int      `json:"page"`        // 页码
	PageSize    int      `json:"pagesize"`    // 每页数量
	Fields      []string `json:"fields"`      // 搜索字段
}

// NewSearchService 创建搜索服务
func NewSearchService(db *database.DB, cache cache.Cache, config *config.Config) *SearchService {
	return &SearchService{
		db:            db,
		cache:         cache,
		config:        config,
		articleModel:  model.NewArticleModel(db),
		productModel:  model.NewProductModel(db),
		downloadModel: model.NewDownloadModel(db),
	}
}

// Search 搜索
func (s *SearchService) Search(options *SearchOptions) ([]*SearchResult, int, error) {
	// 检查缓存
	cacheKey := s.generateCacheKey(options)
	if s.config.Cache.EnableSearchCache {
		if cached, ok := s.cache.Get(cacheKey); ok {
			if result, ok := cached.([]*SearchResult); ok {
				// 获取总数
				totalKey := cacheKey + "_total"
				if totalCached, ok := s.cache.Get(totalKey); ok {
					if total, ok := totalCached.(int); ok {
						return result, total, nil
					}
				}
			}
		}
	}

	// 构建查询
	qb := database.NewQueryBuilder(s.db, "archives")
	qb.Select("a.id", "a.title", "a.description", "a.pubdate", "a.typeid", "a.channel", "t.typename")

	// 添加左连接
	qb.LeftJoin(s.db.TableName("arctype")+" AS t", "a.typeid = t.id")

	// 添加条件
	qb.Where("a.arcrank > -1")

	// 添加频道类型条件
	if options.ChannelType > 0 {
		qb.Where("a.channel = ?", options.ChannelType)
	}

	// 添加栏目条件
	if options.TypeID > 0 {
		qb.Where("a.typeid = ?", options.TypeID)
	}

	// 添加关键词条件
	if options.Keyword != "" {
		// 分割关键词
		keywords := strings.Split(options.Keyword, " ")
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if keyword == "" {
				continue
			}

			// 构建搜索条件
			conditions := make([]string, 0, len(options.Fields))
			values := make([]interface{}, 0, len(options.Fields))

			// 添加搜索字段
			for _, field := range options.Fields {
				conditions = append(conditions, field+" LIKE ?")
				values = append(values, "%"+keyword+"%")
			}

			// 添加条件
			if len(conditions) > 0 {
				qb.Where("("+strings.Join(conditions, " OR ")+")", values...)
			}
		}
	}

	// 添加时间范围条件
	if options.TimeRange != "" {
		switch options.TimeRange {
		case "day":
			qb.Where("a.pubdate >= ?", time.Now().AddDate(0, 0, -1))
		case "week":
			qb.Where("a.pubdate >= ?", time.Now().AddDate(0, 0, -7))
		case "month":
			qb.Where("a.pubdate >= ?", time.Now().AddDate(0, -1, 0))
		case "year":
			qb.Where("a.pubdate >= ?", time.Now().AddDate(-1, 0, 0))
		}
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询搜索总数失败", "error", err)
		return nil, 0, err
	}

	// 添加排序
	switch options.OrderBy {
	case "pubdate":
		qb.OrderBy("a.pubdate DESC")
	case "click":
		qb.OrderBy("a.click DESC")
	case "id":
		qb.OrderBy("a.id DESC")
	default:
		qb.OrderBy("a.pubdate DESC")
	}

	// 添加分页
	offset := (options.Page - 1) * options.PageSize
	qb.Limit(options.PageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询搜索结果失败", "error", err)
		return nil, 0, err
	}

	// 转换为搜索结果
	searchResults := make([]*SearchResult, 0, len(results))
	for _, result := range results {
		searchResult := &SearchResult{}
		searchResult.ID, _ = result["id"].(int64)
		searchResult.Title, _ = result["title"].(string)
		searchResult.Description, _ = result["description"].(string)

		// 处理日期
		if pubdate, ok := result["pubdate"].(time.Time); ok {
			searchResult.PubDate = pubdate
		}

		searchResult.TypeID, _ = result["typeid"].(int64)
		searchResult.TypeName, _ = result["typename"].(string)

		// 处理频道类型
		if channel, ok := result["channel"].(int64); ok {
			searchResult.ChannelType = int(channel)
		}

		// 生成URL
		switch searchResult.ChannelType {
		case 1: // 文章
			searchResult.URL = fmt.Sprintf("/article/%d.html", searchResult.ID)
		case 2: // 产品
			searchResult.URL = fmt.Sprintf("/product/%d.html", searchResult.ID)
		case 3: // 下载
			searchResult.URL = fmt.Sprintf("/download/%d.html", searchResult.ID)
		default:
			searchResult.URL = fmt.Sprintf("/article/%d.html", searchResult.ID)
		}

		// 生成高亮
		if options.Keyword != "" {
			// 分割关键词
			keywords := strings.Split(options.Keyword, " ")
			highlight := searchResult.Description
			if highlight == "" {
				// 获取内容
				var content string
				switch searchResult.ChannelType {
				case 1: // 文章
					article, err := s.articleModel.GetByID(searchResult.ID)
					if err == nil && article != nil {
						content = article.Body
					}
				case 2: // 产品
					product, err := s.productModel.GetByID(searchResult.ID)
					if err == nil && product != nil {
						content = product.Body
					}
				case 3: // 下载
					download, err := s.downloadModel.GetByID(searchResult.ID)
					if err == nil && download != nil {
						content = download.Body
					}
				}

				// 截取内容
				if content != "" {
					// 查找关键词位置
					pos := -1
					for _, keyword := range keywords {
						keyword = strings.TrimSpace(keyword)
						if keyword == "" {
							continue
						}

						// 查找关键词位置
						p := strings.Index(strings.ToLower(content), strings.ToLower(keyword))
						if p >= 0 && (pos < 0 || p < pos) {
							pos = p
						}
					}

					// 截取内容
					if pos >= 0 {
						start := pos - 50
						if start < 0 {
							start = 0
						}
						end := pos + 100
						if end > len(content) {
							end = len(content)
						}
						highlight = content[start:end]
					} else {
						// 截取前200个字符
						if len(content) > 200 {
							highlight = content[:200]
						} else {
							highlight = content
						}
					}
				}
			}

			// 高亮关键词
			for _, keyword := range keywords {
				keyword = strings.TrimSpace(keyword)
				if keyword == "" {
					continue
				}

				// 替换关键词
				highlight = strings.Replace(highlight, keyword, "<span class=\"highlight\">"+keyword+"</span>", -1)
			}

			searchResult.Highlight = highlight
		}

		searchResults = append(searchResults, searchResult)
	}

	// 缓存搜索结果
	if s.config.Cache.EnableSearchCache {
		cache.SafeSet(s.cache, cacheKey, searchResults, time.Duration(s.config.Cache.SearchCacheTime)*time.Second)
		cache.SafeSet(s.cache, cacheKey+"_total", total, time.Duration(s.config.Cache.SearchCacheTime)*time.Second)
	}

	return searchResults, total, nil
}

// generateCacheKey 生成缓存键
func (s *SearchService) generateCacheKey(options *SearchOptions) string {
	return fmt.Sprintf("search:%s:%d:%d:%s:%s:%d:%d:%s",
		options.Keyword,
		options.ChannelType,
		options.TypeID,
		options.OrderBy,
		options.TimeRange,
		options.Page,
		options.PageSize,
		strings.Join(options.Fields, ","),
	)
}

// ClearCache 清除缓存
func (s *SearchService) ClearCache() {
	// 清除搜索缓存
	s.cache.Delete("search:*")
}
