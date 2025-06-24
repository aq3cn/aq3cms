package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// CommentController 评论API控制器
type CommentController struct {
	*BaseController
}

// NewCommentController 创建评论API控制器
func NewCommentController(db *database.DB, cache cache.Cache, config *config.Config) *CommentController {
	return &CommentController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// List 评论列表
func (c *CommentController) List(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 获取文章ID
	aid, err := c.GetInt64Param(r, "aid")
	if err != nil || aid <= 0 {
		c.Error(w, 400, "Invalid article ID")
		return
	}

	// 获取查询参数
	page := c.GetQueryInt(r, "page", 1)
	pageSize := c.GetQueryInt(r, "pagesize", 20)

	// 验证参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取评论列表
	comments, err := c.commentModel.GetByAID(aid)
	if err != nil {
		logger.Error("获取评论列表失败", "error", err)
		c.Error(w, 500, "Failed to get comments")
		return
	}

	// 计算总数
	total := len(comments)

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
		"comments":   comments,
		"pagination": pagination,
	})
}

// Create 创建评论
func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 解析请求体
	var commentData struct {
		AID     int64  `json:"aid"`
		Content string `json:"content"`
		ParentID int64 `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&commentData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证必填字段
	if commentData.AID <= 0 || commentData.Content == "" {
		c.Error(w, 400, "Missing required fields")
		return
	}

	// 获取会员
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		c.Error(w, 404, "Member not found")
		return
	}

	// 获取文章
	article, err := c.articleModel.GetByID(commentData.AID)
	if err != nil {
		logger.Error("获取文章失败", "id", commentData.AID, "error", err)
		c.Error(w, 404, "Article not found")
		return
	}

	// 创建评论
	isCheck := 0
	if c.config.Site.CommentAutoCheck {
		isCheck = 1
	}

	comment := &model.Comment{
		AID:       commentData.AID,
		TypeID:    article.TypeID,
		Username:  member.Username,
		MID:       memberID,
		IP:        r.RemoteAddr,
		IsCheck:   isCheck,
		Dtime:     time.Now(),
		ParentID:  commentData.ParentID,
		Content:   security.FilterXSS(commentData.Content),
		GoodCount: 0,
		BadCount:  0,
		UserFace:  member.Avatar,
	}

	// 保存评论
	id, err := c.commentModel.Create(comment)
	if err != nil {
		logger.Error("创建评论失败", "error", err)
		c.Error(w, 500, "Failed to create comment")
		return
	}

	// 更新文章评论数
	c.articleModel.IncrementCommentCount(commentData.AID)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"id":      id,
		"message": "Comment created successfully",
		"is_check": isCheck == 1,
	})
}

// Delete 删除评论
func (c *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取评论ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid comment ID")
		return
	}

	// 获取评论
	comment, err := c.commentModel.GetByID(id)
	if err != nil {
		logger.Error("获取评论失败", "id", id, "error", err)
		c.Error(w, 404, "Comment not found")
		return
	}

	// 检查权限
	if comment.MID != memberID {
		c.Error(w, 403, "Permission denied")
		return
	}

	// 删除评论
	err = c.commentModel.Delete(id)
	if err != nil {
		logger.Error("删除评论失败", "error", err)
		c.Error(w, 500, "Failed to delete comment")
		return
	}

	// 更新文章评论数
	c.articleModel.DecrementCommentCount(comment.AID)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Comment deleted successfully",
	})
}

// Vote 评论投票
func (c *CommentController) Vote(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取评论ID
	id, err := c.GetInt64Param(r, "id")
	if err != nil || id <= 0 {
		c.Error(w, 400, "Invalid comment ID")
		return
	}

	// 解析请求体
	var voteData struct {
		Type string `json:"type"` // good or bad
	}
	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证投票类型
	if voteData.Type != "good" && voteData.Type != "bad" {
		c.Error(w, 400, "Invalid vote type")
		return
	}

	// 获取评论
	comment, err := c.commentModel.GetByID(id)
	if err != nil {
		logger.Error("获取评论失败", "id", id, "error", err)
		c.Error(w, 404, "Comment not found")
		return
	}

	// 检查是否已投票
	cacheKey := "comment_vote:" + strconv.FormatInt(memberID, 10) + ":" + strconv.FormatInt(id, 10)
	if _, ok := c.cache.Get(cacheKey); ok {
		c.Error(w, 400, "Already voted")
		return
	}

	// 更新评论投票
	if voteData.Type == "good" {
		comment.GoodCount++
	} else {
		comment.BadCount++
	}

	// 保存评论
	err = c.commentModel.Update(comment)
	if err != nil {
		logger.Error("更新评论失败", "error", err)
		c.Error(w, 500, "Failed to update comment")
		return
	}

	// 缓存投票记录
	c.cache.Set(cacheKey, true, time.Hour*24*7)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Vote successful",
		"good_count": comment.GoodCount,
		"bad_count": comment.BadCount,
	})
}
