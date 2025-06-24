package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CommentController 评论控制器
type CommentController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	commentModel    *model.CommentModel
	articleModel    *model.ArticleModel
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
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 评论管理首页
func (c *CommentController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取评论统计信息
	totalComments, err := c.commentModel.GetTotalCount()
	if err != nil {
		logger.Error("获取评论总数失败", "error", err)
		totalComments = 0
	}

	pendingComments, err := c.commentModel.GetPendingCount()
	if err != nil {
		logger.Error("获取待审核评论数失败", "error", err)
		pendingComments = 0
	}

	approvedComments, err := c.commentModel.GetApprovedCount()
	if err != nil {
		logger.Error("获取已审核评论数失败", "error", err)
		approvedComments = 0
	}

	rejectedComments, err := c.commentModel.GetRejectedCount()
	if err != nil {
		logger.Error("获取已拒绝评论数失败", "error", err)
		rejectedComments = 0
	}

	// 获取最新评论
	latestComments, err := c.commentModel.GetLatest(10)
	if err != nil {
		logger.Error("获取最新评论失败", "error", err)
		latestComments = []*model.Comment{}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":          adminID,
		"AdminName":        adminName,
		"TotalComments":    totalComments,
		"PendingComments":  pendingComments,
		"ApprovedComments": approvedComments,
		"RejectedComments": rejectedComments,
		"LatestComments":   latestComments,
		"CurrentMenu":      "comment",
		"PageTitle":        "评论管理",
	}

	// 渲染模板
	tplFile := "admin/comment_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染评论管理首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// List 评论列表
func (c *CommentController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	isCheckStr := r.URL.Query().Get("ischeck")
	keyword := r.URL.Query().Get("keyword")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pagesize")

	// 解析参数
	isCheck := -1
	if isCheckStr != "" {
		isCheck, _ = strconv.Atoi(isCheckStr)
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

	// 获取评论列表
	comments, total, err := c.commentModel.GetList(isCheck, keyword, page, pageSize)
	if err != nil {
		logger.Error("获取评论列表失败", "error", err)
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

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Comments":    comments,
		"Pagination":  pagination,
		"IsCheck":     isCheck,
		"Keyword":     keyword,
		"CurrentMenu": "comment",
		"PageTitle":   "评论管理",
	}

	// 渲染模板
	tplFile := "admin/comment_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染评论列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Detail 评论详情
func (c *CommentController) Detail(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取评论ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// 获取评论
	comment, err := c.commentModel.GetByID(id)
	if err != nil {
		logger.Error("获取评论失败", "id", id, "error", err)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// 获取文章
	article, err := c.articleModel.GetByID(comment.AID)
	if err != nil {
		logger.Error("获取文章失败", "id", comment.AID, "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Comment":     comment,
		"Article":     article,
		"CurrentMenu": "comment",
		"PageTitle":   "评论详情",
	}

	// 渲染模板
	tplFile := "admin/comment_detail.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染评论详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Approve 审核评论
func (c *CommentController) Approve(w http.ResponseWriter, r *http.Request) {
	// 获取评论ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// 审核评论
	err = c.commentModel.UpdateStatus(id, 1)
	if err != nil {
		logger.Error("审核评论失败", "id", id, "error", err)
		http.Error(w, "Failed to approve comment", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论审核成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/comment/list", http.StatusFound)
	}
}

// Reject 拒绝评论
func (c *CommentController) Reject(w http.ResponseWriter, r *http.Request) {
	// 获取评论ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// 拒绝评论
	err = c.commentModel.UpdateStatus(id, -1)
	if err != nil {
		logger.Error("拒绝评论失败", "id", id, "error", err)
		http.Error(w, "Failed to reject comment", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论拒绝成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/comment/list", http.StatusFound)
	}
}

// Delete 删除评论
func (c *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取评论ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// 删除评论
	err = c.commentModel.Delete(id)
	if err != nil {
		logger.Error("删除评论失败", "id", id, "error", err)
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/comment/list", http.StatusFound)
	}
}

// BatchApprove 批量审核评论
func (c *CommentController) BatchApprove(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取评论ID列表
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No comments selected", http.StatusBadRequest)
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

	// 批量审核评论
	for _, id := range ids {
		err := c.commentModel.UpdateStatus(id, 1)
		if err != nil {
			logger.Error("审核评论失败", "id", id, "error", err)
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论批量审核成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/comment/list", http.StatusFound)
	}
}

// BatchDelete 批量删除评论
func (c *CommentController) BatchDelete(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取评论ID列表
	idsStr := r.FormValue("ids")
	if idsStr == "" {
		http.Error(w, "No comments selected", http.StatusBadRequest)
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

	// 批量删除评论
	for _, id := range ids {
		err := c.commentModel.Delete(id)
		if err != nil {
			logger.Error("删除评论失败", "id", id, "error", err)
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "评论批量删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/comment/list", http.StatusFound)
	}
}
