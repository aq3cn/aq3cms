package admin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// I18nController 国际化控制器
type I18nController struct {
	db          *database.DB
	cache       cache.Cache
	config      *config.Config
	i18nService *service.I18nService
	templateService *service.TemplateService
}

// NewI18nController 创建国际化控制器
func NewI18nController(db *database.DB, cache cache.Cache, config *config.Config) *I18nController {
	return &I18nController{
		db:          db,
		cache:       cache,
		config:      config,
		i18nService: service.NewI18nService(db, cache, config),
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 语言列表
func (c *I18nController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取语言列表
	langs := c.i18nService.GetAvailableLangs()

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Langs":       langs,
		"DefaultLang": c.i18nService.GetDefaultLang(),
		"CurrentMenu": "i18n",
		"PageTitle":   "语言管理",
	}

	// 渲染模板
	tplFile := "admin/i18n_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染语言列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Add 添加语言页面
func (c *I18nController) Add(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "i18n",
		"PageTitle":   "添加语言",
	}

	// 渲染模板
	tplFile := "admin/i18n_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加语言模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoAdd 处理添加语言
func (c *I18nController) DoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	lang := r.FormValue("lang")
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 安全检查
	lang = security.SanitizeFilename(lang)
	if lang == "" {
		http.Error(w, "Invalid language code", http.StatusBadRequest)
		return
	}

	// 检查语言是否已存在
	langs := c.i18nService.GetLangs()
	for _, l := range langs {
		if l == lang {
			http.Error(w, "Language already exists", http.StatusBadRequest)
			return
		}
	}

	// 创建语言文件
	err := c.i18nService.CreateLangFile(lang)
	if err != nil {
		logger.Error("创建语言文件失败", "lang", lang, "error", err)
		http.Error(w, "Failed to create language file", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "语言添加成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/i18n/list", http.StatusFound)
	}
}

// Edit 编辑语言页面
func (c *I18nController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取语言代码
	vars := mux.Vars(r)
	lang := vars["lang"]
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 获取语言消息
	messages, err := c.i18nService.GetLangMessages(lang)
	if err != nil {
		logger.Error("获取语言消息失败", "lang", lang, "error", err)
		http.Error(w, "Failed to get language messages", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Lang":        lang,
		"Messages":    messages,
		"CurrentMenu": "i18n",
		"PageTitle":   "编辑语言",
	}

	// 渲染模板
	tplFile := "admin/i18n_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑语言模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑语言
func (c *I18nController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取语言代码
	lang := r.FormValue("lang")
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 获取消息
	messagesJSON := r.FormValue("messages")
	if messagesJSON == "" {
		http.Error(w, "Missing messages", http.StatusBadRequest)
		return
	}

	// 解析消息
	var messages map[string]string
	err := json.Unmarshal([]byte(messagesJSON), &messages)
	if err != nil {
		logger.Error("解析消息失败", "error", err)
		http.Error(w, "Invalid messages format", http.StatusBadRequest)
		return
	}

	// 更新语言消息
	err = c.i18nService.UpdateLangMessages(lang, messages)
	if err != nil {
		logger.Error("更新语言消息失败", "lang", lang, "error", err)
		http.Error(w, "Failed to update language messages", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "语言更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/i18n/list", http.StatusFound)
	}
}

// Delete 删除语言
func (c *I18nController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取语言代码
	vars := mux.Vars(r)
	lang := vars["lang"]
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 检查是否是默认语言
	if lang == c.i18nService.GetDefaultLang() {
		http.Error(w, "Cannot delete default language", http.StatusBadRequest)
		return
	}

	// 删除语言
	err := c.i18nService.DeleteLang(lang)
	if err != nil {
		logger.Error("删除语言失败", "lang", lang, "error", err)
		http.Error(w, "Failed to delete language", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "语言删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/i18n/list", http.StatusFound)
	}
}

// SetDefault 设置默认语言
func (c *I18nController) SetDefault(w http.ResponseWriter, r *http.Request) {
	// 获取语言代码
	vars := mux.Vars(r)
	lang := vars["lang"]
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 设置默认语言
	c.i18nService.SetDefaultLang(lang)

	// 更新配置
	c.config.Site.DefaultLang = lang
	err := c.config.Save()
	if err != nil {
		logger.Error("保存配置失败", "error", err)
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "默认语言设置成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/i18n/list", http.StatusFound)
	}
}

// Import 导入语言页面
func (c *I18nController) Import(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "i18n",
		"PageTitle":   "导入语言",
	}

	// 渲染模板
	tplFile := "admin/i18n_import.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染导入语言模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoImport 处理导入语言
func (c *I18nController) DoImport(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取语言代码
	lang := r.FormValue("lang")
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 安全检查
	lang = security.SanitizeFilename(lang)
	if lang == "" {
		http.Error(w, "Invalid language code", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Error("获取上传文件失败", "error", err)
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "lang_*.json")
	if err != nil {
		logger.Error("创建临时文件失败", "error", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 复制文件内容
	_, err = io.Copy(tempFile, file)
	if err != nil {
		logger.Error("复制文件内容失败", "error", err)
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}

	// 导入语言
	err = c.i18nService.ImportLang(lang, tempFile.Name())
	if err != nil {
		logger.Error("导入语言失败", "lang", lang, "error", err)
		http.Error(w, "Failed to import language", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "语言导入成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/i18n/list", http.StatusFound)
	}
}

// Export 导出语言
func (c *I18nController) Export(w http.ResponseWriter, r *http.Request) {
	// 获取语言代码
	vars := mux.Vars(r)
	lang := vars["lang"]
	if lang == "" {
		http.Error(w, "Missing language code", http.StatusBadRequest)
		return
	}

	// 获取语言目录
	langDir := filepath.Join(c.config.Template.Dir, "lang")

	// 读取语言文件
	filePath := filepath.Join(langDir, lang+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("读取语言文件失败", "lang", lang, "error", err)
		http.Error(w, "Failed to read language file", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+lang+".json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	// 输出文件内容
	w.Write(data)
}
