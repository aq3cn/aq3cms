package service

import (
	"fmt"
	"net/http"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// AdminService 管理员服务
type AdminService struct {
	db         *database.DB
	cache      cache.Cache
	config     *config.Config
	adminModel *model.AdminModel
}

// NewAdminService 创建管理员服务
func NewAdminService(db *database.DB, cache cache.Cache, config *config.Config) *AdminService {
	return &AdminService{
		db:         db,
		cache:      cache,
		config:     config,
		adminModel: model.NewAdminModel(db),
	}
}

// Login 管理员登录
func (s *AdminService) Login(username, password string) (*model.Admin, error) {
	// 获取管理员
	admin, err := s.adminModel.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if admin == nil {
		logger.Error("用户名或密码错误")
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if !security.CheckPassword(password, admin.Password) {
		logger.Error("用户名或密码错误")
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查管理员状态
	if admin.Status != 1 {
		logger.Error("账号已禁用")
		return nil, fmt.Errorf("账号已禁用")
	}

	return admin, nil
}

// UpdateLoginInfo 更新登录信息
func (s *AdminService) UpdateLoginInfo(id int64, r *http.Request) error {
	// 获取管理员
	admin, err := s.adminModel.GetByID(id)
	if err != nil {
		return err
	}

	if admin == nil {
		logger.Error("管理员不存在")
		return fmt.Errorf("管理员不存在")
	}

	// 更新登录信息
	admin.LastLogin = time.Now()
	admin.LastIP = r.RemoteAddr
	admin.LoginCount++

	// 更新管理员
	return s.adminModel.Update(admin)
}

// GetAdminList 获取管理员列表
func (s *AdminService) GetAdminList() ([]*model.Admin, error) {
	return s.adminModel.GetAll()
}

// GetAdminByID 根据ID获取管理员
func (s *AdminService) GetAdminByID(id int64) (*model.Admin, error) {
	return s.adminModel.GetByID(id)
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(admin *model.Admin) (int64, error) {
	return s.adminModel.Create(admin)
}

// UpdateAdmin 更新管理员
func (s *AdminService) UpdateAdmin(admin *model.Admin) error {
	return s.adminModel.Update(admin)
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(id int64) error {
	return s.adminModel.Delete(id)
}

// ChangePassword 修改密码
func (s *AdminService) ChangePassword(id int64, oldPassword, newPassword string) error {
	return s.adminModel.ChangePassword(id, oldPassword, newPassword)
}
