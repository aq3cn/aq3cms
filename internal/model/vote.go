package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// Vote 投票
type Vote struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`      // 投票标题
	Description string   `json:"description"` // 投票描述
	StartTime  time.Time `json:"starttime"`  // 开始时间
	EndTime    time.Time `json:"endtime"`    // 结束时间
	IsMulti    int       `json:"ismulti"`    // 是否多选：0否，1是
	MaxChoices int       `json:"maxchoices"` // 最多可选几项
	Status     int       `json:"status"`     // 状态：0禁用，1启用
	TotalCount int       `json:"totalcount"` // 总投票数
	CreateTime time.Time `json:"createtime"` // 创建时间
	UpdateTime time.Time `json:"updatetime"` // 更新时间
}

// VoteOption 投票选项
type VoteOption struct {
	ID      int64  `json:"id"`
	VoteID  int64  `json:"voteid"`  // 投票ID
	Title   string `json:"title"`   // 选项标题
	Count   int    `json:"count"`   // 投票数
	OrderID int    `json:"orderid"` // 排序ID
}

// VoteLog 投票日志
type VoteLog struct {
	ID         int64     `json:"id"`
	VoteID     int64     `json:"voteid"`     // 投票ID
	OptionID   int64     `json:"optionid"`   // 选项ID
	MemberID   int64     `json:"memberid"`   // 会员ID
	IP         string    `json:"ip"`         // IP地址
	CreateTime time.Time `json:"createtime"` // 创建时间
}

// VoteModel 投票模型
type VoteModel struct {
	db *database.DB
}

// NewVoteModel 创建投票模型
func NewVoteModel(db *database.DB) *VoteModel {
	return &VoteModel{
		db: db,
	}
}

// GetByID 根据ID获取投票
func (m *VoteModel) GetByID(id int64) (*Vote, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取投票失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("vote not found: %d", id)
	}

	// 转换为投票
	vote := &Vote{}
	vote.ID, _ = result["id"].(int64)
	vote.Title, _ = result["title"].(string)
	vote.Description, _ = result["description"].(string)
	vote.StartTime, _ = result["starttime"].(time.Time)
	vote.EndTime, _ = result["endtime"].(time.Time)
	vote.IsMulti, _ = result["ismulti"].(int)
	vote.MaxChoices, _ = result["maxchoices"].(int)
	vote.Status, _ = result["status"].(int)
	vote.TotalCount, _ = result["totalcount"].(int)
	vote.CreateTime, _ = result["createtime"].(time.Time)
	vote.UpdateTime, _ = result["updatetime"].(time.Time)

	return vote, nil
}

// GetAll 获取所有投票
func (m *VoteModel) GetAll(page, pageSize int) ([]*Vote, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote")
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有投票失败", "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取投票总数失败", "error", err)
		return nil, 0, err
	}

	// 转换为投票列表
	votes := make([]*Vote, 0, len(results))
	for _, result := range results {
		vote := &Vote{}
		vote.ID, _ = result["id"].(int64)
		vote.Title, _ = result["title"].(string)
		vote.Description, _ = result["description"].(string)
		vote.StartTime, _ = result["starttime"].(time.Time)
		vote.EndTime, _ = result["endtime"].(time.Time)
		vote.IsMulti, _ = result["ismulti"].(int)
		vote.MaxChoices, _ = result["maxchoices"].(int)
		vote.Status, _ = result["status"].(int)
		vote.TotalCount, _ = result["totalcount"].(int)
		vote.CreateTime, _ = result["createtime"].(time.Time)
		vote.UpdateTime, _ = result["updatetime"].(time.Time)
		votes = append(votes, vote)
	}

	return votes, total, nil
}

// Create 创建投票
func (m *VoteModel) Create(vote *Vote) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	vote.CreateTime = now
	vote.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("vote")+" (title, description, starttime, endtime, ismulti, maxchoices, status, totalcount, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		vote.Title, vote.Description, vote.StartTime, vote.EndTime, vote.IsMulti, vote.MaxChoices, vote.Status, vote.TotalCount, vote.CreateTime, vote.UpdateTime,
	)
	if err != nil {
		logger.Error("创建投票失败", "error", err)
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

// Update 更新投票
func (m *VoteModel) Update(vote *Vote) error {
	// 设置更新时间
	vote.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("vote")+" SET title = ?, description = ?, starttime = ?, endtime = ?, ismulti = ?, maxchoices = ?, status = ?, totalcount = ?, updatetime = ? WHERE id = ?",
		vote.Title, vote.Description, vote.StartTime, vote.EndTime, vote.IsMulti, vote.MaxChoices, vote.Status, vote.TotalCount, vote.UpdateTime, vote.ID,
	)
	if err != nil {
		logger.Error("更新投票失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除投票
func (m *VoteModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("vote")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除投票失败", "error", err)
		return err
	}

	return nil
}

// IncrementTotalCount 增加总投票数
func (m *VoteModel) IncrementTotalCount(id int64, count int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("vote")+" SET totalcount = totalcount + ? WHERE id = ?",
		count, id,
	)
	if err != nil {
		logger.Error("增加总投票数失败", "error", err)
		return err
	}

	return nil
}

// VoteOptionModel 投票选项模型
type VoteOptionModel struct {
	db *database.DB
}

// NewVoteOptionModel 创建投票选项模型
func NewVoteOptionModel(db *database.DB) *VoteOptionModel {
	return &VoteOptionModel{
		db: db,
	}
}

// GetByID 根据ID获取投票选项
func (m *VoteOptionModel) GetByID(id int64) (*VoteOption, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote_option")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取投票选项失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("vote option not found: %d", id)
	}

	// 转换为投票选项
	option := &VoteOption{}
	option.ID, _ = result["id"].(int64)
	option.VoteID, _ = result["voteid"].(int64)
	option.Title, _ = result["title"].(string)
	option.Count, _ = result["count"].(int)
	option.OrderID, _ = result["orderid"].(int)

	return option, nil
}

// GetByVoteID 根据投票ID获取投票选项
func (m *VoteOptionModel) GetByVoteID(voteID int64) ([]*VoteOption, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote_option")
	qb.Where("voteid = ?", voteID)
	qb.OrderBy("orderid ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取投票选项失败", "voteid", voteID, "error", err)
		return nil, err
	}

	// 转换为投票选项列表
	options := make([]*VoteOption, 0, len(results))
	for _, result := range results {
		option := &VoteOption{}
		option.ID, _ = result["id"].(int64)
		option.VoteID, _ = result["voteid"].(int64)
		option.Title, _ = result["title"].(string)
		option.Count, _ = result["count"].(int)
		option.OrderID, _ = result["orderid"].(int)
		options = append(options, option)
	}

	return options, nil
}

// Create 创建投票选项
func (m *VoteOptionModel) Create(option *VoteOption) (int64, error) {
	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("vote_option")+" (voteid, title, count, orderid) VALUES (?, ?, ?, ?)",
		option.VoteID, option.Title, option.Count, option.OrderID,
	)
	if err != nil {
		logger.Error("创建投票选项失败", "error", err)
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

// Update 更新投票选项
func (m *VoteOptionModel) Update(option *VoteOption) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("vote_option")+" SET title = ?, count = ?, orderid = ? WHERE id = ?",
		option.Title, option.Count, option.OrderID, option.ID,
	)
	if err != nil {
		logger.Error("更新投票选项失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除投票选项
func (m *VoteOptionModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("vote_option")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除投票选项失败", "error", err)
		return err
	}

	return nil
}

// DeleteByVoteID 根据投票ID删除投票选项
func (m *VoteOptionModel) DeleteByVoteID(voteID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("vote_option")+" WHERE voteid = ?", voteID)
	if err != nil {
		logger.Error("删除投票选项失败", "voteid", voteID, "error", err)
		return err
	}

	return nil
}

// IncrementCount 增加投票数
func (m *VoteOptionModel) IncrementCount(id int64) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("vote_option")+" SET count = count + 1 WHERE id = ?",
		id,
	)
	if err != nil {
		logger.Error("增加投票数失败", "error", err)
		return err
	}

	return nil
}

// VoteLogModel 投票日志模型
type VoteLogModel struct {
	db *database.DB
}

// NewVoteLogModel 创建投票日志模型
func NewVoteLogModel(db *database.DB) *VoteLogModel {
	return &VoteLogModel{
		db: db,
	}
}

// GetByID 根据ID获取投票日志
func (m *VoteLogModel) GetByID(id int64) (*VoteLog, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote_log")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取投票日志失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("vote log not found: %d", id)
	}

	// 转换为投票日志
	log := &VoteLog{}
	log.ID, _ = result["id"].(int64)
	log.VoteID, _ = result["voteid"].(int64)
	log.OptionID, _ = result["optionid"].(int64)
	log.MemberID, _ = result["memberid"].(int64)
	log.IP, _ = result["ip"].(string)
	log.CreateTime, _ = result["createtime"].(time.Time)

	return log, nil
}

// GetByVoteID 根据投票ID获取投票日志
func (m *VoteLogModel) GetByVoteID(voteID int64, page, pageSize int) ([]*VoteLog, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote_log")
	qb.Where("voteid = ?", voteID)
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取投票日志失败", "voteid", voteID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取投票日志总数失败", "voteid", voteID, "error", err)
		return nil, 0, err
	}

	// 转换为投票日志列表
	logs := make([]*VoteLog, 0, len(results))
	for _, result := range results {
		log := &VoteLog{}
		log.ID, _ = result["id"].(int64)
		log.VoteID, _ = result["voteid"].(int64)
		log.OptionID, _ = result["optionid"].(int64)
		log.MemberID, _ = result["memberid"].(int64)
		log.IP, _ = result["ip"].(string)
		log.CreateTime, _ = result["createtime"].(time.Time)
		logs = append(logs, log)
	}

	return logs, total, nil
}

// Create 创建投票日志
func (m *VoteLogModel) Create(log *VoteLog) (int64, error) {
	// 设置创建时间
	log.CreateTime = time.Now()

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("vote_log")+" (voteid, optionid, memberid, ip, createtime) VALUES (?, ?, ?, ?, ?)",
		log.VoteID, log.OptionID, log.MemberID, log.IP, log.CreateTime,
	)
	if err != nil {
		logger.Error("创建投票日志失败", "error", err)
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

// Delete 删除投票日志
func (m *VoteLogModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("vote_log")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除投票日志失败", "error", err)
		return err
	}

	return nil
}

// DeleteByVoteID 根据投票ID删除投票日志
func (m *VoteLogModel) DeleteByVoteID(voteID int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("vote_log")+" WHERE voteid = ?", voteID)
	if err != nil {
		logger.Error("删除投票日志失败", "voteid", voteID, "error", err)
		return err
	}

	return nil
}

// CheckVoted 检查是否已投票
func (m *VoteLogModel) CheckVoted(voteID int64, memberID int64, ip string) (bool, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "vote_log")
	qb.Where("voteid = ?", voteID)
	if memberID > 0 {
		qb.Where("memberid = ?", memberID)
	} else {
		qb.Where("ip = ?", ip)
	}

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("检查是否已投票失败", "error", err)
		return false, err
	}

	return count > 0, nil
}
