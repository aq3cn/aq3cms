package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// Comment 评论模型
type Comment struct {
	ID          int64     `json:"id"`
	AID         int64     `json:"aid"`         // 文章ID
	TypeID      int64     `json:"typeid"`      // 栏目ID
	Username    string    `json:"username"`    // 用户名
	MID         int64     `json:"mid"`         // 会员ID
	IP          string    `json:"ip"`          // IP地址
	IsCheck     int       `json:"ischeck"`     // 是否审核
	Dtime       time.Time `json:"dtime"`       // 评论时间
	Content     string    `json:"content"`     // 评论内容
	ParentID    int64     `json:"parentid"`    // 父评论ID
	Score       int       `json:"score"`       // 评分
	GoodCount   int       `json:"goodcount"`   // 点赞数
	BadCount    int       `json:"badcount"`    // 踩数
	UserFace    string    `json:"userface"`    // 用户头像
	ChannelType int       `json:"channeltype"` // 频道类型
}

// CommentModel 评论模型操作
type CommentModel struct {
	db *database.DB
}

// NewCommentModel 创建评论模型
func NewCommentModel(db *database.DB) *CommentModel {
	return &CommentModel{
		db: db,
	}
}

// GetByID 根据ID获取评论
func (m *CommentModel) GetByID(id int64) (*Comment, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询评论失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("评论不存在")
	}

	// 转换为评论对象
	comment := &Comment{}
	comment.ID, _ = result["id"].(int64)
	comment.AID, _ = result["aid"].(int64)
	comment.TypeID, _ = result["typeid"].(int64)
	comment.Username, _ = result["username"].(string)
	comment.MID, _ = result["mid"].(int64)
	comment.IP, _ = result["ip"].(string)

	// 处理整数字段
	if ischeck, ok := result["ischeck"].(int64); ok {
		comment.IsCheck = int(ischeck)
	}

	// 处理日期
	if dtime, ok := result["dtime"].(time.Time); ok {
		comment.Dtime = dtime
	}

	comment.Content, _ = result["content"].(string)
	comment.ParentID, _ = result["parentid"].(int64)

	// 处理整数字段
	if score, ok := result["score"].(int64); ok {
		comment.Score = int(score)
	}
	if goodcount, ok := result["goodcount"].(int64); ok {
		comment.GoodCount = int(goodcount)
	}
	if badcount, ok := result["badcount"].(int64); ok {
		comment.BadCount = int(badcount)
	}

	comment.UserFace, _ = result["userface"].(string)

	// 处理整数字段
	if channeltype, ok := result["channeltype"].(int64); ok {
		comment.ChannelType = int(channeltype)
	}

	return comment, nil
}

// GetListByAID 获取文章评论列表
func (m *CommentModel) GetListByAID(aid int64, page, pageSize int) ([]*Comment, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.Where("aid = ?", aid)
	qb.Where("ischeck = 1")

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询评论总数失败", "aid", aid, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("dtime DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询评论列表失败", "aid", aid, "error", err)
		return nil, 0, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}

// GetRecentComments 获取最新评论
func (m *CommentModel) GetRecentComments(limit int) ([]*Comment, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.Where("ischeck = 1")
	qb.OrderBy("dtime DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询最新评论失败", "error", err)
		return nil, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// Create 创建评论
func (m *CommentModel) Create(comment *Comment) (int64, error) {
	// 清理评论内容，防止XSS攻击
	comment.Content = security.FilterXSS(comment.Content)

	// 构建数据
	data := map[string]interface{}{
		"aid":         comment.AID,
		"typeid":      comment.TypeID,
		"username":    comment.Username,
		"mid":         comment.MID,
		"ip":          comment.IP,
		"ischeck":     comment.IsCheck,
		"dtime":       comment.Dtime,
		"content":     comment.Content,
		"parentid":    comment.ParentID,
		"score":       comment.Score,
		"goodcount":   comment.GoodCount,
		"badcount":    comment.BadCount,
		"userface":    comment.UserFace,
		"channeltype": comment.ChannelType,
	}

	// 执行插入
	qb := database.NewQueryBuilder(m.db, "feedback")
	id, err := qb.Insert(data)
	if err != nil {
		logger.Error("创建评论失败", "error", err)
		return 0, err
	}

	return id, nil
}

// Delete 删除评论
func (m *CommentModel) Delete(id int64) error {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("id = ?", id)

	// 执行删除
	_, err := qb.Delete()
	if err != nil {
		logger.Error("删除评论失败", "id", id, "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新评论状态
func (m *CommentModel) UpdateStatus(id int64, isCheck int) error {
	// 构建数据
	data := map[string]interface{}{
		"ischeck": isCheck,
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("id = ?", id)

	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("更新评论状态失败", "id", id, "error", err)
		return err
	}

	return nil
}

// UpdateGoodCount 更新点赞数
func (m *CommentModel) UpdateGoodCount(id int64, increment int) error {
	// 执行更新
	_, err := m.db.Execute("UPDATE "+m.db.TableName("feedback")+" SET goodcount = goodcount + ? WHERE id = ?", increment, id)
	if err != nil {
		logger.Error("更新评论点赞数失败", "id", id, "error", err)
		return err
	}

	return nil
}

// UpdateBadCount 更新踩数
func (m *CommentModel) UpdateBadCount(id int64, increment int) error {
	// 执行更新
	_, err := m.db.Execute("UPDATE "+m.db.TableName("feedback")+" SET badcount = badcount + ? WHERE id = ?", increment, id)
	if err != nil {
		logger.Error("更新评论踩数失败", "id", id, "error", err)
		return err
	}

	return nil
}

// GetByAID 根据文章ID获取评论
func (m *CommentModel) GetByAID(aid int64) ([]*Comment, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.Where("aid = ?", aid)
	qb.Where("ischeck = 1")
	qb.OrderBy("dtime DESC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询文章评论失败", "aid", aid, "error", err)
		return nil, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// Update 更新评论
func (m *CommentModel) Update(comment *Comment) error {
	// 清理评论内容，防止XSS攻击
	comment.Content = security.FilterXSS(comment.Content)

	// 构建数据
	data := map[string]interface{}{
		"aid":         comment.AID,
		"typeid":      comment.TypeID,
		"username":    comment.Username,
		"mid":         comment.MID,
		"ip":          comment.IP,
		"ischeck":     comment.IsCheck,
		"dtime":       comment.Dtime,
		"content":     comment.Content,
		"parentid":    comment.ParentID,
		"score":       comment.Score,
		"goodcount":   comment.GoodCount,
		"badcount":    comment.BadCount,
		"userface":    comment.UserFace,
		"channeltype": comment.ChannelType,
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("id = ?", comment.ID)

	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("更新评论失败", "id", comment.ID, "error", err)
		return err
	}

	return nil
}

// GetList 获取评论列表
func (m *CommentModel) GetList(isCheck int, keyword string, page, pageSize int) ([]*Comment, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")

	// 添加条件
	if isCheck != -1 {
		qb.Where("ischeck = ?", isCheck)
	}

	if keyword != "" {
		qb.Where("(username LIKE ? OR content LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取评论总数失败", "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("dtime DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取评论列表失败", "error", err)
		return nil, 0, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}

// GetCount 获取评论总数
func (m *CommentModel) GetCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取评论总数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetTotalCount 获取评论总数（别名方法）
func (m *CommentModel) GetTotalCount() (int, error) {
	return m.GetCount()
}

// GetPendingCount 获取待审核评论数
func (m *CommentModel) GetPendingCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("ischeck = 0") // 待审核状态

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取待审核评论数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetApprovedCount 获取已审核评论数
func (m *CommentModel) GetApprovedCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("ischeck = 1") // 已审核状态

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取已审核评论数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetRejectedCount 获取已拒绝评论数
func (m *CommentModel) GetRejectedCount() (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Where("ischeck = -1") // 已拒绝状态

	// 执行查询
	count, err := qb.Count()
	if err != nil {
		logger.Error("获取已拒绝评论数失败", "error", err)
		return 0, err
	}

	return count, nil
}

// GetLatest 获取最新评论
func (m *CommentModel) GetLatest(limit int) ([]*Comment, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.OrderBy("dtime DESC")
	qb.Limit(limit)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取最新评论失败", "error", err)
		return nil, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetByMemberID 获取会员评论
func (m *CommentModel) GetByMemberID(memberID int64, page, pageSize int) ([]*Comment, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "feedback")
	qb.Select("*")
	qb.Where("mid = ?", memberID)

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取会员评论总数失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("dtime DESC")
	qb.Limit(pageSize, offset)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取会员评论列表失败", "memberID", memberID, "error", err)
		return nil, 0, err
	}

	// 转换为评论对象
	comments := make([]*Comment, 0, len(results))
	for _, result := range results {
		comment := &Comment{}
		comment.ID, _ = result["id"].(int64)
		comment.AID, _ = result["aid"].(int64)
		comment.TypeID, _ = result["typeid"].(int64)
		comment.Username, _ = result["username"].(string)
		comment.MID, _ = result["mid"].(int64)
		comment.IP, _ = result["ip"].(string)

		// 处理整数字段
		if ischeck, ok := result["ischeck"].(int64); ok {
			comment.IsCheck = int(ischeck)
		}

		// 处理日期
		if dtime, ok := result["dtime"].(time.Time); ok {
			comment.Dtime = dtime
		}

		comment.Content, _ = result["content"].(string)
		comment.ParentID, _ = result["parentid"].(int64)

		// 处理整数字段
		if score, ok := result["score"].(int64); ok {
			comment.Score = int(score)
		}
		if goodcount, ok := result["goodcount"].(int64); ok {
			comment.GoodCount = int(goodcount)
		}
		if badcount, ok := result["badcount"].(int64); ok {
			comment.BadCount = int(badcount)
		}

		comment.UserFace, _ = result["userface"].(string)

		// 处理整数字段
		if channeltype, ok := result["channeltype"].(int64); ok {
			comment.ChannelType = int(channeltype)
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}
