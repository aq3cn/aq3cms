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
)

// MessageController 消息控制器
type MessageController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	messageModel    *model.MessageModel
	memberModel     *model.MemberModel
	templateService *service.TemplateService
}

// NewMessageController 创建消息控制器
func NewMessageController(db *database.DB, cache cache.Cache, config *config.Config) *MessageController {
	return &MessageController{
		db:              db,
		cache:           cache,
		config:          config,
		messageModel:    model.NewMessageModel(db),
		memberModel:     model.NewMemberModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Inbox 收件箱
func (c *MessageController) Inbox(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取收件箱
	pageSize := 10
	messages, total, err := c.messageModel.GetInbox(memberID, page, pageSize)
	if err != nil {
		logger.Error("获取收件箱失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get inbox", http.StatusInternalServerError)
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

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get member info", http.StatusInternalServerError)
		return
	}

	// 获取未读消息数量
	unreadCount, err := c.messageModel.GetUnreadCount(memberID)
	if err != nil {
		logger.Error("获取未读消息数量失败", "memberid", memberID, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"Messages":    messages,
		"Pagination":  pagination,
		"UnreadCount": unreadCount,
		"PageTitle":   "收件箱 - " + c.config.Site.Name,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/inbox.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染收件箱模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Outbox 发件箱
func (c *MessageController) Outbox(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取页码
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取发件箱
	pageSize := 10
	messages, total, err := c.messageModel.GetOutbox(memberID, page, pageSize)
	if err != nil {
		logger.Error("获取发件箱失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get outbox", http.StatusInternalServerError)
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

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get member info", http.StatusInternalServerError)
		return
	}

	// 获取未读消息数量
	unreadCount, err := c.messageModel.GetUnreadCount(memberID)
	if err != nil {
		logger.Error("获取未读消息数量失败", "memberid", memberID, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"Messages":    messages,
		"Pagination":  pagination,
		"UnreadCount": unreadCount,
		"PageTitle":   "发件箱 - " + c.config.Site.Name,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/outbox.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染发件箱模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Read 阅读消息
func (c *MessageController) Read(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取消息ID
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// 获取消息
	message, err := c.messageModel.GetByID(id)
	if err != nil {
		logger.Error("获取消息失败", "id", id, "error", err)
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// 检查权限
	if message.ToID != memberID && message.FromID != memberID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	// 如果是收件人，标记为已读
	if message.ToID == memberID && message.IsRead == 0 {
		err = c.messageModel.MarkAsRead(id, memberID)
		if err != nil {
			logger.Error("标记消息为已读失败", "id", id, "memberid", memberID, "error", err)
		}
	}

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get member info", http.StatusInternalServerError)
		return
	}

	// 获取未读消息数量
	unreadCount, err := c.messageModel.GetUnreadCount(memberID)
	if err != nil {
		logger.Error("获取未读消息数量失败", "memberid", memberID, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"Message":     message,
		"UnreadCount": unreadCount,
		"PageTitle":   "阅读消息 - " + c.config.Site.Name,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/message_read.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染阅读消息模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Send 发送消息页面
func (c *MessageController) Send(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 获取收件人ID
	toIDStr := r.URL.Query().Get("toid")
	var toID int64
	var toMember *model.Member
	if toIDStr != "" {
		toID, _ = strconv.ParseInt(toIDStr, 10, 64)
		if toID > 0 {
			// 获取收件人信息
			var err error
			toMember, err = c.memberModel.GetByID(toID)
			if err != nil {
				logger.Error("获取收件人信息失败", "toid", toID, "error", err)
			}
		}
	}

	// 获取会员信息
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员信息失败", "memberid", memberID, "error", err)
		http.Error(w, "Failed to get member info", http.StatusInternalServerError)
		return
	}

	// 获取未读消息数量
	unreadCount, err := c.messageModel.GetUnreadCount(memberID)
	if err != nil {
		logger.Error("获取未读消息数量失败", "memberid", memberID, "error", err)
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Member":      member,
		"ToMember":    toMember,
		"UnreadCount": unreadCount,
		"PageTitle":   "发送消息 - " + c.config.Site.Name,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/member/message_send.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染发送消息模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoSend 处理发送消息
func (c *MessageController) DoSend(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	toIDStr := r.FormValue("toid")
	toID, err := strconv.ParseInt(toIDStr, 10, 64)
	if err != nil || toID <= 0 {
		http.Error(w, "Invalid recipient ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Message title cannot be empty", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
		return
	}

	// 检查收件人是否存在
	toMember, err := c.memberModel.GetByID(toID)
	if err != nil || toMember == nil {
		logger.Error("收件人不存在", "toid", toID, "error", err)
		http.Error(w, "Recipient not found", http.StatusBadRequest)
		return
	}

	// 创建消息
	message := &model.Message{
		FromID:   memberID,
		ToID:     toID,
		Title:    title,
		Content:  content,
		SendTime: time.Now(),
		IsRead:   0,
		FromDel:  0,
		ToDel:    0,
	}

	// 保存消息
	_, err = c.messageModel.Create(message)
	if err != nil {
		logger.Error("发送消息失败", "error", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "消息发送成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/member/outbox", http.StatusFound)
	}
}

// Delete 删除消息
func (c *MessageController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取会员ID
	memberID := middleware.GetMemberID(r)
	if memberID == 0 {
		http.Redirect(w, r, "/member/login", http.StatusFound)
		return
	}

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取消息ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// 获取消息
	message, err := c.messageModel.GetByID(id)
	if err != nil {
		logger.Error("获取消息失败", "id", id, "error", err)
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// 获取操作类型
	action := r.FormValue("action")
	if action != "inbox" && action != "outbox" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// 检查权限
	if action == "inbox" && message.ToID != memberID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	if action == "outbox" && message.FromID != memberID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	// 删除消息
	var err2 error
	if action == "inbox" {
		err2 = c.messageModel.DeleteFromInbox(id, memberID)
	} else {
		err2 = c.messageModel.DeleteFromOutbox(id, memberID)
	}

	if err2 != nil {
		logger.Error("删除消息失败", "id", id, "action", action, "error", err2)
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	// 如果发件人和收件人都已删除，完全删除消息
	if (action == "inbox" && message.FromDel == 1) || (action == "outbox" && message.ToDel == 1) {
		err = c.messageModel.DeleteCompletely(id)
		if err != nil {
			logger.Error("完全删除消息失败", "id", id, "error", err)
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "消息删除成功",
		})
	} else {
		// 普通表单提交
		if action == "inbox" {
			http.Redirect(w, r, "/member/inbox", http.StatusFound)
		} else {
			http.Redirect(w, r, "/member/outbox", http.StatusFound)
		}
	}
}
