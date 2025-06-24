package tags

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"aq3cms/config"
	"aq3cms/pkg/logger"
)

// IncludeTag 包含标签处理器
type IncludeTag struct {
	Config *config.TemplateConfig
}

// Handle 处理标签
func (t *IncludeTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 获取文件名
	filename := ""
	for k, v := range attrs {
		if k == "file" {
			filename = v
			break
		}
	}

	if filename == "" {
		return "", fmt.Errorf("包含标签缺少file属性")
	}

	// 尝试多个路径
	paths := []string{
		filepath.Join(t.Config.Dir, filename),
		filepath.Join("templets", filename),
		filepath.Join("templates", filename),
		filename,
	}

	var contentBytes []byte
	var err error
	var foundPath string

	// 尝试读取文件
	for _, p := range paths {
		contentBytes, err = ioutil.ReadFile(p)
		if err == nil {
			foundPath = p
			logger.Info("成功读取包含文件", "file", foundPath)
			break
		}
	}

	// 如果所有路径都失败，返回错误
	if err != nil {
		logger.Error("读取包含文件失败", "file", filename, "尝试路径", paths, "error", err)
		return "", fmt.Errorf("无法找到包含文件: %s", filename)
	}

	return string(contentBytes), nil
}
