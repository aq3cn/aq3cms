package model

import (
	"encoding/json"
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// PaymentMethod 支付方式
type PaymentMethod struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`        // 支付方式名称
	Code        string    `json:"code"`        // 支付方式代码
	Description string    `json:"description"` // 描述
	Config      string    `json:"config"`      // 配置，JSON格式
	Icon        string    `json:"icon"`        // 图标
	OrderID     int       `json:"orderid"`     // 排序ID
	Status      int       `json:"status"`      // 状态：0禁用，1启用
	CreateTime  time.Time `json:"createtime"`  // 创建时间
	UpdateTime  time.Time `json:"updatetime"`  // 更新时间
}

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID             int64     `json:"id"`
	OrderNo        string    `json:"orderno"`        // 订单号
	MemberID       int64     `json:"memberid"`       // 会员ID
	Amount         float64   `json:"amount"`         // 金额
	PaymentMethod  string    `json:"paymentmethod"`  // 支付方式
	PaymentOrderNo string    `json:"paymentorderno"` // 支付平台订单号
	Status         int       `json:"status"`         // 状态：0未支付，1已支付，2已取消，3已退款
	Type           int       `json:"type"`           // 类型：0充值，1购买，2其他
	RelatedID      int64     `json:"relatedid"`      // 关联ID
	RelatedType    string    `json:"relatedtype"`    // 关联类型
	Remark         string    `json:"remark"`         // 备注
	IP             string    `json:"ip"`             // IP地址
	CreateTime     time.Time `json:"createtime"`     // 创建时间
	UpdateTime     time.Time `json:"updatetime"`     // 更新时间
	PayTime        time.Time `json:"paytime"`        // 支付时间
	ExtraData      string    `json:"extradata"`      // 额外数据，JSON格式
}

// PaymentMethodModel 支付方式模型
type PaymentMethodModel struct {
	db *database.DB
}

// NewPaymentMethodModel 创建支付方式模型
func NewPaymentMethodModel(db *database.DB) *PaymentMethodModel {
	return &PaymentMethodModel{
		db: db,
	}
}

// GetByID 根据ID获取支付方式
func (m *PaymentMethodModel) GetByID(id int64) (*PaymentMethod, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_method")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取支付方式失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("payment method not found: %d", id)
	}

	// 转换为支付方式
	method := &PaymentMethod{}
	method.ID, _ = result["id"].(int64)
	method.Name, _ = result["name"].(string)
	method.Code, _ = result["code"].(string)
	method.Description, _ = result["description"].(string)
	method.Config, _ = result["config"].(string)
	method.Icon, _ = result["icon"].(string)
	method.OrderID, _ = result["orderid"].(int)
	method.Status, _ = result["status"].(int)
	method.CreateTime, _ = result["createtime"].(time.Time)
	method.UpdateTime, _ = result["updatetime"].(time.Time)

	return method, nil
}

// GetByCode 根据代码获取支付方式
func (m *PaymentMethodModel) GetByCode(code string) (*PaymentMethod, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_method")
	qb.Where("code = ?", code)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取支付方式失败", "code", code, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("payment method not found: %s", code)
	}

	// 转换为支付方式
	method := &PaymentMethod{}
	method.ID, _ = result["id"].(int64)
	method.Name, _ = result["name"].(string)
	method.Code, _ = result["code"].(string)
	method.Description, _ = result["description"].(string)
	method.Config, _ = result["config"].(string)
	method.Icon, _ = result["icon"].(string)
	method.OrderID, _ = result["orderid"].(int)
	method.Status, _ = result["status"].(int)
	method.CreateTime, _ = result["createtime"].(time.Time)
	method.UpdateTime, _ = result["updatetime"].(time.Time)

	return method, nil
}

// GetAll 获取所有支付方式
func (m *PaymentMethodModel) GetAll(status int) ([]*PaymentMethod, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_method")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("orderid ASC, id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有支付方式失败", "error", err)
		return nil, err
	}

	// 转换为支付方式列表
	methods := make([]*PaymentMethod, 0, len(results))
	for _, result := range results {
		method := &PaymentMethod{}
		method.ID, _ = result["id"].(int64)
		method.Name, _ = result["name"].(string)
		method.Code, _ = result["code"].(string)
		method.Description, _ = result["description"].(string)
		method.Config, _ = result["config"].(string)
		method.Icon, _ = result["icon"].(string)
		method.OrderID, _ = result["orderid"].(int)
		method.Status, _ = result["status"].(int)
		method.CreateTime, _ = result["createtime"].(time.Time)
		method.UpdateTime, _ = result["updatetime"].(time.Time)
		methods = append(methods, method)
	}

	return methods, nil
}

// Create 创建支付方式
func (m *PaymentMethodModel) Create(method *PaymentMethod) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	method.CreateTime = now
	method.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("payment_method")+" (name, code, description, config, icon, orderid, status, createtime, updatetime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		method.Name, method.Code, method.Description, method.Config, method.Icon, method.OrderID, method.Status, method.CreateTime, method.UpdateTime,
	)
	if err != nil {
		logger.Error("创建支付方式失败", "error", err)
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

// Update 更新支付方式
func (m *PaymentMethodModel) Update(method *PaymentMethod) error {
	// 设置更新时间
	method.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_method")+" SET name = ?, code = ?, description = ?, config = ?, icon = ?, orderid = ?, status = ?, updatetime = ? WHERE id = ?",
		method.Name, method.Code, method.Description, method.Config, method.Icon, method.OrderID, method.Status, method.UpdateTime, method.ID,
	)
	if err != nil {
		logger.Error("更新支付方式失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除支付方式
func (m *PaymentMethodModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("payment_method")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除支付方式失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新支付方式状态
func (m *PaymentMethodModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_method")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新支付方式状态失败", "error", err)
		return err
	}

	return nil
}

// GetConfig 获取配置
func (m *PaymentMethodModel) GetConfig(id int64) (map[string]interface{}, error) {
	// 获取支付方式
	method, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 解析配置
	var config map[string]interface{}
	if method.Config != "" {
		err = json.Unmarshal([]byte(method.Config), &config)
		if err != nil {
			logger.Error("解析配置失败", "error", err)
			return nil, err
		}
	} else {
		config = make(map[string]interface{})
	}

	return config, nil
}

// SetConfig 设置配置
func (m *PaymentMethodModel) SetConfig(id int64, config map[string]interface{}) error {
	// 获取支付方式
	method, err := m.GetByID(id)
	if err != nil {
		return err
	}

	// 序列化配置
	configData, err := json.Marshal(config)
	if err != nil {
		logger.Error("序列化配置失败", "error", err)
		return err
	}

	// 更新配置
	method.Config = string(configData)
	return m.Update(method)
}

// PaymentOrderModel 支付订单模型
type PaymentOrderModel struct {
	db *database.DB
}

// NewPaymentOrderModel 创建支付订单模型
func NewPaymentOrderModel(db *database.DB) *PaymentOrderModel {
	return &PaymentOrderModel{
		db: db,
	}
}

// GetByID 根据ID获取支付订单
func (m *PaymentOrderModel) GetByID(id int64) (*PaymentOrder, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_order")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取支付订单失败", "id", id, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("payment order not found: %d", id)
	}

	// 转换为支付订单
	order := &PaymentOrder{}
	order.ID, _ = result["id"].(int64)
	order.OrderNo, _ = result["orderno"].(string)
	order.MemberID, _ = result["memberid"].(int64)
	order.Amount, _ = result["amount"].(float64)
	order.PaymentMethod, _ = result["paymentmethod"].(string)
	order.PaymentOrderNo, _ = result["paymentorderno"].(string)
	order.Status, _ = result["status"].(int)
	order.Type, _ = result["type"].(int)
	order.RelatedID, _ = result["relatedid"].(int64)
	order.RelatedType, _ = result["relatedtype"].(string)
	order.Remark, _ = result["remark"].(string)
	order.IP, _ = result["ip"].(string)
	order.CreateTime, _ = result["createtime"].(time.Time)
	order.UpdateTime, _ = result["updatetime"].(time.Time)
	order.PayTime, _ = result["paytime"].(time.Time)
	order.ExtraData, _ = result["extradata"].(string)

	return order, nil
}

// GetByOrderNo 根据订单号获取支付订单
func (m *PaymentOrderModel) GetByOrderNo(orderNo string) (*PaymentOrder, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_order")
	qb.Where("orderno = ?", orderNo)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("获取支付订单失败", "orderno", orderNo, "error", err)
		return nil, err
	}

	// 检查结果
	if result == nil {
		return nil, fmt.Errorf("payment order not found: %s", orderNo)
	}

	// 转换为支付订单
	order := &PaymentOrder{}
	order.ID, _ = result["id"].(int64)
	order.OrderNo, _ = result["orderno"].(string)
	order.MemberID, _ = result["memberid"].(int64)
	order.Amount, _ = result["amount"].(float64)
	order.PaymentMethod, _ = result["paymentmethod"].(string)
	order.PaymentOrderNo, _ = result["paymentorderno"].(string)
	order.Status, _ = result["status"].(int)
	order.Type, _ = result["type"].(int)
	order.RelatedID, _ = result["relatedid"].(int64)
	order.RelatedType, _ = result["relatedtype"].(string)
	order.Remark, _ = result["remark"].(string)
	order.IP, _ = result["ip"].(string)
	order.CreateTime, _ = result["createtime"].(time.Time)
	order.UpdateTime, _ = result["updatetime"].(time.Time)
	order.PayTime, _ = result["paytime"].(time.Time)
	order.ExtraData, _ = result["extradata"].(string)

	return order, nil
}

// GetByMemberID 根据会员ID获取支付订单
func (m *PaymentOrderModel) GetByMemberID(memberID int64, status int, page, pageSize int) ([]*PaymentOrder, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_order")
	qb.Where("memberid = ?", memberID)
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取支付订单失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取支付订单总数失败", "memberid", memberID, "error", err)
		return nil, 0, err
	}

	// 转换为支付订单列表
	orders := make([]*PaymentOrder, 0, len(results))
	for _, result := range results {
		order := &PaymentOrder{}
		order.ID, _ = result["id"].(int64)
		order.OrderNo, _ = result["orderno"].(string)
		order.MemberID, _ = result["memberid"].(int64)
		order.Amount, _ = result["amount"].(float64)
		order.PaymentMethod, _ = result["paymentmethod"].(string)
		order.PaymentOrderNo, _ = result["paymentorderno"].(string)
		order.Status, _ = result["status"].(int)
		order.Type, _ = result["type"].(int)
		order.RelatedID, _ = result["relatedid"].(int64)
		order.RelatedType, _ = result["relatedtype"].(string)
		order.Remark, _ = result["remark"].(string)
		order.IP, _ = result["ip"].(string)
		order.CreateTime, _ = result["createtime"].(time.Time)
		order.UpdateTime, _ = result["updatetime"].(time.Time)
		order.PayTime, _ = result["paytime"].(time.Time)
		order.ExtraData, _ = result["extradata"].(string)
		orders = append(orders, order)
	}

	return orders, total, nil
}

// GetAll 获取所有支付订单
func (m *PaymentOrderModel) GetAll(status int, page, pageSize int) ([]*PaymentOrder, int, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "payment_order")
	if status >= 0 {
		qb.Where("status = ?", status)
	}
	qb.OrderBy("id DESC")
	qb.Limit(pageSize)
	qb.Offset((page - 1) * pageSize)

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("获取所有支付订单失败", "error", err)
		return nil, 0, err
	}

	// 获取总数
	total, err := qb.Count()
	if err != nil {
		logger.Error("获取支付订单总数失败", "error", err)
		return nil, 0, err
	}

	// 转换为支付订单列表
	orders := make([]*PaymentOrder, 0, len(results))
	for _, result := range results {
		order := &PaymentOrder{}
		order.ID, _ = result["id"].(int64)
		order.OrderNo, _ = result["orderno"].(string)
		order.MemberID, _ = result["memberid"].(int64)
		order.Amount, _ = result["amount"].(float64)
		order.PaymentMethod, _ = result["paymentmethod"].(string)
		order.PaymentOrderNo, _ = result["paymentorderno"].(string)
		order.Status, _ = result["status"].(int)
		order.Type, _ = result["type"].(int)
		order.RelatedID, _ = result["relatedid"].(int64)
		order.RelatedType, _ = result["relatedtype"].(string)
		order.Remark, _ = result["remark"].(string)
		order.IP, _ = result["ip"].(string)
		order.CreateTime, _ = result["createtime"].(time.Time)
		order.UpdateTime, _ = result["updatetime"].(time.Time)
		order.PayTime, _ = result["paytime"].(time.Time)
		order.ExtraData, _ = result["extradata"].(string)
		orders = append(orders, order)
	}

	return orders, total, nil
}

// Create 创建支付订单
func (m *PaymentOrderModel) Create(order *PaymentOrder) (int64, error) {
	// 设置创建时间和更新时间
	now := time.Now()
	order.CreateTime = now
	order.UpdateTime = now

	// 执行插入
	result, err := m.db.Exec(
		"INSERT INTO "+m.db.TableName("payment_order")+" (orderno, memberid, amount, paymentmethod, paymentorderno, status, type, relatedid, relatedtype, remark, ip, createtime, updatetime, paytime, extradata) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		order.OrderNo, order.MemberID, order.Amount, order.PaymentMethod, order.PaymentOrderNo, order.Status, order.Type, order.RelatedID, order.RelatedType, order.Remark, order.IP, order.CreateTime, order.UpdateTime, order.PayTime, order.ExtraData,
	)
	if err != nil {
		logger.Error("创建支付订单失败", "error", err)
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

// Update 更新支付订单
func (m *PaymentOrderModel) Update(order *PaymentOrder) error {
	// 设置更新时间
	order.UpdateTime = time.Now()

	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_order")+" SET orderno = ?, memberid = ?, amount = ?, paymentmethod = ?, paymentorderno = ?, status = ?, type = ?, relatedid = ?, relatedtype = ?, remark = ?, ip = ?, updatetime = ?, paytime = ?, extradata = ? WHERE id = ?",
		order.OrderNo, order.MemberID, order.Amount, order.PaymentMethod, order.PaymentOrderNo, order.Status, order.Type, order.RelatedID, order.RelatedType, order.Remark, order.IP, order.UpdateTime, order.PayTime, order.ExtraData, order.ID,
	)
	if err != nil {
		logger.Error("更新支付订单失败", "error", err)
		return err
	}

	return nil
}

// Delete 删除支付订单
func (m *PaymentOrderModel) Delete(id int64) error {
	// 执行删除
	_, err := m.db.Exec("DELETE FROM "+m.db.TableName("payment_order")+" WHERE id = ?", id)
	if err != nil {
		logger.Error("删除支付订单失败", "error", err)
		return err
	}

	return nil
}

// UpdateStatus 更新支付订单状态
func (m *PaymentOrderModel) UpdateStatus(id int64, status int) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_order")+" SET status = ?, updatetime = ? WHERE id = ?",
		status, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新支付订单状态失败", "error", err)
		return err
	}

	return nil
}

// UpdatePaymentOrderNo 更新支付平台订单号
func (m *PaymentOrderModel) UpdatePaymentOrderNo(id int64, paymentOrderNo string) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_order")+" SET paymentorderno = ?, updatetime = ? WHERE id = ?",
		paymentOrderNo, time.Now(), id,
	)
	if err != nil {
		logger.Error("更新支付平台订单号失败", "error", err)
		return err
	}

	return nil
}

// UpdatePaid 更新支付状态为已支付
func (m *PaymentOrderModel) UpdatePaid(id int64, paymentOrderNo string) error {
	// 执行更新
	_, err := m.db.Exec(
		"UPDATE "+m.db.TableName("payment_order")+" SET status = ?, paymentorderno = ?, paytime = ?, updatetime = ? WHERE id = ?",
		1, paymentOrderNo, time.Now(), time.Now(), id,
	)
	if err != nil {
		logger.Error("更新支付状态失败", "error", err)
		return err
	}

	return nil
}

// GetExtraData 获取额外数据
func (m *PaymentOrderModel) GetExtraData(id int64) (map[string]interface{}, error) {
	// 获取支付订单
	order, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 解析额外数据
	var extraData map[string]interface{}
	if order.ExtraData != "" {
		err = json.Unmarshal([]byte(order.ExtraData), &extraData)
		if err != nil {
			logger.Error("解析额外数据失败", "error", err)
			return nil, err
		}
	} else {
		extraData = make(map[string]interface{})
	}

	return extraData, nil
}

// SetExtraData 设置额外数据
func (m *PaymentOrderModel) SetExtraData(id int64, extraData map[string]interface{}) error {
	// 获取支付订单
	order, err := m.GetByID(id)
	if err != nil {
		return err
	}

	// 序列化额外数据
	extraDataJSON, err := json.Marshal(extraData)
	if err != nil {
		logger.Error("序列化额外数据失败", "error", err)
		return err
	}

	// 更新额外数据
	order.ExtraData = string(extraDataJSON)
	return m.Update(order)
}

// GenerateOrderNo 生成订单号
func (m *PaymentOrderModel) GenerateOrderNo() string {
	// 生成订单号：当前时间戳 + 4位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := fmt.Sprintf("%04d", now.Nanosecond()%10000)
	return timestamp + random
}
