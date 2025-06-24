package model

import (
	"fmt"
	"strconv"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Visit 访问记录
type Visit struct {
	ID         int64     `json:"id"`
	IP         string    `json:"ip"`         // IP地址
	UserAgent  string    `json:"useragent"`  // 用户代理
	URL        string    `json:"url"`        // 访问URL
	Referer    string    `json:"referer"`    // 来源URL
	MemberID   int64     `json:"memberid"`   // 会员ID，0表示游客
	SessionID  string    `json:"sessionid"`  // 会话ID
	CreateTime time.Time `json:"createtime"` // 创建时间
}

// SearchKeyword 搜索关键词
type SearchKeyword struct {
	ID         int64     `json:"id"`
	Keyword    string    `json:"keyword"`    // 关键词
	Count      int       `json:"count"`      // 搜索次数
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// VisitStats 访问统计
type VisitStats struct {
	Date      string `json:"date"`      // 日期
	PV        int    `json:"pv"`        // 页面浏览量
	UV        int    `json:"uv"`        // 独立访客数
	IP        int    `json:"ip"`        // IP数
	NewVisit  int    `json:"newvisit"`  // 新访客数
	AvgTime   int    `json:"avgtime"`   // 平均访问时长（秒）
	BounceRate int    `json:"bouncerate"` // 跳出率（百分比）
}

// VisitModel 访问记录模型
type VisitModel struct {
	db *database.DB
}

// NewVisitModel 创建访问记录模型
func NewVisitModel(db *database.DB) *VisitModel {
	return &VisitModel{
		db: db,
	}
}

// Create 创建访问记录
func (m *VisitModel) Create(visit *Visit) (int64, error) {
	// 设置创建时间
	visit.CreateTime = time.Now()

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("visit")+" (ip, useragent, url, referer, memberid, sessionid, createtime) VALUES (?, ?, ?, ?, ?, ?, ?)",
		visit.IP, visit.UserAgent, visit.URL, visit.Referer, visit.MemberID, visit.SessionID, visit.CreateTime,
	)
	if err != nil {
		logger.Error("创建访问记录失败", "error", err)
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

// GetByID 根据ID获取访问记录
func (m *VisitModel) GetByID(id int64) (*Visit, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取访问记录失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("visit not found: %d", id)
	}

	// 转换为访问记录
	visit := &Visit{}
	visit.ID, _ = result["id"].(int64)
	visit.IP, _ = result["ip"].(string)
	visit.UserAgent, _ = result["useragent"].(string)
	visit.URL, _ = result["url"].(string)
	visit.Referer, _ = result["referer"].(string)
	visit.MemberID, _ = result["memberid"].(int64)
	visit.SessionID, _ = result["sessionid"].(string)
	visit.CreateTime, _ = result["createtime"].(time.Time)

	return visit, nil
}

// GetBySessionID 根据会话ID获取访问记录
func (m *VisitModel) GetBySessionID(sessionID string, limit int) ([]*Visit, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit")
	qb.Where("sessionid = ?", sessionID)
	qb.OrderBy("id DESC")
	if limit > 0 {
		qb.Limit(limit)
	}

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取访问记录失败", "sessionid", sessionID, "error", err)
		return nil, err
	}

	// 转换为访问记录列表
	visits := make([]*Visit, 0, len(results))
	for _, result := range results {
		visit := &Visit{}
		visit.ID, _ = result["id"].(int64)
		visit.IP, _ = result["ip"].(string)
		visit.UserAgent, _ = result["useragent"].(string)
		visit.URL, _ = result["url"].(string)
		visit.Referer, _ = result["referer"].(string)
		visit.MemberID, _ = result["memberid"].(int64)
		visit.SessionID, _ = result["sessionid"].(string)
		visit.CreateTime, _ = result["createtime"].(time.Time)
		visits = append(visits, visit)
	}

	return visits, nil
}

// GetByMemberID 根据会员ID获取访问记录
func (m *VisitModel) GetByMemberID(memberID int64, limit int) ([]*Visit, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit")
	qb.Where("memberid = ?", memberID)
	qb.OrderBy("id DESC")
	if limit > 0 {
		qb.Limit(limit)
	}

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取访问记录失败", "memberid", memberID, "error", err)
		return nil, err
	}

	// 转换为访问记录列表
	visits := make([]*Visit, 0, len(results))
	for _, result := range results {
		visit := &Visit{}
		visit.ID, _ = result["id"].(int64)
		visit.IP, _ = result["ip"].(string)
		visit.UserAgent, _ = result["useragent"].(string)
		visit.URL, _ = result["url"].(string)
		visit.Referer, _ = result["referer"].(string)
		visit.MemberID, _ = result["memberid"].(int64)
		visit.SessionID, _ = result["sessionid"].(string)
		visit.CreateTime, _ = result["createtime"].(time.Time)
		visits = append(visits, visit)
	}

	return visits, nil
}

// GetByDate 根据日期获取访问记录
func (m *VisitModel) GetByDate(date string, limit int) ([]*Visit, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit")
	qb.Where("DATE(createtime) = ?", date)
	qb.OrderBy("id DESC")
	if limit > 0 {
		qb.Limit(limit)
	}

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取访问记录失败", "date", date, "error", err)
		return nil, err
	}

	// 转换为访问记录列表
	visits := make([]*Visit, 0, len(results))
	for _, result := range results {
		visit := &Visit{}
		visit.ID, _ = result["id"].(int64)
		visit.IP, _ = result["ip"].(string)
		visit.UserAgent, _ = result["useragent"].(string)
		visit.URL, _ = result["url"].(string)
		visit.Referer, _ = result["referer"].(string)
		visit.MemberID, _ = result["memberid"].(int64)
		visit.SessionID, _ = result["sessionid"].(string)
		visit.CreateTime, _ = result["createtime"].(time.Time)
		visits = append(visits, visit)
	}

	return visits, nil
}

// GetByDateRange 根据日期范围获取访问记录
func (m *VisitModel) GetByDateRange(startDate, endDate string, limit int) ([]*Visit, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "visit")
	qb.Where("DATE(createtime) >= ?", startDate)
	qb.Where("DATE(createtime) <= ?", endDate)
	qb.OrderBy("id DESC")
	if limit > 0 {
		qb.Limit(limit)
	}

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取访问记录失败", "startDate", startDate, "endDate", endDate, "error", err)
		return nil, err
	}

	// 转换为访问记录列表
	visits := make([]*Visit, 0, len(results))
	for _, result := range results {
		visit := &Visit{}
		visit.ID, _ = result["id"].(int64)
		visit.IP, _ = result["ip"].(string)
		visit.UserAgent, _ = result["useragent"].(string)
		visit.URL, _ = result["url"].(string)
		visit.Referer, _ = result["referer"].(string)
		visit.MemberID, _ = result["memberid"].(int64)
		visit.SessionID, _ = result["sessionid"].(string)
		visit.CreateTime, _ = result["createtime"].(time.Time)
		visits = append(visits, visit)
	}

	return visits, nil
}

// GetStats 获取统计数据
func (m *VisitModel) GetStats(date string) (*VisitStats, error) {
	// 获取PV
	pvQuery := "SELECT COUNT(*) AS pv FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ?"
	pvResult, err := m.db.Query(pvQuery, date)
	if err != nil {
		logger.Error("获取PV失败", "error", err)
		return nil, err
	}

	var pv int
	if len(pvResult) > 0 {
		if pvVal, ok := pvResult[0]["pv"]; ok {
			switch v := pvVal.(type) {
			case int:
				pv = v
			case int64:
				pv = int(v)
			case string:
				pv, _ = strconv.Atoi(v)
			}
		}
	}

	// 获取UV
	uvQuery := "SELECT COUNT(DISTINCT sessionid) AS uv FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ?"
	uvResult, err := m.db.Query(uvQuery, date)
	if err != nil {
		logger.Error("获取UV失败", "error", err)
		return nil, err
	}

	var uv int
	if len(uvResult) > 0 {
		if uvVal, ok := uvResult[0]["uv"]; ok {
			switch v := uvVal.(type) {
			case int:
				uv = v
			case int64:
				uv = int(v)
			case string:
				uv, _ = strconv.Atoi(v)
			}
		}
	}

	// 获取IP
	ipQuery := "SELECT COUNT(DISTINCT ip) AS ip FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ?"
	ipResult, err := m.db.Query(ipQuery, date)
	if err != nil {
		logger.Error("获取IP失败", "error", err)
		return nil, err
	}

	var ip int
	if len(ipResult) > 0 {
		if ipVal, ok := ipResult[0]["ip"]; ok {
			switch v := ipVal.(type) {
			case int:
				ip = v
			case int64:
				ip = int(v)
			case string:
				ip, _ = strconv.Atoi(v)
			}
		}
	}

	// 获取新访客数
	newVisitQuery := "SELECT COUNT(DISTINCT sessionid) AS newvisit FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ? AND sessionid NOT IN (SELECT DISTINCT sessionid FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) < ?)"
	newVisitResult, err := m.db.Query(newVisitQuery, date, date)
	if err != nil {
		logger.Error("获取新访客数失败", "error", err)
		return nil, err
	}

	var newVisit int
	if len(newVisitResult) > 0 {
		if newVisitVal, ok := newVisitResult[0]["newvisit"]; ok {
			switch v := newVisitVal.(type) {
			case int:
				newVisit = v
			case int64:
				newVisit = int(v)
			case string:
				newVisit, _ = strconv.Atoi(v)
			}
		}
	}

	// 获取平均访问时长
	avgTimeQuery := "SELECT AVG(TIMESTAMPDIFF(SECOND, MIN(createtime), MAX(createtime))) AS avgtime FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ? GROUP BY sessionid"
	avgTimeResult, err := m.db.Query(avgTimeQuery, date)
	if err != nil {
		logger.Error("获取平均访问时长失败", "error", err)
		return nil, err
	}

	var avgTime int
	if len(avgTimeResult) > 0 {
		if avgTimeVal, ok := avgTimeResult[0]["avgtime"]; ok {
			switch v := avgTimeVal.(type) {
			case int:
				avgTime = v
			case int64:
				avgTime = int(v)
			case float64:
				avgTime = int(v)
			case string:
				avgTime, _ = strconv.Atoi(v)
			}
		}
	}

	// 获取跳出率
	bounceRateQuery := "SELECT COUNT(*) * 100 / (SELECT COUNT(DISTINCT sessionid) FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ?) AS bouncerate FROM (SELECT sessionid, COUNT(*) AS count FROM " + m.db.TableName("visit") + " WHERE DATE(createtime) = ? GROUP BY sessionid HAVING count = 1) AS t"
	bounceRateResult, err := m.db.Query(bounceRateQuery, date, date)
	if err != nil {
		logger.Error("获取跳出率失败", "error", err)
		return nil, err
	}

	var bounceRate int
	if len(bounceRateResult) > 0 {
		if bounceRateVal, ok := bounceRateResult[0]["bouncerate"]; ok {
			switch v := bounceRateVal.(type) {
			case int:
				bounceRate = v
			case int64:
				bounceRate = int(v)
			case float64:
				bounceRate = int(v)
			case string:
				bounceRate, _ = strconv.Atoi(v)
			}
		}
	}

	// 返回统计数据
	return &VisitStats{
		Date:       date,
		PV:         pv,
		UV:         uv,
		IP:         ip,
		NewVisit:   newVisit,
		AvgTime:    avgTime,
		BounceRate: bounceRate,
	}, nil
}

// GetStatsRange 获取统计数据范围
func (m *VisitModel) GetStatsRange(startDate, endDate string) ([]*VisitStats, error) {
	// 获取日期范围
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		logger.Error("解析开始日期失败", "error", err)
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		logger.Error("解析结束日期失败", "error", err)
		return nil, err
	}

	// 获取统计数据
	stats := make([]*VisitStats, 0)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		stat, err := m.GetStats(date)
		if err != nil {
			logger.Error("获取统计数据失败", "date", date, "error", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// SearchKeywordModel 搜索关键词模型
type SearchKeywordModel struct {
	db *database.DB
}

// NewSearchKeywordModel 创建搜索关键词模型
func NewSearchKeywordModel(db *database.DB) *SearchKeywordModel {
	return &SearchKeywordModel{
		db: db,
	}
}

// GetByID 根据ID获取搜索关键词
func (m *SearchKeywordModel) GetByID(id int64) (*SearchKeyword, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_keyword")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取搜索关键词失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("search keyword not found: %d", id)
	}

	// 转换为搜索关键词
	keyword := &SearchKeyword{}
	keyword.ID, _ = result["id"].(int64)
	keyword.Keyword, _ = result["keyword"].(string)
	keyword.Count, _ = result["count"].(int)
	keyword.CreateTime, _ = result["createtime"].(time.Time)
	keyword.UpdateTime, _ = result["updatetime"].(time.Time)

	return keyword, nil
}

// GetByKeyword 根据关键词获取搜索关键词
func (m *SearchKeywordModel) GetByKeyword(keyword string) (*SearchKeyword, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_keyword")
	qb.Where("keyword = ?", keyword)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取搜索关键词失败", "keyword", keyword, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("search keyword not found: %s", keyword)
	}

	// 转换为搜索关键词
	searchKeyword := &SearchKeyword{}
	searchKeyword.ID, _ = result["id"].(int64)
	searchKeyword.Keyword, _ = result["keyword"].(string)
	searchKeyword.Count, _ = result["count"].(int)
	searchKeyword.CreateTime, _ = result["createtime"].(time.Time)
	searchKeyword.UpdateTime, _ = result["updatetime"].(time.Time)

	return searchKeyword, nil
}

// GetAll 获取所有搜索关键词
func (m *SearchKeywordModel) GetAll(page, pageSize int) ([]*SearchKeyword, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_keyword")
	qb.OrderBy("count DESC, id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有搜索关键词失败", "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取搜索关键词总数失败", "error", err)
		return nil, 0, err
	}

	// 转换为搜索关键词列表
	keywords := make([]*SearchKeyword, 0, len(results))
	for _, result := range results {
		keyword := &SearchKeyword{}
		keyword.ID, _ = result["id"].(int64)
		keyword.Keyword, _ = result["keyword"].(string)
		keyword.Count, _ = result["count"].(int)
		keyword.CreateTime, _ = result["createtime"].(time.Time)
		keyword.UpdateTime, _ = result["updatetime"].(time.Time)
		keywords = append(keywords, keyword)
	}

	return keywords, total, nil
}

// GetTop 获取热门搜索关键词
func (m *SearchKeywordModel) GetTop(limit int) ([]*SearchKeyword, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "search_keyword")
	qb.OrderBy("count DESC, id DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取热门搜索关键词失败", "error", err)
		return nil, err
	}

	// 转换为搜索关键词列表
	keywords := make([]*SearchKeyword, 0, len(results))
	for _, result := range results {
		keyword := &SearchKeyword{}
		keyword.ID, _ = result["id"].(int64)
		keyword.Keyword, _ = result["keyword"].(string)
		keyword.Count, _ = result["count"].(int)
		keyword.CreateTime, _ = result["createtime"].(time.Time)
		keyword.UpdateTime, _ = result["updatetime"].(time.Time)
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

// Create 创建搜索关键词
func (m *SearchKeywordModel) Create(keyword *SearchKeyword) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	keyword.CreateTime = now
	keyword.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("search_keyword")+" (keyword, count, createtime, updatetime) VALUES (?, ?, ?, ?)",
		keyword.Keyword, keyword.Count, keyword.CreateTime, keyword.UpdateTime,
	)
	if err != nil {
		logger.Error("创建搜索关键词失败", "error", err)
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

// Update 更新搜索关键词
func (m *SearchKeywordModel) Update(keyword *SearchKeyword) error {
	// 设置更新时间
	keyword.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("search_keyword")+" SET keyword = ?, count = ?, updatetime = ? WHERE id = ?",
		keyword.Keyword, keyword.Count, keyword.UpdateTime, keyword.ID,
	)
	if err != nil {
		logger.Error("更新搜索关键词失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除搜索关键词
func (m *SearchKeywordModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("search_keyword")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除搜索关键词失败", "error", err)
		return err
	}

	return nil
}

// IncrementCount 增加搜索次数
func (m *SearchKeywordModel) IncrementCount(keyword string) error {
	// 检查关键词是否存在
	searchKeyword, err := m.GetByKeyword(keyword)
	if err != nil {
		// 关键词不存在，创建
		newKeyword := &SearchKeyword{
			Keyword: keyword,
			Count:   1,
		}
		_, err = m.Create(newKeyword)
		return err
	}

	// 关键词存在，更新
	searchKeyword.Count++
	return m.Update(searchKeyword)
}
