package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FormController 自定义表单控制器
type FormController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	formService     *service.FormService
	formModel       *model.FormModel
	formDataModel   *model.FormDataModel
	templateService *service.TemplateService
}

// NewFormController 创建自定义表单控制器
func NewFormController(db *database.DB, cache cache.Cache, config *config.Config) *FormController {
	return &FormController{
		db:              db,
		cache:           cache,
		config:          config,
		formService:     service.NewFormService(db, cache, config),
		formModel:       model.NewFormModel(db),
		formDataModel:   model.NewFormDataModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 自定义表单列表
func (c *FormController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取自定义表单列表
	page := 1
	pageSize := 100
	forms, _, err := c.formModel.GetList(page, pageSize)
	if err != nil {
		logger.Error("获取自定义表单列表失败", "error", err)
		http.Error(w, "Failed to get forms", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Forms":       forms,
		"CurrentMenu": "form",
		"PageTitle":   "自定义表单管理",
	}

	// 渲染模板
	tplFile := "admin/form_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染自定义表单列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加自定义表单页面
func (c *FormController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "form",
		"PageTitle":   "添加自定义表单",
	}

	// 渲染模板
	tplFile := "admin/form_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加自定义表单模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加自定义表单
func (c *FormController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code") // 用于兼容旧版本
	description := r.FormValue("description")
	fields := r.FormValue("fields")
	template := r.FormValue("template")
	_ = r.FormValue("successurl") // 未使用
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	status, _ := strconv.Atoi(statusStr)

	// 创建自定义表单
	form := &model.Form{
		Title:       name,
		Description: description,
		Template:    template,
		Status:      status,
	}

	// 解析字段
	var formFields []*model.FormField
	if err := json.Unmarshal([]byte(fields), &formFields); err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields data", http.StatusBadRequest)
		return
	}
	form.Fields = formFields

	// 保存自定义表单
	id, err := c.formModel.Create(form)
	if err != nil {
		logger.Error("创建自定义表单失败", "error", err)
		http.Error(w, "Failed to create form", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("form:" + strconv.FormatInt(id, 10))

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/form_list", http.StatusFound)
	}
}

// Edit 编辑自定义表单页面
func (c *FormController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取自定义表单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 获取字段
	fields := form.Fields

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Form":        form,
		"Fields":      fields,
		"CurrentMenu": "form",
		"PageTitle":   "编辑自定义表单",
	}

	// 渲染模板
	tplFile := "admin/form_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑自定义表单模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑自定义表单
func (c *FormController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取自定义表单ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取原自定义表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code") // 用于兼容旧版本
	description := r.FormValue("description")
	fields := r.FormValue("fields")
	template := r.FormValue("template")
	_ = r.FormValue("successurl") // 未使用
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	status, _ := strconv.Atoi(statusStr)

	// 更新自定义表单
	form.Title = name
	form.Description = description
	form.Template = template
	form.Status = status

	// 解析字段
	var formFields []*model.FormField
	if err := json.Unmarshal([]byte(fields), &formFields); err != nil {
		logger.Error("解析字段失败", "error", err)
		http.Error(w, "Invalid fields data", http.StatusBadRequest)
		return
	}
	form.Fields = formFields

	// 保存自定义表单
	err = c.formModel.Update(form)
	if err != nil {
		logger.Error("更新自定义表单失败", "error", err)
		http.Error(w, "Failed to update form", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("form:" + strconv.FormatInt(form.ID, 10))

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/form_list", http.StatusFound)
	}
}

// Delete 删除自定义表单
func (c *FormController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 删除自定义表单
	err = c.formModel.Delete(id)
	if err != nil {
		logger.Error("删除自定义表单失败", "id", id, "error", err)
		http.Error(w, "Failed to delete form", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("form:" + strconv.FormatInt(form.ID, 10))

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/form_list", http.StatusFound)
	}
}

// Preview 预览自定义表单
func (c *FormController) Preview(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 渲染自定义表单
	html, err := c.formService.RenderForm(strconv.FormatInt(form.ID, 10))
	if err != nil {
		logger.Error("渲染自定义表单失败", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	// 返回HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// DataList 自定义表单数据列表
func (c *FormController) DataList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取自定义表单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 获取查询参数
	statusStr := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取自定义表单数据列表
	dataList, total, err := c.formDataModel.GetByFormID(id, status, page, 20)
	if err != nil {
		logger.Error("获取自定义表单数据列表失败", "error", err)
		http.Error(w, "Failed to get form data", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + 20 - 1) / 20
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
		"Form":        form,
		"DataList":    dataList,
		"Status":      status,
		"Pagination":  pagination,
		"CurrentMenu": "form",
		"PageTitle":   "自定义表单数据管理",
	}

	// 渲染模板
	tplFile := "admin/form_data_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染自定义表单数据列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DataDetail 自定义表单数据详情
func (c *FormController) DataDetail(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取自定义表单数据ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form data ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单数据
	formData, dataMap, err := c.formService.GetFormDataDetail(id)
	if err != nil {
		logger.Error("获取自定义表单数据详情失败", "id", id, "error", err)
		http.Error(w, "Form data not found", http.StatusNotFound)
		return
	}

	// 获取自定义表单
	form, err := c.formModel.GetByID(formData.FormID)
	if err != nil {
		logger.Error("获取自定义表单失败", "id", formData.FormID, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 获取字段
	fields := form.Fields

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Form":        form,
		"FormData":    formData,
		"DataMap":     dataMap,
		"Fields":      fields,
		"CurrentMenu": "form",
		"PageTitle":   "自定义表单数据详情",
	}

	// 渲染模板
	tplFile := "admin/form_data_detail.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染自定义表单数据详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DataProcess 处理自定义表单数据
func (c *FormController) DataProcess(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单数据ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form data ID", http.StatusBadRequest)
		return
	}

	// 处理自定义表单数据
	err = c.formService.ProcessFormData(id)
	if err != nil {
		logger.Error("处理自定义表单数据失败", "id", id, "error", err)
		http.Error(w, "Failed to process form data", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单数据处理成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/form_data_list/"+vars["formid"], http.StatusFound)
	}
}

// DataDelete 删除自定义表单数据
func (c *FormController) DataDelete(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单数据ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form data ID", http.StatusBadRequest)
		return
	}

	// 获取自定义表单数据
	formData, err := c.formDataModel.GetByID(id)
	if err != nil {
		logger.Error("获取自定义表单数据失败", "id", id, "error", err)
		http.Error(w, "Form data not found", http.StatusNotFound)
		return
	}

	// 删除自定义表单数据
	err = c.formDataModel.Delete(id)
	if err != nil {
		logger.Error("删除自定义表单数据失败", "id", id, "error", err)
		http.Error(w, "Failed to delete form data", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单数据删除成功",
		})
	} else {
		// 普通表单提交
		formID := strconv.FormatInt(formData.FormID, 10)
		http.Redirect(w, r, "/admin/form_data_list/"+formID, http.StatusFound)
	}
}

// InitForms 初始化自定义表单
func (c *FormController) InitForms(w http.ResponseWriter, r *http.Request) {
	// 初始化默认自定义表单
	err := c.formService.InitDefaultForms()
	if err != nil {
		logger.Error("初始化自定义表单失败", "error", err)
		http.Error(w, "Failed to initialize forms", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "自定义表单初始化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/form_list", http.StatusFound)
	}
}
