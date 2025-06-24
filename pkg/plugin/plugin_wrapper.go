package plugin

import (
	"fmt"
)

// PluginWrapper 插件包装器
type PluginWrapper struct {
	path string
}

// Open 打开插件
func Open(path string) (*PluginWrapper, error) {
	// 这里简化处理，实际应该使用 plugin.Open
	return &PluginWrapper{
		path: path,
	}, nil
}

// Lookup 查找插件符号
func (p *PluginWrapper) Lookup(symbol string) (interface{}, error) {
	// 这里简化处理，实际应该使用 plugin.Lookup
	return nil, fmt.Errorf("plugin symbol not found: %s", symbol)
}
