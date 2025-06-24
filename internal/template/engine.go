package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"aq3cms/config"
	"aq3cms/internal/template/tags"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/logger"
)

// Engine 模板引擎
type Engine struct {
	config      *config.TemplateConfig
	cache       cache.Cache
	funcMap     template.FuncMap
	tagHandlers map[string]TagHandler
	mutex       sync.RWMutex
}

// TagHandler 标签处理器接口
type TagHandler interface {
	Handle(attrs map[string]string, content string, data interface{}) (string, error)
}

// New 创建新的模板引擎
func New(cfg *config.TemplateConfig, cache cache.Cache) *Engine {
	engine := &Engine{
		config:      cfg,
		cache:       cache,
		funcMap:     make(template.FuncMap),
		tagHandlers: make(map[string]TagHandler),
	}

	// 注册内置函数
	engine.registerBuiltinFuncs()

	// 注册内置标签处理器
	engine.registerBuiltinTags()

	return engine
}

// RegisterFunc 注册模板函数
func (e *Engine) RegisterFunc(name string, fn interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.funcMap[name] = fn
}

// RegisterTag 注册标签处理器
func (e *Engine) RegisterTag(name string, handler TagHandler) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.tagHandlers[name] = handler
}

// Render 渲染模板
func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	// 如果是HTTP响应，设置正确的Content-Type
	if httpWriter, ok := w.(http.ResponseWriter); ok {
		httpWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	// 检查缓存
	if e.config.Cache {
		if cached, ok := e.cache.Get("template:" + name); ok {
			if content, ok := cached.(string); ok {
				_, err := io.WriteString(w, content)
				return err
			}
		}
	}

	// 获取模板内容
	content, err := e.getTemplateContent(name)
	if err != nil {
		return err
	}

	// 解析模板标签
	parsedContent, err := e.parseTemplate(content, data)
	if err != nil {
		return err
	}

	// 使用Go标准模板引擎渲染
	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(parsedContent)
	if err != nil {
		return err
	}

	// 如果启用缓存，先渲染到缓冲区
	if e.config.Cache {
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return err
		}

		// 存入缓存
		cache.SafeSet(e.cache, "template:"+name, buf.String(), 0)

		// 写入响应
		_, err = io.Copy(w, &buf)
		return err
	}

	// 直接渲染到响应
	return tmpl.Execute(w, data)
}

// 获取模板内容
func (e *Engine) getTemplateContent(name string) (string, error) {
	// 构建完整路径
	path := filepath.Join(e.config.Dir, name)

	// 读取文件
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %v", err)
	}

	return string(content), nil
}

// 解析模板中的aq3cmsCMS风格标签
func (e *Engine) parseTemplate(content string, data interface{}) (string, error) {
	// 解析 {aq3cms:标签 属性=值}内容{/aq3cms:标签} 格式的标签
	// 不使用反向引用，而是在代码中手动处理标签匹配
	tagPattern := regexp.MustCompile(`{aq3cms:([a-zA-Z0-9_]+)([^}]*)}([\s\S]*?){/aq3cms:([a-zA-Z0-9_]+)}`)

	// 解析 {aq3cms:标签 属性=值/} 格式的自闭合标签
	selfClosingTagPattern := regexp.MustCompile(`{aq3cms:([a-zA-Z0-9_]+)([^}]*)/}`)

	// 解析 {aq3cms:field.字段名/} 格式的标签
	fieldPattern := regexp.MustCompile(`{aq3cms:field\.([a-zA-Z0-9_]+)/}`)

	// 解析 {aq3cms:global.变量名/} 格式的标签
	globalPattern := regexp.MustCompile(`{aq3cms:global\.([a-zA-Z0-9_]+)/}`)

	// 处理复杂标签
	content = tagPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := tagPattern.FindStringSubmatch(match)
		if len(submatches) < 5 {
			return match
		}

		// 检查开始标签和结束标签是否匹配
		startTag := submatches[1]
		endTag := submatches[4]
		if startTag != endTag {
			logger.Warn("标签不匹配", "startTag", startTag, "endTag", endTag)
			// 使用开始标签作为标签名，忽略结束标签不匹配的问题
			endTag = startTag
		}

		tagName := startTag
		attrStr := submatches[2]
		innerContent := submatches[3]

		// 解析属性
		attrs := parseAttributes(attrStr)

		// 查找标签处理器
		e.mutex.RLock()
		handler, exists := e.tagHandlers[tagName]
		e.mutex.RUnlock()

		if !exists {
			logger.Warn("未找到标签处理器", "tag", tagName)
			return match
		}

		// 处理标签
		result, err := handler.Handle(attrs, innerContent, data)
		if err != nil {
			logger.Error("处理标签失败", "tag", tagName, "error", err)
			return match
		}

		return result
	})

	// 处理字段标签
	content = fieldPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := fieldPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		fieldName := submatches[1]
		return fmt.Sprintf("{{.Fields.%s}}", fieldName)
	})

	// 处理全局变量标签
	content = globalPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := globalPattern.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		varName := submatches[1]
		return fmt.Sprintf("{{.Globals.%s}}", varName)
	})

	// 处理自闭合标签
	content = selfClosingTagPattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := selfClosingTagPattern.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		tagName := submatches[1]
		attrStr := submatches[2]

		// 解析属性
		attrs := parseAttributes(attrStr)

		// 查找标签处理器
		e.mutex.RLock()
		handler, exists := e.tagHandlers[tagName]
		e.mutex.RUnlock()

		if !exists {
			logger.Warn("未找到标签处理器", "tag", tagName)
			return match
		}

		// 处理标签
		result, err := handler.Handle(attrs, "", data)
		if err != nil {
			logger.Error("处理标签失败", "tag", tagName, "error", err)
			return match
		}

		return result
	})

	return content, nil
}

// 解析标签属性
func parseAttributes(attrStr string) map[string]string {
	attrs := make(map[string]string)
	attrStr = strings.TrimSpace(attrStr)

	if attrStr == "" {
		return attrs
	}

	// 匹配 name='value' 或 name="value" 格式的属性
	// 分别处理单引号和双引号的情况
	singleQuotePattern := regexp.MustCompile(`([a-zA-Z0-9_]+)\s*=\s*'([^']*)'`)
	doubleQuotePattern := regexp.MustCompile(`([a-zA-Z0-9_]+)\s*=\s*"([^"]*)"`)

	// 先处理单引号属性
	singleMatches := singleQuotePattern.FindAllStringSubmatch(attrStr, -1)
	for _, match := range singleMatches {
		if len(match) >= 3 {
			attrs[match[1]] = match[2]
		}
	}

	// 再处理双引号属性
	doubleMatches := doubleQuotePattern.FindAllStringSubmatch(attrStr, -1)
	for _, match := range doubleMatches {
		if len(match) >= 3 {
			attrs[match[1]] = match[2]
		}
	}

	return attrs
}

// 注册内置函数
func (e *Engine) registerBuiltinFuncs() {
	// 注册所有模板函数
	funcs := TemplateFunctions()
	for name, fn := range funcs {
		e.RegisterFunc(name, fn)
	}
}

// 注册内置标签处理器
func (e *Engine) registerBuiltinTags() {
	// 注册aq3cmsCMS标签处理器
	e.RegisterTag("arclist", &tags.ArcListTag{})
	e.RegisterTag("channel", &tags.ChannelTag{})
	e.RegisterTag("tag", &tags.TagTag{})
	e.RegisterTag("flink", &tags.FLinkTag{})
}
