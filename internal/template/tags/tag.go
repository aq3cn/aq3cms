package tags

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// TagTag 标签标签处理器
type TagTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *TagTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	row := 10
	if rowStr, ok := attrs["row"]; ok {
		if r, err := strconv.Atoi(rowStr); err == nil && r > 0 {
			row = r
		}
	}

	orderby := "count"
	if ob, ok := attrs["orderby"]; ok && ob != "" {
		orderby = ob
	}

	orderway := "desc"
	if ow, ok := attrs["orderway"]; ok && (ow == "asc" || ow == "desc") {
		orderway = ow
	}

	isHot := -1
	if isHotStr, ok := attrs["ishot"]; ok {
		if ih, err := strconv.Atoi(isHotStr); err == nil {
			isHot = ih
		}
	}

	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "tagindex")
	qb.Select("*")

	// 添加条件
	if isHot >= 0 {
		qb.Where("ishot = ?", isHot)
	}

	// 设置排序和限制
	qb.OrderBy(orderby + " " + orderway)
	qb.Limit(row)

	// 执行查询
	tags, err := qb.Get()
	if err != nil {
		logger.Error("查询标签列表失败", "error", err)
		return "", err
	}

	// 如果没有标签，返回空字符串
	if len(tags) == 0 {
		return "", nil
	}

	// 处理每个标签
	var result bytes.Buffer
	for _, tag := range tags {
		// 创建字段映射
		fields := make(map[string]interface{})
		for k, v := range tag {
			fields[k] = v
		}

		// 添加特殊字段
		fields["tagurl"] = getTagUrl(tag)

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
					if funcName == "rand(1, 5)" && fieldName == "id" {
						// 简单实现随机数
						return fmt.Sprintf("%d", (tag["id"].(int64) % 5) + 1)
					}
					// 其他函数可以在这里添加
				}
				return fmt.Sprintf("%v", value)
			}

			return ""
		})

		result.WriteString(itemContent)
	}

	return result.String(), nil
}

// getTagUrl 获取标签URL
func getTagUrl(tag map[string]interface{}) string {
	tagName, _ := tag["tag"].(string)
	return fmt.Sprintf("/tag/%s.html", tagName)
}
