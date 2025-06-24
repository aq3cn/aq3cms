package template

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TemplateFunctions 返回模板函数映射
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		// 字符串处理
		"html":       template.HTMLEscapeString,
		"url":        template.URLQueryEscaper,
		"js":         template.JSEscapeString,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"trim":       strings.TrimSpace,
		"substr":     substrFunc,
		"replace":    strings.Replace,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"split":      strings.Split,
		"join":       strings.Join,
		"stripTags":  stripTags,
		"truncate":   truncate,
		"nl2br":      nl2br,
		"urlEncode":  url.QueryEscape,
		"urlDecode":  urlDecode,
		"htmlDecode": htmlDecode,

		// 数字处理
		"add":      add,
		"sub":      sub,
		"mul":      mul,
		"div":      div,
		"mod":      mod,
		"round":    round,
		"floor":    floor,
		"ceil":     ceil,
		"max":      max,
		"min":      min,
		"abs":      abs,
		"formatNum": formatNum,

		// 日期处理
		"now":         now,
		"date":        dateFormat,
		"dateAdd":     dateAdd,
		"dateSub":     dateSub,
		"dateCompare": dateCompare,
		"timestamp":   timestamp,

		// 条件处理
		"eq":  eq,
		"ne":  ne,
		"lt":  lt,
		"le":  le,
		"gt":  gt,
		"ge":  ge,
		"and": and,
		"or":  or,
		"not": not,
		"if":  ifFunc,

		// 数组处理
		"first":   first,
		"last":    last,
		"slice":   slice,
		"inArray": inArray,
		"length":  length,

		// 其他
		"default": defaultFunc,
		"raw":     raw,
	}
}

// 字符串处理函数

// substrFunc 截取字符串
func substrFunc(s string, start, length int) string {
	runes := []rune(s)
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

// stripTags 去除HTML标签
func stripTags(s string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(s, "")
}

// truncate 截断字符串并添加省略号
func truncate(s string, length int) string {
	runes := []rune(s)
	if len(runes) <= length {
		return s
	}
	return string(runes[:length]) + "..."
}

// nl2br 将换行符转换为<br>标签
func nl2br(s string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(s), "\n", "<br>", -1))
}

// urlDecode URL解码
func urlDecode(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		return s
	}
	return decoded
}

// htmlDecode HTML解码
func htmlDecode(s string) template.HTML {
	return template.HTML(s)
}

// 数字处理函数

// add 加法
func add(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av + bv
}

// sub 减法
func sub(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av - bv
}

// mul 乘法
func mul(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av * bv
}

// div 除法
func div(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	if bv == 0 {
		return 0
	}
	return av / bv
}

// mod 取模
func mod(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	if bv == 0 {
		return 0
	}
	return math.Mod(av, bv)
}

// round 四舍五入
func round(a interface{}, precision int) float64 {
	av := convertToFloat64(a)
	p := math.Pow10(precision)
	return math.Round(av*p) / p
}

// floor 向下取整
func floor(a interface{}) float64 {
	av := convertToFloat64(a)
	return math.Floor(av)
}

// ceil 向上取整
func ceil(a interface{}) float64 {
	av := convertToFloat64(a)
	return math.Ceil(av)
}

// max 取最大值
func max(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return math.Max(av, bv)
}

// min 取最小值
func min(a, b interface{}) interface{} {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return math.Min(av, bv)
}

// abs 取绝对值
func abs(a interface{}) interface{} {
	av := convertToFloat64(a)
	return math.Abs(av)
}

// formatNum 格式化数字
func formatNum(a interface{}, format string) string {
	av := convertToFloat64(a)
	return fmt.Sprintf(format, av)
}

// 日期处理函数

// now 获取当前时间
func now() time.Time {
	return time.Now()
}

// dateFormat 格式化日期
func dateFormat(t interface{}, format string) string {
	var tt time.Time

	switch v := t.(type) {
	case time.Time:
		tt = v
	case int64:
		tt = time.Unix(v, 0)
	case int:
		tt = time.Unix(int64(v), 0)
	case string:
		var err error
		tt, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return ""
		}
	default:
		return ""
	}

	// 替换常用格式
	format = strings.Replace(format, "Y", "2006", -1)
	format = strings.Replace(format, "m", "01", -1)
	format = strings.Replace(format, "d", "02", -1)
	format = strings.Replace(format, "H", "15", -1)
	format = strings.Replace(format, "i", "04", -1)
	format = strings.Replace(format, "s", "05", -1)

	return tt.Format(format)
}

// dateAdd 日期加法
func dateAdd(t interface{}, d string) time.Time {
	var tt time.Time

	switch v := t.(type) {
	case time.Time:
		tt = v
	case int64:
		tt = time.Unix(v, 0)
	case int:
		tt = time.Unix(int64(v), 0)
	case string:
		var err error
		tt, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return time.Now()
		}
	default:
		return time.Now()
	}

	// 解析时间间隔
	duration, err := time.ParseDuration(d)
	if err != nil {
		return tt
	}

	return tt.Add(duration)
}

// dateSub 日期减法
func dateSub(t interface{}, d string) time.Time {
	var tt time.Time

	switch v := t.(type) {
	case time.Time:
		tt = v
	case int64:
		tt = time.Unix(v, 0)
	case int:
		tt = time.Unix(int64(v), 0)
	case string:
		var err error
		tt, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return time.Now()
		}
	default:
		return time.Now()
	}

	// 解析时间间隔
	duration, err := time.ParseDuration(d)
	if err != nil {
		return tt
	}

	return tt.Add(-duration)
}

// dateCompare 日期比较
func dateCompare(t1, t2 interface{}) int {
	var tt1, tt2 time.Time

	switch v := t1.(type) {
	case time.Time:
		tt1 = v
	case int64:
		tt1 = time.Unix(v, 0)
	case int:
		tt1 = time.Unix(int64(v), 0)
	case string:
		var err error
		tt1, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return 0
		}
	default:
		return 0
	}

	switch v := t2.(type) {
	case time.Time:
		tt2 = v
	case int64:
		tt2 = time.Unix(v, 0)
	case int:
		tt2 = time.Unix(int64(v), 0)
	case string:
		var err error
		tt2, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return 0
		}
	default:
		return 0
	}

	if tt1.Before(tt2) {
		return -1
	} else if tt1.After(tt2) {
		return 1
	}
	return 0
}

// timestamp 获取时间戳
func timestamp(t interface{}) int64 {
	var tt time.Time

	switch v := t.(type) {
	case time.Time:
		tt = v
	case string:
		var err error
		tt, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return 0
		}
	default:
		tt = time.Now()
	}

	return tt.Unix()
}

// 条件处理函数

// eq 等于
func eq(a, b interface{}) bool {
	return a == b
}

// ne 不等于
func ne(a, b interface{}) bool {
	return a != b
}

// lt 小于
func lt(a, b interface{}) bool {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av < bv
}

// le 小于等于
func le(a, b interface{}) bool {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av <= bv
}

// gt 大于
func gt(a, b interface{}) bool {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av > bv
}

// ge 大于等于
func ge(a, b interface{}) bool {
	av := convertToFloat64(a)
	bv := convertToFloat64(b)
	return av >= bv
}

// and 逻辑与
func and(a, b bool) bool {
	return a && b
}

// or 逻辑或
func or(a, b bool) bool {
	return a || b
}

// not 逻辑非
func not(a bool) bool {
	return !a
}

// ifFunc 条件判断
func ifFunc(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// 数组处理函数

// first 获取数组第一个元素
func first(arr interface{}) interface{} {
	value := reflect.ValueOf(arr)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil
	}
	if value.Len() == 0 {
		return nil
	}
	return value.Index(0).Interface()
}

// last 获取数组最后一个元素
func last(arr interface{}) interface{} {
	value := reflect.ValueOf(arr)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil
	}
	if value.Len() == 0 {
		return nil
	}
	return value.Index(value.Len() - 1).Interface()
}

// slice 截取数组
func slice(arr interface{}, start, end int) interface{} {
	value := reflect.ValueOf(arr)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil
	}

	size := value.Len()
	if start < 0 {
		start = 0
	}
	if end > size {
		end = size
	}
	if start > end {
		return nil
	}

	return value.Slice(start, end).Interface()
}

// inArray 检查元素是否在数组中
func inArray(needle interface{}, haystack interface{}) bool {
	value := reflect.ValueOf(haystack)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(needle, value.Index(i).Interface()) {
			return true
		}
	}

	return false
}

// length 获取数组长度
func length(arr interface{}) int {
	value := reflect.ValueOf(arr)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array && value.Kind() != reflect.String && value.Kind() != reflect.Map {
		return 0
	}

	return value.Len()
}

// 其他函数

// defaultFunc 默认值
func defaultFunc(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return defaultValue
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return defaultValue
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return defaultValue
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return defaultValue
		}
	case reflect.Bool:
		if !v.Bool() {
			return defaultValue
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return defaultValue
		}
	}

	return value
}

// raw 原始HTML
func raw(s string) template.HTML {
	return template.HTML(s)
}

// 辅助函数

// convertToFloat64 将值转换为float64
func convertToFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}
