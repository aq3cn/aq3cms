package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// UploadController 上传API控制器
type UploadController struct {
	*BaseController
}

// NewUploadController 创建上传API控制器
func NewUploadController(db *database.DB, cache cache.Cache, config *config.Config) *UploadController {
	return &UploadController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// Image 上传图片
func (c *UploadController) Image(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		c.Error(w, 400, "Invalid form data")
		return
	}

	// 获取上传的文件
	file, handler, err := r.FormFile("image")
	if err != nil {
		logger.Error("获取上传文件失败", "error", err)
		c.Error(w, 400, "Failed to get uploaded file")
		return
	}
	defer file.Close()

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := strings.Split(c.config.Upload.AllowedExts, ",")
	allowed := false
	for _, allowedExt := range allowedExts {
		if strings.TrimSpace(allowedExt) == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		c.Error(w, 400, "Invalid file type")
		return
	}

	// 检查文件大小
	if handler.Size > int64(c.config.Upload.MaxSize*1024*1024) {
		c.Error(w, 400, "File too large")
		return
	}

	// 创建上传目录
	uploadDir := c.config.Upload.Dir
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	now := time.Now()
	yearMonth := now.Format("200601")
	uploadPath := filepath.Join(uploadDir, "images", yearMonth)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("创建上传目录失败", "error", err)
		c.Error(w, 500, "Failed to create upload directory")
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s%s", now.Unix(), security.RandomString(8), ext)
	filePath := filepath.Join(uploadPath, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Error("创建目标文件失败", "error", err)
		c.Error(w, 500, "Failed to create destination file")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error("复制文件内容失败", "error", err)
		c.Error(w, 500, "Failed to copy file content")
		return
	}

	// 构建URL
	url := "/" + strings.ReplaceAll(filePath, "\\", "/")

	// 返回数据
	c.Success(w, map[string]interface{}{
		"url":      url,
		"filename": filename,
		"size":     handler.Size,
		"message":  "Image uploaded successfully",
	})
}

// File 上传文件
func (c *UploadController) File(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 解析表单
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		c.Error(w, 400, "Invalid form data")
		return
	}

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error("获取上传文件失败", "error", err)
		c.Error(w, 400, "Failed to get uploaded file")
		return
	}
	defer file.Close()

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	denyExts := strings.Split(c.config.Upload.DenyExts, ",")
	for _, denyExt := range denyExts {
		if strings.TrimSpace(denyExt) == ext {
			c.Error(w, 400, "Invalid file type")
			return
		}
	}

	// 检查文件大小
	if handler.Size > int64(c.config.Upload.MaxSize*1024*1024) {
		c.Error(w, 400, "File too large")
		return
	}

	// 创建上传目录
	uploadDir := c.config.Upload.Dir
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	now := time.Now()
	yearMonth := now.Format("200601")
	uploadPath := filepath.Join(uploadDir, "files", yearMonth)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("创建上传目录失败", "error", err)
		c.Error(w, 500, "Failed to create upload directory")
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s%s", now.Unix(), security.RandomString(8), ext)
	filePath := filepath.Join(uploadPath, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Error("创建目标文件失败", "error", err)
		c.Error(w, 500, "Failed to create destination file")
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error("复制文件内容失败", "error", err)
		c.Error(w, 500, "Failed to copy file content")
		return
	}

	// 构建URL
	url := "/" + strings.ReplaceAll(filePath, "\\", "/")

	// 返回数据
	c.Success(w, map[string]interface{}{
		"url":      url,
		"filename": filename,
		"size":     handler.Size,
		"message":  "File uploaded successfully",
	})
}
