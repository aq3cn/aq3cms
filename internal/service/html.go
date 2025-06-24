package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// HtmlService HTML生成服务
type HtmlService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	templateService *TemplateService
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	productModel    *model.ProductModel
	downloadModel   *model.DownloadModel
	specialModel    *model.SpecialModel
	tagModel        *model.TagModel
	mutex           sync.Mutex
}

// NewHtmlService 创建HTML生成服务
func NewHtmlService(db *database.DB, cache cache.Cache, config *config.Config) *HtmlService {
	templateService := NewTemplateService(db, cache, config)

	return &HtmlService{
		db:              db,
		cache:           cache,
		config:          config,
		templateService: templateService,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		productModel:    model.NewProductModel(db),
		downloadModel:   model.NewDownloadModel(db),
		specialModel:    model.NewSpecialModel(db),
		tagModel:        model.NewTagModel(db),
	}
}

// GenerateIndex 生成首页
func (s *HtmlService) GenerateIndex() error {
	logger.Info("开始生成首页")

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 获取推荐文章
	recommendArticles, _, err := s.articleModel.GetList(0, 1, 10)
	if err != nil {
		logger.Error("获取推荐文章失败", "error", err)
	}

	// 获取热门文章
	hotArticles, _, err := s.articleModel.GetList(0, 1, 10)
	if err != nil {
		logger.Error("获取热门文章失败", "error", err)
	}

	// 获取顶级栏目
	topCategories, err := s.categoryModel.GetTopCategories()
	if err != nil {
		logger.Error("获取顶级栏目失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":           globals,
		"RecommendArticles": recommendArticles,
		"HotArticles":       hotArticles,
		"TopCategories":     topCategories,
		"PageTitle":         s.config.Site.Name,
		"Keywords":          s.config.Site.Keywords,
		"Description":       s.config.Site.Description,
	}

	// 渲染模板
	tplFile := s.config.Template.DefaultTpl + "/index.htm"

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, "index.html")
	if err != nil {
		logger.Error("生成静态首页失败", "error", err)
		return err
	}

	// 生成移动端首页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/index.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		err = s.templateService.GenerateStaticPage(mobileTplFile, data, "m/index.html")
		if err != nil {
			logger.Error("生成移动端静态首页失败", "error", err)
		}
	}

	logger.Info("首页生成完成")
	return nil
}

// GenerateList 生成列表页
func (s *HtmlService) GenerateList(typeid int64) error {
	logger.Info("开始生成列表页", "typeid", typeid)

	// 获取栏目信息
	category, err := s.categoryModel.GetByID(typeid)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", typeid, "error", err)
		return err
	}

	// 获取文章总数
	var total int
	switch category.ChannelType {
	case 1: // 文章
		_, total, err = s.articleModel.GetList(typeid, 1, 1)
	case 2: // 产品
		_, total, err = s.productModel.GetList(typeid, 1, 1)
	case 3: // 下载
		_, total, err = s.downloadModel.GetList(typeid, 1, 1)
	default:
		_, total, err = s.articleModel.GetList(typeid, 1, 1)
	}

	if err != nil {
		logger.Error("获取内容总数失败", "typeid", typeid, "error", err)
		return err
	}

	// 计算总页数
	pageSize := 10
	totalPages := (total + pageSize - 1) / pageSize

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 获取子栏目
	subCategories, err := s.categoryModel.GetChildCategories(typeid)
	if err != nil {
		logger.Error("获取子栏目失败", "typeid", typeid, "error", err)
	}

	// 确定模板文件
	var tplFile string
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	} else {
		tplFile = s.config.Template.DefaultTpl + "/list.htm"
	}

	// 生成每一页
	for page := 1; page <= totalPages; page++ {
		// 获取内容列表
		var articles interface{}
		switch category.ChannelType {
		case 1: // 文章
			articles, _, err = s.articleModel.GetList(typeid, page, pageSize)
		case 2: // 产品
			articles, _, err = s.productModel.GetList(typeid, page, pageSize)
		case 3: // 下载
			articles, _, err = s.downloadModel.GetList(typeid, page, pageSize)
		default:
			articles, _, err = s.articleModel.GetList(typeid, page, pageSize)
		}

		if err != nil {
			logger.Error("获取内容列表失败", "typeid", typeid, "page", page, "error", err)
			continue
		}

		// 计算分页信息
		pagination := map[string]interface{}{
			"CurrentPage": page,
			"TotalPages":  totalPages,
			"TotalItems":  total,
			"HasPrev":     page > 1,
			"HasNext":     page < totalPages,
			"PrevPage":    page - 1,
			"NextPage":    page + 1,
		}

		// 准备模板数据
		data := map[string]interface{}{
			"Globals":       globals,
			"Category":      category,
			"Articles":      articles,
			"Pagination":    pagination,
			"SubCategories": subCategories,
			"PageTitle":     category.TypeName + " - " + s.config.Site.Name,
			"Keywords":      category.Keywords,
			"Description":   category.Description,
		}

		// 生成静态页面
		var staticPath string
		if page == 1 {
			staticPath = fmt.Sprintf("list/%d.html", typeid)
		} else {
			staticPath = fmt.Sprintf("list/%d_%d.html", typeid, page)
		}

		// 确保目录存在
		dir := filepath.Dir(staticPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Error("创建目录失败", "dir", dir, "error", err)
			continue
		}

		// 生成静态页面
		err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
		if err != nil {
			logger.Error("生成静态列表页失败", "typeid", typeid, "page", page, "error", err)
			continue
		}

		// 生成移动端列表页
		mobileTplFile := s.config.Template.DefaultTpl + "/mobile/list.htm"
		if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
			mobileStaticPath := fmt.Sprintf("m/list/%d", typeid)
			if page > 1 {
				mobileStaticPath += fmt.Sprintf("_%d", page)
			}
			mobileStaticPath += ".html"

			// 确保目录存在
			mobileDir := filepath.Dir(mobileStaticPath)
			if err := os.MkdirAll(mobileDir, 0755); err != nil {
				logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
				continue
			}

			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态列表页失败", "typeid", typeid, "page", page, "error", err)
			}
		}
	}

	logger.Info("列表页生成完成", "typeid", typeid)
	return nil
}

// GenerateArticle 生成文章页
func (s *HtmlService) GenerateArticle(id int64) error {
	logger.Info("开始生成文章页", "id", id)

	// 获取文章详情
	article, err := s.articleModel.GetByID(id)
	if err != nil {
		logger.Warn("文章不存在，跳过生成", "id", id, "error", err)
		return nil // 文章不存在时返回nil，不视为错误
	}

	// 检查文章是否为空
	if article == nil {
		logger.Warn("文章为空，跳过生成", "id", id)
		return nil
	}

	// 获取栏目信息
	category, err := s.categoryModel.GetByID(article.TypeID)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", article.TypeID, "error", err)
	}

	// 获取上一篇文章
	prevArticle, err := s.articleModel.GetByID(id - 1)
	if err != nil {
		logger.Error("获取上一篇文章失败", "id", id, "error", err)
	}

	// 获取下一篇文章
	nextArticle, err := s.articleModel.GetByID(id + 1)
	if err != nil {
		logger.Error("获取下一篇文章失败", "id", id, "error", err)
	}

	// 获取相关文章
	relatedArticles, _, err := s.articleModel.GetList(article.TypeID, 1, 10)
	if err != nil {
		logger.Error("获取相关文章失败", "id", id, "error", err)
	}

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":         globals,
		"Article":         article,
		"Category":        category,
		"PrevArticle":     prevArticle,
		"NextArticle":     nextArticle,
		"RelatedArticles": relatedArticles,
		"PageTitle":       article.Title + " - " + s.config.Site.Name,
		"Keywords":        article.Keywords,
		"Description":     article.Description,
	}

	// 确定模板文件
	var tplFile string
	if category != nil && category.ArticleTpl != "" {
		tplFile = category.ArticleTpl
	} else {
		tplFile = s.config.Template.DefaultTpl + "/article.htm"
	}

	// 生成静态页面
	staticPath := fmt.Sprintf("a/%d.html", id)

	// 确保目录存在
	dir := filepath.Dir(staticPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
	if err != nil {
		logger.Error("生成静态文章页失败", "id", id, "error", err)
		return err
	}

	// 生成移动端文章页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/article.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		mobileStaticPath := fmt.Sprintf("m/a/%d.html", id)

		// 确保目录存在
		mobileDir := filepath.Dir(mobileStaticPath)
		if err := os.MkdirAll(mobileDir, 0755); err != nil {
			logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
		} else {
			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态文章页失败", "id", id, "error", err)
			}
		}
	}

	logger.Info("文章页生成完成", "id", id)
	return nil
}

// GenerateProduct 生成产品页
func (s *HtmlService) GenerateProduct(id int64) error {
	logger.Info("开始生成产品页", "id", id)

	// 获取产品详情
	product, err := s.productModel.GetByID(id)
	if err != nil {
		logger.Error("获取产品详情失败", "id", id, "error", err)
		return err
	}

	// 获取栏目信息
	category, err := s.categoryModel.GetByID(product.TypeID)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", product.TypeID, "error", err)
	}

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Product":     product,
		"Category":    category,
		"PageTitle":   product.Title + " - " + s.config.Site.Name,
		"Keywords":    product.Keywords,
		"Description": product.Description,
	}

	// 确定模板文件
	var tplFile string
	if category != nil && category.ArticleTpl != "" {
		tplFile = category.ArticleTpl
	} else {
		tplFile = s.config.Template.DefaultTpl + "/product.htm"
	}

	// 生成静态页面
	staticPath := fmt.Sprintf("p/%d.html", id)

	// 确保目录存在
	dir := filepath.Dir(staticPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
	if err != nil {
		logger.Error("生成静态产品页失败", "id", id, "error", err)
		return err
	}

	// 生成移动端产品页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/product.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		mobileStaticPath := fmt.Sprintf("m/p/%d.html", id)

		// 确保目录存在
		mobileDir := filepath.Dir(mobileStaticPath)
		if err := os.MkdirAll(mobileDir, 0755); err != nil {
			logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
		} else {
			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态产品页失败", "id", id, "error", err)
			}
		}
	}

	logger.Info("产品页生成完成", "id", id)
	return nil
}

// GenerateDownload 生成下载页
func (s *HtmlService) GenerateDownload(id int64) error {
	logger.Info("开始生成下载页", "id", id)

	// 获取下载详情
	download, err := s.downloadModel.GetByID(id)
	if err != nil {
		logger.Error("获取下载详情失败", "id", id, "error", err)
		return err
	}

	// 获取栏目信息
	category, err := s.categoryModel.GetByID(download.TypeID)
	if err != nil {
		logger.Error("获取栏目信息失败", "typeid", download.TypeID, "error", err)
	}

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Download":    download,
		"Category":    category,
		"PageTitle":   download.Title + " - " + s.config.Site.Name,
		"Keywords":    download.Keywords,
		"Description": download.Description,
	}

	// 确定模板文件
	var tplFile string
	if category != nil && category.ArticleTpl != "" {
		tplFile = category.ArticleTpl
	} else {
		tplFile = s.config.Template.DefaultTpl + "/download.htm"
	}

	// 生成静态页面
	staticPath := fmt.Sprintf("d/%d.html", id)

	// 确保目录存在
	dir := filepath.Dir(staticPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
	if err != nil {
		logger.Error("生成静态下载页失败", "id", id, "error", err)
		return err
	}

	// 生成移动端下载页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/download.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		mobileStaticPath := fmt.Sprintf("m/d/%d.html", id)

		// 确保目录存在
		mobileDir := filepath.Dir(mobileStaticPath)
		if err := os.MkdirAll(mobileDir, 0755); err != nil {
			logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
		} else {
			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态下载页失败", "id", id, "error", err)
			}
		}
	}

	logger.Info("下载页生成完成", "id", id)
	return nil
}

// GenerateAllArticles 生成所有文章页
func (s *HtmlService) GenerateAllArticles() error {
	logger.Info("开始生成所有文章页")

	// 获取所有文章
	articles, _, err := s.articleModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取所有文章失败", "error", err)
		return err
	}

	// 提取ID
	ids := make([]int64, 0, len(articles))
	for _, article := range articles {
		ids = append(ids, article.ID)
	}

	// 生成每篇文章
	for _, id := range ids {
		err := s.GenerateArticle(id)
		if err != nil {
			logger.Error("生成文章页失败", "id", id, "error", err)
		}

		// 避免服务器负载过高，每生成一篇文章后暂停一下
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("所有文章页生成完成", "count", len(ids))
	return nil
}

// GenerateAllProducts 生成所有产品页
func (s *HtmlService) GenerateAllProducts() error {
	logger.Info("开始生成所有产品页")

	// 获取所有产品
	products, _, err := s.productModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取所有产品失败", "error", err)
		return err
	}

	// 提取ID
	ids := make([]int64, 0, len(products))
	for _, product := range products {
		ids = append(ids, product.ID)
	}

	// 生成每个产品
	for _, id := range ids {
		err := s.GenerateProduct(id)
		if err != nil {
			logger.Error("生成产品页失败", "id", id, "error", err)
		}

		// 避免服务器负载过高，每生成一个产品后暂停一下
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("所有产品页生成完成", "count", len(ids))
	return nil
}

// GenerateAllDownloads 生成所有下载页
func (s *HtmlService) GenerateAllDownloads() error {
	logger.Info("开始生成所有下载页")

	// 获取所有下载
	downloads, _, err := s.downloadModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取所有下载失败", "error", err)
		return err
	}

	// 提取ID
	ids := make([]int64, 0, len(downloads))
	for _, download := range downloads {
		ids = append(ids, download.ID)
	}

	// 生成每个下载
	for _, id := range ids {
		err := s.GenerateDownload(id)
		if err != nil {
			logger.Error("生成下载页失败", "id", id, "error", err)
		}

		// 避免服务器负载过高，每生成一个下载后暂停一下
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("所有下载页生成完成", "count", len(ids))
	return nil
}

// GenerateAllLists 生成所有列表页
func (s *HtmlService) GenerateAllLists() error {
	logger.Info("开始生成所有列表页")

	// 获取所有栏目
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取所有栏目失败", "error", err)
		return err
	}

	// 提取ID
	ids := make([]int64, 0, len(categories))
	for _, category := range categories {
		ids = append(ids, category.ID)
	}

	// 生成每个栏目
	for _, id := range ids {
		err := s.GenerateList(id)
		if err != nil {
			logger.Error("生成列表页失败", "id", id, "error", err)
		}

		// 避免服务器负载过高，每生成一个栏目后暂停一下
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("所有列表页生成完成", "count", len(ids))
	return nil
}

// GenerateAll 生成所有页面
func (s *HtmlService) GenerateAll() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	logger.Info("开始生成所有页面")

	// 生成首页
	err := s.GenerateIndex()
	if err != nil {
		logger.Error("生成首页失败", "error", err)
	}

	// 生成所有列表页
	err = s.GenerateAllLists()
	if err != nil {
		logger.Error("生成所有列表页失败", "error", err)
	}

	// 生成所有文章页
	err = s.GenerateAllArticles()
	if err != nil {
		logger.Error("生成所有文章页失败", "error", err)
	}

	// 生成所有产品页
	err = s.GenerateAllProducts()
	if err != nil {
		logger.Error("生成所有产品页失败", "error", err)
	}

	// 生成所有下载页
	err = s.GenerateAllDownloads()
	if err != nil {
		logger.Error("生成所有下载页失败", "error", err)
	}

	logger.Info("所有页面生成完成")
	return nil
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
