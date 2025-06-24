package model

import (
	"fmt"
	"strconv"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// VisitLog 访问日志
type VisitLog struct {
	ID         int64     `json:"id"`
	MemberID   int64     `json:"memberid"`   // 会员ID
	ArticleID  int64     `json:"articleid"`  // 文章ID
	CategoryID int64     `json:"categoryid"` // 栏目ID
	IP         string    `json:"ip"`         // IP地址
	UserAgent  string    `json:"useragent"`  // User-Agent
	Referer    string    `json:"referer"`    // 来源
	URL        string    `json:"url"`        // 访问URL
	Device     string    `json:"device"`     // 设备
	Browser    string    `json:"browser"`    // 浏览器
	OS         string    `json:"os"`         // 操作系统
	Region     string    `json:"region"`     // 地区
	VisitTime  time.Time `json:"visittime"`  // 访问时间
}

// VisitLogModel 访问日志模型
type VisitLogModel struct {
	db *database.DB
}

// NewVisitLogModel 创建访问日志模型
func NewVisitLogModel(db *database.DB) *VisitLogModel {
	return &VisitLogModel{
		db: db,
	}
}

// Create 创建访问日志
func (m *VisitLogModel) Create(visitLog *VisitLog) (int64, error) {
	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("visit_log")+" (memberid, articleid, categoryid, ip, useragent, referer, url, device, browser, os, region, visittime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		visitLog.MemberID, visitLog.ArticleID, visitLog.CategoryID, visitLog.IP, visitLog.UserAgent, visitLog.Referer, visitLog.URL, visitLog.Device, visitLog.Browser, visitLog.OS, visitLog.Region, visitLog.VisitTime,
	)
	if err != nil {
		logger.Error("创建访问日志失败", "error", err)
		return 0, err
	}

	// 获取插入ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("获取插入ID失败", "error", err)
		return 0, err
	}

	return id, nil
}

// GetTodayVisits 获取今日访问量
func (m *VisitLogModel) GetTodayVisits() (int, error) {
	// 获取今日开始时间
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")
	qb.Where("visittime >= ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取今日访问量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetYesterdayVisits 获取昨日访问量
func (m *VisitLogModel) GetYesterdayVisits() (int, error) {
	// 获取昨日开始时间和结束时间
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")
	qb.Where("visittime >= ?", yesterdayStart)
	qb.Where("visittime < ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取昨日访问量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetWeekVisits 获取本周访问量
func (m *VisitLogModel) GetWeekVisits() (int, error) {
	// 获取本周开始时间
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")
	qb.Where("visittime >= ?", weekStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本周访问量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetMonthVisits 获取本月访问量
func (m *VisitLogModel) GetMonthVisits() (int, error) {
	// 获取本月开始时间
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")
	qb.Where("visittime >= ?", monthStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本月访问量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetTotalVisits 获取总访问量
func (m *VisitLogModel) GetTotalVisits() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取总访问量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetCategoryVisits 获取栏目访问量
func (m *VisitLogModel) GetCategoryVisits(categoryID int64) (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit_log")
	qb.Where("categoryid = ?", categoryID)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取栏目访问量失败", "categoryid", categoryID, "error", err)
		return 0, err
	}

	return count, nil
}

// GetVisitTrend 获取访问趋势
func (m *VisitLogModel) GetVisitTrend(days int) (map[string]int, error) {
	// 初始化结果
	trend := make(map[string]int)

	// 获取开始时间
	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// 初始化日期
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		trend[date.Format("2006-01-02")] = 0
	}

	// 构建查询
	sql := "SELECT DATE(visittime) as date, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " WHERE visittime >= ? GROUP BY DATE(visittime)"
	results, err := m.db.Query(sql, startDate)
	if err != nil {
		logger.Error("获取访问趋势失败", "error", err)
		return trend, err
	}

	// 处理结果
	for _, row := range results {
		date, ok := row["date"].(string)
		if !ok {
			logger.Error("日期类型转换失败", "date", row["date"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		trend[date] = int(count)
	}

	return trend, nil
}

// GetVisitSources 获取访问来源
func (m *VisitLogModel) GetVisitSources() (map[string]int, error) {
	// 初始化结果
	sources := make(map[string]int)

	// 构建查询
	sql := "SELECT referer, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY referer ORDER BY count DESC LIMIT 10"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问来源失败", "error", err)
		return sources, err
	}

	// 处理结果
	for _, row := range results {
		referer, ok := row["referer"].(string)
		if !ok {
			logger.Error("来源类型转换失败", "referer", row["referer"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		if referer == "" {
			referer = "直接访问"
		}
		sources[referer] = int(count)
	}

	return sources, nil
}

// GetVisitDevices 获取访问设备
func (m *VisitLogModel) GetVisitDevices() (map[string]int, error) {
	// 初始化结果
	devices := make(map[string]int)

	// 构建查询
	sql := "SELECT device, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY device ORDER BY count DESC"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问设备失败", "error", err)
		return devices, err
	}

	// 处理结果
	for _, row := range results {
		device, ok := row["device"].(string)
		if !ok {
			logger.Error("设备类型转换失败", "device", row["device"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		devices[device] = int(count)
	}

	return devices, nil
}

// GetVisitBrowsers 获取访问浏览器
func (m *VisitLogModel) GetVisitBrowsers() (map[string]int, error) {
	// 初始化结果
	browsers := make(map[string]int)

	// 构建查询
	sql := "SELECT browser, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY browser ORDER BY count DESC"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问浏览器失败", "error", err)
		return browsers, err
	}

	// 处理结果
	for _, row := range results {
		browser, ok := row["browser"].(string)
		if !ok {
			logger.Error("浏览器类型转换失败", "browser", row["browser"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		browsers[browser] = int(count)
	}

	return browsers, nil
}

// GetVisitOS 获取访问操作系统
func (m *VisitLogModel) GetVisitOS() (map[string]int, error) {
	// 初始化结果
	os := make(map[string]int)

	// 构建查询
	sql := "SELECT os, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY os ORDER BY count DESC"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问操作系统失败", "error", err)
		return os, err
	}

	// 处理结果
	for _, row := range results {
		osName, ok := row["os"].(string)
		if !ok {
			logger.Error("操作系统类型转换失败", "os", row["os"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		os[osName] = int(count)
	}

	return os, nil
}

// GetVisitRegions 获取访问地区
func (m *VisitLogModel) GetVisitRegions() (map[string]int, error) {
	// 初始化结果
	regions := make(map[string]int)

	// 构建查询
	sql := "SELECT region, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY region ORDER BY count DESC"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问地区失败", "error", err)
		return regions, err
	}

	// 处理结果
	for _, row := range results {
		region, ok := row["region"].(string)
		if !ok {
			logger.Error("地区类型转换失败", "region", row["region"])
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		regions[region] = int(count)
	}

	return regions, nil
}

// GetVisitHours 获取访问时段
func (m *VisitLogModel) GetVisitHours() (map[string]int, error) {
	// 初始化结果
	hours := make(map[string]int)
	for i := 0; i < 24; i++ {
		hours[fmt.Sprintf("%02d", i)] = 0
	}

	// 构建查询
	sql := "SELECT HOUR(visittime) as hour, COUNT(*) as count FROM " + m.db.TableName("visit_log") + " GROUP BY HOUR(visittime) ORDER BY hour"
	results, err := m.db.Query(sql)
	if err != nil {
		logger.Error("获取访问时段失败", "error", err)
		return hours, err
	}

	// 处理结果
	for _, row := range results {
		var hour int
		hourVal, ok := row["hour"]
		if !ok {
			logger.Error("小时字段不存在", "row", row)
			continue
		}

		switch h := hourVal.(type) {
		case int:
			hour = h
		case int64:
			hour = int(h)
		case float64:
			hour = int(h)
		case string:
			var err error
			hour, err = strconv.Atoi(h)
			if err != nil {
				logger.Error("小时类型转换失败", "hour", h, "error", err)
				continue
			}
		default:
			logger.Error("小时类型转换失败", "hour", hourVal)
			continue
		}

		count, ok := row["count"].(int64)
		if !ok {
			// 尝试其他类型转换
			if countStr, ok := row["count"].(string); ok {
				var countInt int
				_, err := fmt.Sscanf(countStr, "%d", &countInt)
				if err != nil {
					logger.Error("计数转换失败", "count", row["count"])
					continue
				}
				count = int64(countInt)
			} else {
				logger.Error("计数类型转换失败", "count", row["count"])
				continue
			}
		}

		hours[fmt.Sprintf("%02d", hour)] = int(count)
	}

	return hours, nil
}
