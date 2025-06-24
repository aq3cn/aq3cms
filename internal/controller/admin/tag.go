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

// TagController 后台标签管理控制器
type TagController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	tagModel        *model.TagModel
	templateService *service.TemplateService
}

// NewTagController 创建标签管理控制器
func NewTagController(db *database.DB, cache cache.Cache, config *config.Config) *TagController {
	return &TagController{
		db:              db,
		cache:           cache,
		config:          config,
		tagModel:        model.NewTagModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 标签管理首页
func (c *TagController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":   adminID,
		"AdminName": adminName,
		"Title":     "标签管理",
	}

	// 渲染模板
	tplFile := "admin/tag_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签管理首页失败", "error", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// List 标签列表
func (c *TagController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))

	// 获取标签列表
	pageSize := 20
	var tags []*model.Tag
	var total int

	if keyword != "" {
		// 暂时使用GetList，后续添加Search方法
		tags, total, err = c.tagModel.GetList(page, pageSize)
	} else {
		tags, total, err = c.tagModel.GetList(page, pageSize)
	}

	if err != nil {
		logger.Error("获取标签列表失败", "error", err)
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":    adminID,
		"AdminName":  adminName,
		"Title":      "标签列表",
		"Tags":       tags,
		"Total":      total,
		"Page":       page,
		"PageSize":   pageSize,
		"TotalPages": totalPages,
		"Keyword":    keyword,
		"HasPrev":    page > 1,
		"HasNext":    page < int(totalPages),
		"PrevPage":   page - 1,
		"NextPage":   page + 1,
	}

	// 渲染模板
	tplFile := "admin/tag_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染标签列表失败", "error", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// Add 添加标签页面
func (c *TagController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":   adminID,
		"AdminName": adminName,
		"Title":     "添加标签",
	}

	// 渲染模板
	tplFile := "admin/tag_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加标签页面失败", "error", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DoAdd 执行添加标签
func (c *TagController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	err := r.ParseForm()
	if err != nil {
		logger.Error("解析表单失败", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	tagName := strings.TrimSpace(r.FormValue("tag"))
	isHotStr := r.FormValue("ishot")
	rankStr := r.FormValue("rank")

	// 验证数据
	if tagName == "" {
		c.showError(w, "标签名不能为空")
		return
	}

	// 检查标签是否已存在
	existingTag, err := c.tagModel.GetByName(tagName)
	if err == nil && existingTag != nil {
		c.showError(w, "标签已存在")
		return
	}

	// 转换数据类型
	isHot := 0
	if isHotStr == "1" {
		isHot = 1
	}

	rank := 0
	if rankStr != "" {
		rank, _ = strconv.Atoi(rankStr)
	}

	// 创建标签对象
	tag := &model.Tag{
		Tag:   tagName,
		Count: 0,
		Rank:  rank,
		IsHot: isHot,
	}

	// 保存标签
	id, err := c.tagModel.Create(tag)
	if err != nil {
		logger.Error("创建标签失败", "error", err)
		c.showError(w, "创建标签失败: "+err.Error())
		return
	}

	logger.Info("标签创建成功", "id", id, "tag", tagName)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "标签创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/tag_list", http.StatusFound)
	}
}

// Edit 编辑标签页面
func (c *TagController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取标签ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("无效的标签ID", "id", idStr, "error", err)
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// 获取标签信息
	tag, err := c.tagModel.GetByID(id)
	if err != nil {
		logger.Error("获取标签失败", "id", id, "error", err)
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":   adminID,
		"AdminName": adminName,
		"Title":     "编辑标签",
		"Tag":       tag,
	}

	// 渲染模板
	tplFile := "admin/tag_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑标签页面失败", "error", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DoEdit 执行编辑标签
func (c *TagController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	err := r.ParseForm()
	if err != nil {
		logger.Error("解析表单失败", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// 获取标签ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("无效的标签ID", "id", idStr, "error", err)
		c.showError(w, "无效的标签ID")
		return
	}

	// 获取表单数据
	tagName := strings.TrimSpace(r.FormValue("tag"))
	isHotStr := r.FormValue("ishot")
	rankStr := r.FormValue("rank")

	// 验证数据
	if tagName == "" {
		c.showError(w, "标签名不能为空")
		return
	}

	// 获取原标签信息
	tag, err := c.tagModel.GetByID(id)
	if err != nil {
		logger.Error("获取标签失败", "id", id, "error", err)
		c.showError(w, "标签不存在")
		return
	}

	// 检查标签名是否被其他标签使用
	if tag.Tag != tagName {
		existingTag, err := c.tagModel.GetByName(tagName)
		if err == nil && existingTag != nil && existingTag.ID != id {
			c.showError(w, "标签名已被其他标签使用")
			return
		}
	}

	// 转换数据类型
	isHot := 0
	if isHotStr == "1" {
		isHot = 1
	}

	rank := 0
	if rankStr != "" {
		rank, _ = strconv.Atoi(rankStr)
	}

	// 更新标签信息
	tag.Tag = tagName
	tag.IsHot = isHot
	tag.Rank = rank

	// 保存标签
	err = c.tagModel.Update(tag)
	if err != nil {
		logger.Error("更新标签失败", "error", err)
		c.showError(w, "更新标签失败: "+err.Error())
		return
	}

	logger.Info("标签更新成功", "id", id, "tag", tagName)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "标签更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/tag_list", http.StatusFound)
	}
}

// Delete 删除标签
func (c *TagController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取标签ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("无效的标签ID", "id", idStr, "error", err)
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// 删除标签
	err = c.tagModel.Delete(id)
	if err != nil {
		logger.Error("删除标签失败", "id", id, "error", err)
		c.showError(w, "删除标签失败: "+err.Error())
		return
	}

	logger.Info("标签删除成功", "id", id)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "标签删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/tag_list", http.StatusFound)
	}
}

// showError 显示错误信息
func (c *TagController) showError(w http.ResponseWriter, message string) {
	if r := w.Header().Get("X-Requested-With"); r == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": message,
		})
	} else {
		// 普通请求
		http.Error(w, message, http.StatusBadRequest)
	}
}
