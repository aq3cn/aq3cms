package tags

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ChannelTag 栏目标签处理器
type ChannelTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *ChannelTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	typeid := attrs["typeid"]
	row := 10
	if rowStr, ok := attrs["row"]; ok {
		if r, err := strconv.Atoi(rowStr); err == nil && r > 0 {
			row = r
		}
	}

	currentstyle := ""
	if cs, ok := attrs["currentstyle"]; ok {
		currentstyle = cs
	}

	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "arctype")
	qb.Select("*")

	// 添加条件
	qb.Where("ishidden = 0")

	if typeid != "" {
		// 获取指定栏目的子栏目
		typeID, err := strconv.Atoi(typeid)
		if err == nil {
			qb.Where("reid = ?", typeID)
		}
	} else {
		// 获取顶级栏目
		qb.Where("reid = 0")
	}

	// 设置排序和限制
	qb.OrderBy("sortrank ASC")
	qb.Limit(row)

	// 执行查询
	channels, err := qb.Get()
	if err != nil {
		logger.Error("查询栏目列表失败", "error", err)
		return "", err
	}

	// 如果没有栏目，返回空字符串
	if len(channels) == 0 {
		return "", nil
	}

	// 获取当前栏目ID
	var currentTypeID int64
	if dataMap, ok := data.(map[string]interface{}); ok {
		if fields, ok := dataMap["Fields"].(map[string]interface{}); ok {
			if tid, ok := fields["typeid"].(int64); ok {
				currentTypeID = tid
			}
		}
	}

	// 处理每个栏目
	var result bytes.Buffer
	for _, channel := range channels {
		// 创建字段映射
		fields := make(map[string]interface{})
		for k, v := range channel {
			fields[k] = v
		}

		// 添加特殊字段
		fields["typeurl"] = getTypeUrl(channel)

		// 判断是否为当前栏目
		channelID, _ := channel["id"].(int64)
		if channelID == currentTypeID {
			fields["currentstyle"] = currentstyle
		} else {
			fields["currentstyle"] = ""
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
					// 这里可以添加函数处理逻辑
					_ = funcName
				}
				return fmt.Sprintf("%v", value)
			}

			return ""
		})

		result.WriteString(itemContent)
	}

	return result.String(), nil
}
