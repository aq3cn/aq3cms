package interfaces

// StatsServiceInterface 统计服务接口
type StatsServiceInterface interface {
	// RecordPageView 记录页面访问
	RecordPageView(path, referer, userAgent, ip string)

	// GetPageViews 获取页面访问量
	GetPageViews(path string, startTime, endTime int64) (int, error)

	// GetUniqueVisitors 获取独立访客数
	GetUniqueVisitors(startTime, endTime int64) (int, error)

	// GetPopularPages 获取热门页面
	GetPopularPages(limit int, startTime, endTime int64) ([]map[string]interface{}, error)

	// GetReferrers 获取来源网站
	GetReferrers(limit int, startTime, endTime int64) ([]map[string]interface{}, error)

	// GetBrowsers 获取浏览器统计
	GetBrowsers(startTime, endTime int64) ([]map[string]interface{}, error)

	// GetOS 获取操作系统统计
	GetOS(startTime, endTime int64) ([]map[string]interface{}, error)

	// GetDevices 获取设备统计
	GetDevices(startTime, endTime int64) ([]map[string]interface{}, error)

	// GetSiteStats 获取站点统计
	GetSiteStats(startTime, endTime int64) (map[string]interface{}, error)

	// GetCategoryStats 获取分类统计
	GetCategoryStats(startTime, endTime int64) (map[string]interface{}, error)

	// GetMemberStats 获取会员统计
	GetMemberStats(startTime, endTime int64) (map[string]interface{}, error)

	// GetVisitStats 获取访问统计
	GetVisitStats(startTime, endTime int64) (map[string]interface{}, error)

	// GetSearchStats 获取搜索统计
	GetSearchStats(startTime, endTime int64) (map[string]interface{}, error)
}
