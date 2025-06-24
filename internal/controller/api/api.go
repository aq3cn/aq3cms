package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/internal/service"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
	"aq3cms/pkg/security"
)

// BaseController API基础控制器
type BaseController struct {
	db              *database.DB
	cache           cache.Cache
	config          *config.Config
	articleModel    *model.ArticleModel
	categoryModel   *model.CategoryModel
	memberModel     *model.MemberModel
	commentModel    *model.CommentModel
	tagModel        *model.TagModel
	specialModel    *model.SpecialModel
	contentModel    *model.ContentModelModel
	templateService *service.TemplateService
	htmlService     *service.HtmlService
	seoService      *service.SEOService
	statsService    *service.StatsService
}

// NewBaseController 创建API基础控制器
func NewBaseController(db *database.DB, cache cache.Cache, config *config.Config) *BaseController {
	return &BaseController{
		db:              db,
		cache:           cache,
		config:          config,
		articleModel:    model.NewArticleModel(db),
		categoryModel:   model.NewCategoryModel(db),
		memberModel:     model.NewMemberModel(db),
		commentModel:    model.NewCommentModel(db),
		tagModel:        model.NewTagModel(db),
		specialModel:    model.NewSpecialModel(db),
		contentModel:    model.NewContentModelModel(db),
		templateService: service.NewTemplateService(db, cache, config),
		htmlService:     service.NewHtmlService(db, cache, config),
		seoService:      service.NewSEOService(db, cache, config),
		statsService:    service.NewStatsService(db, cache, config),
	}
}

// Response API响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func (c *BaseController) Success(w http.ResponseWriter, data interface{}) {
	c.JSON(w, http.StatusOK, &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func (c *BaseController) Error(w http.ResponseWriter, code int, message string) {
	c.JSON(w, http.StatusOK, &Response{
		Code:    code,
		Message: message,
	})
}

// JSON 输出JSON
func (c *BaseController) JSON(w http.ResponseWriter, status int, data interface{}) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// 输出JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("输出JSON失败", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetInt64Param 获取int64参数
func (c *BaseController) GetInt64Param(r *http.Request, name string) (int64, error) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	if !ok {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

// GetIntParam 获取int参数
func (c *BaseController) GetIntParam(r *http.Request, name string) (int, error) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	if !ok {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// GetQueryInt64 获取查询参数int64
func (c *BaseController) GetQueryInt64(r *http.Request, name string, defaultValue int64) int64 {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetQueryInt 获取查询参数int
func (c *BaseController) GetQueryInt(r *http.Request, name string, defaultValue int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetQueryString 获取查询参数string
func (c *BaseController) GetQueryString(r *http.Request, name string, defaultValue string) string {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryBool 获取查询参数bool
func (c *BaseController) GetQueryBool(r *http.Request, name string, defaultValue bool) bool {
	value := r.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetToken 获取Token
func (c *BaseController) GetToken(r *http.Request) string {
	// 从请求头获取
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	// 从查询参数获取
	token = r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// 从表单获取
	token = r.FormValue("token")
	return token
}

// VerifyToken 验证Token
func (c *BaseController) VerifyToken(token string) (int64, error) {
	// 从缓存获取
	cacheKey := "token:" + token
	if cached, ok := c.cache.Get(cacheKey); ok {
		if memberID, ok := cached.(int64); ok {
			return memberID, nil
		}
	}

	// 解析Token
	claims, err := security.ParseToken(token, c.config.Server.JWTSecret)
	if err != nil {
		return 0, err
	}

	// 获取会员ID
	memberID, ok := claims["member_id"].(float64)
	if !ok {
		return 0, nil
	}

	// 缓存Token
	c.cache.Set(cacheKey, int64(memberID), time.Hour*24)

	return int64(memberID), nil
}

// CheckAuth 检查认证
func (c *BaseController) CheckAuth(w http.ResponseWriter, r *http.Request) (int64, bool) {
	// 获取Token
	token := c.GetToken(r)
	if token == "" {
		c.Error(w, 401, "Unauthorized")
		return 0, false
	}

	// 验证Token
	memberID, err := c.VerifyToken(token)
	if err != nil {
		c.Error(w, 401, "Invalid token")
		return 0, false
	}

	// 检查会员是否存在
	member, err := c.memberModel.GetByID(memberID)
	if err != nil || member == nil {
		c.Error(w, 401, "Member not found")
		return 0, false
	}

	// 检查会员状态
	if member.Status != 1 {
		c.Error(w, 403, "Member is disabled")
		return 0, false
	}

	return memberID, true
}

// CheckAPIKey 检查API密钥
func (c *BaseController) CheckAPIKey(w http.ResponseWriter, r *http.Request) bool {
	// 获取API密钥
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		apiKey = r.URL.Query().Get("api_key")
	}

	// 检查API密钥
	if apiKey != c.config.API.Key {
		c.Error(w, 401, "Invalid API key")
		return false
	}

	return true
}

// RecordAPIAccess 记录API访问
func (c *BaseController) RecordAPIAccess(r *http.Request, memberID int64) {
	// 记录访问
	c.statsService.RecordVisit(r, memberID, 0, 0)
}
