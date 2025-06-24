package tags

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// FLinkTag 友情链接标签处理器
type FLinkTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *FLinkTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	row := 10
	if rowStr, ok := attrs["row"]; ok {
		if r, err := strconv.Atoi(rowStr); err == nil && r > 0 {
			row = r
		}
	}

	titlelen := 24
	if titlelenStr, ok := attrs["titlelen"]; ok {
		if tl, err := strconv.Atoi(titlelenStr); err == nil && tl > 0 {
			titlelen = tl
		}
	}

	typeid := 0
	if typeidStr, ok := attrs["typeid"]; ok {
		if tid, err := strconv.Atoi(typeidStr); err == nil {
			typeid = tid
		}
	}

	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "flink")
	qb.Select("*")

	// 添加条件
	qb.Where("ischeck = 1")
	if typeid > 0 {
		qb.Where("typeid = ?", typeid)
	}

	// 设置排序和限制
	qb.OrderBy("id ASC")
	qb.Limit(row)

	// 执行查询
	links, err := qb.Get()
	if err != nil {
		logger.Error("查询友情链接失败", "error", err)
		return "", err
	}

	// 如果没有链接，返回空字符串
	if len(links) == 0 {
		return "", nil
	}

	// 处理每个链接
	var result bytes.Buffer
	for _, link := range links {
		// 创建字段映射
		fields := make(map[string]interface{})
		for k, v := range link {
			fields[k] = v
		}

		// 处理标题长度
		if title, ok := fields["webname"].(string); ok && len(title) > titlelen {
			fields["webname"] = title[:titlelen] + "..."
		}

		// 处理内容中的字段标签
		itemContent := content
		fieldPattern := regexp.MustCompile(`\[field:([a-zA-Z0-9_]+)(?:\s+function="([^"]+)")?\s*/\]`)
		itemContent = fieldPattern.ReplaceAllStringFunc(itemContent, func(match string) string {
			matches := fieldPattern.FindStringSubmatch(match)
			if len(matches) < 2 {
				return match
			}

			fieldName := matches[1]
			if value, ok := fields[fieldName]; ok {
				// 如果有函数，则应用函数
				if len(matches) > 2 && matches[2] != "" {
					funcName := matches[2]
					// 处理函数，这里可以添加更多函数
					if funcName == "substr(0,30)" && fieldName == "webname" {
						if str, ok := value.(string); ok && len(str) > 30 {
							return str[:30] + "..."
						}
					}
				}
				return fmt.Sprintf("%v", value)
			}

			return ""
		})

		result.WriteString(itemContent)
	}

	return result.String(), nil
}
