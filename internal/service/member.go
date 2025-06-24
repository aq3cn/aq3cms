package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// MemberService 会员服务
type MemberService struct {
	db           *database.DB
	cache        cache.Cache
	config       *config.Config
	memberModel  *model.MemberModel
}

// NewMemberService 创建会员服务
func NewMemberService(db *database.DB, cache cache.Cache, config *config.Config) *MemberService {
	return &MemberService{
		db:           db,
		cache:        cache,
		config:       config,
		memberModel:  model.NewMemberModel(db),
	}
}

// Login 会员登录
func (s *MemberService) Login(username, password string) (*model.Member, error) {
	// 获取会员信息
	member, err := s.memberModel.GetByUsername(username)
	if err != nil {
		logger.Error("会员登录失败", "username", username, "error", err)
		return nil, err
	}

	// 验证密码
	// 简单的密码验证，实际应该使用加密算法
	if member.Password != password {
		return nil, fmt.Errorf("密码错误")
	}

	// 更新会员登录信息
	member.LastLogin = time.Now()
	member.LastIP = ""
	member.LoginCount++
	err = s.memberModel.Update(member)
	if err != nil {
		logger.Error("更新会员登录信息失败", "id", member.ID, "error", err)
	}

	return member, nil
}

// GenerateToken 生成会员令牌
func (s *MemberService) GenerateToken(memberID int64) (string, error) {
	// 设置过期时间
	expireSeconds := 86400 // 默认24小时
	expireTime := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()

	// 构建令牌数据
	tokenData := map[string]interface{}{
		"member_id":  memberID,
		"expire_at":  expireTime,
		"created_at": time.Now().Unix(),
	}

	// 序列化令牌数据
	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	// 计算签名
	h := hmac.New(sha256.New, []byte(s.config.Server.JWTSecret))
	h.Write(tokenJSON)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 构建令牌
	token := base64.StdEncoding.EncodeToString(tokenJSON) + "." + signature

	// 缓存令牌
	cacheKey := fmt.Sprintf("member:token:%d", memberID)
	cache.SafeSet(s.cache, cacheKey, token, time.Duration(expireSeconds)*time.Second)

	return token, nil
}

// Register 会员注册
func (s *MemberService) Register(username, password, email, nickname string) (int64, error) {
	// 检查用户名是否存在
	existingMember, err := s.memberModel.GetByUsername(username)
	if err == nil && existingMember != nil {
		return 0, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否存在
	existingMember, err = s.memberModel.GetByEmail(email)
	if err == nil && existingMember != nil {
		return 0, fmt.Errorf("邮箱已存在")
	}

	// 创建会员
	member := &model.Member{
		Username:   username,
		Password:   password, // 实际应该使用加密算法
		Email:      email,
		Sex:        "保密",
		RegTime:    time.Now(),
		RegIP:      "",
		LastLogin:  time.Now(),
		LastIP:     "",
		LoginCount: 1,
		Status:     1,
	}
	return s.memberModel.Create(member)
}

// GetMemberByID 根据ID获取会员信息
func (s *MemberService) GetMemberByID(id int64) (*model.Member, error) {
	// 缓存键
	cacheKey := fmt.Sprintf("member:%d", id)

	// 检查缓存
	if cached, ok := s.cache.Get(cacheKey); ok {
		if member, ok := cached.(*model.Member); ok {
			return member, nil
		}
	}

	// 获取会员信息
	member, err := s.memberModel.GetByID(id)
	if err != nil {
		logger.Error("获取会员信息失败", "id", id, "error", err)
		return nil, err
	}

	// 缓存会员信息
	cache.SafeSet(s.cache, cacheKey, member, time.Duration(3600)*time.Second) // 使用默认过期时间1小时

	return member, nil
}

// UpdateMember 更新会员信息
func (s *MemberService) UpdateMember(member *model.Member) error {
	// 更新会员信息
	err := s.memberModel.Update(member)
	if err != nil {
		logger.Error("更新会员信息失败", "id", member.ID, "error", err)
		return err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("member:%d", member.ID)
	s.cache.Delete(cacheKey)

	return nil
}

// VerifyToken 验证会员令牌
func (s *MemberService) VerifyToken(token string) (int64, error) {
	// 解析令牌
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	// 解码令牌数据
	tokenJSON, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, err
	}

	// 解析令牌数据
	var tokenData map[string]interface{}
	err = json.Unmarshal(tokenJSON, &tokenData)
	if err != nil {
		return 0, err
	}

	// 获取会员ID
	memberID, ok := tokenData["member_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token data: member_id")
	}

	// 获取过期时间
	expireAt, ok := tokenData["expire_at"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token data: expire_at")
	}

	// 检查令牌是否过期
	if time.Now().Unix() > int64(expireAt) {
		return 0, fmt.Errorf("token expired")
	}

	// 计算签名
	h := hmac.New(sha256.New, []byte(s.config.Server.JWTSecret))
	h.Write(tokenJSON)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 验证签名
	if signature != parts[1] {
		return 0, fmt.Errorf("invalid token signature")
	}

	// 从缓存获取令牌
	cacheKey := fmt.Sprintf("member:token:%d", int64(memberID))
	cachedToken, ok := s.cache.Get(cacheKey)
	if !ok {
		return 0, fmt.Errorf("token not found in cache")
	}

	// 验证令牌
	if cachedToken != token {
		return 0, fmt.Errorf("token mismatch")
	}

	return int64(memberID), nil
}
