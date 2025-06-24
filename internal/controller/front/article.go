package front

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ArticleController 前台文章控制器
type ArticleController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryService *service.CategoryService
	commentModel    *model.CommentModel
	memberService   *service.MemberService
	templateService *service.TemplateService
}

// NewArticleController 创建前台文章控制器
func NewArticleController(db *database.DB, cache cache.Cache, config *config.Config) *ArticleController {
	return &ArticleController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryService: service.NewCategoryService(db, cache, config),
		commentModel:    model.NewCommentModel(db),
		memberService:   service.NewMemberService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Show 显示文章
func (c *ArticleController) Show(w http.ResponseWriter, r *http.Request) {
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

	// 检查文章状态
	if article.ArcRank < 0 {
		http.Error(w, "Article is not published", http.StatusForbidden)
		return
	}

	// 增加文章浏览量
	go c.articleModel.IncrementClick(id)

	// 获取文章栏目
	category, err := c.categoryService.GetCategoryByID(article.TypeID)
	if err != nil {
		logger.Error("获取文章栏目失败", "id", article.TypeID, "error", err)
	}

	// 获取文章标签
	var tags []string
	if article.Keywords != "" {
		tags = strings.Split(article.Keywords, ",")
	}

	// 获取相关文章
	relatedArticles, _, err := c.articleModel.GetByTypeID(article.TypeID, 1, 5)
	if err != nil {
		logger.Error("获取相关文章失败", "error", err)
	}

	// 获取上一篇文章
	prevArticle, err := c.articleModel.GetPrevArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取上一篇文章失败", "error", err)
	}

	// 获取下一篇文章
	nextArticle, err := c.articleModel.GetNextArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取下一篇文章失败", "error", err)
	}

	// 获取文章评论
	comments, total, err := c.commentModel.GetListByAID(id, 1, 100)
	if err != nil {
		logger.Error("获取文章评论失败", "error", err)
	}

	// 广告相关
	var contentBanner string
	var sideBanner string

	// 简化处理，实际应该从广告服务获取
	contentBanner = ""
	sideBanner = ""

	// 获取会员信息
	memberID := middleware.GetMemberID(r)
	var member *model.Member
	if memberID > 0 {
		member, _ = c.memberService.GetMemberByID(memberID)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Article":        article,
		"Category":       category,
		"Tags":           tags,
		"RelatedArticles": relatedArticles,
		"PrevArticle":    prevArticle,
		"NextArticle":    nextArticle,
		"Comments":       comments,
		"CommentTotal":   total,
		"ContentBanner":  contentBanner,
		"SideBanner":     sideBanner,
		"Member":         member,
		"PageTitle":      article.Title,
	}

	// 渲染模板
	tplFile := "article.htm"
	if category != nil && category.ArticleTpl != "" {
		tplFile = category.ArticleTpl
	}
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染文章模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AddComment 添加评论
func (c *ArticleController) AddComment(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	articleIDStr := r.FormValue("articleid")
	content := r.FormValue("content")
	parentIDStr := r.FormValue("parentid")

	// 验证必填字段
	if articleIDStr == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	articleID, _ := strconv.ParseInt(articleIDStr, 10, 64)
	parentID, _ := strconv.ParseInt(parentIDStr, 10, 64)

	// 获取会员信息
	memberID := middleware.GetMemberID(r)
	if memberID <= 0 {
		http.Error(w, "Please login first", http.StatusForbidden)
		return
	}

	// 添加评论
	comment := &model.Comment{
		AID:      articleID,
		MID:      memberID,
		Content:  content,
		ParentID: parentID,
		IP:       r.RemoteAddr,
		Dtime:    time.Now(),
	}
	id, err := c.commentModel.Create(comment)
	if err != nil {
		logger.Error("添加评论失败", "error", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论添加成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/article/"+articleIDStr, http.StatusFound)
	}
}

// Detail 显示文章详情
func (c *ArticleController) Detail(w http.ResponseWriter, r *http.Request) {
	// 调用 Show 方法
	c.Show(w, r)
}

// Like 点赞文章
func (c *ArticleController) Like(w http.ResponseWriter, r *http.Request) {
	// 获取文章ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 增加文章点赞数 - 简化处理
	err = c.articleModel.IncrementClick(id)
	if err != nil {
		logger.Error("增加文章点赞数失败", "error", err)
		http.Error(w, "Failed to like article", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "点赞成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/article/"+idStr, http.StatusFound)
	}
}

// List 显示文章列表
func (c *ArticleController) List(w http.ResponseWriter, r *http.Request) {
	// 获取栏目ID
	vars := mux.Vars(r)
	typeidStr := vars["typeid"]
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// 获取分页参数
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	pageSize := 10

	// 获取栏目
	category, err := c.categoryService.GetCategoryByID(typeid)
	if err != nil {
		logger.Error("获取栏目失败", "id", typeid, "error", err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	// 获取文章列表
	articles, total, err := c.articleModel.GetByTypeID(typeid, page, pageSize)
	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
		http.Error(w, "Failed to get articles", http.StatusInternalServerError)
		return
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
		"Category":    category,
		"Articles":    articles,
		"Pagination":  pagination,
		"PageTitle":   category.TypeName,
	}

	// 渲染模板
	tplFile := "list.htm"
	if category.ListTpl != "" {
		tplFile = category.ListTpl
	}
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}