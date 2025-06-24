package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// Message 消息模型
type Message struct {
	ID        int64     `json:"id"`
	FromID    int64     `json:"fromid"`    // 发送者ID
	ToID      int64     `json:"toid"`      // 接收者ID
	Title     string    `json:"title"`     // 消息标题
	Content   string    `json:"content"`   // 消息内容
	SendTime  time.Time `json:"sendtime"`  // 发送时间
	ReadTime  time.Time `json:"readtime"`  // 阅读时间
	IsRead    int       `json:"isread"`    // 是否已读
	FromDel   int       `json:"fromdel"`   // 发送者是否删除
	ToDel     int       `json:"todel"`     // 接收者是否删除
	FromName  string    `json:"fromname"`  // 发送者名称
	ToName    string    `json:"toname"`    // 接收者名称
}

// MessageModel 消息模型操作
type MessageModel struct {
	db *database.DB
}

// NewMessageModel 创建消息模型
func NewMessageModel(db *database.DB) *MessageModel {
	return &MessageModel{
		db: db,
	}
}

// GetByID 根据ID获取消息
func (m *MessageModel) GetByID(id int64) (*Message, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Select("m.*, f.uname as fromname, t.uname as toname")
	qb.LeftJoin(m.db.TableName("member")+" AS f", "m.fromid = f.mid")
	qb.LeftJoin(m.db.TableName("member")+" AS t", "m.toid = t.mid")
	qb.Where("m.id = ?", id)
	
	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询消息失败", "id", id, "error", err)
		return nil, err
	}
	
	if result == nil {
		return nil, fmt.Errorf("消息不存在")
	}
	
	// 转换为消息对象
	message := &Message{}
	message.ID, _ = result["id"].(int64)
	message.FromID, _ = result["fromid"].(int64)
	message.ToID, _ = result["toid"].(int64)
	message.Title, _ = result["title"].(string)
	message.Content, _ = result["content"].(string)
	
	// 处理日期
	if sendtime, ok := result["sendtime"].(time.Time); ok {
		message.SendTime = sendtime
	}
	if readtime, ok := result["readtime"].(time.Time); ok {
		message.ReadTime = readtime
	}
	
	// 处理整数字段
	if isread, ok := result["isread"].(int64); ok {
		message.IsRead = int(isread)
	}
	if fromdel, ok := result["fromdel"].(int64); ok {
		message.FromDel = int(fromdel)
	}
	if todel, ok := result["todel"].(int64); ok {
		message.ToDel = int(todel)
	}
	
	message.FromName, _ = result["fromname"].(string)
	message.ToName, _ = result["toname"].(string)
	
	return message, nil
}

// GetInbox 获取收件箱
func (m *MessageModel) GetInbox(memberID int64, page, pageSize int) ([]*Message, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Select("m.*, f.uname as fromname, t.uname as toname")
	qb.LeftJoin(m.db.TableName("member")+" AS f", "m.fromid = f.mid")
	qb.LeftJoin(m.db.TableName("member")+" AS t", "m.toid = t.mid")
	qb.Where("m.toid = ?", memberID)
	qb.Where("m.todel = 0")
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询收件箱总数失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}
	
	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("m.sendtime DESC")
	qb.Limit(pageSize, offset)
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询收件箱失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}
	
	// 转换为消息对象
	messages := make([]*Message, 0, len(results))
	for _, result := range results {
		message := &Message{}
		message.ID, _ = result["id"].(int64)
		message.FromID, _ = result["fromid"].(int64)
		message.ToID, _ = result["toid"].(int64)
		message.Title, _ = result["title"].(string)
		message.Content, _ = result["content"].(string)
		
		// 处理日期
		if sendtime, ok := result["sendtime"].(time.Time); ok {
			message.SendTime = sendtime
		}
		if readtime, ok := result["readtime"].(time.Time); ok {
			message.ReadTime = readtime
		}
		
		// 处理整数字段
		if isread, ok := result["isread"].(int64); ok {
			message.IsRead = int(isread)
		}
		if fromdel, ok := result["fromdel"].(int64); ok {
			message.FromDel = int(fromdel)
		}
		if todel, ok := result["todel"].(int64); ok {
			message.ToDel = int(todel)
		}
		
		message.FromName, _ = result["fromname"].(string)
		message.ToName, _ = result["toname"].(string)
		
		messages = append(messages, message)
	}
	
	return messages, total, nil
}

// GetOutbox 获取发件箱
func (m *MessageModel) GetOutbox(memberID int64, page, pageSize int) ([]*Message, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Select("m.*, f.uname as fromname, t.uname as toname")
	qb.LeftJoin(m.db.TableName("member")+" AS f", "m.fromid = f.mid")
	qb.LeftJoin(m.db.TableName("member")+" AS t", "m.toid = t.mid")
	qb.Where("m.fromid = ?", memberID)
	qb.Where("m.fromdel = 0")
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询发件箱总数失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}
	
	// 设置分页
	offset := (page - 1) * pageSize
	qb.OrderBy("m.sendtime DESC")
	qb.Limit(pageSize, offset)
	
	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询发件箱失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}
	
	// 转换为消息对象
	messages := make([]*Message, 0, len(results))
	for _, result := range results {
		message := &Message{}
		message.ID, _ = result["id"].(int64)
		message.FromID, _ = result["fromid"].(int64)
		message.ToID, _ = result["toid"].(int64)
		message.Title, _ = result["title"].(string)
		message.Content, _ = result["content"].(string)
		
		// 处理日期
		if sendtime, ok := result["sendtime"].(time.Time); ok {
			message.SendTime = sendtime
		}
		if readtime, ok := result["readtime"].(time.Time); ok {
			message.ReadTime = readtime
		}
		
		// 处理整数字段
		if isread, ok := result["isread"].(int64); ok {
			message.IsRead = int(isread)
		}
		if fromdel, ok := result["fromdel"].(int64); ok {
			message.FromDel = int(fromdel)
		}
		if todel, ok := result["todel"].(int64); ok {
			message.ToDel = int(todel)
		}
		
		message.FromName, _ = result["fromname"].(string)
		message.ToName, _ = result["toname"].(string)
		
		messages = append(messages, message)
	}
	
	return messages, total, nil
}

// GetUnreadCount 获取未读消息数量
func (m *MessageModel) GetUnreadCount(memberID int64) (int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Where("toid = ?", memberID)
	qb.Where("isread = 0")
	qb.Where("todel = 0")
	
	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("查询未读消息数量失败", "memberid", memberID, "error", err)
		return 0, err
	}
	
	return total, nil
}

// Create 创建消息
func (m *MessageModel) Create(message *Message) (int64, error) {
	// 清理消息内容，防止XSS攻击
	message.Title = security.CleanHTML(message.Title)
	message.Content = security.CleanHTML(message.Content)
	
	// 设置默认值
	if message.SendTime.IsZero() {
		message.SendTime = time.Now()
	}
	
	// 构建数据
	data := map[string]interface{}{
		"fromid":   message.FromID,
		"toid":     message.ToID,
		"title":    message.Title,
		"content":  message.Content,
		"sendtime": message.SendTime,
		"isread":   message.IsRead,
		"fromdel":  message.FromDel,
		"todel":    message.ToDel,
	}
	
	// 执行插入
	qb := database.NewQueryBuilder(m.db, "member_msg")
	id, err := qb.Insert(data)
	if err != nil {
		logger.Error("创建消息失败", "error", err)
		return 0, err
	}
	
	return id, nil
}

// MarkAsRead 标记为已读
func (m *MessageModel) MarkAsRead(id, memberID int64) error {
	// 构建数据
	data := map[string]interface{}{
		"isread":   1,
		"readtime": time.Now(),
	}
	
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Where("id = ?", id)
	qb.Where("toid = ?", memberID)
	
	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("标记消息为已读失败", "id", id, "memberid", memberID, "error", err)
		return err
	}
	
	return nil
}

// DeleteFromInbox 从收件箱删除
func (m *MessageModel) DeleteFromInbox(id, memberID int64) error {
	// 构建数据
	data := map[string]interface{}{
		"todel": 1,
	}
	
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Where("id = ?", id)
	qb.Where("toid = ?", memberID)
	
	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("从收件箱删除消息失败", "id", id, "memberid", memberID, "error", err)
		return err
	}
	
	return nil
}

// DeleteFromOutbox 从发件箱删除
func (m *MessageModel) DeleteFromOutbox(id, memberID int64) error {
	// 构建数据
	data := map[string]interface{}{
		"fromdel": 1,
	}
	
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Where("id = ?", id)
	qb.Where("fromid = ?", memberID)
	
	// 执行更新
	_, err := qb.Update(data)
	if err != nil {
		logger.Error("从发件箱删除消息失败", "id", id, "memberid", memberID, "error", err)
		return err
	}
	
	return nil
}

// DeleteCompletely 完全删除消息
func (m *MessageModel) DeleteCompletely(id int64) error {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "member_msg")
	qb.Where("id = ?", id)
	qb.Where("fromdel = 1")
	qb.Where("todel = 1")
	
	// 执行删除
	_, err := qb.Delete()
	if err != nil {
		logger.Error("完全删除消息失败", "id", id, "error", err)
		return err
	}
	
	return nil
}
