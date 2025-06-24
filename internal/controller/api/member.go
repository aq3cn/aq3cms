package api

import (
	"encoding/json"
	"net/http"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// MemberController 会员API控制器
type MemberController struct {
	*BaseController
}

// NewMemberController 创建会员API控制器
func NewMemberController(db *database.DB, cache cache.Cache, config *config.Config) *MemberController {
	return &MemberController{
		BaseController: NewBaseController(db, cache, config),
	}
}

// Login 会员登录
func (c *MemberController) Login(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 解析请求体
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证必填字段
	if loginData.Username == "" || loginData.Password == "" {
		c.Error(w, 400, "Missing required fields")
		return
	}

	// 获取会员
	member, err := c.memberModel.GetByUsername(loginData.Username)
	if err != nil || member == nil {
		c.Error(w, 401, "Invalid username or password")
		return
	}

	// 验证密码
	if !security.CheckPassword(loginData.Password, member.Password) {
		c.Error(w, 401, "Invalid username or password")
		return
	}

	// 检查会员状态
	if member.Status != 1 {
		c.Error(w, 403, "Account is disabled")
		return
	}

	// 更新登录信息
	member.LastLogin = time.Now()
	member.LastIP = r.RemoteAddr
	member.LoginCount++
	c.memberModel.Update(member)

	// 生成Token
	token, err := security.GenerateToken(map[string]interface{}{
		"member_id": member.ID,
		"username":  member.Username,
	}, c.config.Server.JWTSecret, 60*60*24*7) // 7天过期
	if err != nil {
		logger.Error("生成Token失败", "error", err)
		c.Error(w, 500, "Failed to generate token")
		return
	}

	// 缓存Token
	cacheKey := "token:" + token
	c.cache.Set(cacheKey, member.ID, time.Hour*24*7)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"token":   token,
		"member":  member,
		"message": "Login successful",
	})
}

// Register 会员注册
func (c *MemberController) Register(w http.ResponseWriter, r *http.Request) {
	// 记录API访问
	c.RecordAPIAccess(r, 0)

	// 解析请求体
	var registerData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Mobile   string `json:"mobile"`
	}
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证必填字段
	if registerData.Username == "" || registerData.Password == "" || registerData.Email == "" {
		c.Error(w, 400, "Missing required fields")
		return
	}

	// 检查用户名是否已存在
	existingMember, err := c.memberModel.GetByUsername(registerData.Username)
	if err == nil && existingMember != nil {
		c.Error(w, 400, "Username already exists")
		return
	}

	// 检查邮箱是否已存在
	existingMember, err = c.memberModel.GetByEmail(registerData.Email)
	if err == nil && existingMember != nil {
		c.Error(w, 400, "Email already exists")
		return
	}

	// 创建会员
	member := &model.Member{
		Username: registerData.Username,
		Password: security.HashPassword(registerData.Password),
		Email:    registerData.Email,
		Mobile:   registerData.Mobile,
		MType:    0,
		// Sex字段暂时跳过，避免数据库兼容性问题
		Status:    1,
		RegTime:   time.Now(),
		RegIP:     r.RemoteAddr,
		LastLogin: time.Time{},
		LastIP:    "",
		Money:     0,
		Score:     0,
	}

	// 保存会员
	id, err := c.memberModel.Create(member)
	if err != nil {
		logger.Error("创建会员失败", "error", err)
		c.Error(w, 500, "Failed to create member")
		return
	}

	// 生成Token
	token, err := security.GenerateToken(map[string]interface{}{
		"member_id": id,
		"username":  member.Username,
	}, c.config.Server.JWTSecret, 60*60*24*7) // 7天过期
	if err != nil {
		logger.Error("生成Token失败", "error", err)
		c.Error(w, 500, "Failed to generate token")
		return
	}

	// 缓存Token
	cacheKey := "token:" + token
	c.cache.Set(cacheKey, id, time.Hour*24*7)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"token":   token,
		"member":  member,
		"message": "Registration successful",
	})
}

// Profile 会员资料
func (c *MemberController) Profile(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取会员
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		c.Error(w, 404, "Member not found")
		return
	}

	// 获取会员文章
	articles, _, err := c.articleModel.GetByMemberID(memberID, 1, 10)
	if err != nil {
		logger.Error("获取会员文章失败", "id", memberID, "error", err)
	}

	// 获取会员评论
	comments, _, err := c.commentModel.GetByMemberID(memberID, 1, 10)
	if err != nil {
		logger.Error("获取会员评论失败", "id", memberID, "error", err)
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"member":   member,
		"articles": articles,
		"comments": comments,
	})
}

// UpdateProfile 更新会员资料
func (c *MemberController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取会员
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		c.Error(w, 404, "Member not found")
		return
	}

	// 解析请求体
	var updateData struct {
		Email  string `json:"email"`
		Mobile string `json:"mobile"`
		Sex    string `json:"sex"`
		QQ     string `json:"qq"`
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 检查邮箱是否已存在
	if updateData.Email != "" && updateData.Email != member.Email {
		existingMember, err := c.memberModel.GetByEmail(updateData.Email)
		if err == nil && existingMember != nil && existingMember.ID != memberID {
			c.Error(w, 400, "Email already exists")
			return
		}
	}

	// 更新会员资料
	if updateData.Email != "" {
		member.Email = updateData.Email
	}
	if updateData.Mobile != "" {
		member.Mobile = updateData.Mobile
	}
	if updateData.Sex != "" {
		member.Sex = updateData.Sex
	}

	if updateData.QQ != "" {
		member.QQ = updateData.QQ
	}
	if updateData.Avatar != "" {
		member.Avatar = updateData.Avatar
	}

	// 保存会员
	err = c.memberModel.Update(member)
	if err != nil {
		logger.Error("更新会员失败", "error", err)
		c.Error(w, 500, "Failed to update member")
		return
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Profile updated successfully",
	})
}

// ChangePassword 修改密码
func (c *MemberController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// 检查认证
	memberID, ok := c.CheckAuth(w, r)
	if !ok {
		return
	}

	// 记录API访问
	c.RecordAPIAccess(r, memberID)

	// 获取会员
	member, err := c.memberModel.GetByID(memberID)
	if err != nil {
		logger.Error("获取会员失败", "id", memberID, "error", err)
		c.Error(w, 404, "Member not found")
		return
	}

	// 解析请求体
	var passwordData struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
		c.Error(w, 400, "Invalid request body")
		return
	}

	// 验证必填字段
	if passwordData.OldPassword == "" || passwordData.NewPassword == "" {
		c.Error(w, 400, "Missing required fields")
		return
	}

	// 验证旧密码
	if !security.CheckPassword(passwordData.OldPassword, member.Password) {
		c.Error(w, 400, "Invalid old password")
		return
	}

	// 更新密码
	member.Password = security.HashPassword(passwordData.NewPassword)

	// 保存会员
	err = c.memberModel.Update(member)
	if err != nil {
		logger.Error("更新会员失败", "error", err)
		c.Error(w, 500, "Failed to update member")
		return
	}

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// Logout 会员退出
func (c *MemberController) Logout(w http.ResponseWriter, r *http.Request) {
	// 获取Token
	token := c.GetToken(r)
	if token == "" {
		c.Error(w, 400, "Missing token")
		return
	}

	// 删除缓存
	cacheKey := "token:" + token
	c.cache.Delete(cacheKey)

	// 返回数据
	c.Success(w, map[string]interface{}{
		"message": "Logout successful",
	})
}
