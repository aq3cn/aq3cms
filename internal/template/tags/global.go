package tags

import (
	"fmt"
)

// GlobalTag 全局变量标签处理器
type GlobalTag struct {
	Globals map[string]interface{}
}

// Handle 处理标签
func (t *GlobalTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 获取变量名
	name := ""
	for k, v := range attrs {
		if k == "name" {
			name = v
			break
		}
	}
	
	if name == "" {
		return "", fmt.Errorf("全局变量标签缺少name属性")
	}
	
	// 获取全局变量值
	value, ok := t.Globals[name]
	if !ok {
		return "", nil
	}
	
	// 返回变量值
	return fmt.Sprintf("%v", value), nil
}
