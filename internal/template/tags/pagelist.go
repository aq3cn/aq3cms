package tags

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PageListTag 分页标签处理器
type PageListTag struct {}

// Handle 处理标签
func (t *PageListTag) Handle(attrs map[string]string, content string, data interface{}) (string, error) {
	// 获取分页数据
	pagination, err := getPaginationData(data)
	if err != nil {
		return "", err
	}

	// 获取属性
	listSize := 10
	if sizeStr, ok := attrs["listsize"]; ok {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			listSize = size
		}
	}

	listItem := ""
	if item, ok := attrs["listitem"]; ok {
		listItem = item
	} else {
		listItem = "index,pre,pageno,next,end"
	}

	listStyle := ""
	if style, ok := attrs["liststyle"]; ok {
		listStyle = style
	}

	// 解析列表项
	items := strings.Split(listItem, ",")

	// 生成分页HTML
	var result bytes.Buffer

	// 添加样式
	if listStyle != "" {
		result.WriteString(fmt.Sprintf("<div class=\"%s\">\n", listStyle))
	} else {
		result.WriteString("<div class=\"pagelist\">\n")
	}

	// 处理每个列表项
	for _, item := range items {
		item = strings.TrimSpace(item)

		switch item {
		case "index":
			// 首页
			if pagination.CurrentPage > 1 {
				result.WriteString(fmt.Sprintf("<a href=\"%s\">首页</a>\n", getPageUrl(pagination.PageUrl, 1)))
			} else {
				result.WriteString("<span class=\"disabled\">首页</span>\n")
			}
		case "pre":
			// 上一页
			if pagination.HasPrev {
				result.WriteString(fmt.Sprintf("<a href=\"%s\">上一页</a>\n", getPageUrl(pagination.PageUrl, pagination.PrevPage)))
			} else {
				result.WriteString("<span class=\"disabled\">上一页</span>\n")
			}
		case "pageno":
			// 页码
			startPage := pagination.CurrentPage - listSize/2
			if startPage < 1 {
				startPage = 1
			}
			endPage := startPage + listSize - 1
			if endPage > pagination.TotalPages {
				endPage = pagination.TotalPages
				startPage = endPage - listSize + 1
				if startPage < 1 {
					startPage = 1
				}
			}

			for i := startPage; i <= endPage; i++ {
				if i == pagination.CurrentPage {
					result.WriteString(fmt.Sprintf("<span class=\"current\">%d</span>\n", i))
				} else {
					result.WriteString(fmt.Sprintf("<a href=\"%s\">%d</a>\n", getPageUrl(pagination.PageUrl, i), i))
				}
			}
		case "next":
			// 下一页
			if pagination.HasNext {
				result.WriteString(fmt.Sprintf("<a href=\"%s\">下一页</a>\n", getPageUrl(pagination.PageUrl, pagination.NextPage)))
			} else {
				result.WriteString("<span class=\"disabled\">下一页</span>\n")
			}
		case "end":
			// 末页
			if pagination.CurrentPage < pagination.TotalPages {
				result.WriteString(fmt.Sprintf("<a href=\"%s\">末页</a>\n", getPageUrl(pagination.PageUrl, pagination.TotalPages)))
			} else {
				result.WriteString("<span class=\"disabled\">末页</span>\n")
			}
		case "info":
			// 信息
			result.WriteString(fmt.Sprintf("<span class=\"pageinfo\">共 <strong>%d</strong> 页 <strong>%d</strong> 条</span>\n",
				pagination.TotalPages, pagination.TotalItems))
		}
	}

	result.WriteString("</div>\n")

	return result.String(), nil
}

// 分页数据结构
type paginationData struct {
	CurrentPage int
	TotalPages  int
	TotalItems  int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	PageUrl     string
}

// 获取分页数据
func getPaginationData(data interface{}) (*paginationData, error) {
	// 从模板数据中获取分页信息
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无法获取模板数据")
	}

	pagination, ok := dataMap["Pagination"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无法获取分页数据")
	}

	// 构建分页数据
	result := &paginationData{}

	if currentPage, ok := pagination["CurrentPage"].(int); ok {
		result.CurrentPage = currentPage
	} else {
		result.CurrentPage = 1
	}

	if totalPages, ok := pagination["TotalPages"].(int); ok {
		result.TotalPages = totalPages
	} else {
		result.TotalPages = 1
	}

	if totalItems, ok := pagination["TotalItems"].(int); ok {
		result.TotalItems = totalItems
	}

	if hasPrev, ok := pagination["HasPrev"].(bool); ok {
		result.HasPrev = hasPrev
	}

	if hasNext, ok := pagination["HasNext"].(bool); ok {
		result.HasNext = hasNext
	}

	if prevPage, ok := pagination["PrevPage"].(int); ok {
		result.PrevPage = prevPage
	} else {
		result.PrevPage = result.CurrentPage - 1
		if result.PrevPage < 1 {
			result.PrevPage = 1
		}
	}

	if nextPage, ok := pagination["NextPage"].(int); ok {
		result.NextPage = nextPage
	} else {
		result.NextPage = result.CurrentPage + 1
		if result.NextPage > result.TotalPages {
			result.NextPage = result.TotalPages
		}
	}

	// 获取当前URL
	result.PageUrl = getCurrentUrl(dataMap)

	return result, nil
}

// 获取当前URL
func getCurrentUrl(data map[string]interface{}) string {
	// 尝试从请求中获取URL
	if request, ok := data["Request"].(map[string]interface{}); ok {
		if url, ok := request["URL"].(string); ok {
			return url
		}
	}

	// 尝试从分类中获取URL
	if category, ok := data["Category"].(map[string]interface{}); ok {
		if typeid, ok := category["id"].(int64); ok {
			return fmt.Sprintf("/list/%d.html", typeid)
		}
	}

	// 默认URL
	return "/list/0.html"
}

// 获取分页URL
func getPageUrl(baseUrl string, page int) string {
	// 检查URL中是否已有page参数
	if strings.Contains(baseUrl, "page=") {
		re := regexp.MustCompile(`page=\d+`)
		return re.ReplaceAllString(baseUrl, fmt.Sprintf("page=%d", page))
	}

	// 检查URL中是否有其他参数
	if strings.Contains(baseUrl, "?") {
		return baseUrl + fmt.Sprintf("&page=%d", page)
	}

	// 检查URL是否以.html结尾
	if strings.HasSuffix(baseUrl, ".html") {
		// 替换为带页码的URL
		if page == 1 {
			return baseUrl
		}
		return strings.Replace(baseUrl, ".html", fmt.Sprintf("_%d.html", page), 1)
	}

	// 添加页码参数
	return baseUrl + fmt.Sprintf("?page=%d", page)
}
