package tags

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ArcListTag 文章列表标签处理器
type ArcListTag struct {
	DB *database.DB
}

// Handle 处理标签
func (t *ArcListTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 解析属性
	typeid := attrs["typeid"]
	row := 10
	if rowStr, ok := attrs["row"]; ok {
		if r, err := strconv.Atoi(rowStr); err == nil && r > 0 {
			row = r
		}
	}

	orderby := "id"
	if ob, ok := attrs["orderby"]; ok && ob != "" {
		orderby = ob
	}

	orderway := "desc"
	if ow, ok := attrs["orderway"]; ok && (ow == "asc" || ow == "desc") {
		orderway = ow
	}

	// 构建查询
	qb := database.NewQueryBuilder(t.DB, "archives AS a")
	qb.Select("a.*", "t.typename", "t.typedir", "ad.body")
	qb.LeftJoin(t.DB.TableName("arctype")+" AS t", "a.typeid = t.id")
	qb.LeftJoin(t.DB.TableName("addonarticle")+" AS ad", "a.id = ad.aid")

	// 添加条件
	qb.Where("a.arcrank > -1")

	if typeid != "" {
		if strings.Contains(typeid, ",") {
			// 多个栏目ID
			qb.Where("a.typeid IN (" + typeid + ")")
		} else {
			// 单个栏目ID
			typeID, err := strconv.Atoi(typeid)
			if err == nil {
				qb.Where("a.typeid = ?", typeID)
			}
		}
	}

	// 设置排序和限制
	qb.OrderBy("a." + orderby + " " + orderway)
	qb.Limit(row)

	// 执行查询
	articles, err := qb.Get()
	if err != nil {
		logger.Error("查询文章列表失败", "error", err)
		return "", err
	}

	// 如果没有文章，返回空字符串
	if len(articles) == 0 {
		return "", nil
	}

	// 处理每篇文章
	var result bytes.Buffer
	for _, article := range articles {
		// 创建字段映射
		fields := make(map[string]interface{})
		for k, v := range article {
			fields[k] = v
		}

		// 添加特殊字段
		fields["arcurl"] = getArcUrl(article)
		fields["typeurl"] = getTypeUrl(article)

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
					// 处理日期格式化函数
					if fieldName == "pubdate" && funcName == "date('Y-m-d',@me)" {
						if timestamp, ok := value.(int64); ok {
							return time.Unix(timestamp, 0).Format("2006-01-02")
						}
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

// getArcUrl 获取文章URL
func getArcUrl(article map[string]interface{}) string {
	id, _ := article["id"].(int64)
	return fmt.Sprintf("/article/%d.html", id)
}

// getTypeUrl 获取栏目URL
func getTypeUrl(article map[string]interface{}) string {
	typeid, _ := article["typeid"].(int64)
	return fmt.Sprintf("/list/%d.html", typeid)
}
