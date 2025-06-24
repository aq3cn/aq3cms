package admin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/plugin"
)

// PluginController 插件控制器
type PluginController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	pluginManager   *plugin.Manager
	templateService *service.TemplateService
}

// NewPluginController 创建插件控制器
func NewPluginController(db *database.DB, cache cache.Cache, config *config.Config, pluginManager *plugin.Manager) *PluginController {
	return &PluginController{
		db:              db,
		cache:           cache,
		config:          config,
		pluginManager:   pluginManager,
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// List 插件列表
func (c *PluginController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取插件列表
	plugins := c.pluginManager.GetPlugins()
	pluginInfos := c.pluginManager.GetPluginInfos()

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Plugins":     plugins,
		"PluginInfos": pluginInfos,
		"CurrentMenu": "plugin",
		"PageTitle":   "插件管理",
	}

	// 渲染模板
	tplFile := "admin/plugin_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染插件列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Enable 启用插件
func (c *PluginController) Enable(w http.ResponseWriter, r *http.Request) {
	// 获取插件名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Missing plugin name", http.StatusBadRequest)
		return
	}

	// 启用插件
	err := c.pluginManager.EnablePlugin(name)
	if err != nil {
		logger.Error("启用插件失败", "name", name, "error", err)
		http.Error(w, "Failed to enable plugin", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "插件启用成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/plugin/list", http.StatusFound)
	}
}

// Disable 禁用插件
func (c *PluginController) Disable(w http.ResponseWriter, r *http.Request) {
	// 获取插件名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Missing plugin name", http.StatusBadRequest)
		return
	}

	// 禁用插件
	err := c.pluginManager.DisablePlugin(name)
	if err != nil {
		logger.Error("禁用插件失败", "name", name, "error", err)
		http.Error(w, "Failed to disable plugin", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "插件禁用成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/plugin/list", http.StatusFound)
	}
}

// Config 插件配置页面
func (c *PluginController) Config(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取插件名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Missing plugin name", http.StatusBadRequest)
		return
	}

	// 获取插件
	plugin, err := c.pluginManager.GetPlugin(name)
	if err != nil {
		logger.Error("获取插件失败", "name", name, "error", err)
		http.Error(w, "Plugin not found", http.StatusNotFound)
		return
	}

	// 获取插件信息
	info, err := c.pluginManager.GetPluginInfo(name)
	if err != nil {
		logger.Error("获取插件信息失败", "name", name, "error", err)
		http.Error(w, "Plugin info not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Plugin":      plugin,
		"PluginInfo":  info,
		"CurrentMenu": "plugin",
		"PageTitle":   "插件配置 - " + plugin.Name(),
	}

	// 渲染模板
	tplFile := "admin/plugin_config.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染插件配置模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoConfig 处理插件配置
func (c *PluginController) DoConfig(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取插件名称
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Missing plugin name", http.StatusBadRequest)
		return
	}

	// 获取插件配置
	configJSON := r.FormValue("config")
	if configJSON == "" {
		http.Error(w, "Missing plugin config", http.StatusBadRequest)
		return
	}

	// 验证JSON格式
	var config json.RawMessage
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		logger.Error("解析插件配置失败", "error", err)
		http.Error(w, "Invalid config format", http.StatusBadRequest)
		return
	}

	// 更新插件配置
	err := c.pluginManager.UpdatePluginConfig(name, config)
	if err != nil {
		logger.Error("更新插件配置失败", "name", name, "error", err)
		http.Error(w, "Failed to update plugin config", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "插件配置保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/plugin/list", http.StatusFound)
	}
}

// Upload 上传插件页面
func (c *PluginController) Upload(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "plugin",
		"PageTitle":   "上传插件",
	}

	// 渲染模板
	tplFile := "admin/plugin_upload.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染上传插件模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoUpload 处理上传插件
func (c *PluginController) DoUpload(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	file, handler, err := r.FormFile("plugin")
	if err != nil {
		logger.Error("获取上传文件失败", "error", err)
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 检查文件扩展名
	if filepath.Ext(handler.Filename) != ".so" {
		http.Error(w, "Invalid plugin file", http.StatusBadRequest)
		return
	}

	// 创建插件目录
	if _, err := os.Stat(c.config.Plugin.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.config.Plugin.Dir, 0755); err != nil {
			logger.Error("创建插件目录失败", "error", err)
			http.Error(w, "Failed to create plugin directory", http.StatusInternalServerError)
			return
		}
	}

	// 创建目标文件
	dst, err := os.Create(filepath.Join(c.config.Plugin.Dir, handler.Filename))
	if err != nil {
		logger.Error("创建目标文件失败", "error", err)
		http.Error(w, "Failed to create destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error("复制文件内容失败", "error", err)
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}

	// 加载插件
	pluginPath := filepath.Join(c.config.Plugin.Dir, handler.Filename)
	if err := c.pluginManager.LoadPlugins(); err != nil {
		logger.Error("加载插件失败", "path", pluginPath, "error", err)
		http.Error(w, "Failed to load plugin", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "插件上传成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/plugin/list", http.StatusFound)
	}
}

// Delete 删除插件
func (c *PluginController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取插件名称
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Missing plugin name", http.StatusBadRequest)
		return
	}

	// 获取插件信息
	info, err := c.pluginManager.GetPluginInfo(name)
	if err != nil {
		logger.Error("获取插件信息失败", "name", name, "error", err)
		http.Error(w, "Plugin info not found", http.StatusNotFound)
		return
	}

	// 禁用插件
	if info.Enabled {
		if err := c.pluginManager.DisablePlugin(name); err != nil {
			logger.Error("禁用插件失败", "name", name, "error", err)
			http.Error(w, "Failed to disable plugin", http.StatusInternalServerError)
			return
		}
	}

	// 获取插件文件
	files, err := ioutil.ReadDir(c.config.Plugin.Dir)
	if err != nil {
		logger.Error("读取插件目录失败", "error", err)
		http.Error(w, "Failed to read plugin directory", http.StatusInternalServerError)
		return
	}

	// 查找插件文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 检查文件扩展名
		if filepath.Ext(file.Name()) != ".so" {
			continue
		}

		// 加载插件
		pluginPath := filepath.Join(c.config.Plugin.Dir, file.Name())
		p, err := plugin.Open(pluginPath)
		if err != nil {
			continue
		}

		// 获取插件实例
		symPlugin, err := p.Lookup("Plugin")
		if err != nil {
			continue
		}

		// 转换为插件接口
		plugin, ok := symPlugin.(plugin.Plugin)
		if !ok {
			continue
		}

		// 检查插件名称
		if plugin.Name() == name {
			// 删除插件文件
			if err := os.Remove(pluginPath); err != nil {
				logger.Error("删除插件文件失败", "path", pluginPath, "error", err)
				http.Error(w, "Failed to delete plugin file", http.StatusInternalServerError)
				return
			}
			break
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "插件删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/plugin/list", http.StatusFound)
	}
}
