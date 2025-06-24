package service

import (
	"encoding/json"
	"fmt"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// PaymentService 支付服务
type PaymentService struct {
	db                 *database.DB
	cache              cache.Cache
	config             *config.Config
	paymentMethodModel *model.PaymentMethodModel
	paymentOrderModel  *model.PaymentOrderModel
	memberModel        *model.MemberModel
}

// NewPaymentService 创建支付服务
func NewPaymentService(db *database.DB, cache cache.Cache, config *config.Config) *PaymentService {
	return &PaymentService{
		db:                 db,
		cache:              cache,
		config:             config,
		paymentMethodModel: model.NewPaymentMethodModel(db),
		paymentOrderModel:  model.NewPaymentOrderModel(db),
		memberModel:        model.NewMemberModel(db),
	}
}

// GetPaymentMethods 获取支付方式
func (s *PaymentService) GetPaymentMethods() ([]*model.PaymentMethod, error) {
	// 从缓存获取
	cacheKey := "payment:methods"
	if cached, ok := s.cache.Get(cacheKey); ok {
		if methods, ok := cached.([]*model.PaymentMethod); ok {
			return methods, nil
		}
	}

	// 获取支付方式
	methods, err := s.paymentMethodModel.GetAll(1)
	if err != nil {
		return nil, err
	}

	// 缓存支付方式
	cache.SafeSet(s.cache, cacheKey, methods, time.Hour)

	return methods, nil
}

// GetPaymentMethod 获取支付方式
func (s *PaymentService) GetPaymentMethod(code string) (*model.PaymentMethod, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("payment:method:%s", code)
	if cached, ok := s.cache.Get(cacheKey); ok {
		if method, ok := cached.(*model.PaymentMethod); ok {
			return method, nil
		}
	}

	// 获取支付方式
	method, err := s.paymentMethodModel.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// 缓存支付方式
	cache.SafeSet(s.cache, cacheKey, method, time.Hour)

	return method, nil
}

// CreateOrder 创建订单
func (s *PaymentService) CreateOrder(memberID int64, amount float64, paymentMethod string, orderType int, relatedID int64, relatedType string, remark string, ip string, extraData map[string]interface{}) (*model.PaymentOrder, error) {
	// 获取支付方式
	method, err := s.GetPaymentMethod(paymentMethod)
	if err != nil {
		return nil, err
	}

	// 检查支付方式状态
	if method.Status != 1 {
		return nil, fmt.Errorf("payment method is disabled")
	}

	// 序列化额外数据
	extraDataJSON := ""
	if extraData != nil {
		extraDataBytes, err := json.Marshal(extraData)
		if err != nil {
			logger.Error("序列化额外数据失败", "error", err)
			return nil, err
		}
		extraDataJSON = string(extraDataBytes)
	}

	// 创建订单
	order := &model.PaymentOrder{
		OrderNo:       s.paymentOrderModel.GenerateOrderNo(),
		MemberID:      memberID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        0,
		Type:          orderType,
		RelatedID:     relatedID,
		RelatedType:   relatedType,
		Remark:        remark,
		IP:            ip,
		ExtraData:     extraDataJSON,
	}

	// 保存订单
	id, err := s.paymentOrderModel.Create(order)
	if err != nil {
		return nil, err
	}

	// 获取订单
	order, err = s.paymentOrderModel.GetByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder 获取订单
func (s *PaymentService) GetOrder(orderNo string) (*model.PaymentOrder, error) {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *PaymentService) UpdateOrderStatus(orderNo string, status int) error {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 更新订单状态
	err = s.paymentOrderModel.UpdateStatus(order.ID, status)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrderPaid 更新订单为已支付
func (s *PaymentService) UpdateOrderPaid(orderNo string, paymentOrderNo string) error {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != 0 {
		return fmt.Errorf("order status is not unpaid")
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 更新订单为已支付
	err = s.paymentOrderModel.UpdatePaid(order.ID, paymentOrderNo)
	if err != nil {
		return err
	}

	// 处理订单类型
	switch order.Type {
	case 0: // 充值
		// 更新会员余额
		err = s.memberModel.UpdateMoney(order.MemberID, order.Amount)
		if err != nil {
			return err
		}
	case 1: // 购买
		// 处理购买逻辑
		// 这里需要根据RelatedType和RelatedID处理不同的购买逻辑
		// 暂时不实现
	case 2: // 其他
		// 处理其他逻辑
		// 暂时不实现
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// CancelOrder 取消订单
func (s *PaymentService) CancelOrder(orderNo string) error {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != 0 {
		return fmt.Errorf("order status is not unpaid")
	}

	// 更新订单状态为已取消
	err = s.paymentOrderModel.UpdateStatus(order.ID, 2)
	if err != nil {
		return err
	}

	return nil
}

// RefundOrder 退款订单
func (s *PaymentService) RefundOrder(orderNo string) error {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != 1 {
		return fmt.Errorf("order status is not paid")
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 更新订单状态为已退款
	err = s.paymentOrderModel.UpdateStatus(order.ID, 3)
	if err != nil {
		return err
	}

	// 处理订单类型
	switch order.Type {
	case 0: // 充值
		// 更新会员余额
		err = s.memberModel.UpdateMoney(order.MemberID, -order.Amount)
		if err != nil {
			return err
		}
	case 1: // 购买
		// 处理退款逻辑
		// 这里需要根据RelatedType和RelatedID处理不同的退款逻辑
		// 暂时不实现
	case 2: // 其他
		// 处理其他逻辑
		// 暂时不实现
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// GetPaymentURL 获取支付URL
func (s *PaymentService) GetPaymentURL(orderNo string) (string, error) {
	// 获取订单
	order, err := s.paymentOrderModel.GetByOrderNo(orderNo)
	if err != nil {
		return "", err
	}

	// 检查订单状态
	if order.Status != 0 {
		return "", fmt.Errorf("order status is not unpaid")
	}

	// 获取支付方式
	method, err := s.GetPaymentMethod(order.PaymentMethod)
	if err != nil {
		return "", err
	}

	// 获取支付方式配置
	_, err = s.paymentMethodModel.GetConfig(method.ID)
	if err != nil {
		return "", err
	}

	// 根据支付方式生成支付URL
	switch method.Code {
	case "alipay":
		// 生成支付宝支付URL
		// 这里需要集成支付宝SDK
		// 暂时返回一个示例URL
		return fmt.Sprintf("/payment/alipay?orderno=%s", orderNo), nil
	case "wechat":
		// 生成微信支付URL
		// 这里需要集成微信支付SDK
		// 暂时返回一个示例URL
		return fmt.Sprintf("/payment/wechat?orderno=%s", orderNo), nil
	case "paypal":
		// 生成PayPal支付URL
		// 这里需要集成PayPal SDK
		// 暂时返回一个示例URL
		return fmt.Sprintf("/payment/paypal?orderno=%s", orderNo), nil
	default:
		return "", fmt.Errorf("unsupported payment method: %s", method.Code)
	}
}

// InitDefaultMethods 初始化默认支付方式
func (s *PaymentService) InitDefaultMethods() error {
	// 默认支付方式
	defaultMethods := []*model.PaymentMethod{
		{
			Name:        "支付宝",
			Code:        "alipay",
			Description: "支付宝支付",
			Config:      `{"appId":"","privateKey":"","publicKey":"","gatewayUrl":"https://openapi.alipay.com/gateway.do"}`,
			Icon:        "static/images/payment/alipay.png",
			OrderID:     1,
			Status:      1,
		},
		{
			Name:        "微信支付",
			Code:        "wechat",
			Description: "微信支付",
			Config:      `{"appId":"","mchId":"","key":"","certPath":"","keyPath":""}`,
			Icon:        "static/images/payment/wechat.png",
			OrderID:     2,
			Status:      1,
		},
		{
			Name:        "PayPal",
			Code:        "paypal",
			Description: "PayPal支付",
			Config:      `{"clientId":"","secret":"","mode":"sandbox"}`,
			Icon:        "static/images/payment/paypal.png",
			OrderID:     3,
			Status:      1,
		},
	}

	// 创建默认支付方式
	for _, method := range defaultMethods {
		// 检查支付方式是否已存在
		existingMethod, err := s.paymentMethodModel.GetByCode(method.Code)
		if err == nil && existingMethod != nil {
			// 支付方式已存在，跳过
			continue
		}

		// 创建支付方式
		_, err = s.paymentMethodModel.Create(method)
		if err != nil {
			logger.Error("创建默认支付方式失败", "code", method.Code, "error", err)
			return err
		}
	}

	// 清除缓存
	s.cache.Delete("payment:methods")

	return nil
}
