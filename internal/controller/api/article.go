package api

import (
	"encoding/json"
	"net/http"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ArticleController 文章API控制器
type ArticleController struct {
	*BaseController
}

// NewArticleController 创建文章API控制器
func NewArticleController(db *database.DB, cache cache.Cache, config *config.Config) *ArticleController {
	return &ArticleController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// List 文章列表
func (c *ArticleController) List(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取查询参数
	typeID := c.GetQueryInt64(r, "typeid", 0)
	page := c.GetQueryInt(r, "page", 1)
	pageSize := c.GetQueryInt(r, "pagesize", 20)
	keyword := c.GetQueryString(r, "keyword", "")
	_ = c.GetQueryInt(r, "flag", 0) // 未使用的变量
	_ = c.GetQueryString(r, "orderby", "pubdate") // 未使用的变量
	_ = c.GetQueryString(r, "ordertype", "desc") // 未使用的变量

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取文章列表
	var articles []*model.Article
	var total int
	var err error

	if keyword != "" {
		// 搜索文章
		articles, total, err = c.articleModel.Search(keyword, page, pageSize)
	} else if typeID > 0 {
		// 获取栏目文章
		articles, total, err = c.articleModel.GetByTypeID(typeID, page, pageSize)
	} else {
		// 获取所有文章
		articles, total, err = c.articleModel.GetList(0, page, pageSize)
	}

	if err != nil {
		logger.Error("获取文章列表失败", "error", err)
		c.Error(w, 500, "Failed to get articles")
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := map[string]interface{}{
		"current_page": page,
		"total_pages":  totalPages,
		"total_items":  total,
		"page_size":    pageSize,
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"articles":   articles,
		"pagination": pagination,
	})
}

// Detail 文章详情
func (c *ArticleController) Detail(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取文章ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid article ID")
		return
	}

	// 获取文章
	article, err := c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		c.Error(w, 404, "Article not found")
		return
	}

	// 更新点击量
	c.articleModel.IncrementClick(id)

	// 获取栏目
	category, err := c.categoryModel.GetByID(article.TypeID)
	if err != nil {
		logger.Error("获取栏目失败", "id", article.TypeID, "error", err)
	}

	// 获取标签
	tags, err := c.tagModel.GetByAID(id)
	if err != nil {
		logger.Error("获取标签失败", "id", id, "error", err)
	}

	// 获取相关文章
	relatedArticles, err := c.articleModel.GetRelatedArticles(article.Keywords, id, 5)
	if err != nil {
		logger.Error("获取相关文章失败", "id", id, "error", err)
	}

	// 获取上一篇文章
	prevArticle, err := c.articleModel.GetPrevArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取上一篇文章失败", "id", id, "error", err)
	}

	// 获取下一篇文章
	nextArticle, err := c.articleModel.GetNextArticle(id, article.TypeID)
	if err != nil {
		logger.Error("获取下一篇文章失败", "id", id, "error", err)
	}

	// 获取评论
	comments, err := c.commentModel.GetByAID(id)
	if err != nil {
		logger.Error("获取评论失败", "id", id, "error", err)
	}

	// 获取扩展模型内容
	var modelContent map[string]interface{}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"article":         article,
		"category":        category,
		"tags":            tags,
		"related":         relatedArticles,
		"prev":            prevArticle,
		"next":            nextArticle,
		"comments":        comments,
		"model_content":   modelContent,
	})
}

// Create 创建文章
func (c *ArticleController) Create(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 解析请求体
	var article model.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证必填字段
	if article.Title == "" || article.TypeID <= 0 {
		c.Error(w, 400, "Missing required fields")
		return
	}

	// 设置默认值
	article.Click = 0
	article.IsTop = 0
	article.IsRecommend = 0
	article.IsHot = 0
	article.Status = 0 // 待审核
	article.PubDate = time.Now()
	article.SendDate = time.Now()
	article.UpdateDate = time.Now()

	// 创建文章
	id, err := c.articleModel.Create(&article)
	if err != nil {
		logger.Error("创建文章失败", "error", err)
		c.Error(w, 500, "Failed to create article")
		return
	}

	// 处理标签
	if article.Tags != "" {
		c.tagModel.UpdateTags(id, article.Tags)
	}

	// 处理扩展模型内容

	// 返回数据
	c.Success(w, map[string]interface{}{
		"id":      id,
		"message": "Article created successfully",
	})
}

// Update 更新文章
func (c *ArticleController) Update(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取文章ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid article ID")
		return
	}

	// 获取原文章
	article, err := c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		c.Error(w, 404, "Article not found")
		return
	}

	// 检查权限
	if article.MemberID != memberID {
		c.Error(w, 403, "Permission denied")
		return
	}

	// 解析请求体
	var updateArticle model.Article
	if err := json.NewDecoder(r.Body).Decode(&updateArticle); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 更新文章
	article.Title = updateArticle.Title
	article.ShortTitle = updateArticle.ShortTitle
	article.Color = updateArticle.Color
	article.Source = updateArticle.Source
	article.Author = updateArticle.Author
	article.LitPic = updateArticle.LitPic
	article.Description = updateArticle.Description
	article.Keywords = updateArticle.Keywords
	article.Content = updateArticle.Content
	article.UpdateDate = time.Now()

	// 保存文章
	err = c.articleModel.Update(article)
	if err != nil {
		logger.Error("更新文章失败", "error", err)
		c.Error(w, 500, "Failed to update article")
		return
	}

	// 处理标签
	if updateArticle.Tags != "" {
		c.tagModel.UpdateTags(id, updateArticle.Tags)
	}

	// 处理扩展模型内容

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Article updated successfully",
	})
}

// Delete 删除文章
func (c *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取文章ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid article ID")
		return
	}

	// 获取原文章
	_, err = c.articleModel.GetByID(id)
	if err != nil {
		logger.Error("获取文章失败", "id", id, "error", err)
		c.Error(w, 404, "Article not found")
		return
	}

	// 检查权限

	// 删除文章
	err = c.articleModel.Delete(id)
	if err != nil {
		logger.Error("删除文章失败", "error", err)
		c.Error(w, 500, "Failed to delete article")
		return
	}

	// 删除标签关联
	c.tagModel.DeleteByAID(id)

	// 删除扩展模型内容

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Article deleted successfully",
	})
}
