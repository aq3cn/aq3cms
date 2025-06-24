package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// ScoreRule 积分规则
type ScoreRule struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 规则名称
	Code        string    `json:"code"`        // 规则代码
	Score       int       `json:"score"`       // 积分值
	MaxTimes    int       `json:"maxtimes"`    // 最大次数，0表示不限制
	CycleType   int       `json:"cycletype"`   // 周期类型：0不限，1每天，2每周，3每月，4每年
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	Description string    `json:"description"` // 描述
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// ScoreLog 积分日志
type ScoreLog struct {
	ID         int64     `json:"id"`
	MemberID   int64     `json:"memberid"`   // 会员ID
	RuleID     int64     `json:"ruleid"`     // 规则ID
	Score      int       `json:"score"`      // 积分值
	Remark     string    `json:"remark"`     // 备注
	IP         string    `json:"ip"`         // IP地址
	CreateTime time.Time `json:"createtime"` // 创建时间
}

// ScoreRuleModel 积分规则模型
type ScoreRuleModel struct {
	db *database.DB
}

// NewScoreRuleModel 创建积分规则模型
func NewScoreRuleModel(db *database.DB) *ScoreRuleModel {
	return &ScoreRuleModel{
		db: db,
	}
}

// GetByID 根据ID获取积分规则
func (m *ScoreRuleModel) GetByID(id int64) (*ScoreRule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_rule")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取积分规则失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("score rule not found: %d", id)
	}

	// 转换为积分规则
	rule := &ScoreRule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.Code, _ = result["code"].(string)
	rule.Score, _ = result["score"].(int)
	rule.MaxTimes, _ = result["maxtimes"].(int)
	rule.CycleType, _ = result["cycletype"].(int)
	rule.Status, _ = result["status"].(int)
	rule.Description, _ = result["description"].(string)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)

	return rule, nil
}

// GetByCode 根据代码获取积分规则
func (m *ScoreRuleModel) GetByCode(code string) (*ScoreRule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_rule")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取积分规则失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("score rule not found: %s", code)
	}

	// 转换为积分规则
	rule := &ScoreRule{}
	rule.ID, _ = result["id"].(int64)
	rule.Name, _ = result["name"].(string)
	rule.Code, _ = result["code"].(string)
	rule.Score, _ = result["score"].(int)
	rule.MaxTimes, _ = result["maxtimes"].(int)
	rule.CycleType, _ = result["cycletype"].(int)
	rule.Status, _ = result["status"].(int)
	rule.Description, _ = result["description"].(string)
	rule.CreateTime, _ = result["createtime"].(time.Time)
	rule.UpdateTime, _ = result["updatetime"].(time.Time)

	return rule, nil
}

// GetAll 获取所有积分规则
func (m *ScoreRuleModel) GetAll(status int) ([]*ScoreRule, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_rule")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有积分规则失败", "error", err)
		return nil, err
	}

	// 转换为积分规则列表
	rules := make([]*ScoreRule, 0, len(results))
	for _, result := range results {
		rule := &ScoreRule{}
		rule.ID, _ = result["id"].(int64)
		rule.Name, _ = result["name"].(string)
		rule.Code, _ = result["code"].(string)
		rule.Score, _ = result["score"].(int)
		rule.MaxTimes, _ = result["maxtimes"].(int)
		rule.CycleType, _ = result["cycletype"].(int)
		rule.Status, _ = result["status"].(int)
		rule.Description, _ = result["description"].(string)
		rule.CreateTime, _ = result["createtime"].(time.Time)
		rule.UpdateTime, _ = result["updatetime"].(time.Time)
		rules = append(rules, rule)
	}

	return rules, nil
}

// Create 创建积分规则
func (m *ScoreRuleModel) Create(rule *ScoreRule) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	rule.CreateTime = now
	rule.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("score_rule")+" (name, code, score, maxtimes, cycletype, status, description, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		rule.Name, rule.Code, rule.Score, rule.MaxTimes, rule.CycleType, rule.Status, rule.Description, rule.CreateTime, rule.UpdateTime,
	)
	if err != nil {
		logger.Error("创建积分规则失败", "error", err)
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

// Update 更新积分规则
func (m *ScoreRuleModel) Update(rule *ScoreRule) error {
	// 设置更新时间
	rule.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("score_rule")+" SET name = ?, code = ?, score = ?, maxtimes = ?, cycletype = ?, status = ?, description = ?, updatetime = ? WHERE id = ?",
		rule.Name, rule.Code, rule.Score, rule.MaxTimes, rule.CycleType, rule.Status, rule.Description, rule.UpdateTime, rule.ID,
	)
	if err != nil {
		logger.Error("更新积分规则失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除积分规则
func (m *ScoreRuleModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("score_rule")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除积分规则失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新积分规则状态
func (m *ScoreRuleModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("score_rule")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新积分规则状态失败", "error", err)
		return err
	}

	return nil
}

// ScoreLogModel 积分日志模型
type ScoreLogModel struct {
	db *database.DB
}

// NewScoreLogModel 创建积分日志模型
func NewScoreLogModel(db *database.DB) *ScoreLogModel {
	return &ScoreLogModel{
		db: db,
	}
}

// GetByID 根据ID获取积分日志
func (m *ScoreLogModel) GetByID(id int64) (*ScoreLog, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取积分日志失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("score log not found: %d", id)
	}

	// 转换为积分日志
	log := &ScoreLog{}
	log.ID, _ = result["id"].(int64)
	log.MemberID, _ = result["memberid"].(int64)
	log.RuleID, _ = result["ruleid"].(int64)
	log.Score, _ = result["score"].(int)
	log.Remark, _ = result["remark"].(string)
	log.IP, _ = result["ip"].(string)
	log.CreateTime, _ = result["createtime"].(time.Time)

	return log, nil
}

// GetByMemberID 根据会员ID获取积分日志
func (m *ScoreLogModel) GetByMemberID(memberID int64, page, pageSize int) ([]*ScoreLog, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取积分日志失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取积分日志总数失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}

	// 转换为积分日志列表
	logs := make([]*ScoreLog, 0, len(results))
	for _, result := range results {
		log := &ScoreLog{}
		log.ID, _ = result["id"].(int64)
		log.MemberID, _ = result["memberid"].(int64)
		log.RuleID, _ = result["ruleid"].(int64)
		log.Score, _ = result["score"].(int)
		log.Remark, _ = result["remark"].(string)
		log.IP, _ = result["ip"].(string)
		log.CreateTime, _ = result["createtime"].(time.Time)
		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetByRuleID 根据规则ID获取积分日志
func (m *ScoreLogModel) GetByRuleID(ruleID int64, page, pageSize int) ([]*ScoreLog, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("ruleid = ?", ruleID)
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取积分日志失败", "ruleid", ruleID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取积分日志总数失败", "ruleid", ruleID, "error", err)
		return nil, 0, err
	}

	// 转换为积分日志列表
	logs := make([]*ScoreLog, 0, len(results))
	for _, result := range results {
		log := &ScoreLog{}
		log.ID, _ = result["id"].(int64)
		log.MemberID, _ = result["memberid"].(int64)
		log.RuleID, _ = result["ruleid"].(int64)
		log.Score, _ = result["score"].(int)
		log.Remark, _ = result["remark"].(string)
		log.IP, _ = result["ip"].(string)
		log.CreateTime, _ = result["createtime"].(time.Time)
		logs = append(logs, log)
	}

	return logs, total, nil
}

// Create 创建积分日志
func (m *ScoreLogModel) Create(log *ScoreLog) (int64, error) {
	// 设置创建时间
	log.CreateTime = time.Now()

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("score_log")+" (memberid, ruleid, score, remark, ip, createtime) VALUES (?, ?, ?, ?, ?, ?)",
		log.MemberID, log.RuleID, log.Score, log.Remark, log.IP, log.CreateTime,
	)
	if err != nil {
		logger.Error("创建积分日志失败", "error", err)
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

// Delete 删除积分日志
func (m *ScoreLogModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("score_log")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除积分日志失败", "error", err)
		return err
	}

	return nil
}

// DeleteByMemberID 根据会员ID删除积分日志
func (m *ScoreLogModel) DeleteByMemberID(memberID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("score_log")+" WHERE memberid = ?", memberID)
	if err != nil {
		logger.Error("删除积分日志失败", "memberid", memberID, "error", err)
		return err
	}

	return nil
}

// DeleteByRuleID 根据规则ID删除积分日志
func (m *ScoreLogModel) DeleteByRuleID(ruleID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("score_log")+" WHERE ruleid = ?", ruleID)
	if err != nil {
		logger.Error("删除积分日志失败", "ruleid", ruleID, "error", err)
		return err
	}

	return nil
}

// GetTodayCount 获取今日次数
func (m *ScoreLogModel) GetTodayCount(memberID int64, ruleID int64) (int, error) {
	// 获取今日开始时间
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.Where("ruleid = ?", ruleID)
	qb.Where("createtime >= ?", todayStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取今日次数失败", "memberid", memberID, "ruleid", ruleID, "error", err)
		return 0, err
	}

	return count, nil
}

// GetWeekCount 获取本周次数
func (m *ScoreLogModel) GetWeekCount(memberID int64, ruleID int64) (int, error) {
	// 获取本周开始时间
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.Where("ruleid = ?", ruleID)
	qb.Where("createtime >= ?", weekStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本周次数失败", "memberid", memberID, "ruleid", ruleID, "error", err)
		return 0, err
	}

	return count, nil
}

// GetMonthCount 获取本月次数
func (m *ScoreLogModel) GetMonthCount(memberID int64, ruleID int64) (int, error) {
	// 获取本月开始时间
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.Where("ruleid = ?", ruleID)
	qb.Where("createtime >= ?", monthStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本月次数失败", "memberid", memberID, "ruleid", ruleID, "error", err)
		return 0, err
	}

	return count, nil
}

// GetYearCount 获取本年次数
func (m *ScoreLogModel) GetYearCount(memberID int64, ruleID int64) (int, error) {
	// 获取本年开始时间
	now := time.Now()
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.Where("ruleid = ?", ruleID)
	qb.Where("createtime >= ?", yearStart)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取本年次数失败", "memberid", memberID, "ruleid", ruleID, "error", err)
		return 0, err
	}

	return count, nil
}

// GetTotalCount 获取总次数
func (m *ScoreLogModel) GetTotalCount(memberID int64, ruleID int64) (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "score_log")
	qb.Where("memberid = ?", memberID)
	qb.Where("ruleid = ?", ruleID)

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取总次数失败", "memberid", memberID, "ruleid", ruleID, "error", err)
		return 0, err
	}

	return count, nil
}
