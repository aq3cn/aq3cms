package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// SEORule SEO规则
type SEORule struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 规则名称
	Type        int       `json:"type"`        // 类型：0首页，1栏目页，2文章页，3标签页，4自定义
	Pattern     string    `json:"pattern"`     // 匹配模式
	Title       string    `json:"title"`       // 标题模板
	Keywords    string    `json:"keywords"`    // 关键词模板
	Description string    `json:"description"` // 描述模板
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// SEOKeyword SEO关键词
type SEOKeyword struct {
	ID         int64     `json:"id"`
	Keyword    string    `json:"keyword"`    // 关键词
	URL        string    `json:"url"`        // 链接
	Weight     int       `json:"weight"`     // 权重
	Status     int       `json:"status"`     // 状态：0禁用，1启用
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// SEORuleModel SEO规则模型
type SEORuleModel struct {
	db *database.DB
}

// NewSEORuleModel 创建SEO规则模型
func NewSEORuleModel(db *database.DB) *SEORuleModel {
	return &SEORuleModel{
		db: db,
	}
}

// GetByID 根据ID获取SEO规则
func (m *SEORuleModel) GetByID(id int64) (*SEORule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_rule")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取SEO规则失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("seo rule not found: %d", id)
	}

	// 转换为SEO规则
	rule := &SEORule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.Type, _ = result["type"].(int)
	rule.Pattern, _ = result["pattern"].(string)
	rule.Title, _ = result["title"].(string)
	rule.Keywords, _ = result["keywords"].(string)
	rule.Description, _ = result["description"].(string)
	rule.Status, _ = result["status"].(int)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)

	return rule, nil
}

// GetByType 根据类型获取SEO规则
func (m *SEORuleModel) GetByType(ruleType int) (*SEORule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_rule")
	qb.Where("type = ?", ruleType)
	qb.Where("status = ?", 1)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取SEO规则失败", "type", ruleType, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("seo rule not found: %d", ruleType)
	}

	// 转换为SEO规则
	rule := &SEORule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.Type, _ = result["type"].(int)
	rule.Pattern, _ = result["pattern"].(string)
	rule.Title, _ = result["title"].(string)
	rule.Keywords, _ = result["keywords"].(string)
	rule.Description, _ = result["description"].(string)
	rule.Status, _ = result["status"].(int)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)

	return rule, nil
}

// GetByPattern 根据匹配模式获取SEO规则
func (m *SEORuleModel) GetByPattern(pattern string) (*SEORule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_rule")
	qb.Where("pattern = ?", pattern)
	qb.Where("status = ?", 1)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取SEO规则失败", "pattern", pattern, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("seo rule not found: %s", pattern)
	}

	// 转换为SEO规则
	rule := &SEORule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.Type, _ = result["type"].(int)
	rule.Pattern, _ = result["pattern"].(string)
	rule.Title, _ = result["title"].(string)
	rule.Keywords, _ = result["keywords"].(string)
	rule.Description, _ = result["description"].(string)
	rule.Status, _ = result["status"].(int)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)

	return rule, nil
}

// GetAll 获取所有SEO规则
func (m *SEORuleModel) GetAll(status int) ([]*SEORule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_rule")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("type ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有SEO规则失败", "error", err)
		return nil, err
	}

	// 转换为SEO规则列表
	rules := make([]*SEORule, 0, len(results))
	for _, result := range results {
		rule := &SEORule{}
		rule.ID, _ = result["id"].(int64)
		rule.Name, _ = result["name"].(string)
		rule.Type, _ = result["type"].(int)
		rule.Pattern, _ = result["pattern"].(string)
		rule.Title, _ = result["title"].(string)
		rule.Keywords, _ = result["keywords"].(string)
		rule.Description, _ = result["description"].(string)
		rule.Status, _ = result["status"].(int)
		rule.CreateTime, _ = result["createtime"].(time.Time)
		rule.UpdateTime, _ = result["updatetime"].(time.Time)
		rules = append(rules, rule)
	}

	return rules, nil
}

// Create 创建SEO规则
func (m *SEORuleModel) Create(rule *SEORule) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	rule.CreateTime = now
	rule.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("seo_rule")+" (name, type, pattern, title, keywords, description, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rule.Name, rule.Type, rule.Pattern, rule.Title, rule.Keywords, rule.Description, rule.Status, rule.CreateTime, rule.UpdateTime,
	)
	if err != nil {
		logger.Error("创建SEO规则失败", "error", err)
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

// Update 更新SEO规则
func (m *SEORuleModel) Update(rule *SEORule) error {
	// 设置更新时间
	rule.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("seo_rule")+" SET name = ?, type = ?, pattern = ?, title = ?, keywords = ?, description = ?, status = ?, updatetime = ? WHERE id = ?",
		rule.Name, rule.Type, rule.Pattern, rule.Title, rule.Keywords, rule.Description, rule.Status, rule.UpdateTime, rule.ID,
	)
	if err != nil {
		logger.Error("更新SEO规则失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除SEO规则
func (m *SEORuleModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("seo_rule")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除SEO规则失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新SEO规则状态
func (m *SEORuleModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("seo_rule")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新SEO规则状态失败", "error", err)
		return err
	}

	return nil
}

// SEOKeywordModel SEO关键词模型
type SEOKeywordModel struct {
	db *database.DB
}

// NewSEOKeywordModel 创建SEO关键词模型
func NewSEOKeywordModel(db *database.DB) *SEOKeywordModel {
	return &SEOKeywordModel{
		db: db,
	}
}

// GetByID 根据ID获取SEO关键词
func (m *SEOKeywordModel) GetByID(id int64) (*SEOKeyword, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_keyword")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取SEO关键词失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("seo keyword not found: %d", id)
	}

	// 转换为SEO关键词
	keyword := &SEOKeyword{}
	keyword.ID, _ = result["id"].(int64)
	keyword.Keyword, _ = result["keyword"].(string)
	keyword.URL, _ = result["url"].(string)
	keyword.Weight, _ = result["weight"].(int)
	keyword.Status, _ = result["status"].(int)
	keyword.CreateTime, _ = result["createtime"].(time.Time)
	keyword.UpdateTime, _ = result["updatetime"].(time.Time)

	return keyword, nil
}

// GetByKeyword 根据关键词获取SEO关键词
func (m *SEOKeywordModel) GetByKeyword(keyword string) (*SEOKeyword, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_keyword")
	qb.Where("keyword = ?", keyword)
	qb.Where("status = ?", 1)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取SEO关键词失败", "keyword", keyword, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("seo keyword not found: %s", keyword)
	}

	// 转换为SEO关键词
	seoKeyword := &SEOKeyword{}
	seoKeyword.ID, _ = result["id"].(int64)
	seoKeyword.Keyword, _ = result["keyword"].(string)
	seoKeyword.URL, _ = result["url"].(string)
	seoKeyword.Weight, _ = result["weight"].(int)
	seoKeyword.Status, _ = result["status"].(int)
	seoKeyword.CreateTime, _ = result["createtime"].(time.Time)
	seoKeyword.UpdateTime, _ = result["updatetime"].(time.Time)

	return seoKeyword, nil
}

// GetAll 获取所有SEO关键词
func (m *SEOKeywordModel) GetAll(status int, page, pageSize int) ([]*SEOKeyword, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "seo_keyword")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("weight DESC, id ASC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有SEO关键词失败", "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取SEO关键词总数失败", "error", err)
		return nil, 0, err
	}

	// 转换为SEO关键词列表
	keywords := make([]*SEOKeyword, 0, len(results))
	for _, result := range results {
		keyword := &SEOKeyword{}
		keyword.ID, _ = result["id"].(int64)
		keyword.Keyword, _ = result["keyword"].(string)
		keyword.URL, _ = result["url"].(string)
		keyword.Weight, _ = result["weight"].(int)
		keyword.Status, _ = result["status"].(int)
		keyword.CreateTime, _ = result["createtime"].(time.Time)
		keyword.UpdateTime, _ = result["updatetime"].(time.Time)
		keywords = append(keywords, keyword)
	}

	return keywords, total, nil
}

// Create 创建SEO关键词
func (m *SEOKeywordModel) Create(keyword *SEOKeyword) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	keyword.CreateTime = now
	keyword.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("seo_keyword")+" (keyword, url, weight, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?)",
		keyword.Keyword, keyword.URL, keyword.Weight, keyword.Status, keyword.CreateTime, keyword.UpdateTime,
	)
	if err != nil {
		logger.Error("创建SEO关键词失败", "error", err)
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

// Update 更新SEO关键词
func (m *SEOKeywordModel) Update(keyword *SEOKeyword) error {
	// 设置更新时间
	keyword.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("seo_keyword")+" SET keyword = ?, url = ?, weight = ?, status = ?, updatetime = ? WHERE id = ?",
		keyword.Keyword, keyword.URL, keyword.Weight, keyword.Status, keyword.UpdateTime, keyword.ID,
	)
	if err != nil {
		logger.Error("更新SEO关键词失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除SEO关键词
func (m *SEOKeywordModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("seo_keyword")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除SEO关键词失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新SEO关键词状态
func (m *SEOKeywordModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("seo_keyword")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新SEO关键词状态失败", "error", err)
		return err
	}

	return nil
}

// BatchCreate 批量创建SEO关键词
func (m *SEOKeywordModel) BatchCreate(keywords []*SEOKeyword) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 设置时间
	now := time.Now()

	// 批量插入
	for _, keyword := range keywords {
		keyword.CreateTime = now
		keyword.UpdateTime = now

		_, err = tx.Exec(
			"INSERT INTO "+m.db.TableName("seo_keyword")+" (keyword, url, weight, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?)",
			keyword.Keyword, keyword.URL, keyword.Weight, keyword.Status, keyword.CreateTime, keyword.UpdateTime,
		)
		if err != nil {
			logger.Error("创建SEO关键词失败", "error", err)
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// BatchDelete 批量删除SEO关键词
func (m *SEOKeywordModel) BatchDelete(ids []int64) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 批量删除
	for _, id := range ids {
		_, err = tx.Exec("DELETE FROM "+m.db.TableName("seo_keyword")+" WHERE id = ?", id)
		if err != nil {
			logger.Error("删除SEO关键词失败", "error", err)
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// BatchUpdateStatus 批量更新SEO关键词状态
func (m *SEOKeywordModel) BatchUpdateStatus(ids []int64, status int) error {
	// 开始事务
	tx, err := m.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 设置时间
	now := time.Now()

	// 批量更新
	for _, id := range ids {
		_, err = tx.Exec(
			"UPDATE "+m.db.TableName("seo_keyword")+" SET status = ?, updatetime = ? WHERE id = ?",
			status, now, id,
		)
		if err != nil {
			logger.Error("更新SEO关键词状态失败", "error", err)
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}
