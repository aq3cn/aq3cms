package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SearchLog 搜索日志
type SearchLog struct {
	ID          int64     `json:"id"`
	Keyword     string    `json:"keyword"`     // 搜索关键词
	IP          string    `json:"ip"`          // IP地址
	ResultCount int       `json:"resultcount"` // 结果数量
	SearchTime  time.Time `json:"searchtime"`  // 搜索时间
}

// SearchLogModel 搜索日志模型
type SearchLogModel struct {
	db *database.DB
}

// NewSearchLogModel 创建搜索日志模型
func NewSearchLogModel(db *database.DB) *SearchLogModel {
	return &SearchLogModel{
		db: db,
	}
}

// Create 创建搜索日志
func (m *SearchLogModel) Create(searchLog *SearchLog) (int64, error) {
	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("search_log")+" (keyword, ip, resultcount, searchtime) VALUES (?, ?, ?, ?)",
		searchLog.Keyword, searchLog.IP, searchLog.ResultCount, searchLog.SearchTime,
	)
	if err != nil {
		logger.Error("创建搜索日志失败", "error", err)
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

// GetHotSearches 获取热门搜索
func (m *SearchLogModel) GetHotSearches(limit int) ([]map[string]interface{}, error) {
	// 构建查询
	sql := "SELECT keyword, COUNT(*) as count FROM " + m.db.TableName("search_log") + " GROUP BY keyword ORDER BY count DESC LIMIT ?"
	results, err := m.db.Query(sql, limit)
	if err != nil {
		logger.Error("获取热门搜索失败", "error", err)
		return nil, err
	}

	// 处理结果
	hotSearches := make([]map[string]interface{}, 0, len(results))
	for _, row := range results {
		hotSearches = append(hotSearches, map[string]interface{}{
			"Keyword": row["keyword"],
			"Count":   row["count"],
		})
	}

	return hotSearches, nil
}

// GetNoResultSearches 获取无结果搜索
func (m *SearchLogModel) GetNoResultSearches(limit int) ([]map[string]interface{}, error) {
	// 构建查询
	sql := "SELECT keyword, COUNT(*) as count FROM " + m.db.TableName("search_log") + " WHERE resultcount = 0 GROUP BY keyword ORDER BY count DESC LIMIT ?"
	results, err := m.db.Query(sql, limit)
	if err != nil {
		logger.Error("获取无结果搜索失败", "error", err)
		return nil, err
	}

	// 处理结果
	noResultSearches := make([]map[string]interface{}, 0, len(results))
	for _, row := range results {
		noResultSearches = append(noResultSearches, map[string]interface{}{
			"Keyword": row["keyword"],
			"Count":   row["count"],
		})
	}

	return noResultSearches, nil
}

// GetSearchTrend 获取搜索趋势
func (m *SearchLogModel) GetSearchTrend(days int) (map[string]int, error) {
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
	sql := "SELECT DATE(searchtime) as date, COUNT(*) as count FROM " + m.db.TableName("search_log") + " WHERE searchtime >= ? GROUP BY DATE(searchtime)"
	results, err := m.db.Query(sql, startDate)
	if err != nil {
		logger.Error("获取搜索趋势失败", "error", err)
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

// GetTodaySearches 获取今日搜索量
func (m *SearchLogModel) GetTodaySearches() (int, error) {
	// 获取今日开始时间
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_log")
	qb.Where("searchtime >= ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取今日搜索量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetYesterdaySearches 获取昨日搜索量
func (m *SearchLogModel) GetYesterdaySearches() (int, error) {
	// 获取昨日开始时间和结束时间
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_log")
	qb.Where("searchtime >= ?", yesterdayStart)
	qb.Where("searchtime < ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取昨日搜索量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetWeekSearches 获取本周搜索量
func (m *SearchLogModel) GetWeekSearches() (int, error) {
	// 获取本周开始时间
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_log")
	qb.Where("searchtime >= ?", weekStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本周搜索量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetMonthSearches 获取本月搜索量
func (m *SearchLogModel) GetMonthSearches() (int, error) {
	// 获取本月开始时间
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_log")
	qb.Where("searchtime >= ?", monthStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本月搜索量失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetTotalSearches 获取总搜索量
func (m *SearchLogModel) GetTotalSearches() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_log")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取总搜索量失败", "error", err)
		return 0, err
	}

	return count, nil
}
