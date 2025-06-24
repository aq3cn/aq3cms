package tags

import (
	"bytes"
	"fmt"
	"strings"

	"aq3cms/internal/interfaces"
)

// I18nTag 国际化标签
type I18nTag struct {
	I18nService interfaces.I18nServiceInterface
}

// Parse 解析标签
func (t *I18nTag) Parse(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取语言
	lang := params["lang"]
	if lang == "" {
		// 从上下文获取语言
		if langCtx, ok := ctx["Lang"].(string); ok {
			lang = langCtx
		} else {
			// 使用默认语言
			lang = t.I18nService.GetDefaultLang()
		}
	}

	// 获取键
	key := params["key"]
	if key == "" {
		return "", fmt.Errorf("i18n tag requires key parameter")
	}

	// 获取参数
	args := make([]interface{}, 0)
	for i := 1; ; i++ {
		argKey := fmt.Sprintf("arg%d", i)
		if arg, ok := params[argKey]; ok {
			args = append(args, arg)
		} else {
			break
		}
	}

	// 翻译
	return t.I18nService.T(lang, key, args...), nil
}

// ParseBlock 解析块标签
func (t *I18nTag) ParseBlock(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取语言
	lang := params["lang"]
	if lang == "" {
		// 从上下文获取语言
		if langCtx, ok := ctx["Lang"].(string); ok {
			lang = langCtx
		} else {
			// 使用默认语言
			lang = t.I18nService.GetDefaultLang()
		}
	}

	// 分割内容
	parts := strings.Split(content, "|")
	if len(parts) == 1 {
		// 只有一个部分，直接返回
		return parts[0], nil
	}

	// 获取可用语言
	langs := t.I18nService.GetLangs()

	// 查找匹配的语言
	for i, l := range langs {
		if l == lang && i < len(parts) {
			return parts[i], nil
		}
	}

	// 没有找到匹配的语言，返回第一个部分
	return parts[0], nil
}

// GetLangSelector 获取语言选择器
func (t *I18nTag) GetLangSelector(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取当前语言
	currentLang := ""
	if langCtx, ok := ctx["Lang"].(string); ok {
		currentLang = langCtx
	} else {
		currentLang = t.I18nService.GetDefaultLang()
	}

	// 获取可用语言
	langs := t.I18nService.GetAvailableLangs()

	// 构建语言选择器
	var buf bytes.Buffer
	buf.WriteString(`<div class="lang-selector">`)
	for _, lang := range langs {
		if lang["Code"] == currentLang {
			buf.WriteString(fmt.Sprintf(`<span class="lang-item active">%s</span>`, lang["Name"]))
		} else {
			buf.WriteString(fmt.Sprintf(`<a href="?lang=%s" class="lang-item">%s</a>`, lang["Code"], lang["Name"]))
		}
	}
	buf.WriteString(`</div>`)

	return buf.String(), nil
}
