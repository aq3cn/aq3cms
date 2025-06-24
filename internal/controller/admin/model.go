package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// ModelController 模型控制器
type ModelController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	modelModel      *model.ContentModelModel
	templateService *service.TemplateService
}

// NewModelController 创建模型控制器
func NewModelController(db *database.DB, cache cache.Cache, config *config.Config) *ModelController {
	return &ModelController{
		db:              db,
		cache:           cache,
		config:          config,
		modelModel:      model.NewContentModelModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 模型列表
func (c *ModelController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模型列表
	models, err := c.modelModel.GetAll()
	if err != nil {
		logger.Error("获取模型列表失败", "error", err)
		http.Error(w, "Failed to get models", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Models":      models,
		"CurrentMenu": "model",
		"PageTitle":   "内容模型管理",
	}

	// 渲染模板
	tplFile := "admin/model_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模型列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加模型页面
func (c *ModelController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "model",
		"PageTitle":   "添加内容模型",
	}

	// 渲染模板
	tplFile := "admin/model_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加模型模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加模型
func (c *ModelController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	tableName := r.FormValue("tablename")
	description := r.FormValue("description")
	stateStr := r.FormValue("state")
	fieldsJSON := r.FormValue("fields")

	// 验证必填字段
	if name == "" || tableName == "" || fieldsJSON == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	var fields []model.Field
	err := json.Unmarshal([]byte(fieldsJSON), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 验证字段
	for _, field := range fields {
		if field.Name == "" || field.Title == "" || field.Type == "" {
			http.Error(w, "Invalid field definition", http.StatusBadRequest)
			return
		}
	}

	// 解析状态
	state := 1
	if stateStr != "" {
		state, _ = strconv.Atoi(stateStr)
	}

	// 创建模型
	contentModel := &model.ContentModel{
		Name:        name,
		TableName:   tableName,
		Description: description,
		State:       state,
		Fields:      fieldsJSON,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	// 保存模型
	id, err := c.modelModel.Create(contentModel)
	if err != nil {
		logger.Error("创建模型失败", "error", err)
		http.Error(w, "Failed to create model", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "模型创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/list", http.StatusFound)
	}
}

// Edit 编辑模型页面
func (c *ModelController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模型ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(id)
	if err != nil {
		logger.Error("获取模型失败", "id", id, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Model":       contentModel,
		"CurrentMenu": "model",
		"PageTitle":   "编辑内容模型",
	}

	// 渲染模板
	tplFile := "admin/model_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑模型模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑模型
func (c *ModelController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取模型ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取原模型
	contentModel, err := c.modelModel.GetByID(id)
	if err != nil {
		logger.Error("获取模型失败", "id", id, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	description := r.FormValue("description")
	stateStr := r.FormValue("state")
	fieldsJSON := r.FormValue("fields")

	// 验证必填字段
	if name == "" || fieldsJSON == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(fieldsJSON), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 验证字段
	for _, field := range fields {
		if field.Name == "" || field.Title == "" || field.Type == "" {
			http.Error(w, "Invalid field definition", http.StatusBadRequest)
			return
		}
	}

	// 解析状态
	state := 1
	if stateStr != "" {
		state, _ = strconv.Atoi(stateStr)
	}

	// 更新模型
	contentModel.Name = name
	contentModel.Description = description
	contentModel.State = state
	contentModel.Fields = fieldsJSON
	contentModel.UpdateTime = time.Now()

	// 保存模型
	err = c.modelModel.Update(contentModel)
	if err != nil {
		logger.Error("更新模型失败", "error", err)
		http.Error(w, "Failed to update model", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "模型更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/list", http.StatusFound)
	}
}

// Delete 删除模型
func (c *ModelController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取模型ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 删除模型
	err = c.modelModel.Delete(id)
	if err != nil {
		logger.Error("删除模型失败", "id", id, "error", err)
		http.Error(w, "Failed to delete model", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "模型删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/list", http.StatusFound)
	}
}

// Fields 获取模型字段
func (c *ModelController) Fields(w http.ResponseWriter, r *http.Request) {
	// 获取模型ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(id)
	if err != nil {
		logger.Error("获取模型失败", "id", id, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 返回字段
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(contentModel.Fields))
}

// Content 内容列表
func (c *ModelController) Content(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模型ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(id)
	if err != nil {
		logger.Error("获取模型失败", "id", id, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 获取查询参数
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pagesize")

	// 解析参数
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

	// 获取内容列表
	qb := database.NewQueryBuilder(c.db, contentModel.TableName)
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)
	contents, err := qb.Get()
	if err != nil {
		logger.Error("获取内容列表失败", "error", err)
		http.Error(w, "Failed to get contents", http.StatusInternalServerError)
		return
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取内容总数失败", "error", err)
		http.Error(w, "Failed to get content count", http.StatusInternalServerError)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(contentModel.Fields), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
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
		"Model":       contentModel,
		"Fields":      fields,
		"Contents":    contents,
		"Pagination":  pagination,
		"CurrentMenu": "model",
		"PageTitle":   "内容管理 - " + contentModel.Name,
	}

	// 渲染模板
	tplFile := "admin/model_content.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染内容列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ContentAdd 添加内容页面
func (c *ModelController) ContentAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模型ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(id)
	if err != nil {
		logger.Error("获取模型失败", "id", id, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(contentModel.Fields), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Model":       contentModel,
		"Fields":      fields,
		"CurrentMenu": "model",
		"PageTitle":   "添加内容 - " + contentModel.Name,
	}

	// 渲染模板
	tplFile := "admin/model_content_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加内容模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ContentDoAdd 处理添加内容
func (c *ModelController) ContentDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取模型ID
	modelIDStr := r.FormValue("modelid")
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取文章ID
	aidStr := r.FormValue("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(modelID)
	if err != nil {
		logger.Error("获取模型失败", "id", modelID, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(contentModel.Fields), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 构建数据
	data := make(map[string]interface{})
	for _, field := range fields {
		value := r.FormValue(field.Name)
		data[field.Name] = value
	}

	// 保存内容
	err = c.modelModel.SaveContent(modelID, aid, data)
	if err != nil {
		logger.Error("保存内容失败", "error", err)
		http.Error(w, "Failed to save content", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "内容保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/content/"+modelIDStr, http.StatusFound)
	}
}

// ContentEdit 编辑内容页面
func (c *ModelController) ContentEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模型ID和内容ID
	vars := mux.Vars(r)
	modelIDStr := vars["id"]
	aidStr := vars["aid"]
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(modelID)
	if err != nil {
		logger.Error("获取模型失败", "id", modelID, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 获取内容
	content, err := c.modelModel.GetContent(modelID, aid)
	if err != nil {
		logger.Error("获取内容失败", "error", err)
		http.Error(w, "Content not found", http.StatusNotFound)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(contentModel.Fields), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Model":       contentModel,
		"Fields":      fields,
		"Content":     content,
		"CurrentMenu": "model",
		"PageTitle":   "编辑内容 - " + contentModel.Name,
	}

	// 渲染模板
	tplFile := "admin/model_content_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑内容模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ContentDoEdit 处理编辑内容
func (c *ModelController) ContentDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取模型ID
	modelIDStr := r.FormValue("modelid")
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}

	// 获取文章ID
	aidStr := r.FormValue("aid")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// 获取模型
	contentModel, err := c.modelModel.GetByID(modelID)
	if err != nil {
		logger.Error("获取模型失败", "id", modelID, "error", err)
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}

	// 解析字段
	var fields []model.Field
	err = json.Unmarshal([]byte(contentModel.Fields), &fields)
	if err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields format", http.StatusBadRequest)
		return
	}

	// 构建数据
	data := make(map[string]interface{})
	for _, field := range fields {
		value := r.FormValue(field.Name)
		data[field.Name] = value
	}

	// 保存内容
	err = c.modelModel.SaveContent(modelID, aid, data)
	if err != nil {
		logger.Error("保存内容失败", "error", err)
		http.Error(w, "Failed to save content", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "内容保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/content/"+modelIDStr, http.StatusFound)
	}
}

// ContentDelete 删除内容
func (c *ModelController) ContentDelete(w http.ResponseWriter, r *http.Request) {
	// 获取模型ID和内容ID
	vars := mux.Vars(r)
	modelIDStr := vars["id"]
	aidStr := vars["aid"]
	modelID, err := strconv.ParseInt(modelIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid model ID", http.StatusBadRequest)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	// 删除内容
	err = c.modelModel.DeleteContent(modelID, aid)
	if err != nil {
		logger.Error("删除内容失败", "error", err)
		http.Error(w, "Failed to delete content", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "内容删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/model/content/"+modelIDStr, http.StatusFound)
	}
}
