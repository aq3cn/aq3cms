package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Template 模板
type Template struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 模板名称
	Path        string    `json:"path"`        // 模板路径
	Type        int       `json:"type"`        // 类型：0系统，1用户
	Description string    `json:"description"` // 描述
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// TemplateFile 模板文件
type TemplateFile struct {
	Name         string    `json:"name"`         // 文件名
	Path         string    `json:"path"`         // 文件路径
	RelativePath string    `json:"relativepath"` // 相对路径
	IsDir        bool      `json:"isdir"`        // 是否目录
	Size         int64     `json:"size"`         // 文件大小
	ModTime      time.Time `json:"modtime"`      // 修改时间
}

// TemplateModel 模板模型
type TemplateModel struct {
	db         *database.DB
	templateDir string
}

// NewTemplateModel 创建模板模型
func NewTemplateModel(db *database.DB, templateDir string) *TemplateModel {
	return &TemplateModel{
		db:         db,
		templateDir: templateDir,
	}
}

// GetByID 根据ID获取模板
func (m *TemplateModel) GetByID(id int64) (*Template, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "template")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取模板失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("template not found: %d", id)
	}

	// 转换为模板
	template := &Template{}
	template.ID, _ = result["id"].(int64)
	template.Name, _ = result["name"].(string)
	template.Path, _ = result["path"].(string)
	template.Type, _ = result["type"].(int)
	template.Description, _ = result["description"].(string)
	template.Status, _ = result["status"].(int)
	template.CreateTime, _ = result["createtime"].(time.Time)
	template.UpdateTime, _ = result["updatetime"].(time.Time)

	return template, nil
}

// GetByPath 根据路径获取模板
func (m *TemplateModel) GetByPath(path string) (*Template, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "template")
	qb.Where("path = ?", path)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取模板失败", "path", path, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("template not found: %s", path)
	}

	// 转换为模板
	template := &Template{}
	template.ID, _ = result["id"].(int64)
	template.Name, _ = result["name"].(string)
	template.Path, _ = result["path"].(string)
	template.Type, _ = result["type"].(int)
	template.Description, _ = result["description"].(string)
	template.Status, _ = result["status"].(int)
	template.CreateTime, _ = result["createtime"].(time.Time)
	template.UpdateTime, _ = result["updatetime"].(time.Time)

	return template, nil
}

// GetAll 获取所有模板
func (m *TemplateModel) GetAll(status int) ([]*Template, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "template")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有模板失败", "error", err)
		return nil, err
	}

	// 转换为模板列表
	templates := make([]*Template, 0, len(results))
	for _, result := range results {
		template := &Template{}
		template.ID, _ = result["id"].(int64)
		template.Name, _ = result["name"].(string)
		template.Path, _ = result["path"].(string)
		template.Type, _ = result["type"].(int)
		template.Description, _ = result["description"].(string)
		template.Status, _ = result["status"].(int)
		template.CreateTime, _ = result["createtime"].(time.Time)
		template.UpdateTime, _ = result["updatetime"].(time.Time)
		templates = append(templates, template)
	}

	return templates, nil
}

// Create 创建模板
func (m *TemplateModel) Create(template *Template) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	template.CreateTime = now
	template.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("template")+" (name, path, type, description, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?)",
		template.Name, template.Path, template.Type, template.Description, template.Status, template.CreateTime, template.UpdateTime,
	)
	if err != nil {
		logger.Error("创建模板失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新模板
func (m *TemplateModel) Update(template *Template) error {
	// 设置更新时间
	template.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("template")+" SET name = ?, path = ?, type = ?, description = ?, status = ?, updatetime = ? WHERE id = ?",
		template.Name, template.Path, template.Type, template.Description, template.Status, template.UpdateTime, template.ID,
	)
	if err != nil {
		logger.Error("更新模板失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除模板
func (m *TemplateModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("template")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除模板失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新模板状态
func (m *TemplateModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("template")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新模板状态失败", "error", err)
		return err
	}

	return nil
}

// GetTemplateFiles 获取模板文件
func (m *TemplateModel) GetTemplateFiles(path string) ([]*TemplateFile, error) {
	// 构建完整路径
	fullPath := filepath.Join(m.templateDir, path)

	// 检查路径是否存在
	_, err := os.Stat(fullPath)
	if err != nil {
		logger.Error("获取模板文件失败", "path", path, "error", err)
		return nil, err
	}

	// 读取目录
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		logger.Error("读取目录失败", "path", path, "error", err)
		return nil, err
	}

	// 转换为模板文件列表
	templateFiles := make([]*TemplateFile, 0, len(files))
	for _, file := range files {
		// 跳过隐藏文件
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// 构建相对路径
		relativePath := filepath.Join(path, file.Name())

		// 创建模板文件
		templateFile := &TemplateFile{
			Name:         file.Name(),
			Path:         filepath.Join(fullPath, file.Name()),
			RelativePath: relativePath,
			IsDir:        file.IsDir(),
			Size:         file.Size(),
			ModTime:      file.ModTime(),
		}
		templateFiles = append(templateFiles, templateFile)
	}

	return templateFiles, nil
}

// GetTemplateFileContent 获取模板文件内容
func (m *TemplateModel) GetTemplateFileContent(path string) (string, error) {
	// 构建完整路径
	fullPath := filepath.Join(m.templateDir, path)

	// 检查路径是否存在
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		logger.Error("获取模板文件内容失败", "path", path, "error", err)
		return "", err
	}

	// 检查是否是目录
	if fileInfo.IsDir() {
		return "", fmt.Errorf("path is a directory: %s", path)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		logger.Error("读取文件内容失败", "path", path, "error", err)
		return "", err
	}

	return string(content), nil
}

// SaveTemplateFileContent 保存模板文件内容
func (m *TemplateModel) SaveTemplateFileContent(path string, content string) error {
	// 构建完整路径
	fullPath := filepath.Join(m.templateDir, path)

	// 检查路径是否存在
	fileInfo, err := os.Stat(fullPath)
	if err == nil && fileInfo.IsDir() {
		return fmt.Errorf("path is a directory: %s", path)
	}

	// 创建目录
	dir := filepath.Dir(fullPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 写入文件内容
	err = ioutil.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		logger.Error("写入文件内容失败", "path", path, "error", err)
		return err
	}

	return nil
}

// CreateTemplateFile 创建模板文件
func (m *TemplateModel) CreateTemplateFile(path string, isDir bool) error {
	// 构建完整路径
	fullPath := filepath.Join(m.templateDir, path)

	// 检查路径是否存在
	_, err := os.Stat(fullPath)
	if err == nil {
		return fmt.Errorf("path already exists: %s", path)
	}

	// 创建目录或文件
	if isDir {
		// 创建目录
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			logger.Error("创建目录失败", "path", path, "error", err)
			return err
		}
	} else {
		// 创建目录
		dir := filepath.Dir(fullPath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error("创建目录失败", "dir", dir, "error", err)
			return err
		}

		// 创建文件
		file, err := os.Create(fullPath)
		if err != nil {
			logger.Error("创建文件失败", "path", path, "error", err)
			return err
		}
		defer file.Close()
	}

	return nil
}

// DeleteTemplateFile 删除模板文件
func (m *TemplateModel) DeleteTemplateFile(path string) error {
	// 构建完整路径
	fullPath := filepath.Join(m.templateDir, path)

	// 检查路径是否存在
	_, err := os.Stat(fullPath)
	if err != nil {
		logger.Error("删除模板文件失败", "path", path, "error", err)
		return err
	}

	// 删除文件或目录
	err = os.RemoveAll(fullPath)
	if err != nil {
		logger.Error("删除文件或目录失败", "path", path, "error", err)
		return err
	}

	return nil
}

// RenameTemplateFile 重命名模板文件
func (m *TemplateModel) RenameTemplateFile(oldPath, newPath string) error {
	// 构建完整路径
	oldFullPath := filepath.Join(m.templateDir, oldPath)
	newFullPath := filepath.Join(m.templateDir, newPath)

	// 检查旧路径是否存在
	_, err := os.Stat(oldFullPath)
	if err != nil {
		logger.Error("重命名模板文件失败", "oldPath", oldPath, "error", err)
		return err
	}

	// 检查新路径是否存在
	_, err = os.Stat(newFullPath)
	if err == nil {
		return fmt.Errorf("path already exists: %s", newPath)
	}

	// 创建目录
	dir := filepath.Dir(newFullPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 重命名文件或目录
	err = os.Rename(oldFullPath, newFullPath)
	if err != nil {
		logger.Error("重命名文件或目录失败", "oldPath", oldPath, "newPath", newPath, "error", err)
		return err
	}

	return nil
}
