package tags

import (
	"bytes"
	"fmt"

	"aq3cms/internal/interfaces"
)

// SEOTag SEO标签
type SEOTag struct {
	SEOService interfaces.SEOServiceInterface
}

// Parse 解析标签
func (t *SEOTag) Parse(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取类型
	tagType := params["type"]
	if tagType == "" {
		return "", fmt.Errorf("seo tag requires type parameter")
	}

	switch tagType {
	case "meta":
		return t.generateMetaTags(params, ctx)
	case "opengraph":
		return t.generateOpenGraphTags(params, ctx)
	case "twitter":
		return t.generateTwitterCardTags(params, ctx)
	case "canonical":
		return t.generateCanonicalLink(params, ctx)
	case "alternate":
		return t.generateAlternateLinks(params, ctx)
	case "sitemap":
		return t.generateSitemap(params, ctx)
	case "robots":
		return t.generateRobotsTxt(params, ctx)
	default:
		return "", fmt.Errorf("unknown seo tag type: %s", tagType)
	}
}

// ParseBlock 解析块标签
func (t *SEOTag) ParseBlock(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// SEO标签不支持块标签
	return content, nil
}

// 生成Meta标签
func (t *SEOTag) generateMetaTags(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取参数
	title := params["title"]
	keywords := params["keywords"]
	description := params["description"]

	// 获取Meta标签
	metaTags := t.SEOService.GetMetaTags(title, keywords, description)

	// 构建HTML
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`<title>%s</title>`, metaTags["title"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="keywords" content="%s">`, metaTags["keywords"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="description" content="%s">`, metaTags["description"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="author" content="%s">`, metaTags["author"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="generator" content="%s">`, metaTags["generator"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="robots" content="%s">`, metaTags["robots"]))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`<meta name="viewport" content="%s">`, metaTags["viewport"]))

	return buf.String(), nil
}

// 生成Open Graph标签
func (t *SEOTag) generateOpenGraphTags(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取参数
	title := params["title"]
	description := params["description"]
	url := params["url"]
	image := params["image"]

	// 获取Open Graph标签
	ogTags := t.SEOService.GetOpenGraphTags(title, description, url, image)

	// 构建HTML
	var buf bytes.Buffer
	for name, content := range ogTags {
		buf.WriteString(fmt.Sprintf(`<meta property="%s" content="%s">`, name, content))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// 生成Twitter Card标签
func (t *SEOTag) generateTwitterCardTags(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取参数
	title := params["title"]
	description := params["description"]
	image := params["image"]

	// 获取Twitter Card标签
	twitterTags := t.SEOService.GetTwitterCardTags(title, description, image)

	// 构建HTML
	var buf bytes.Buffer
	for name, content := range twitterTags {
		buf.WriteString(fmt.Sprintf(`<meta name="%s" content="%s">`, name, content))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// 生成规范链接
func (t *SEOTag) generateCanonicalLink(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取参数
	path := params["path"]

	// 获取规范URL
	canonicalURL := t.SEOService.GetCanonicalURL(path)

	// 构建HTML
	return fmt.Sprintf(`<link rel="canonical" href="%s">`, canonicalURL), nil
}

// 生成备用链接
func (t *SEOTag) generateAlternateLinks(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取参数
	path := params["path"]

	// 获取备用URL
	alternateURLs := t.SEOService.GetAlternateURLs(path)

	// 构建HTML
	var buf bytes.Buffer
	for lang, url := range alternateURLs {
		buf.WriteString(fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="%s">`, lang, url))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// 生成站点地图
func (t *SEOTag) generateSitemap(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 生成站点地图
	sitemap, err := t.SEOService.GenerateSitemap()
	if err != nil {
		return "", err
	}

	return sitemap, nil
}

// 生成robots.txt
func (t *SEOTag) generateRobotsTxt(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 生成robots.txt
	robotsTxt := t.SEOService.GenerateRobotsTxt()

	return robotsTxt, nil
}
