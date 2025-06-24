package front

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FormController 前台自定义表单控制器
type FormController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	formService     *service.FormService
	templateService *service.TemplateService
}

// NewFormController 创建前台自定义表单控制器
func NewFormController(db *database.DB, cache cache.Cache, config *config.Config) *FormController {
	return &FormController{
		db:              db,
		cache:           cache,
		config:          config,
		formService:     service.NewFormService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Show 显示自定义表单
func (c *FormController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单代码
	vars := mux.Vars(r)
	code := vars["code"]

	// 获取自定义表单
	form, err := c.formService.GetForm(code)
	if err != nil {
		logger.Error("获取自定义表单失败", "code", code, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 检查表单状态
	if form.Status != 1 {
		http.Error(w, "Form is disabled", http.StatusForbidden)
		return
	}

	// 渲染自定义表单
	html, err := c.formService.RenderForm(code)
	if err != nil {
		logger.Error("渲染自定义表单失败", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Form":      form,
		"FormHTML":  html,
		"PageTitle": form.Title,
	}

	// 渲染模板
	tplFile := "form.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染自定义表单模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Submit 提交自定义表单
func (c *FormController) Submit(w http.ResponseWriter, r *http.Request) {
	// 获取自定义表单代码
	vars := mux.Vars(r)
	code := vars["code"]

	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	formData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) == 1 {
			formData[key] = values[0]
		} else {
			formData[key] = values
		}
	}

	// 提交自定义表单
	err := c.formService.SubmitForm(code, formData, r.RemoteAddr)
	if err != nil {
		logger.Error("提交自定义表单失败", "error", err)
		http.Error(w, "Failed to submit form", http.StatusInternalServerError)
		return
	}

	// 获取自定义表单
	form, err := c.formService.GetForm(code)
	if err != nil {
		logger.Error("获取自定义表单失败", "code", code, "error", err)
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "表单提交成功",
			"url":     form.SuccessURL,
		})
	} else {
		// 普通表单提交
		if form.SuccessURL != "" {
			http.Redirect(w, r, form.SuccessURL, http.StatusFound)
		} else {
			http.Redirect(w, r, "/form/success", http.StatusFound)
		}
	}
}

// Success 表单提交成功页面
func (c *FormController) Success(w http.ResponseWriter, r *http.Request) {
	// 准备模板数据
	data := map[string]interface{}{
		"PageTitle": "表单提交成功",
	}

	// 渲染模板
	tplFile := "form_success.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染表单提交成功模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
