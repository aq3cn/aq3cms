package model

import (
	"encoding/json"
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// CollectRule 采集规则
type CollectRule struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 规则名称
	SourceType  int       `json:"sourcetype"`  // 来源类型：0列表，1RSS
	SourceURL   string    `json:"sourceurl"`   // 来源URL
	StartPage   int       `json:"startpage"`   // 开始页码
	EndPage     int       `json:"endpage"`     // 结束页码
	PageRule    string    `json:"pagerule"`    // 分页规则
	ListRule    string    `json:"listrule"`    // 列表规则
	TitleRule   string    `json:"titlerule"`   // 标题规则
	ContentRule string    `json:"contentrule"` // 内容规则
	TypeID      int64     `json:"typeid"`      // 栏目ID
	ModelID     int64     `json:"modelid"`     // 模型ID
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	LastTime    time.Time `json:"lasttime"`    // 最后采集时间
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
	FieldRules  string    `json:"fieldrules"`  // 字段规则，JSON格式
}

// FieldRule 字段规则
type FieldRule struct {
	Field    string `json:"field"`    // 字段名
	Rule     string `json:"rule"`     // 规则
	Required bool   `json:"required"` // 是否必填
	Default  string `json:"default"`  // 默认值
}

// CollectItem 采集项目
type CollectItem struct {
	ID         int64     `json:"id"`
	RuleID     int64     `json:"ruleid"`     // 规则ID
	Title      string    `json:"title"`      // 标题
	URL        string    `json:"url"`        // 来源URL
	Content    string    `json:"content"`    // 内容
	Status     int       `json:"status"`     // 状态：0未处理，1已处理，2已发布
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
	FieldData  string    `json:"fielddata"`  // 字段数据，JSON格式
}

// CollectRuleModel 采集规则模型
type CollectRuleModel struct {
	db *database.DB
}

// NewCollectRuleModel 创建采集规则模型
func NewCollectRuleModel(db *database.DB) *CollectRuleModel {
	return &CollectRuleModel{
		db: db,
	}
}

// GetByID 根据ID获取采集规则
func (m *CollectRuleModel) GetByID(id int64) (*CollectRule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "collect_rule")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取采集规则失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("collect rule not found: %d", id)
	}

	// 转换为采集规则
	rule := &CollectRule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.SourceType, _ = result["sourcetype"].(int)
	rule.SourceURL, _ = result["sourceurl"].(string)
	rule.StartPage, _ = result["startpage"].(int)
	rule.EndPage, _ = result["endpage"].(int)
	rule.PageRule, _ = result["pagerule"].(string)
	rule.ListRule, _ = result["listrule"].(string)
	rule.TitleRule, _ = result["titlerule"].(string)
	rule.ContentRule, _ = result["contentrule"].(string)
	rule.TypeID, _ = result["typeid"].(int64)
	rule.ModelID, _ = result["modelid"].(int64)
	rule.Status, _ = result["status"].(int)
	rule.LastTime, _ = result["lasttime"].(time.Time)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)
	rule.FieldRules, _ = result["fieldrules"].(string)

	return rule, nil
}

// GetAll 获取所有采集规则
func (m *CollectRuleModel) GetAll() ([]*CollectRule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "collect_rule")
	qb.OrderBy("id DESC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有采集规则失败", "error", err)
		return nil, err
	}

	// 转换为采集规则列表
	rules := make([]*CollectRule, 0, len(results))
	for _, result := range results {
		rule := &CollectRule{}
		rule.ID, _ = result["id"].(int64)
		rule.Name, _ = result["name"].(string)
		rule.SourceType, _ = result["sourcetype"].(int)
		rule.SourceURL, _ = result["sourceurl"].(string)
		rule.StartPage, _ = result["startpage"].(int)
		rule.EndPage, _ = result["endpage"].(int)
		rule.PageRule, _ = result["pagerule"].(string)
		rule.ListRule, _ = result["listrule"].(string)
		rule.TitleRule, _ = result["titlerule"].(string)
		rule.ContentRule, _ = result["contentrule"].(string)
		rule.TypeID, _ = result["typeid"].(int64)
		rule.ModelID, _ = result["modelid"].(int64)
		rule.Status, _ = result["status"].(int)
		rule.LastTime, _ = result["lasttime"].(time.Time)
		rule.CreateTime, _ = result["createtime"].(time.Time)
		rule.UpdateTime, _ = result["updatetime"].(time.Time)
		rule.FieldRules, _ = result["fieldrules"].(string)
		rules = append(rules, rule)
	}

	return rules, nil
}

// Create 创建采集规则
func (m *CollectRuleModel) Create(rule *CollectRule) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	rule.CreateTime = now
	rule.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("collect_rule")+" (name, sourcetype, sourceurl, startpage, endpage, pagerule, listrule, titlerule, contentrule, typeid, modelid, status, lasttime, createtime, updatetime, fieldrules) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rule.Name, rule.SourceType, rule.SourceURL, rule.StartPage, rule.EndPage, rule.PageRule, rule.ListRule, rule.TitleRule, rule.ContentRule, rule.TypeID, rule.ModelID, rule.Status, rule.LastTime, rule.CreateTime, rule.UpdateTime, rule.FieldRules,
	)
	if err != nil {
		logger.Error("创建采集规则失败", "error", err)
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

// Update 更新采集规则
func (m *CollectRuleModel) Update(rule *CollectRule) error {
	// 设置更新时间
	rule.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("collect_rule")+" SET name = ?, sourcetype = ?, sourceurl = ?, startpage = ?, endpage = ?, pagerule = ?, listrule = ?, titlerule = ?, contentrule = ?, typeid = ?, modelid = ?, status = ?, lasttime = ?, updatetime = ?, fieldrules = ? WHERE id = ?",
		rule.Name, rule.SourceType, rule.SourceURL, rule.StartPage, rule.EndPage, rule.PageRule, rule.ListRule, rule.TitleRule, rule.ContentRule, rule.TypeID, rule.ModelID, rule.Status, rule.LastTime, rule.UpdateTime, rule.FieldRules, rule.ID,
	)
	if err != nil {
		logger.Error("更新采集规则失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除采集规则
func (m *CollectRuleModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("collect_rule")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除采集规则失败", "error", err)
		return err
	}

	return nil
}

// UpdateLastTime 更新最后采集时间
func (m *CollectRuleModel) UpdateLastTime(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("collect_rule")+" SET lasttime = ? WHERE id = ?",
		time.Now(), id,
	)
	if err != nil {
		logger.Error("更新最后采集时间失败", "error", err)
		return err
	}

	return nil
}

// CollectItemModel 采集项目模型
type CollectItemModel struct {
	db *database.DB
}

// NewCollectItemModel 创建采集项目模型
func NewCollectItemModel(db *database.DB) *CollectItemModel {
	return &CollectItemModel{
		db: db,
	}
}

// GetByID 根据ID获取采集项目
func (m *CollectItemModel) GetByID(id int64) (*CollectItem, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "collect_item")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取采集项目失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("collect item not found: %d", id)
	}

	// 转换为采集项目
	item := &CollectItem{}
	item.ID, _ = result["id"].(int64)
	item.RuleID, _ = result["ruleid"].(int64)
	item.Title, _ = result["title"].(string)
	item.URL, _ = result["url"].(string)
	item.Content, _ = result["content"].(string)
	item.Status, _ = result["status"].(int)
	item.CreateTime, _ = result["createtime"].(time.Time)
	item.UpdateTime, _ = result["updatetime"].(time.Time)
	item.FieldData, _ = result["fielddata"].(string)

	return item, nil
}

// GetByRuleID 根据规则ID获取采集项目
func (m *CollectItemModel) GetByRuleID(ruleID int64, status int, page, pageSize int) ([]*CollectItem, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "collect_item")
	qb.Where("ruleid = ?", ruleID)
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取采集项目失败", "ruleid", ruleID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取采集项目总数失败", "ruleid", ruleID, "error", err)
		return nil, 0, err
	}

	// 转换为采集项目列表
	items := make([]*CollectItem, 0, len(results))
	for _, result := range results {
		item := &CollectItem{}
		item.ID, _ = result["id"].(int64)
		item.RuleID, _ = result["ruleid"].(int64)
		item.Title, _ = result["title"].(string)
		item.URL, _ = result["url"].(string)
		item.Content, _ = result["content"].(string)
		item.Status, _ = result["status"].(int)
		item.CreateTime, _ = result["createtime"].(time.Time)
		item.UpdateTime, _ = result["updatetime"].(time.Time)
		item.FieldData, _ = result["fielddata"].(string)
		items = append(items, item)
	}

	return items, total, nil
}

// Create 创建采集项目
func (m *CollectItemModel) Create(item *CollectItem) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	item.CreateTime = now
	item.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("collect_item")+" (ruleid, title, url, content, status, createtime, updatetime, fielddata) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		item.RuleID, item.Title, item.URL, item.Content, item.Status, item.CreateTime, item.UpdateTime, item.FieldData,
	)
	if err != nil {
		logger.Error("创建采集项目失败", "error", err)
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

// Update 更新采集项目
func (m *CollectItemModel) Update(item *CollectItem) error {
	// 设置更新时间
	item.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("collect_item")+" SET title = ?, content = ?, status = ?, updatetime = ?, fielddata = ? WHERE id = ?",
		item.Title, item.Content, item.Status, item.UpdateTime, item.FieldData, item.ID,
	)
	if err != nil {
		logger.Error("更新采集项目失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除采集项目
func (m *CollectItemModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("collect_item")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除采集项目失败", "error", err)
		return err
	}

	return nil
}

// DeleteByRuleID 根据规则ID删除采集项目
func (m *CollectItemModel) DeleteByRuleID(ruleID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("collect_item")+" WHERE ruleid = ?", ruleID)
	if err != nil {
		logger.Error("删除采集项目失败", "ruleid", ruleID, "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新采集项目状态
func (m *CollectItemModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("collect_item")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新采集项目状态失败", "error", err)
		return err
	}

	return nil
}

// GetFieldData 获取字段数据
func (m *CollectItemModel) GetFieldData(id int64) (map[string]interface{}, error) {
	// 获取采集项目
	item, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 解析字段数据
	var fieldData map[string]interface{}
	if item.FieldData != "" {
		err = json.Unmarshal([]byte(item.FieldData), &fieldData)
		if err != nil {
			logger.Error("解析字段数据失败", "error", err)
			return nil, err
		}
	} else {
		fieldData = make(map[string]interface{})
	}

	return fieldData, nil
}

// SetFieldData 设置字段数据
func (m *CollectItemModel) SetFieldData(id int64, fieldData map[string]interface{}) error {
	// 获取采集项目
	item, err := m.GetByID(id)
	if err != nil {
		return err
	}

	// 序列化字段数据
	data, err := json.Marshal(fieldData)
	if err != nil {
		logger.Error("序列化字段数据失败", "error", err)
		return err
	}

	// 更新字段数据
	item.FieldData = string(data)
	return m.Update(item)
}
