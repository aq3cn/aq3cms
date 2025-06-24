package model

import (
	"fmt"
	"time"

	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// Admin 管理员模型
type Admin struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	RealName   string    `json:"realname"`
	Rank       int       `json:"rank"`
	Status     int       `json:"status"`
	LastLogin  time.Time `json:"lastlogin"`
	LastIP     string    `json:"lastip"`
	LoginCount int       `json:"logincount"`
}

// AdminModel 管理员模型操作
type AdminModel struct {
	db *database.DB
}

// NewAdminModel 创建管理员模型
func NewAdminModel(db *database.DB) *AdminModel {
	return &AdminModel{
		db: db,
	}
}

// GetByID 根据ID获取管理员
func (m *AdminModel) GetByID(id int64) (*Admin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Select("*")
	qb.Where("id = ?", id)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询管理员失败", "id", id, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为管理员对象
	admin := &Admin{}
	admin.ID, _ = result["id"].(int64)
	admin.Username, _ = result["username"].(string)
	admin.Password, _ = result["password"].(string)
	admin.Email, _ = result["email"].(string)
	admin.RealName, _ = result["realname"].(string)

	// 处理整数字段
	if rank, ok := result["rank"].(int64); ok {
		admin.Rank = int(rank)
	}
	if status, ok := result["status"].(int64); ok {
		admin.Status = int(status)
	}

	// 处理日期
	if lastlogin, ok := result["lastlogin"].(time.Time); ok {
		admin.LastLogin = lastlogin
	}

	admin.LastIP, _ = result["lastip"].(string)

	// 处理整数字段
	if logincount, ok := result["logincount"].(int64); ok {
		admin.LoginCount = int(logincount)
	}

	return admin, nil
}

// GetByUsername 根据用户名获取管理员
func (m *AdminModel) GetByUsername(username string) (*Admin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Select("*")
	qb.Where("username = ?", username)

	// 执行查询
	result, err := qb.First()
	if err != nil {
		logger.Error("查询管理员失败", "username", username, "error", err)
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	// 转换为管理员对象
	admin := &Admin{}
	admin.ID, _ = result["id"].(int64)
	admin.Username, _ = result["username"].(string)
	admin.Password, _ = result["password"].(string)
	admin.Email, _ = result["email"].(string)
	admin.RealName, _ = result["realname"].(string)

	// 处理整数字段
	if rank, ok := result["rank"].(int64); ok {
		admin.Rank = int(rank)
	}
	if status, ok := result["status"].(int64); ok {
		admin.Status = int(status)
	}

	// 处理日期
	if lastlogin, ok := result["lastlogin"].(time.Time); ok {
		admin.LastLogin = lastlogin
	}

	admin.LastIP, _ = result["lastip"].(string)

	// 处理整数字段
	if logincount, ok := result["logincount"].(int64); ok {
		admin.LoginCount = int(logincount)
	}

	return admin, nil
}

// GetAll 获取所有管理员
func (m *AdminModel) GetAll() ([]*Admin, error) {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Select("*")
	qb.OrderBy("id ASC")

	// 执行查询
	results, err := qb.Get()
	if err != nil {
		logger.Error("查询所有管理员失败", "error", err)
		return nil, err
	}

	// 转换为管理员对象
	admins := make([]*Admin, 0, len(results))
	for _, result := range results {
		admin := &Admin{}
		admin.ID, _ = result["id"].(int64)
		admin.Username, _ = result["username"].(string)
		admin.Password, _ = result["password"].(string)
		admin.Email, _ = result["email"].(string)
		admin.RealName, _ = result["realname"].(string)

		// 处理整数字段
		if rank, ok := result["rank"].(int64); ok {
			admin.Rank = int(rank)
		}
		if status, ok := result["status"].(int64); ok {
			admin.Status = int(status)
		}

		// 处理日期
		if lastlogin, ok := result["lastlogin"].(time.Time); ok {
			admin.LastLogin = lastlogin
		}

		admin.LastIP, _ = result["lastip"].(string)

		// 处理整数字段
		if logincount, ok := result["logincount"].(int64); ok {
			admin.LoginCount = int(logincount)
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

// Create 创建管理员
func (m *AdminModel) Create(admin *Admin) (int64, error) {
	// 检查用户名是否已存在
	existingAdmin, err := m.GetByUsername(admin.Username)
	if err == nil && existingAdmin != nil {
		logger.Error("用户名已存在")
		return 0, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	admin.Password = security.HashPassword(admin.Password)

	// 构建数据
	data := map[string]interface{}{
		"username":   admin.Username,
		"password":   admin.Password,
		"email":      admin.Email,
		"realname":   admin.RealName,
		"rank":       admin.Rank,
		"status":     admin.Status,
		"lastlogin":  admin.LastLogin,
		"lastip":     admin.LastIP,
		"logincount": admin.LoginCount,
	}

	// 执行插入
	qb := database.NewQueryBuilder(m.db, "admin")
	id, err := qb.Insert(data)
	if err != nil {
		logger.Error("创建管理员失败", "username", admin.Username, "error", err)
		return 0, err
	}

	return id, nil
}

// Update 更新管理员
func (m *AdminModel) Update(admin *Admin) error {
	// 检查用户名是否已存在
	existingAdmin, err := m.GetByUsername(admin.Username)
	if err == nil && existingAdmin != nil && existingAdmin.ID != admin.ID {
		logger.Error("用户名已存在")
		return fmt.Errorf("用户名已存在")
	}

	// 构建数据
	data := map[string]interface{}{
		"email":      admin.Email,
		"realname":   admin.RealName,
		"rank":       admin.Rank,
		"status":     admin.Status,
		"lastlogin":  admin.LastLogin,
		"lastip":     admin.LastIP,
		"logincount": admin.LoginCount,
	}

	// 如果密码不为空，则更新密码
	if admin.Password != "" {
		data["password"] = security.HashPassword(admin.Password)
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Where("id = ?", admin.ID)

	// 执行更新
	_, err = qb.Update(data)
	if err != nil {
		logger.Error("更新管理员失败", "id", admin.ID, "error", err)
		return err
	}

	return nil
}

// Delete 删除管理员
func (m *AdminModel) Delete(id int64) error {
	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Where("id = ?", id)

	// 执行删除
	_, err := qb.Delete()
	if err != nil {
		logger.Error("删除管理员失败", "id", id, "error", err)
		return err
	}

	return nil
}

// ChangePassword 修改密码
func (m *AdminModel) ChangePassword(id int64, oldPassword, newPassword string) error {
	// 获取管理员
	admin, err := m.GetByID(id)
	if err != nil {
		return err
	}

	if admin == nil {
		logger.Error("管理员不存在")
		return fmt.Errorf("管理员不存在")
	}

	// 验证旧密码
	if !security.CheckPassword(oldPassword, admin.Password) {
		logger.Error("旧密码错误")
		return fmt.Errorf("旧密码错误")
	}

	// 更新密码
	admin.Password = security.HashPassword(newPassword)

	// 构建数据
	data := map[string]interface{}{
		"password": admin.Password,
	}

	// 构建查询
	qb := database.NewQueryBuilder(m.db, "admin")
	qb.Where("id = ?", id)

	// 执行更新
	_, err = qb.Update(data)
	if err != nil {
		logger.Error("修改密码失败", "id", id, "error", err)
		return err
	}

	return nil
}

// InitDefaultAdmin 初始化默认管理员
func (m *AdminModel) InitDefaultAdmin() error {
	// 检查是否已存在管理员
	existingAdmin, err := m.GetByUsername("admin")
	if err == nil && existingAdmin != nil {
		logger.Info("默认管理员已存在，跳过初始化")
		return nil
	}

	// 创建默认管理员
	admin := &Admin{
		Username:   "admin",
		Password:   "admin123",
		Email:      "admin@example.com",
		RealName:   "系统管理员",
		Rank:       10, // 超级管理员
		Status:     1,  // 启用
		LastLogin:  time.Now(),
		LastIP:     "127.0.0.1",
		LoginCount: 0,
	}

	id, err := m.Create(admin)
	if err != nil {
		logger.Error("创建默认管理员失败", "error", err)
		return err
	}

	logger.Info("默认管理员创建成功", "id", id, "username", admin.Username)
	return nil
}
