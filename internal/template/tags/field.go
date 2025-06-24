package tags

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FieldTag 字段标签处理器
type FieldTag struct {}

// Handle 处理标签
func (t *FieldTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 获取字段名
	name := ""
	for k, v := range attrs {
		if k == "name" {
			name = v
			break
		}
	}

	if name == "" {
		return "", fmt.Errorf("字段标签缺少name属性")
	}

	// 获取字段值
	value := getFieldValue(name, data)
	if value == nil {
		return "", nil
	}

	// 处理函数
	if function, ok := attrs["function"]; ok {
		value = applyFunction(function, value)
	}

	// 返回字段值
	return fmt.Sprintf("%v", value), nil
}

// getFieldValue 获取字段值
func getFieldValue(name string, data interface{}) interface{} {
	// 处理嵌套字段，如 article.title
	parts := strings.Split(name, ".")

	// 获取数据
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	// 处理第一级
	current := dataMap

	// 处理嵌套字段
	for i, part := range parts {
		if i == len(parts)-1 {
			// 最后一级，返回值
			if value, ok := current[part]; ok {
				return value
			}
			return nil
		}

		// 中间级，继续查找
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// applyFunction 应用函数处理
func applyFunction(function string, value interface{}) interface{} {
	// 解析函数名和参数
	parts := strings.Split(function, "(")
	if len(parts) != 2 {
		return value
	}

	funcName := strings.TrimSpace(parts[0])
	params := strings.TrimRight(parts[1], ")")
	paramList := strings.Split(params, ",")

	// 处理不同的函数
	switch funcName {
	case "substring":
		// 截取字符串
		if len(paramList) >= 2 {
			str := fmt.Sprintf("%v", value)
			start, _ := strconv.Atoi(strings.TrimSpace(paramList[0]))
			length, _ := strconv.Atoi(strings.TrimSpace(paramList[1]))

			runes := []rune(str)
			if start < 0 {
				start = 0
			}
			if start > len(runes) {
				return ""
			}
			end := start + length
			if end > len(runes) {
				end = len(runes)
			}
			return string(runes[start:end])
		}
	case "strftime":
		// 格式化时间
		if len(paramList) >= 1 {
			format := strings.Trim(strings.TrimSpace(paramList[0]), "'\"")

			var timestamp int64
			switch v := value.(type) {
			case int64:
				timestamp = v
			case int:
				timestamp = int64(v)
			case string:
				t, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return value
				}
				timestamp = t
			default:
				return value
			}

			// 替换常用格式
			format = strings.Replace(format, "Y", "2006", -1)
			format = strings.Replace(format, "m", "01", -1)
			format = strings.Replace(format, "d", "02", -1)
			format = strings.Replace(format, "H", "15", -1)
			format = strings.Replace(format, "i", "04", -1)
			format = strings.Replace(format, "s", "05", -1)

			return time.Unix(timestamp, 0).Format(format)
		}
	}

	return value
}
