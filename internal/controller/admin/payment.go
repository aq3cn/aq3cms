package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/middleware"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// PaymentController 支付控制器
type PaymentController struct {
	db                 *database.DB
	cache              cache.Cache
	config             *config.Config
	paymentService     *service.PaymentService
	paymentMethodModel *model.PaymentMethodModel
	paymentOrderModel  *model.PaymentOrderModel
	memberModel        *model.MemberModel
	templateService    *service.TemplateService
}

// NewPaymentController 创建支付控制器
func NewPaymentController(db *database.DB, cache cache.Cache, config *config.Config) *PaymentController {
	return &PaymentController{
		db:                 db,
		cache:              cache,
		config:             config,
		paymentService:     service.NewPaymentService(db, cache, config),
		paymentMethodModel: model.NewPaymentMethodModel(db),
		paymentOrderModel:  model.NewPaymentOrderModel(db),
		memberModel:        model.NewMemberModel(db),
		templateService:    service.NewTemplateService(db, cache, config),
	}
}

// MethodList 支付方式列表
func (c *PaymentController) MethodList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取支付方式列表
	methods, err := c.paymentMethodModel.GetAll(-1)
	if err != nil {
		logger.Error("获取支付方式列表失败", "error", err)
		http.Error(w, "Failed to get payment methods", http.StatusInternalServerError)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Methods":     methods,
		"CurrentMenu": "payment",
		"PageTitle":   "支付方式管理",
	}

	// 渲染模板
	tplFile := "admin/payment_method_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染支付方式列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// MethodAdd 添加支付方式页面
func (c *PaymentController) MethodAdd(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"CurrentMenu": "payment",
		"PageTitle":   "添加支付方式",
	}

	// 渲染模板
	tplFile := "admin/payment_method_add.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染添加支付方式模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// MethodDoAdd 处理添加支付方式
func (c *PaymentController) MethodDoAdd(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	config := r.FormValue("config")
	icon := r.FormValue("icon")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)

	// 创建支付方式
	method := &model.PaymentMethod{
		Name:        name,
		Code:        code,
		Description: description,
		Config:      config,
		Icon:        icon,
		OrderID:     orderID,
		Status:      status,
	}

	// 保存支付方式
	id, err := c.paymentMethodModel.Create(method)
	if err != nil {
		logger.Error("创建支付方式失败", "error", err)
		http.Error(w, "Failed to create payment method", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("payment:methods")
	c.cache.Delete("payment:method:" + code)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付方式创建成功",
			"id":      id,
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_method_list", http.StatusFound)
	}
}

// MethodEdit 编辑支付方式页面
func (c *PaymentController) MethodEdit(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取支付方式ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment method ID", http.StatusBadRequest)
		return
	}

	// 获取支付方式
	method, err := c.paymentMethodModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付方式失败", "id", id, "error", err)
		http.Error(w, "Payment method not found", http.StatusNotFound)
		return
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Method":      method,
		"CurrentMenu": "payment",
		"PageTitle":   "编辑支付方式",
	}

	// 渲染模板
	tplFile := "admin/payment_method_edit.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染编辑支付方式模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// MethodDoEdit 处理编辑支付方式
func (c *PaymentController) MethodDoEdit(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// 获取支付方式ID
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment method ID", http.StatusBadRequest)
		return
	}

	// 获取原支付方式
	method, err := c.paymentMethodModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付方式失败", "id", id, "error", err)
		http.Error(w, "Payment method not found", http.StatusNotFound)
		return
	}

	// 获取表单数据
	name := r.FormValue("name")
	code := r.FormValue("code")
	description := r.FormValue("description")
	config := r.FormValue("config")
	icon := r.FormValue("icon")
	orderIDStr := r.FormValue("orderid")
	statusStr := r.FormValue("status")

	// 验证必填字段
	if name == "" || code == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 解析数值
	orderID, _ := strconv.Atoi(orderIDStr)
	status, _ := strconv.Atoi(statusStr)

	// 更新支付方式
	method.Name = name
	method.Code = code
	method.Description = description
	method.Config = config
	method.Icon = icon
	method.OrderID = orderID
	method.Status = status

	// 保存支付方式
	err = c.paymentMethodModel.Update(method)
	if err != nil {
		logger.Error("更新支付方式失败", "error", err)
		http.Error(w, "Failed to update payment method", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("payment:methods")
	c.cache.Delete("payment:method:" + code)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付方式更新成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_method_list", http.StatusFound)
	}
}

// MethodDelete 删除支付方式
func (c *PaymentController) MethodDelete(w http.ResponseWriter, r *http.Request) {
	// 获取支付方式ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment method ID", http.StatusBadRequest)
		return
	}

	// 获取支付方式
	method, err := c.paymentMethodModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付方式失败", "id", id, "error", err)
		http.Error(w, "Payment method not found", http.StatusNotFound)
		return
	}

	// 删除支付方式
	err = c.paymentMethodModel.Delete(id)
	if err != nil {
		logger.Error("删除支付方式失败", "id", id, "error", err)
		http.Error(w, "Failed to delete payment method", http.StatusInternalServerError)
		return
	}

	// 清除缓存
	c.cache.Delete("payment:methods")
	c.cache.Delete("payment:method:" + method.Code)

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付方式删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_method_list", http.StatusFound)
	}
}

// OrderList 支付订单列表
func (c *PaymentController) OrderList(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取查询参数
	statusStr := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")

	// 解析参数
	status := -1
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	// 获取支付订单列表
	orders, total, err := c.paymentOrderModel.GetAll(status, page, 20)
	if err != nil {
		logger.Error("获取支付订单列表失败", "error", err)
		http.Error(w, "Failed to get payment orders", http.StatusInternalServerError)
		return
	}

	// 计算分页信息
	totalPages := (total + 20 - 1) / 20
	pagination := map[string]interface{}{
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"TotalItems":  total,
		"HasPrev":     page > 1,
		"HasNext":     page < totalPages,
		"PrevPage":    page - 1,
		"NextPage":    page + 1,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Orders":      orders,
		"Status":      status,
		"Pagination":  pagination,
		"CurrentMenu": "payment",
		"PageTitle":   "支付订单管理",
	}

	// 渲染模板
	tplFile := "admin/payment_order_list.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染支付订单列表模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// OrderDetail 支付订单详情
func (c *PaymentController) OrderDetail(w http.ResponseWriter, r *http.Request) {
	// 获取管理员信息
	adminID := middleware.GetAdminID(r)
	adminName := middleware.GetAdminName(r)

	// 获取支付订单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment order ID", http.StatusBadRequest)
		return
	}

	// 获取支付订单
	order, err := c.paymentOrderModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付订单失败", "id", id, "error", err)
		http.Error(w, "Payment order not found", http.StatusNotFound)
		return
	}

	// 获取会员
	var member *model.Member
	if order.MemberID > 0 {
		member, err = c.memberModel.GetByID(order.MemberID)
		if err != nil {
			logger.Error("获取会员失败", "id", order.MemberID, "error", err)
		}
	}

	// 获取额外数据
	extraData, err := c.paymentOrderModel.GetExtraData(id)
	if err != nil {
		logger.Error("获取额外数据失败", "error", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"AdminID":     adminID,
		"AdminName":   adminName,
		"Order":       order,
		"Member":      member,
		"ExtraData":   extraData,
		"CurrentMenu": "payment",
		"PageTitle":   "支付订单详情",
	}

	// 渲染模板
	tplFile := "admin/payment_order_detail.htm"
	if err := c.templateService.Render(w, tplFile, data); err != nil {
		logger.Error("渲染支付订单详情模板失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// OrderCancel 取消支付订单
func (c *PaymentController) OrderCancel(w http.ResponseWriter, r *http.Request) {
	// 获取支付订单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment order ID", http.StatusBadRequest)
		return
	}

	// 获取支付订单
	order, err := c.paymentOrderModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付订单失败", "id", id, "error", err)
		http.Error(w, "Payment order not found", http.StatusNotFound)
		return
	}

	// 取消支付订单
	err = c.paymentService.CancelOrder(order.OrderNo)
	if err != nil {
		logger.Error("取消支付订单失败", "id", id, "error", err)
		http.Error(w, "Failed to cancel payment order", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付订单取消成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_order_list", http.StatusFound)
	}
}

// OrderRefund 退款支付订单
func (c *PaymentController) OrderRefund(w http.ResponseWriter, r *http.Request) {
	// 获取支付订单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment order ID", http.StatusBadRequest)
		return
	}

	// 获取支付订单
	order, err := c.paymentOrderModel.GetByID(id)
	if err != nil {
		logger.Error("获取支付订单失败", "id", id, "error", err)
		http.Error(w, "Payment order not found", http.StatusNotFound)
		return
	}

	// 退款支付订单
	err = c.paymentService.RefundOrder(order.OrderNo)
	if err != nil {
		logger.Error("退款支付订单失败", "id", id, "error", err)
		http.Error(w, "Failed to refund payment order", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付订单退款成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_order_list", http.StatusFound)
	}
}

// OrderDelete 删除支付订单
func (c *PaymentController) OrderDelete(w http.ResponseWriter, r *http.Request) {
	// 获取支付订单ID
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment order ID", http.StatusBadRequest)
		return
	}

	// 删除支付订单
	err = c.paymentOrderModel.Delete(id)
	if err != nil {
		logger.Error("删除支付订单失败", "id", id, "error", err)
		http.Error(w, "Failed to delete payment order", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付订单删除成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_order_list", http.StatusFound)
	}
}

// InitMethods 初始化支付方式
func (c *PaymentController) InitMethods(w http.ResponseWriter, r *http.Request) {
	// 初始化默认支付方式
	err := c.paymentService.InitDefaultMethods()
	if err != nil {
		logger.Error("初始化支付方式失败", "error", err)
		http.Error(w, "Failed to initialize payment methods", http.StatusInternalServerError)
		return
	}

	// 返回成功信息
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		// AJAX请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "支付方式初始化成功",
		})
	} else {
		// 普通表单提交
		http.Redirect(w, r, "/admin/payment_method_list", http.StatusFound)
	}
}
