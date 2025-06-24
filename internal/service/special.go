package service

import (
	"fmt"
	"os"
	"path/filepath"

	"aq3cms/pkg/logger"
)

// GenerateSpecial 生成专题页
func (s *HtmlService) GenerateSpecial(id int64) error {
	logger.Info("开始生成专题页", "id", id)

	// 获取专题详情
	special, err := s.specialModel.GetByID(id)
	if err != nil {
		logger.Error("获取专题详情失败", "id", id, "error", err)
		return err
	}

	// 获取专题文章
	articles, _, err := s.specialModel.GetArticles(id, 1, 100)
	if err != nil {
		logger.Error("获取专题文章失败", "id", id, "error", err)
	}

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"Special":     special,
		"Articles":    articles,
		"PageTitle":   special.Title + " - " + s.config.Site.Name,
		"Keywords":    special.Keywords,
		"Description": special.Description,
	}

	// 确定模板文件
	tplFile := s.config.Template.DefaultTpl + "/special.htm"

	// 生成静态页面
	staticPath := fmt.Sprintf("special/%d.html", id)

	// 确保目录存在
	dir := filepath.Dir(staticPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
	if err != nil {
		logger.Error("生成静态专题页失败", "id", id, "error", err)
		return err
	}

	// 生成移动端专题页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/special.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		mobileStaticPath := fmt.Sprintf("m/special/%d.html", id)

		// 确保目录存在
		mobileDir := filepath.Dir(mobileStaticPath)
		if err := os.MkdirAll(mobileDir, 0755); err != nil {
			logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
		} else {
			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态专题页失败", "id", id, "error", err)
			}
		}
	}

	logger.Info("专题页生成完成", "id", id)
	return nil
}

// GenerateTag 生成标签页
func (s *HtmlService) GenerateTag(tagName string) error {
	logger.Info("开始生成标签页", "tagName", tagName)

	// 获取标签文章
	articles, _, err := s.articleModel.GetByTag(tagName, 1, 100)
	if err != nil {
		logger.Error("获取标签文章失败", "tagName", tagName, "error", err)
		return err
	}

	// 获取全局变量
	globals := s.templateService.GetGlobals()

	// 准备模板数据
	data := map[string]interface{}{
		"Globals":     globals,
		"TagName":     tagName,
		"Articles":    articles,
		"PageTitle":   tagName + " - " + s.config.Site.Name,
		"Keywords":    tagName,
		"Description": "标签 " + tagName + " 下的所有文章",
	}

	// 确定模板文件
	tplFile := s.config.Template.DefaultTpl + "/tag.htm"

	// 生成静态页面
	staticPath := fmt.Sprintf("tag/%s.html", tagName)

	// 确保目录存在
	dir := filepath.Dir(staticPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("创建目录失败", "dir", dir, "error", err)
		return err
	}

	// 生成静态页面
	err = s.templateService.GenerateStaticPage(tplFile, data, staticPath)
	if err != nil {
		logger.Error("生成静态标签页失败", "tagName", tagName, "error", err)
		return err
	}

	// 生成移动端标签页
	mobileTplFile := s.config.Template.DefaultTpl + "/mobile/tag.htm"
	if fileExists(filepath.Join(s.config.Template.Dir, mobileTplFile)) {
		mobileStaticPath := fmt.Sprintf("m/tag/%s.html", tagName)

		// 确保目录存在
		mobileDir := filepath.Dir(mobileStaticPath)
		if err := os.MkdirAll(mobileDir, 0755); err != nil {
			logger.Error("创建移动端目录失败", "dir", mobileDir, "error", err)
		} else {
			err = s.templateService.GenerateStaticPage(mobileTplFile, data, mobileStaticPath)
			if err != nil {
				logger.Error("生成移动端静态标签页失败", "tagName", tagName, "error", err)
			}
		}
	}

	logger.Info("标签页生成完成", "tagName", tagName)
	return nil
}
