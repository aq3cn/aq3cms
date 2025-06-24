package security

import (
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	// 危险标签正则表达式
	dangerousTags = regexp.MustCompile(`(?i)<(script|iframe|object|embed|form|style|link|meta|base|applet|param|layer|frameset|frame|ilayer|bgsound|title|xml|svg)`)
	
	// 危险属性正则表达式
	dangerousAttrs = regexp.MustCompile(`(?i)(on\w+|style|formaction|xmlns|xlink:href)=["']?[^>"']*["']?`)
	
	// JavaScript URL正则表达式
	jsURLs = regexp.MustCompile(`(?i)javascript:`)
	
	// 数据URL正则表达式
	dataURLs = regexp.MustCompile(`(?i)data:`)
	
	// 注释正则表达式
	comments = regexp.MustCompile(`<!--.*?-->`)
	
	// 允许的标签和属性
	allowedTags = map[string][]string{
		"a":          {"href", "target", "title", "class", "id", "rel"},
		"abbr":       {"title", "class", "id"},
		"address":    {"class", "id"},
		"area":       {"shape", "coords", "href", "alt", "target"},
		"article":    {"class", "id"},
		"aside":      {"class", "id"},
		"audio":      {"autoplay", "controls", "loop", "preload", "src", "class", "id"},
		"b":          {"class", "id"},
		"bdi":        {"dir", "class", "id"},
		"bdo":        {"dir", "class", "id"},
		"blockquote": {"cite", "class", "id"},
		"br":         {"class", "id"},
		"caption":    {"class", "id"},
		"cite":       {"class", "id"},
		"code":       {"class", "id"},
		"col":        {"align", "valign", "span", "width", "class", "id"},
		"colgroup":   {"align", "valign", "span", "width", "class", "id"},
		"dd":         {"class", "id"},
		"del":        {"datetime", "class", "id"},
		"details":    {"open", "class", "id"},
		"dfn":        {"class", "id"},
		"div":        {"class", "id"},
		"dl":         {"class", "id"},
		"dt":         {"class", "id"},
		"em":         {"class", "id"},
		"figcaption": {"class", "id"},
		"figure":     {"class", "id"},
		"footer":     {"class", "id"},
		"h1":         {"class", "id"},
		"h2":         {"class", "id"},
		"h3":         {"class", "id"},
		"h4":         {"class", "id"},
		"h5":         {"class", "id"},
		"h6":         {"class", "id"},
		"header":     {"class", "id"},
		"hr":         {"class", "id"},
		"i":          {"class", "id"},
		"img":        {"src", "alt", "title", "width", "height", "class", "id"},
		"ins":        {"datetime", "class", "id"},
		"li":         {"class", "id"},
		"mark":       {"class", "id"},
		"nav":        {"class", "id"},
		"ol":         {"class", "id"},
		"p":          {"class", "id"},
		"pre":        {"class", "id"},
		"q":          {"cite", "class", "id"},
		"s":          {"class", "id"},
		"section":    {"class", "id"},
		"small":      {"class", "id"},
		"span":       {"class", "id"},
		"strong":     {"class", "id"},
		"sub":        {"class", "id"},
		"summary":    {"class", "id"},
		"sup":        {"class", "id"},
		"table":      {"width", "border", "align", "valign", "class", "id"},
		"tbody":      {"align", "valign", "class", "id"},
		"td":         {"width", "rowspan", "colspan", "align", "valign", "class", "id"},
		"tfoot":      {"align", "valign", "class", "id"},
		"th":         {"width", "rowspan", "colspan", "align", "valign", "class", "id"},
		"thead":      {"align", "valign", "class", "id"},
		"time":       {"datetime", "class", "id"},
		"tr":         {"align", "valign", "class", "id"},
		"u":          {"class", "id"},
		"ul":         {"class", "id"},
		"video":      {"autoplay", "controls", "loop", "preload", "src", "height", "width", "class", "id"},
		"wbr":        {"class", "id"},
	}
)

// CleanHTML 清理HTML，防止XSS攻击
func CleanHTML(input string) string {
	// 如果输入为空，直接返回
	if input == "" {
		return input
	}
	
	// 移除注释
	input = comments.ReplaceAllString(input, "")
	
	// 移除危险标签
	input = dangerousTags.ReplaceAllStringFunc(input, func(match string) string {
		return "&lt;" + match[1:]
	})
	
	// 移除危险属性
	input = dangerousAttrs.ReplaceAllString(input, "")
	
	// 移除JavaScript URL
	input = jsURLs.ReplaceAllString(input, "invalid:")
	
	// 移除数据URL
	input = dataURLs.ReplaceAllString(input, "invalid:")
	
	return input
}

// StripTags 去除所有HTML标签
func StripTags(input string) string {
	// 如果输入为空，直接返回
	if input == "" {
		return input
	}
	
	// 匹配所有HTML标签
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	
	// 去除所有标签
	return tagPattern.ReplaceAllString(input, "")
}

// EscapeHTML 转义HTML特殊字符
func EscapeHTML(input string) string {
	return html.EscapeString(input)
}

// EscapeURL 转义URL
func EscapeURL(input string) string {
	return url.QueryEscape(input)
}

// SanitizeFilename 清理文件名
func SanitizeFilename(filename string) string {
	// 移除路径分隔符和其他危险字符
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, ":", "")
	filename = strings.ReplaceAll(filename, "*", "")
	filename = strings.ReplaceAll(filename, "?", "")
	filename = strings.ReplaceAll(filename, "\"", "")
	filename = strings.ReplaceAll(filename, "<", "")
	filename = strings.ReplaceAll(filename, ">", "")
	filename = strings.ReplaceAll(filename, "|", "")
	
	return filename
}
