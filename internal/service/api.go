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
)

// APIService API服务
type APIService struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	memberModel     *model.MemberModel
	formService     *FormService
	commentModel    *model.CommentModel
	articleService  *ArticleService
	memberService   *MemberService
}

// NewAPIService 创建API服务
func NewAPIService(db *database.DB, cache cache.Cache, config *config.Config) *APIService {
	articleModel := model.NewArticleModel(db)
	categoryModel := model.NewCategoryModel(db)
	memberModel := model.NewMemberModel(db)

	return &APIService{
		db:             db,
		cache:          cache,
		config:         config,
		articleModel:   articleModel,
		categoryModel:  categoryModel,
		memberModel:    memberModel,
		formService:    NewFormService(db, cache, config),
		commentModel:   model.NewCommentModel(db),
		articleService: NewArticleService(db, cache, config),
		memberService:  NewMemberService(db, cache, config),
	}
}

// GenerateToken 生成API令牌
func (s *APIService) GenerateToken(appID, appSecret string, expireSeconds int) (string, error) {
	// 验证应用ID和密钥
	if appID != s.config.API.Key || appSecret != s.config.API.Secret {
		return "", fmt.Errorf("invalid app id or app secret")
	}

	// 设置过期时间
	if expireSeconds <= 0 {
		expireSeconds = 3600 // 默认1小时
	}
	expireTime := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()

	// 构建令牌数据
	tokenData := map[string]interface{}{
		"app_id":     appID,
		"expire_at":  expireTime,
		"created_at": time.Now().Unix(),
	}

	// 序列化令牌数据
	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	// 计算签名
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write(tokenJSON)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 构建令牌
	token := base64.StdEncoding.EncodeToString(tokenJSON) + "." + signature

	// 缓存令牌
	cacheKey := "api:token:" + appID
	cache.SafeSet(s.cache, cacheKey, token, time.Duration(expireSeconds)*time.Second)

	return token, nil
}

// VerifyToken 验证API令牌
func (s *APIService) VerifyToken(token string) (bool, error) {
	// 解析令牌
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid token format")
	}

	// 解码令牌数据
	tokenJSON, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	// 解析令牌数据
	var tokenData map[string]interface{}
	err = json.Unmarshal(tokenJSON, &tokenData)
	if err != nil {
		return false, err
	}

	// 获取应用ID
	appID, ok := tokenData["app_id"].(string)
	if !ok {
		return false, fmt.Errorf("invalid token data: app_id")
	}

	// 获取过期时间
	expireAt, ok := tokenData["expire_at"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid token data: expire_at")
	}

	// 检查令牌是否过期
	if time.Now().Unix() > int64(expireAt) {
		return false, fmt.Errorf("token expired")
	}

	// 获取应用密钥
	appSecret := s.config.API.Secret

	// 计算签名
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write(tokenJSON)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 验证签名
	if signature != parts[1] {
		return false, fmt.Errorf("invalid token signature")
	}

	// 从缓存获取令牌
	cacheKey := "api:token:" + appID
	cachedToken, ok := s.cache.Get(cacheKey)
	if !ok {
		return false, fmt.Errorf("token not found in cache")
	}

	// 验证令牌
	if cachedToken != token {
		return false, fmt.Errorf("token mismatch")
	}

	return true, nil
}

// GetArticles 获取文章列表
func (s *APIService) GetArticles(categoryID int64, page, pageSize int, orderBy string) ([]*model.Article, int, error) {
	// 获取文章列表
	return s.articleModel.GetList(categoryID, page, pageSize)
}

// GetArticle 获取文章详情
func (s *APIService) GetArticle(id int64) (*model.Article, error) {
	// 获取文章详情
	return s.articleModel.GetByID(id)
}

// GetCategories 获取栏目列表
func (s *APIService) GetCategories(parentID int64) ([]*model.Category, error) {
	// 获取栏目列表
	if parentID > 0 {
		return s.categoryModel.GetChildCategories(parentID)
	}
	return s.categoryModel.GetTopCategories()
}

// GetCategory 获取栏目详情
func (s *APIService) GetCategory(id int64) (*model.Category, error) {
	// 获取栏目详情
	return s.categoryModel.GetByID(id)
}

// Login 会员登录
func (s *APIService) Login(username, password string) (string, error) {
	// 会员登录
	member, err := s.memberService.Login(username, password)
	if err != nil {
		return "", err
	}

	// 生成会员令牌
	token, err := s.memberService.GenerateToken(member.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register 会员注册
func (s *APIService) Register(username, password, email, nickname string) (int64, error) {
	// 会员注册
	return s.memberService.Register(username, password, email, nickname)
}

// GetMember 获取会员信息
func (s *APIService) GetMember(id int64) (*model.Member, error) {
	// 获取会员信息
	return s.memberService.GetMemberByID(id)
}

// UpdateMember 更新会员信息
func (s *APIService) UpdateMember(id int64, data map[string]interface{}) error {
	// 获取会员信息
	member, err := s.memberService.GetMemberByID(id)
	if err != nil {
		return err
	}

	// 更新会员信息
	if username, ok := data["username"].(string); ok {
		member.Username = username
	}
	if email, ok := data["email"].(string); ok {
		member.Email = email
	}
	if avatar, ok := data["avatar"].(string); ok {
		member.Avatar = avatar
	}
	if mobile, ok := data["mobile"].(string); ok {
		member.Mobile = mobile
	}
	if sex, ok := data["sex"].(string); ok {
		member.Sex = sex
	}
	if qq, ok := data["qq"].(string); ok {
		member.QQ = qq
	}

	// 保存会员信息
	return s.memberService.UpdateMember(member)
}

// SubmitForm 提交表单
func (s *APIService) SubmitForm(code string, formData map[string]interface{}, ip string) error {
	// 提交表单
	return s.formService.SubmitForm(code, formData, ip)
}

// Search 搜索
func (s *APIService) Search(keyword string, page, pageSize int) ([]*model.Article, int, error) {
	// 搜索
	return s.articleService.SearchArticles(keyword, page, pageSize)
}

// GetTags 获取标签列表
func (s *APIService) GetTags() ([]string, error) {
	// 获取标签列表
	return s.articleService.GetAllTags()
}

// GetArticlesByTag 获取标签文章
func (s *APIService) GetArticlesByTag(tag string, page, pageSize int) ([]*model.Article, int, error) {
	// 获取标签文章
	return s.articleService.GetArticlesByTag(tag, page, pageSize)
}

// GetComments 获取评论列表
func (s *APIService) GetComments(articleID int64, page, pageSize int) ([]*model.Comment, int, error) {
	// 获取评论列表
	return s.commentModel.GetListByAID(articleID, page, pageSize)
}

// AddComment 添加评论
func (s *APIService) AddComment(articleID, memberID int64, content string, parentID int64) (int64, error) {
	// 添加评论
	comment := &model.Comment{
		AID:      articleID,
		MID:      memberID,
		Content:  content,
		ParentID: parentID,
		Dtime:    time.Now(),
		IsCheck:  1,
	}
	return s.commentModel.Create(comment)
}

// GetSiteConfig 获取站点配置
func (s *APIService) GetSiteConfig() map[string]interface{} {
	// 获取站点配置
	return map[string]interface{}{
		"name":        s.config.Site.Name,
		"keywords":    s.config.Site.Keywords,
		"description": s.config.Site.Description,
		"url":         s.config.Site.URL,
		"email":       s.config.Site.Email,
		"icp":         s.config.Site.ICP,
		"copyright":   s.config.Site.CopyRight,
	}
}
