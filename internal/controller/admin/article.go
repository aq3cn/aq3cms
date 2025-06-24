package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
	"github.com/gorilla/mux"
)

// ArticleController 文章控制器
type ArticleController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	tagModel        *model.TagModel
	htmlService     *service.HtmlService
	templateService *service.TemplateService
}

// NewArticleController 创建文章控制器
func NewArticleController(db *database.DB, cache cache.Cache, config *config.Config) *ArticleController {
	return &ArticleController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		tagModel:        model.NewTagModel(db),
		htmlService:     service.NewHtmlService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 文章管理主页
func (c *ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取文章统计信息
	totalArticles, err := c.articleModel.GetTotalCount()
	if err != nil {
		logger.Error("获取文章总数失败", "error", err)
		totalArticles = 0
	}

	// 获取今日新增文章数
	todayArticles, err := c.articleModel.GetTodayCount()
	if err != nil {
		logger.Error("获取今日文章数失败", "error", err)
		todayArticles = 0
	}

	// 获取待审核文章数
	pendingArticles, err := c.articleModel.GetPendingCount()
	if err != nil {
		logger.Error("获取待审核文章数失败", "error", err)
		pendingArticles = 0
	}

	// 获取最新文章列表（前10条）
	latestArticles, _, err := c.articleModel.GetList(0, 1, 10)
	if err != nil {
		logger.Error("获取最新文章失败", "error", err)
		latestArticles = []*model.Article{}
	}

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		categories = []*model.Category{}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":         adminID,
		"AdminName":       adminName,
		"TotalArticles":   totalArticles,
		"TodayArticles":   todayArticles,
		"PendingArticles": pendingArticles,
		"LatestArticles":  latestArticles,
		"Categories":      categories,
		"CurrentMenu":     "article",
		"PageTitle":       "文章管理",
	}

	// 渲染模板
	tplFile := "admin/article_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章管理主页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// List 文章列表
func (c *ArticleController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	typeidStr := r.URL.Query().Get("typeid")
	keyword := r.URL.Query().Get("keyword")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pagesize")

	// 解析参数
	typeid := int64(0)
	if typeidStr != "" {
		typeid, _ = strconv.ParseInt(typeidStr, 10, 64)
	}

	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	pageSize := 20
	if pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
		if pageSize < 1 {
			pageSize = 20
		}
	}

	// 获取文章列表
	var articles []*model.Article
	var total int
	var err error

	if keyword != "" {
		// 搜索文章
		articles, total, err = c.articleModel.Search(keyword, page, pageSize)
	} else {
		// 获取文章列表
		articles, total, err = c.articleModel.GetList(typeid, page, pageSize)
	}

	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
	}

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
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
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Articles":    articles,
		"Categories":  categories,
		"Pagination":  pagination,
		"TypeID":      typeid,
		"Keyword":     keyword,
		"CurrentMenu": "article",
		"PageTitle":   "文章管理",
	}

	// 渲染模板
	tplFile := "admin/article_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加文章页面
func (c *ArticleController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Categories":  categories,
		"CurrentMenu": "article",
		"PageTitle":   "添加文章",
	}

	// 渲染模板
	tplFile := "admin/article_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加文章模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加文章
func (c *ArticleController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	typeidStr := r.FormValue("typeid")
	title := r.FormValue("title")
	shortTitle := r.FormValue("shorttitle")
	color := r.FormValue("color")
	writer := r.FormValue("writer")
	source := r.FormValue("source")
	litpic := r.FormValue("litpic")
	keywords := r.FormValue("keywords")
	description := r.FormValue("description")
	body := r.FormValue("body")
	isTopStr := r.FormValue("istop")
	isRecommendStr := r.FormValue("isrecommend")
	isHotStr := r.FormValue("ishot")
	arcRankStr := r.FormValue("arcrank")
	filename := r.FormValue("filename")
	tags := r.FormValue("tags")

	// 验证必填字段
	if typeidStr == "" || title == "" || body == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	isTop := 0
	if isTopStr == "1" {
		isTop = 1
	}

	isRecommend := 0
	if isRecommendStr == "1" {
		isRecommend = 1
	}

	isHot := 0
	if isHotStr == "1" {
		isHot = 1
	}

	arcRank := 0
	if arcRankStr != "" {
		arcRank, _ = strconv.Atoi(arcRankStr)
	}

	// 处理上传的缩略图
	file, header, err := r.FormFile("litpic_upload")
	if err == nil {
		defer file.Close()

		// 生成文件名
		filename := time.Now().Format("20060102150405") + "_" + security.SanitizeFilename(header.Filename)
		filepath := "uploads/images/" + filename

		// 保存文件
		if err := saveUploadedFile(file, filepath); err != nil {
			logger.Error("保存上传文件失败", "error", err)
		} else {
			litpic = "/" + filepath
		}
	}

	// 创建文章
	article := &model.Article{
		TypeID:      typeid,
		Title:       title,
		ShortTitle:  shortTitle,
		Color:       color,
		Writer:      writer,
		Source:      source,
		LitPic:      litpic,
		PubDate:     time.Now(),
		SendDate:    time.Now(),
		Keywords:    keywords,
		Description: description,
		Filename:    filename,
		IsTop:       isTop,
		IsRecommend: isRecommend,
		IsHot:       isHot,
		ArcRank:     arcRank,
		Click:       0,
		Body:        body,
	}

	// 保存文章
	id, err := c.articleModel.Create(article)
	if err != nil {
		logger.Error("创建文章失败", "error", err)
		http.Error(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	// 处理标签
	if tags != "" {
		err = c.tagModel.UpdateArticleTags(id, tags)
		if err != nil {
			logger.Error("更新文章标签失败", "error", err)
		}
	}

	// 生成静态页面
	if c.config.Site.StaticArticle {
		go c.htmlService.GenerateArticle(id)
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "文章创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/article_list", http.StatusFound)
	}
}

// Edit 编辑文章页面
func (c *ArticleController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取文章ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取文章
	article, err := c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// 获取栏目列表
	categories, err := c.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
	}

	// 获取文章标签
	tags, err := c.tagModel.GetArticleTags(id)
	if err != nil {
		logger.Error("获取文章标签失败", "id", id, "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Article":     article,
		"Categories":  categories,
		"Tags":        strings.Join(tags, ","),
		"CurrentMenu": "article",
		"PageTitle":   "编辑文章",
	}

	// 渲染模板
	tplFile := "admin/article_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑文章模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑文章
func (c *ArticleController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取文章ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取原文章
	article, err := c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	typeidStr := r.FormValue("typeid")
	title := r.FormValue("title")
	shortTitle := r.FormValue("shorttitle")
	color := r.FormValue("color")
	writer := r.FormValue("writer")
	source := r.FormValue("source")
	litpic := r.FormValue("litpic")
	keywords := r.FormValue("keywords")
	description := r.FormValue("description")
	body := r.FormValue("body")
	isTopStr := r.FormValue("istop")
	isRecommendStr := r.FormValue("isrecommend")
	isHotStr := r.FormValue("ishot")
	arcRankStr := r.FormValue("arcrank")
	filename := r.FormValue("filename")
	tags := r.FormValue("tags")

	// 验证必填字段
	if typeidStr == "" || title == "" || body == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	isTop := 0
	if isTopStr == "1" {
		isTop = 1
	}

	isRecommend := 0
	if isRecommendStr == "1" {
		isRecommend = 1
	}

	isHot := 0
	if isHotStr == "1" {
		isHot = 1
	}

	arcRank := 0
	if arcRankStr != "" {
		arcRank, _ = strconv.Atoi(arcRankStr)
	}

	// 处理上传的缩略图
	file, header, err := r.FormFile("litpic_upload")
	if err == nil {
		defer file.Close()

		// 生成文件名
		filename := time.Now().Format("20060102150405") + "_" + security.SanitizeFilename(header.Filename)
		filepath := "uploads/images/" + filename

		// 保存文件
		if err := saveUploadedFile(file, filepath); err != nil {
			logger.Error("保存上传文件失败", "error", err)
		} else {
			litpic = "/" + filepath
		}
	}

	// 更新文章
	article.TypeID = typeid
	article.Title = title
	article.ShortTitle = shortTitle
	article.Color = color
	article.Writer = writer
	article.Source = source
	if litpic != "" {
		article.LitPic = litpic
	}
	article.Keywords = keywords
	article.Description = description
	article.Filename = filename
	article.IsTop = isTop
	article.IsRecommend = isRecommend
	article.IsHot = isHot
	article.ArcRank = arcRank
	article.Body = body

	// 保存文章
	err = c.articleModel.Update(article)
	if err != nil {
		logger.Error("更新文章失败", "error", err)
		http.Error(w, "Failed to update article", http.StatusInternalServerError)
		return
	}

	// 处理标签
	err = c.tagModel.UpdateArticleTags(id, tags)
	if err != nil {
		logger.Error("更新文章标签失败", "error", err)
	}

	// 生成静态页面
	if c.config.Site.StaticArticle {
		go c.htmlService.GenerateArticle(id)
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "文章更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/article_list", http.StatusFound)
	}
}

// Delete 删除文章
func (c *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取文章ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 删除文章
	err = c.articleModel.Delete(id)
	if err != nil {
		logger.Error("删除文章失败", "id", id, "error", err)
		http.Error(w, "Failed to delete article", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "文章删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/article_list", http.StatusFound)
	}
}

// BatchDelete 批量删除文章
func (c *ArticleController) BatchDelete(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取文章ID列表
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No articles selected", http.StatusBadRequest)
		return
	}

	// 分割ID列表
	idStrs := strings.Split(idsStr, ",")
	ids := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	// 批量删除文章
	for _, id := range ids {
		err := c.articleModel.Delete(id)
		if err != nil {
			logger.Error("删除文章失败", "id", id, "error", err)
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "文章批量删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/article_list", http.StatusFound)
	}
}

// saveUploadedFile 保存上传文件
func saveUploadedFile(file io.Reader, dst string) error {
	// 创建目录
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 复制文件内容
	_, err = io.Copy(out, file)
	return err
}
