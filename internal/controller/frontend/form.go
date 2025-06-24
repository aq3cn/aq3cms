package frontend

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FormController 表单控制器
type FormController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	formModel       *model.FormModel
	templateService *service.TemplateService
}

// NewFormController 创建表单控制器
func NewFormController(db *database.DB, cache cache.Cache, config *config.Config) *FormController {
	return &FormController{
		db:              db,
		cache:           cache,
		config:          config,
		formModel:       model.NewFormModel(db),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Show 显示表单
func (c *FormController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 检查表单状态
	if form.Status != 1 {
		http.Error(w, "Form is closed", http.StatusForbidden)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Form":        form,
		"PageTitle":   form.Title + " - " + c.config.Site.Name,
		"Keywords":    form.Title + "," + c.config.Site.Keywords,
		"Description": form.Description,
	}

	// 确定模板文件
	var tplFile string
	if form.Template != "" {
		tplFile = form.Template
	} else {
		tplFile = c.config.Template.DefaultTpl + "/form.htm"
	}

	// 渲染模板
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染表单模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Submit 提交表单
func (c *FormController) Submit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单ID
	idStr := r.FormValue("formid")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 检查表单状态
	if form.Status != 1 {
		http.Error(w, "Form is closed", http.StatusForbidden)
		return
	}

	// 收集表单数据
	formData := make(map[string]string)
	for _, field := range form.Fields {
		value := r.FormValue(field.FieldName)
		
		// 检查必填字段
		if field.IsRequired == 1 && value == "" {
			http.Error(w, field.FieldTitle+" is required", http.StatusBadRequest)
			return
		}
		
		formData[field.FieldName] = value
	}

	// 提交表单
	_, err = c.formModel.SubmitForm(id, formData, r.RemoteAddr)
	if err != nil {
		logger.Error("提交表单失败", "id", id, "error", err)
		http.Error(w, "Failed to submit form", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "表单提交成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/form/success/"+idStr, http.StatusFound)
	}
}

// Success 表单提交成功
func (c *FormController) Success(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid form ID", http.StatusBadRequest)
		return
	}

	// 获取表单
	form, err := c.formModel.GetByID(id)
	if err != nil {
		logger.Error("获取表单失败", "id", id, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 获取全局变量
	globals := c.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Form":        form,
		"PageTitle":   "提交成功 - " + form.Title + " - " + c.config.Site.Name,
		"Keywords":    form.Title + "," + c.config.Site.Keywords,
		"Description": form.Description,
	}

	// 渲染模板
	tplFile := c.config.Template.DefaultTpl + "/form_success.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染表单成功模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
