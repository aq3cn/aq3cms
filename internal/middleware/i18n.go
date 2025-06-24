package middleware

import (
	"context"
	"net/http"
	"strings"

	"aq3cms/pkg/i18n"
)

// I18nMiddleware 多语言中间件
func I18nMiddleware(i18n *i18n.I18n) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取语言
			lang := getLang(r, i18n)

			// 将语言保存到请求上下文
			ctx := r.Context()
			ctx = context.WithValue(ctx, "lang", lang)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getLang 获取语言
func getLang(r *http.Request, i18n *i18n.I18n) string {
	// 从Cookie获取语言
	cookie, err := r.Cookie("lang")
	if err == nil && cookie.Value != "" {
		// 检查语言是否有效
		for _, lang := range i18n.GetLangs() {
			if lang == cookie.Value {
				return lang
			}
		}
	}

	// 从Accept-Language头获取语言
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		// 解析Accept-Language
		langs := parseAcceptLanguage(acceptLang)
		for _, lang := range langs {
			// 检查语言是否有效
			for _, validLang := range i18n.GetLangs() {
				if strings.HasPrefix(lang, validLang) {
					return validLang
				}
			}
		}
	}

	// 使用默认语言
	return i18n.GetDefaultLang()
}

// parseAcceptLanguage 解析Accept-Language头
func parseAcceptLanguage(acceptLang string) []string {
	// 分割语言
	parts := strings.Split(acceptLang, ",")
	langs := make([]string, 0, len(parts))

	// 解析每个语言
	for _, part := range parts {
		// 分割语言和权重
		langPart := strings.Split(part, ";")
		lang := strings.TrimSpace(langPart[0])
		langs = append(langs, lang)
	}

	return langs
}

// GetLang 从请求上下文获取语言
func GetLang(r *http.Request) string {
	// 从请求上下文获取语言
	if lang, ok := r.Context().Value("lang").(string); ok {
		return lang
	}

	// 默认返回空字符串
	return ""
}
