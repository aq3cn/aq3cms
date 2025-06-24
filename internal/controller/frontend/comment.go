/*
 * @Author: xxx@xxx.com
 * @Date: 2025-05-01 10:19:36
 * @LastEditors: xxx@xxx.com
 * @LastEditTime: 2025-05-01 10:37:17
 * @FilePath: \aq3cms\aq3cms\internal\controller\frontend\comment.go
 * @Description:
 *
 * Copyright (c) 2022 by xxx@xxx.com, All Rights Reserved.
 */
package frontend

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// CommentController 评论控制器
type CommentController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	commentModel    *model.CommentModel
	articleModel    *model.ArticleModel
	memberModel     *model.MemberModel
	templateService *service.TemplateService
}

// NewCommentController 创建评论控制器
func NewCommentController(db *database.DB, cache cache.Cache, config *config.Config) *CommentController {
	return &CommentController{
		db:              db,
		cache:           cache,
		config:          config,
		commentModel:    model.NewCommentModel(db),
		articleModel:    model.NewArticleModel(db),
		memberModel:     model.NewMemberModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 评论列表
func (c *CommentController) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aidStr := vars["aid"]
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取评论列表
	pageSize := 10
	comments, total, err := c.commentModel.GetListByAID(aid, page, pageSize)
	if err != nil {
		logger.Error("获取评论列表失败", "aid", aid, "error", err)
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
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

	// 获取文章信息
	article, err := c.articleModel.GetByID(aid)
	if err != nil {
		logger.Error("获取文章信息失败", "aid", aid, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":    globals,
		"Comments":   comments,
		"Pagination": pagination,
		"Article":    article,
		"PageTitle":  "评论列表 - " + c.config.Site.Name,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/comment_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染评论列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Post 发表评论
func (c *CommentController) Post(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	aidStr := r.FormValue("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	typeidStr := r.FormValue("typeid")
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	// 清理评论内容，防止XSS攻击
	content = security.CleanHTML(content)

	// 获取会员信息
	var mid int64
	var username string
	var userface string

	if middleware.IsMemberLoggedIn(r) {
		mid = middleware.GetMemberID(r)
		username = middleware.GetMemberName(r)

		// 获取会员信息
		member, err := c.memberModel.GetByID(mid)
		if err == nil && member != nil {
			userface = member.Face
		}
	} else {
		// 匿名评论
		username = r.FormValue("username")
		if username == "" {
			username = "游客"
		}
	}

	// 获取父评论ID
	parentIDStr := r.FormValue("parentid")
	var parentID int64
	if parentIDStr != "" {
		parentID, _ = strconv.ParseInt(parentIDStr, 10, 64)
	}

	// 获取评分
	scoreStr := r.FormValue("score")
	score := 0
	if scoreStr != "" {
		score, _ = strconv.Atoi(scoreStr)
		if score < 1 || score > 5 {
			score = 5
		}
	}

	// 创建评论
	isCheck := 0
	if c.config.Site.CommentAutoCheck {
		isCheck = 1
	}

	comment := &model.Comment{
		AID:         aid,
		TypeID:      typeid,
		Username:    username,
		MID:         mid,
		IP:          r.RemoteAddr,
		IsCheck:     isCheck, // 根据配置决定是否自动审核
		Dtime:       time.Now(),
		Content:     content,
		ParentID:    parentID,
		Score:       score,
		GoodCount:   0,
		BadCount:    0,
		UserFace:    userface,
		ChannelType: 1, // 默认为文章频道
	}

	// 保存评论
	_, err = c.commentModel.Create(comment)
	if err != nil {
		logger.Error("保存评论失败", "error", err)
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论提交成功，" + (map[bool]string{true: "已发布", false: "等待审核"})[c.config.Site.CommentAutoCheck],
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/article/"+aidStr+".html", http.StatusFound)
	}
}

// Vote 评论投票
func (c *CommentController) Vote(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")
	if action != "good" && action != "bad" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// 检查是否已投票
	cookie, err := r.Cookie("comment_vote_" + idStr)
	if err == nil && cookie.Value != "" {
		// 已投票
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "您已经投过票了",
		})
		return
	}

	// 更新评论投票数
	var err2 error
	if action == "good" {
		err2 = c.commentModel.UpdateGoodCount(id, 1)
	} else {
		err2 = c.commentModel.UpdateBadCount(id, 1)
	}

	if err2 != nil {
		logger.Error("更新评论投票数失败", "id", id, "action", action, "error", err2)
		http.Error(w, "Failed to update vote count", http.StatusInternalServerError)
		return
	}

	// 设置Cookie，防止重复投票
	http.SetCookie(w, &http.Cookie{
		Name:     "comment_vote_" + idStr,
		Value:    "1",
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	// 返回成功信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "投票成功",
	})
}

// Reply 回复评论
func (c *CommentController) Reply(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	aidStr := r.FormValue("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	typeidStr := r.FormValue("typeid")
	typeid, err := strconv.ParseInt(typeidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	parentIDStr := r.FormValue("parentid")
	parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid parent comment ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Reply content cannot be empty", http.StatusBadRequest)
		return
	}

	// 清理评论内容，防止XSS攻击
	content = security.CleanHTML(content)

	// 获取会员信息
	var mid int64
	var username string
	var userface string

	if middleware.IsMemberLoggedIn(r) {
		mid = middleware.GetMemberID(r)
		username = middleware.GetMemberName(r)

		// 获取会员信息
		member, err := c.memberModel.GetByID(mid)
		if err == nil && member != nil {
			userface = member.Face
		}
	} else {
		// 匿名评论
		username = r.FormValue("username")
		if username == "" {
			username = "游客"
		}
	}

	// 创建评论
	isCheck := 0
	if c.config.Site.CommentAutoCheck {
		isCheck = 1
	}

	comment := &model.Comment{
		AID:         aid,
		TypeID:      typeid,
		Username:    username,
		MID:         mid,
		IP:          r.RemoteAddr,
		IsCheck:     isCheck, // 根据配置决定是否自动审核
		Dtime:       time.Now(),
		Content:     content,
		ParentID:    parentID,
		Score:       0,
		GoodCount:   0,
		BadCount:    0,
		UserFace:    userface,
		ChannelType: 1, // 默认为文章频道
	}

	// 保存评论
	_, err = c.commentModel.Create(comment)
	if err != nil {
		logger.Error("保存回复失败", "error", err)
		http.Error(w, "Failed to save reply", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "回复提交成功，" + (map[bool]string{true: "已发布", false: "等待审核"})[c.config.Site.CommentAutoCheck],
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/article/"+aidStr+".html", http.StatusFound)
	}
}
