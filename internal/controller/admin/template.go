package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
	"github.com/gorilla/mux"
)

// TemplateController 模板控制器
type TemplateController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	templateService *service.TemplateService
}

// NewTemplateController 创建模板控制器
func NewTemplateController(db *database.DB, cache cache.Cache, config *config.Config) *TemplateController {
	return &TemplateController{
		db:              db,
		cache:           cache,
		config:          config,
		templateService: service.NewTemplateService(db, cache, config),
	}
}

// Index 模板管理首页
func (c *TemplateController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模板目录
	templateDir := c.config.Template.Dir
	if templateDir == "" {
		templateDir = "templets"
	}

	// 获取模板统计信息
	templateStats, err := c.getTemplateStats(templateDir)
	if err != nil {
		logger.Error("获取模板统计信息失败", "error", err)
		templateStats = map[string]int{
			"total_files":   0,
			"html_files":    0,
			"css_files":     0,
			"js_files":      0,
			"image_files":   0,
			"directories":   0,
			"last_modified": 0,
		}
	}

	// 获取最近修改的模板文件
	recentFiles, err := c.getRecentFiles(templateDir, 10)
	if err != nil {
		logger.Error("获取最近修改文件失败", "error", err)
		recentFiles = []map[string]interface{}{}
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":       adminID,
		"AdminName":     adminName,
		"TemplateDir":   templateDir,
		"TemplateStats": templateStats,
		"RecentFiles":   recentFiles,
		"CurrentMenu":   "template",
		"PageTitle":     "模板管理",
	}

	// 渲染模板
	tplFile := "admin/template_index.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模板管理首页模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// List 模板列表
func (c *TemplateController) List(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取模板目录
	templateDir := c.config.Template.Dir
	if templateDir == "" {
		templateDir = "templates"
	}

	// 获取当前目录
	currentDir := r.URL.Query().Get("dir")
	if currentDir == "" {
		currentDir = templateDir
	} else {
		// 安全检查，确保目录在模板目录下
		if !strings.HasPrefix(currentDir, templateDir) {
			currentDir = templateDir
		}
	}

	// 获取目录内容
	files, err := ioutil.ReadDir(currentDir)
	if err != nil {
		logger.Error("读取模板目录失败", "dir", currentDir, "error", err)
		http.Error(w, "Failed to read template directory", http.StatusInternalServerError)
		return
	}

	// 构建文件列表
	fileList := make([]map[string]interface{}, 0, len(files))
	for _, file := range files {
		fileInfo := map[string]interface{}{
			"Name":    file.Name(),
			"IsDir":   file.IsDir(),
			"Size":    file.Size(),
			"ModTime": file.ModTime(),
			"Path":    filepath.Join(currentDir, file.Name()),
		}
		fileList = append(fileList, fileInfo)
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(templateDir, currentDir)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Files":       fileList,
		"CurrentDir":  currentDir,
		"TemplateDir": templateDir,
		"Breadcrumbs": breadcrumbs,
		"CurrentMenu": "template",
		"PageTitle":   "模板管理",
	}

	// 渲染模板
	tplFile := "admin/template_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染模板列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Edit 编辑模板
func (c *TemplateController) Edit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取文件路径
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "Missing file path", http.StatusBadRequest)
		return
	}

	// 安全检查，确保文件在模板目录下
	templateDir := c.config.Template.Dir
	if !strings.HasPrefix(filePath, templateDir) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error("获取文件信息失败", "file", filePath, "error", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// 检查是否是目录
	if fileInfo.IsDir() {
		http.Error(w, "Cannot edit directory", http.StatusBadRequest)
		return
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("读取文件内容失败", "file", filePath, "error", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(templateDir, filepath.Dir(filePath))

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"FilePath":    filePath,
		"FileName":    filepath.Base(filePath),
		"FileContent": string(content),
		"FileSize":    fileInfo.Size(),
		"FileModTime": fileInfo.ModTime(),
		"Breadcrumbs": breadcrumbs,
		"CurrentMenu": "template",
		"PageTitle":   "编辑模板",
	}

	// 渲染模板
	tplFile := "admin/template_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑模板模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoEdit 处理编辑模板
func (c *TemplateController) DoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	filePath := r.FormValue("file_path")
	content := r.FormValue("content")

	// 安全检查，确保文件在模板目录下
	templateDir := c.config.Template.Dir
	if !strings.HasPrefix(filePath, templateDir) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		logger.Error("获取文件信息失败", "file", filePath, "error", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// 写入文件内容
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		logger.Error("写入文件内容失败", "file", filePath, "error", err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	// 清除模板缓存
	c.templateService.ClearCache()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "模板保存成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/template_list?dir="+filepath.Dir(filePath), http.StatusFound)
	}
}

// Create 创建模板
func (c *TemplateController) Create(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取当前目录
	currentDir := r.URL.Query().Get("dir")
	if currentDir == "" {
		currentDir = c.config.Template.Dir
	}

	// 安全检查，确保目录在模板目录下
	templateDir := c.config.Template.Dir
	if !strings.HasPrefix(currentDir, templateDir) {
		currentDir = templateDir
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(templateDir, currentDir)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentDir":  currentDir,
		"Breadcrumbs": breadcrumbs,
		"CurrentMenu": "template",
		"PageTitle":   "创建模板",
	}

	// 渲染模板
	tplFile := "admin/template_create.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染创建模板模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoCreate 处理创建模板
func (c *TemplateController) DoCreate(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	currentDir := r.FormValue("current_dir")
	fileName := r.FormValue("file_name")
	content := r.FormValue("content")
	fileType := r.FormValue("file_type")

	// 安全检查，确保目录在模板目录下
	templateDir := c.config.Template.Dir
	if !strings.HasPrefix(currentDir, templateDir) {
		http.Error(w, "Invalid directory", http.StatusBadRequest)
		return
	}

	// 安全检查，文件名
	fileName = security.SanitizeFilename(fileName)
	if fileName == "" {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	// 添加文件扩展名
	if fileType == "directory" {
		// 创建目录
		dirPath := filepath.Join(currentDir, fileName)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			logger.Error("创建目录失败", "dir", dirPath, "error", err)
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}
	} else {
		// 创建文件
		if !strings.Contains(fileName, ".") {
			fileName += ".htm"
		}
		filePath := filepath.Join(currentDir, fileName)

		// 检查文件是否已存在
		if _, err := os.Stat(filePath); err == nil {
			http.Error(w, "File already exists", http.StatusBadRequest)
			return
		}

		// 写入文件内容
		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			logger.Error("写入文件内容失败", "file", filePath, "error", err)
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		message := "文件创建成功"
		if fileType == "directory" {
			message = "目录创建成功"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": message,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/template_list?dir="+currentDir, http.StatusFound)
	}
}

// Delete 删除模板
func (c *TemplateController) Delete(w http.ResponseWriter, r *http.Request) {
	// 获取文件路径
	vars := mux.Vars(r)
	filePath := vars["path"]
	if filePath == "" {
		http.Error(w, "Missing file path", http.StatusBadRequest)
		return
	}

	// 安全检查，确保文件在模板目录下
	templateDir := c.config.Template.Dir
	if !strings.HasPrefix(filePath, templateDir) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logger.Error("获取文件信息失败", "file", filePath, "error", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// 删除文件或目录
	if fileInfo.IsDir() {
		// 删除目录
		err = os.RemoveAll(filePath)
	} else {
		// 删除文件
		err = os.Remove(filePath)
	}

	if err != nil {
		logger.Error("删除文件失败", "file", filePath, "error", err)
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	// 清除模板缓存
	c.templateService.ClearCache()

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		message := "文件删除成功"
		if fileInfo.IsDir() {
			message = "目录删除成功"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": message,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/aq3cms/template_list?dir="+filepath.Dir(filePath), http.StatusFound)
	}
}

// 构建面包屑导航
func buildBreadcrumbs(templateDir, currentDir string) []map[string]string {
	// 替换反斜杠为正斜杠
	templateDir = filepath.ToSlash(templateDir)
	currentDir = filepath.ToSlash(currentDir)

	// 分割路径
	parts := strings.Split(strings.TrimPrefix(currentDir, templateDir), "/")

	// 构建面包屑
	breadcrumbs := make([]map[string]string, 0, len(parts)+1)

	// 添加根目录
	breadcrumbs = append(breadcrumbs, map[string]string{
		"Name": "模板根目录",
		"Path": templateDir,
	})

	// 添加子目录
	path := templateDir
	for _, part := range parts {
		if part == "" {
			continue
		}
		path = filepath.Join(path, part)
		breadcrumbs = append(breadcrumbs, map[string]string{
			"Name": part,
			"Path": path,
		})
	}

	return breadcrumbs
}

// getTemplateStats 获取模板统计信息
func (c *TemplateController) getTemplateStats(templateDir string) (map[string]int, error) {
	stats := map[string]int{
		"total_files": 0,
		"html_files":  0,
		"css_files":   0,
		"js_files":    0,
		"image_files": 0,
		"directories": 0,
	}

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续遍历
		}

		if info.IsDir() {
			stats["directories"]++
		} else {
			stats["total_files"]++
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".htm", ".html":
				stats["html_files"]++
			case ".css":
				stats["css_files"]++
			case ".js":
				stats["js_files"]++
			case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg":
				stats["image_files"]++
			}
		}
		return nil
	})

	return stats, err
}

// getRecentFiles 获取最近修改的文件
func (c *TemplateController) getRecentFiles(templateDir string, limit int) ([]map[string]interface{}, error) {
	type fileInfo struct {
		Path    string
		Name    string
		Size    int64
		ModTime time.Time
		IsDir   bool
	}

	var files []fileInfo

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误，继续遍历
		}

		// 只处理文件，不处理目录
		if !info.IsDir() {
			files = append(files, fileInfo{
				Path:    path,
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按修改时间排序（最新的在前）
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i].ModTime.Before(files[j].ModTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// 限制返回数量
	if len(files) > limit {
		files = files[:limit]
	}

	// 转换为map格式
	result := make([]map[string]interface{}, len(files))
	for i, file := range files {
		result[i] = map[string]interface{}{
			"Path":    file.Path,
			"Name":    file.Name,
			"Size":    file.Size,
			"ModTime": file.ModTime,
			"IsDir":   file.IsDir,
		}
	}

	return result, nil
}
