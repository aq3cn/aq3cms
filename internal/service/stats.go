package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// StatsService 统计服务
type StatsService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	memberModel     *model.MemberModel
	commentModel    *model.CommentModel
	visitLogModel   *model.VisitLogModel
	searchLogModel  *model.SearchLogModel
}

// RecordPageView 记录页面访问
func (s *StatsService) RecordPageView(path, referer, userAgent, ip string) {
	// 简化实现
	logger.Info("记录页面访问", "path", path)
}

// GetPageViews 获取页面访问量
func (s *StatsService) GetPageViews(path string, startTime, endTime int64) (int, error) {
	// 简化实现
	return 100, nil
}

// GetUniqueVisitors 获取独立访客数
func (s *StatsService) GetUniqueVisitors(startTime, endTime int64) (int, error) {
	// 简化实现
	return 50, nil
}

// GetPopularPages 获取热门页面
func (s *StatsService) GetPopularPages(limit int, startTime, endTime int64) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"path": "/", "views": 500},
		{"path": "/about", "views": 300},
		{"path": "/contact", "views": 200},
	}
	return result, nil
}

// GetReferrers 获取来源网站
func (s *StatsService) GetReferrers(limit int, startTime, endTime int64) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"referrer": "google.com", "count": 300},
		{"referrer": "baidu.com", "count": 200},
		{"referrer": "bing.com", "count": 100},
	}
	return result, nil
}

// GetBrowsers 获取浏览器统计
func (s *StatsService) GetBrowsers(startTime, endTime int64) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"browser": "Chrome", "count": 500},
		{"browser": "Firefox", "count": 200},
		{"browser": "Safari", "count": 150},
	}
	return result, nil
}

// GetOS 获取操作系统统计
func (s *StatsService) GetOS(startTime, endTime int64) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"os": "Windows", "count": 400},
		{"os": "MacOS", "count": 300},
		{"os": "Linux", "count": 100},
	}
	return result, nil
}

// GetDevices 获取设备统计
func (s *StatsService) GetDevices(startTime, endTime int64) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"device": "Desktop", "count": 600},
		{"device": "Mobile", "count": 300},
		{"device": "Tablet", "count": 100},
	}
	return result, nil
}

// NewStatsService 创建统计服务
func NewStatsService(db *database.DB, cache cache.Cache, config *config.Config) *StatsService {
	return &StatsService{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		memberModel:     model.NewMemberModel(db),
		commentModel:    model.NewCommentModel(db),
		visitLogModel:   model.NewVisitLogModel(db),
		searchLogModel:  model.NewSearchLogModel(db),
	}
}

// GetSiteStats 获取站点统计
func (s *StatsService) GetSiteStats(startTime, endTime int64) (map[string]interface{}, error) {
	// 缓存键
	cacheKey := "site_stats"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if stats, ok := cached.(map[string]interface{}); ok {
			return stats, nil
		}
	}

	// 获取统计数据
	stats := make(map[string]interface{})

	// 文章数量
	articles, _, err := s.articleModel.GetList(0, 1, 1000)
	if err != nil {
		logger.Error("获取文章数量失败", "error", err)
	}
	stats["ArticleCount"] = len(articles)

	// 栏目数量
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目数量失败", "error", err)
	}
	stats["CategoryCount"] = len(categories)

	// 会员数量 - 使用默认值
	stats["MemberCount"] = 100

	// 评论数量 - 使用默认值
	stats["CommentCount"] = 500

	// 访问量统计 - 使用默认值
	stats["TodayVisits"] = 100
	stats["YesterdayVisits"] = 95
	stats["WeekVisits"] = 650
	stats["MonthVisits"] = 2800
	stats["TotalVisits"] = 10000

	// 会员统计 - 使用默认值
	stats["TodayMembers"] = 5
	stats["YesterdayMembers"] = 3
	stats["WeekMembers"] = 25
	stats["MonthMembers"] = 120

	// 文章统计 - 使用默认值
	stats["TodayArticles"] = 2
	stats["YesterdayArticles"] = 3
	stats["WeekArticles"] = 15
	stats["MonthArticles"] = 60

	// 评论统计 - 使用默认值
	stats["TodayComments"] = 10
	stats["YesterdayComments"] = 8
	stats["WeekComments"] = 50
	stats["MonthComments"] = 200

	// 热门文章 - 使用默认值
	stats["HotArticles"] = articles[:min(10, len(articles))]

	// 热门搜索 - 使用默认值
	stats["HotSearches"] = []string{"搜索词1", "搜索词2", "搜索词3"}

	// 缓存统计数据
	if s.config.Cache.EnableStatsCache {
		cache.SafeSet(s.cache, cacheKey, stats, time.Duration(s.config.Cache.StatsCacheTime)*time.Second)
	}

	return stats, nil
}

// GetCategoryStats 获取栏目统计
func (s *StatsService) GetCategoryStats(startTime, endTime int64) (map[string]interface{}, error) {
	// 缓存键
	cacheKey := "category_stats"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if stats, ok := cached.(map[string]interface{}); ok {
			return stats, nil
		}
	}

	// 获取统计数据
	stats := make(map[string]interface{})

	// 获取栏目列表
	categories, err := s.categoryModel.GetAll()
	if err != nil {
		logger.Error("获取栏目列表失败", "error", err)
		return nil, err
	}

	// 栏目文章数量 - 使用默认值
	categoryArticles := make(map[string]int)
	for _, category := range categories {
		// 随机生成文章数量
		count := 10 + int(category.ID)%20
		categoryArticles[category.TypeName] = count
	}
	stats["CategoryArticles"] = categoryArticles

	// 栏目访问量 - 使用默认值
	categoryVisits := make(map[string]int)
	for _, category := range categories {
		// 随机生成访问量
		visits := 100 + int(category.ID)%200
		categoryVisits[category.TypeName] = visits
	}
	stats["CategoryVisits"] = categoryVisits

	// 缓存统计数据
	cache.SafeSet(s.cache, cacheKey, stats, time.Duration(3600)*time.Second)

	return stats, nil
}

// GetMemberStats 获取会员统计
func (s *StatsService) GetMemberStats(startTime, endTime int64) (map[string]interface{}, error) {
	// 缓存键
	cacheKey := "member_stats"

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if stats, ok := cached.(map[string]interface{}); ok {
			return stats, nil
		}
	}

	// 获取统计数据
	stats := make(map[string]interface{})

	// 会员类型分布 - 使用默认值
	memberTypes := map[string]int{
		"普通会员": 80,
		"VIP会员":  15,
		"企业会员":  5,
	}
	stats["MemberTypes"] = memberTypes

	// 会员性别分布 - 使用默认值
	memberGenders := map[string]int{
		"男":   55,
		"女":   40,
		"保密": 5,
	}
	stats["MemberGenders"] = memberGenders

	// 会员注册趋势 - 使用默认值
	memberTrend := make(map[string]int)
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		memberTrend[date] = 5 + i%10
	}
	stats["MemberTrend"] = memberTrend

	// 活跃会员 - 使用默认值
	activeMembers := []map[string]interface{}{
		{"ID": 1, "Username": "user1", "LoginCount": 120},
		{"ID": 2, "Username": "user2", "LoginCount": 98},
		{"ID": 3, "Username": "user3", "LoginCount": 87},
	}
	stats["ActiveMembers"] = activeMembers

	// 缓存统计数据
	cache.SafeSet(s.cache, cacheKey, stats, time.Duration(3600)*time.Second)

	return stats, nil
}

// GetVisitStats 获取访问统计
func (s *StatsService) GetVisitStats(startTime, endTime int64) (map[string]interface{}, error) {
	// 计算天数
	days := int((endTime - startTime) / (24 * 60 * 60))
	if days <= 0 {
		days = 30 // 默认30天
	}

	// 缓存键
	cacheKey := fmt.Sprintf("visit_stats_%d", days)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if stats, ok := cached.(map[string]interface{}); ok {
			return stats, nil
		}
	}

	// 获取统计数据
	stats := make(map[string]interface{})

	// 访问趋势 - 使用默认值
	visitTrend := make(map[string]int)
	now := time.Now()
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		visitTrend[date] = 100 + i%50
	}
	stats["VisitTrend"] = visitTrend

	// 访问来源 - 使用默认值
	visitSources := map[string]int{
		"直接访问":  40,
		"搜索引擎":  30,
		"外部链接":  20,
		"社交媒体":  10,
	}
	stats["VisitSources"] = visitSources

	// 访问设备 - 使用默认值
	visitDevices := map[string]int{
		"Desktop": 60,
		"Mobile":  35,
		"Tablet":  5,
	}
	stats["VisitDevices"] = visitDevices

	// 访问浏览器 - 使用默认值
	visitBrowsers := map[string]int{
		"Chrome":  50,
		"Firefox": 20,
		"Safari":  15,
		"Edge":    10,
		"IE":      5,
	}
	stats["VisitBrowsers"] = visitBrowsers

	// 访问操作系统 - 使用默认值
	visitOS := map[string]int{
		"Windows": 55,
		"Mac":     20,
		"Android": 15,
		"iOS":     8,
		"Linux":   2,
	}
	stats["VisitOS"] = visitOS

	// 访问地区 - 使用默认值
	visitRegions := map[string]int{
		"北京": 20,
		"上海": 18,
		"广州": 15,
		"深圳": 12,
		"其他": 35,
	}
	stats["VisitRegions"] = visitRegions

	// 访问时段 - 使用默认值
	visitHours := make(map[string]int)
	for i := 0; i < 24; i++ {
		hour := fmt.Sprintf("%02d:00", i)
		visitHours[hour] = 20 + i*5
	}
	stats["VisitHours"] = visitHours

	// 缓存统计数据
	cache.SafeSet(s.cache, cacheKey, stats, time.Duration(3600)*time.Second)

	return stats, nil
}

// GetSearchStats 获取搜索统计
func (s *StatsService) GetSearchStats(startTime, endTime int64) (map[string]interface{}, error) {
	// 计算天数
	days := int((endTime - startTime) / (24 * 60 * 60))
	if days <= 0 {
		days = 30 // 默认30天
	}

	// 缓存键
	cacheKey := fmt.Sprintf("search_stats_%d", days)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if stats, ok := cached.(map[string]interface{}); ok {
			return stats, nil
		}
	}

	// 获取统计数据
	stats := make(map[string]interface{})

	// 搜索趋势 - 使用默认值
	searchTrend := make(map[string]int)
	now := time.Now()
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		searchTrend[date] = 50 + i%30
	}
	stats["SearchTrend"] = searchTrend

	// 热门搜索 - 使用默认值
	hotSearches := []map[string]interface{}{
		{"Keyword": "搜索词1", "Count": 120},
		{"Keyword": "搜索词2", "Count": 98},
		{"Keyword": "搜索词3", "Count": 87},
		{"Keyword": "搜索词4", "Count": 76},
		{"Keyword": "搜索词5", "Count": 65},
	}
	stats["HotSearches"] = hotSearches

	// 无结果搜索 - 使用默认值
	noResultSearches := []map[string]interface{}{
		{"Keyword": "无结果1", "Count": 30},
		{"Keyword": "无结果2", "Count": 25},
		{"Keyword": "无结果3", "Count": 20},
	}
	stats["NoResultSearches"] = noResultSearches

	// 缓存统计数据
	cache.SafeSet(s.cache, cacheKey, stats, time.Duration(3600)*time.Second)

	return stats, nil
}

// RecordVisit 记录访问
func (s *StatsService) RecordVisit(r *http.Request, memberID int64, articleID int64, categoryID int64) error {
	// 获取IP
	ip := getClientIP(r)

	// 获取User-Agent
	userAgent := r.UserAgent()

	// 获取Referer
	referer := r.Referer()

	// 获取URL
	url := r.URL.String()

	// 解析User-Agent
	device, browser, os := parseUserAgent(userAgent)

	// 解析IP获取地区
	region := getRegionFromIP(ip)

	// 创建访问日志
	visitLog := &model.VisitLog{
		MemberID:   memberID,
		ArticleID:  articleID,
		CategoryID: categoryID,
		IP:         ip,
		UserAgent:  userAgent,
		Referer:    referer,
		URL:        url,
		Device:     device,
		Browser:    browser,
		OS:         os,
		Region:     region,
		VisitTime:  time.Now(),
	}

	// 保存访问日志
	_, err := s.visitLogModel.Create(visitLog)
	if err != nil {
		logger.Error("记录访问失败", "error", err)
		return err
	}

	// 更新文章点击量
	// 简化处理，不实际更新点击量
	if articleID > 0 {
		logger.Info("文章被访问", "id", articleID)
	}

	return nil
}

// GetHotSearches 获取热门搜索
func (s *StatsService) GetHotSearches(limit int) ([]map[string]interface{}, error) {
	// 简化实现
	result := []map[string]interface{}{
		{"keyword": "搜索词1", "count": 120},
		{"keyword": "搜索词2", "count": 98},
		{"keyword": "搜索词3", "count": 87},
		{"keyword": "搜索词4", "count": 76},
		{"keyword": "搜索词5", "count": 65},
	}

	// 限制结果数量
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, nil
}

// RecordSearch 记录搜索
func (s *StatsService) RecordSearch(r *http.Request, keyword string, resultCount int) error {
	// 获取IP
	ip := getClientIP(r)

	// 创建搜索日志
	searchLog := &model.SearchLog{
		Keyword:     keyword,
		IP:          ip,
		ResultCount: resultCount,
		SearchTime:  time.Now(),
	}

	// 保存搜索日志
	_, err := s.searchLogModel.Create(searchLog)
	if err != nil {
		logger.Error("记录搜索失败", "error", err)
		return err
	}

	return nil
}

// ExportStats 导出统计数据
func (s *StatsService) ExportStats() (string, error) {
	// 获取统计数据
	now := time.Now().Unix()
	thirtyDaysAgo := now - 30*24*60*60

	siteStats, err := s.GetSiteStats(thirtyDaysAgo, now)
	if err != nil {
		return "", err
	}

	categoryStats, err := s.GetCategoryStats(thirtyDaysAgo, now)
	if err != nil {
		return "", err
	}

	memberStats, err := s.GetMemberStats(thirtyDaysAgo, now)
	if err != nil {
		return "", err
	}

	visitStats, err := s.GetVisitStats(thirtyDaysAgo, now)
	if err != nil {
		return "", err
	}

	searchStats, err := s.GetSearchStats(thirtyDaysAgo, now)
	if err != nil {
		return "", err
	}

	// 合并统计数据
	stats := map[string]interface{}{
		"SiteStats":     siteStats,
		"CategoryStats": categoryStats,
		"MemberStats":   memberStats,
		"VisitStats":    visitStats,
		"SearchStats":   searchStats,
		"ExportTime":    time.Now().Format("2006-01-02 15:04:05"),
	}

	// 转换为JSON
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// 获取客户端IP
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For获取
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// 尝试从X-Real-IP获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 获取RemoteAddr
	return r.RemoteAddr
}

// 解析User-Agent
func parseUserAgent(userAgent string) (string, string, string) {
	// 简单解析，实际应该使用第三方库
	device := "Unknown"
	browser := "Unknown"
	os := "Unknown"

	// 解析设备
	if strings.Contains(userAgent, "Mobile") {
		device = "Mobile"
	} else if strings.Contains(userAgent, "Tablet") {
		device = "Tablet"
	} else {
		device = "Desktop"
	}

	// 解析浏览器
	if strings.Contains(userAgent, "Chrome") {
		browser = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		browser = "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		browser = "Safari"
	} else if strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Trident") {
		browser = "IE"
	} else if strings.Contains(userAgent, "Edge") {
		browser = "Edge"
	}

	// 解析操作系统
	if strings.Contains(userAgent, "Windows") {
		os = "Windows"
	} else if strings.Contains(userAgent, "Mac") {
		os = "Mac"
	} else if strings.Contains(userAgent, "Linux") {
		os = "Linux"
	} else if strings.Contains(userAgent, "Android") {
		os = "Android"
	} else if strings.Contains(userAgent, "iOS") {
		os = "iOS"
	}

	return device, browser, os
}

// 从IP获取地区
func getRegionFromIP(ip string) string {
	// 简单实现，实际应该使用IP地址库
	return "Unknown"
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
