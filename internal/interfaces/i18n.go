/*
 * @Author: xxx@xxx.com
 * @Date: 2025-05-01 22:01:01
 * @LastEditors: xxx@xxx.com
 * @LastEditTime: 2025-05-01 22:02:08
 * @FilePath: \aq3cms\aq3cms\internal\interfaces\i18n.go
 * @Description:
 *
 * Copyright (c) 2022 by xxx@xxx.com, All Rights Reserved.
 */
package interfaces

// I18nServiceInterface 国际化服务接口
type I18nServiceInterface interface {
	// T 翻译
	T(lang, key string, args ...interface{}) string

	// GetDefaultLang 获取默认语言
	GetDefaultLang() string

	// GetLangs 获取可用语言代码
	GetLangs() []string

	// GetAvailableLangs 获取可用语言信息
	GetAvailableLangs() []map[string]string
}
