package tags

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"aq3cms/internal/interfaces"
)

// StatsTag 统计标签
type StatsTag struct {
	StatsService interfaces.StatsServiceInterface
}

// Parse 解析标签
func (t *StatsTag) Parse(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取类型
	statsType := params["type"]
	if statsType == "" {
		return "", fmt.Errorf("stats tag requires type parameter")
	}

	switch statsType {
	case "site":
		return t.getSiteStats(params, ctx)
	case "category":
		return t.getCategoryStats(params, ctx)
	case "member":
		return t.getMemberStats(params, ctx)
	case "visit":
		return t.getVisitStats(params, ctx)
	case "search":
		return t.getSearchStats(params, ctx)
	case "chart":
		return t.generateChart(params, ctx)
	default:
		return "", fmt.Errorf("unknown stats tag type: %s", statsType)
	}
}

// ParseBlock 解析块标签
func (t *StatsTag) ParseBlock(content string, params map[string]string, ctx map[string]interface{}) (string, error) {
	// 统计标签不支持块标签
	return content, nil
}

// 获取站点统计
func (t *StatsTag) getSiteStats(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取时间范围
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	endTime := time.Now().Unix()

	// 获取站点统计
	stats, err := t.StatsService.GetSiteStats(startTime, endTime)
	if err != nil {
		return "", err
	}

	// 获取要显示的字段
	field := params["field"]
	if field == "" {
		// 返回所有统计数据的JSON
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 返回指定字段
	if value, ok := stats[field]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "0", nil
}

// 获取栏目统计
func (t *StatsTag) getCategoryStats(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取时间范围
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	endTime := time.Now().Unix()

	// 获取栏目统计
	stats, err := t.StatsService.GetCategoryStats(startTime, endTime)
	if err != nil {
		return "", err
	}

	// 获取要显示的字段
	field := params["field"]
	if field == "" {
		// 返回所有统计数据的JSON
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 返回指定字段
	if value, ok := stats[field]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "0", nil
}

// 获取会员统计
func (t *StatsTag) getMemberStats(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取时间范围
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	endTime := time.Now().Unix()

	// 获取会员统计
	stats, err := t.StatsService.GetMemberStats(startTime, endTime)
	if err != nil {
		return "", err
	}

	// 获取要显示的字段
	field := params["field"]
	if field == "" {
		// 返回所有统计数据的JSON
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 返回指定字段
	if value, ok := stats[field]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "0", nil
}

// 获取访问统计
func (t *StatsTag) getVisitStats(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取天数
	days := 30
	if daysStr, ok := params["days"]; ok {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// 计算时间范围
	startTime := time.Now().AddDate(0, 0, -days).Unix()
	endTime := time.Now().Unix()

	// 获取访问统计
	stats, err := t.StatsService.GetVisitStats(startTime, endTime)
	if err != nil {
		return "", err
	}

	// 获取要显示的字段
	field := params["field"]
	if field == "" {
		// 返回所有统计数据的JSON
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 返回指定字段
	if value, ok := stats[field]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "0", nil
}

// 获取搜索统计
func (t *StatsTag) getSearchStats(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取天数
	days := 30
	if daysStr, ok := params["days"]; ok {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// 计算时间范围
	startTime := time.Now().AddDate(0, 0, -days).Unix()
	endTime := time.Now().Unix()

	// 获取搜索统计
	stats, err := t.StatsService.GetSearchStats(startTime, endTime)
	if err != nil {
		return "", err
	}

	// 获取要显示的字段
	field := params["field"]
	if field == "" {
		// 返回所有统计数据的JSON
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// 返回指定字段
	if value, ok := stats[field]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "0", nil
}

// 生成图表
func (t *StatsTag) generateChart(params map[string]string, ctx map[string]interface{}) (string, error) {
	// 获取图表类型
	chartType := params["chart_type"]
	if chartType == "" {
		return "", fmt.Errorf("chart requires chart_type parameter")
	}

	// 获取数据源
	dataSource := params["data_source"]
	if dataSource == "" {
		return "", fmt.Errorf("chart requires data_source parameter")
	}

	// 获取图表ID
	chartID := params["id"]
	if chartID == "" {
		chartID = "chart_" + chartType + "_" + dataSource
	}

	// 获取图表标题
	chartTitle := params["title"]
	if chartTitle == "" {
		chartTitle = "Chart"
	}

	// 获取图表宽度和高度
	width := "100%"
	if w, ok := params["width"]; ok {
		width = w
	}

	height := "400px"
	if h, ok := params["height"]; ok {
		height = h
	}

	// 获取数据
	var data interface{}
	var err error

	// 获取时间范围
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	endTime := time.Now().Unix()

	// 如果指定了天数，则使用指定的天数
	if daysStr, ok := params["days"]; ok {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			startTime = time.Now().AddDate(0, 0, -days).Unix()
		}
	}

	switch dataSource {
	case "site":
		stats, err := t.StatsService.GetSiteStats(startTime, endTime)
		if err != nil {
			return "", err
		}
		data = stats
	case "category":
		stats, err := t.StatsService.GetCategoryStats(startTime, endTime)
		if err != nil {
			return "", err
		}
		data = stats
	case "member":
		stats, err := t.StatsService.GetMemberStats(startTime, endTime)
		if err != nil {
			return "", err
		}
		data = stats
	case "visit":
		stats, err := t.StatsService.GetVisitStats(startTime, endTime)
		if err != nil {
			return "", err
		}
		data = stats
	case "search":
		stats, err := t.StatsService.GetSearchStats(startTime, endTime)
		if err != nil {
			return "", err
		}
		data = stats
	default:
		return "", fmt.Errorf("unknown data source: %s", dataSource)
	}

	// 转换数据为JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 构建图表HTML
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`<div id="%s" style="width: %s; height: %s;"></div>`, chartID, width, height))
	buf.WriteString("\n")
	buf.WriteString("<script>\n")
	buf.WriteString("document.addEventListener('DOMContentLoaded', function() {\n")
	buf.WriteString(fmt.Sprintf("  var chartData = %s;\n", dataJSON))
	buf.WriteString(fmt.Sprintf("  var chartType = '%s';\n", chartType))
	buf.WriteString(fmt.Sprintf("  var chartTitle = '%s';\n", chartTitle))
	buf.WriteString(fmt.Sprintf("  var chartID = '%s';\n", chartID))
	buf.WriteString("  initChart(chartID, chartType, chartTitle, chartData);\n")
	buf.WriteString("});\n")
	buf.WriteString("</script>\n")

	return buf.String(), nil
}
