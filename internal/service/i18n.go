package service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"aq3cms/config"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/i18n"
	"aq3cms/pkg/logger"
)

// I18nService 国际化服务
type I18nService struct {
	db     *database.DB
	cache  cache.Cache
	config *config.Config
	i18n   *i18n.I18n
}

// NewI18nService 创建国际化服务
func NewI18nService(db *database.DB, cache cache.Cache, config *config.Config) *I18nService {
	// 初始化国际化
	langDir := filepath.Join(config.Template.Dir, "lang")
	i18nInstance, err := i18n.New(config.Site.DefaultLang, langDir)
	if err != nil {
		logger.Error("初始化国际化失败", "error", err)
		// 使用默认语言
		i18nInstance, _ = i18n.New("zh-cn", langDir)
	}

	return &I18nService{
		db:     db,
		cache:  cache,
		config: config,
		i18n:   i18nInstance,
	}
}

// T 翻译
func (s *I18nService) T(lang, key string, args ...interface{}) string {
	return s.i18n.T(lang, key, args...)
}

// GetLangs 获取所有语言
func (s *I18nService) GetLangs() []string {
	return s.i18n.GetLangs()
}

// GetDefaultLang 获取默认语言
func (s *I18nService) GetDefaultLang() string {
	return s.i18n.GetDefaultLang()
}

// SetDefaultLang 设置默认语言
func (s *I18nService) SetDefaultLang(lang string) {
	s.i18n.SetDefaultLang(lang)
}

// AddLang 添加语言
func (s *I18nService) AddLang(lang string, messages map[string]string) {
	s.i18n.AddLang(lang, messages)
}

// AddMessage 添加翻译
func (s *I18nService) AddMessage(lang, key, message string) {
	s.i18n.AddMessage(lang, key, message)
}

// SaveLang 保存语言
func (s *I18nService) SaveLang(lang string) error {
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	return s.i18n.SaveLang(lang, langDir)
}

// ImportLang 导入语言
func (s *I18nService) ImportLang(lang, filePath string) error {
	// 读取文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析JSON
	var messages map[string]string
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return err
	}

	// 添加语言
	s.i18n.AddLang(lang, messages)

	// 保存语言
	return s.SaveLang(lang)
}

// ExportLang 导出语言
func (s *I18nService) ExportLang(lang, filePath string) error {
	// 获取语言目录
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	
	// 读取语言文件
	data, err := ioutil.ReadFile(filepath.Join(langDir, lang+".json"))
	if err != nil {
		return err
	}

	// 创建目录
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 写入文件
	return ioutil.WriteFile(filePath, data, 0644)
}

// CreateLangFile 创建语言文件
func (s *I18nService) CreateLangFile(lang string) error {
	// 获取语言目录
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	
	// 检查目录是否存在
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(langDir, 0755); err != nil {
			return err
		}
	}

	// 创建空的语言文件
	filePath := filepath.Join(langDir, lang+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建文件
		emptyMessages := map[string]string{}
		data, err := json.MarshalIndent(emptyMessages, "", "  ")
		if err != nil {
			return err
		}
		return ioutil.WriteFile(filePath, data, 0644)
	}

	return nil
}

// GetLangMessages 获取语言消息
func (s *I18nService) GetLangMessages(lang string) (map[string]string, error) {
	// 获取语言目录
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	
	// 读取语言文件
	data, err := ioutil.ReadFile(filepath.Join(langDir, lang+".json"))
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var messages map[string]string
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateLangMessages 更新语言消息
func (s *I18nService) UpdateLangMessages(lang string, messages map[string]string) error {
	// 添加语言
	s.i18n.AddLang(lang, messages)

	// 保存语言
	return s.SaveLang(lang)
}

// DeleteLang 删除语言
func (s *I18nService) DeleteLang(lang string) error {
	// 获取语言目录
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	
	// 删除语言文件
	filePath := filepath.Join(langDir, lang+".json")
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}

	return nil
}

// GetAvailableLangs 获取可用语言
func (s *I18nService) GetAvailableLangs() []map[string]string {
	// 获取语言目录
	langDir := filepath.Join(s.config.Template.Dir, "lang")
	
	// 获取语言文件列表
	files, err := ioutil.ReadDir(langDir)
	if err != nil {
		return []map[string]string{}
	}

	// 构建语言列表
	langs := make([]map[string]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 检查文件扩展名
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// 获取语言代码
		lang := strings.TrimSuffix(file.Name(), ".json")

		// 添加语言
		langs = append(langs, map[string]string{
			"Code": lang,
			"Name": getLangName(lang),
		})
	}

	return langs
}

// 获取语言名称
func getLangName(lang string) string {
	switch lang {
	case "zh-cn":
		return "简体中文"
	case "zh-tw":
		return "繁體中文"
	case "en":
		return "English"
	case "ja":
		return "日本語"
	case "ko":
		return "한국어"
	case "fr":
		return "Français"
	case "de":
		return "Deutsch"
	case "es":
		return "Español"
	case "it":
		return "Italiano"
	case "ru":
		return "Русский"
	default:
		return lang
	}
}
