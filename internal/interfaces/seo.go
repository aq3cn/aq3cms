/*
 * @Author: xxx@xxx.com
 * @Date: 2025-05-01 22:01:21
 * @LastEditors: xxx@xxx.com
 * @LastEditTime: 2025-05-01 22:03:47
 * @FilePath: \aq3cms\aq3cms\internal\interfaces\seo.go
 * @Description:
 *
 * Copyright (c) 2022 by xxx@xxx.com, All Rights Reserved.
 */
package interfaces

// SEOServiceInterface SEO服务接口
type SEOServiceInterface interface {
	// GetMetaTags 获取Meta标签
	GetMetaTags(title, keywords, description string) map[string]string

	// GetOpenGraphTags 获取Open Graph标签
	GetOpenGraphTags(title, description, url, image string) map[string]string

	// GetTwitterCardTags 获取Twitter Card标签
	GetTwitterCardTags(title, description, image string) map[string]string

	// GetCanonicalURL 获取规范URL
	GetCanonicalURL(path string) string

	// GetAlternateURLs 获取备用URL
	GetAlternateURLs(path string) map[string]string

	// GetAvailableLangs 获取可用语言
	GetAvailableLangs() []map[string]string

	// GenerateSitemap 生成站点地图
	GenerateSitemap() (string, error)

	// GenerateRobotsTxt 生成robots.txt
	GenerateRobotsTxt() string
}
