package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"aq3cms/pkg/logger"
)

// I18n 国际化
type I18n struct {
	defaultLang string
	langs       map[string]map[string]string
	mutex       sync.RWMutex
}

// New 创建国际化
func New(defaultLang, langDir string) (*I18n, error) {
	i18n := &I18n{
		defaultLang: defaultLang,
		langs:       make(map[string]map[string]string),
	}

	// 加载语言文件
	err := i18n.loadLangs(langDir)
	if err != nil {
		return nil, err
	}

	return i18n, nil
}

// NewMemory 创建内存中的国际化实例
func NewMemory(defaultLang string) *I18n {
	return &I18n{
		defaultLang: defaultLang,
		langs:       make(map[string]map[string]string),
	}
}

// loadLangs 加载语言文件
func (i *I18n) loadLangs(langDir string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 检查目录是否存在
	_, err := os.Stat(langDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("language directory not found: %s", langDir)
	}

	// 遍历目录
	files, err := ioutil.ReadDir(langDir)
	if err != nil {
		return err
	}

	// 加载语言文件
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

		// 读取文件
		data, err := ioutil.ReadFile(filepath.Join(langDir, file.Name()))
		if err != nil {
			logger.Error("读取语言文件失败", "file", file.Name(), "error", err)
			continue
		}

		// 解析JSON
		var messages map[string]string
		err = json.Unmarshal(data, &messages)
		if err != nil {
			logger.Error("解析语言文件失败", "file", file.Name(), "error", err)
			continue
		}

		// 保存语言
		i.langs[lang] = messages
	}

	// 检查默认语言是否存在
	if _, ok := i.langs[i.defaultLang]; !ok {
		if len(i.langs) > 0 {
			// 使用第一个语言作为默认语言
			for lang := range i.langs {
				i.defaultLang = lang
				break
			}
		} else {
			return fmt.Errorf("no language files found")
		}
	}

	return nil
}

// T 翻译
func (i *I18n) T(lang, key string, args ...interface{}) string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 检查语言是否存在
	messages, ok := i.langs[lang]
	if !ok {
		// 使用默认语言
		messages = i.langs[i.defaultLang]
	}

	// 获取翻译
	message, ok := messages[key]
	if !ok {
		// 检查默认语言
		if lang != i.defaultLang {
			if defaultMessages, ok := i.langs[i.defaultLang]; ok {
				if defaultMessage, ok := defaultMessages[key]; ok {
					message = defaultMessage
				} else {
					message = key
				}
			} else {
				message = key
			}
		} else {
			message = key
		}
	}

	// 格式化翻译
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return message
}

// GetLangs 获取所有语言
func (i *I18n) GetLangs() []string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	langs := make([]string, 0, len(i.langs))
	for lang := range i.langs {
		langs = append(langs, lang)
	}

	return langs
}

// GetDefaultLang 获取默认语言
func (i *I18n) GetDefaultLang() string {
	return i.defaultLang
}

// SetDefaultLang 设置默认语言
func (i *I18n) SetDefaultLang(lang string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if _, ok := i.langs[lang]; ok {
		i.defaultLang = lang
	}
}

// AddLang 添加语言
func (i *I18n) AddLang(lang string, messages map[string]string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.langs[lang] = messages
}

// AddMessage 添加翻译
func (i *I18n) AddMessage(lang, key, message string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if _, ok := i.langs[lang]; !ok {
		i.langs[lang] = make(map[string]string)
	}

	i.langs[lang][key] = message
}

// SaveLang 保存语言
func (i *I18n) SaveLang(lang, langDir string) error {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// 检查语言是否存在
	messages, ok := i.langs[lang]
	if !ok {
		return fmt.Errorf("language not found: %s", lang)
	}

	// 检查目录是否存在
	_, err := os.Stat(langDir)
	if os.IsNotExist(err) {
		// 创建目录
		err = os.MkdirAll(langDir, 0755)
		if err != nil {
			return err
		}
	}

	// 序列化JSON
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	err = ioutil.WriteFile(filepath.Join(langDir, lang+".json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
