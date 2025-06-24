package middleware

import (
	"net/http"
	"time"

	"aq3cms/config"
	"aq3cms/pkg/logger"
	"github.com/gorilla/sessions"
)

var (
	// 会话存储
	sessionStore *sessions.CookieStore
)

// InitAuth 初始化认证中间件
func InitAuth(cfg *config.Config) {
	sessionStore = sessions.NewCookieStore([]byte(cfg.Server.SessionSecret))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7天
		HttpOnly: true,
		Secure:   false,                // 开发环境使用HTTP
		SameSite: http.SameSiteLaxMode, // 使用Lax模式
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取会话
		session, err := sessionStore.Get(r, "admin-session")
		if err != nil {
			logger.Error("获取管理员会话失败", "error", err)
			http.Redirect(w, r, "/aq3cms/login", http.StatusFound)
			return
		}

		// 检查是否已登录
		if session.Values["admin_id"] == nil {
			http.Redirect(w, r, "/aq3cms/login", http.StatusFound)
			return
		}

		// 检查管理员权限
		if rank, ok := session.Values["admin_rank"].(int); !ok || rank < 1 {
			http.Redirect(w, r, "/aq3cms/login", http.StatusFound)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// MemberAuth 会员认证中间件
func MemberAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取会话
		session, err := sessionStore.Get(r, "member-session")
		if err != nil {
			logger.Error("获取会员会话失败", "error", err)
			http.Redirect(w, r, "/member/login", http.StatusFound)
			return
		}

		// 检查是否已登录
		if session.Values["member_id"] == nil {
			http.Redirect(w, r, "/member/login", http.StatusFound)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// GetAdminID 获取管理员ID
func GetAdminID(r *http.Request) int64 {
	session, err := sessionStore.Get(r, "admin-session")
	if err != nil {
		return 0
	}

	if id, ok := session.Values["admin_id"].(int64); ok {
		return id
	}

	return 0
}

// GetMemberID 获取会员ID
func GetMemberID(r *http.Request) int64 {
	session, err := sessionStore.Get(r, "member-session")
	if err != nil {
		return 0
	}

	if id, ok := session.Values["member_id"].(int64); ok {
		return id
	}

	return 0
}

// GetAdminName 获取管理员名称
func GetAdminName(r *http.Request) string {
	session, err := sessionStore.Get(r, "admin-session")
	if err != nil {
		return ""
	}

	if name, ok := session.Values["admin_name"].(string); ok {
		return name
	}

	return ""
}

// GetMemberName 获取会员名称
func GetMemberName(r *http.Request) string {
	session, err := sessionStore.Get(r, "member-session")
	if err != nil {
		return ""
	}

	if name, ok := session.Values["member_name"].(string); ok {
		return name
	}

	return ""
}

// IsAdminLoggedIn 检查管理员是否已登录
func IsAdminLoggedIn(r *http.Request) bool {
	return GetAdminID(r) > 0
}

// IsMemberLoggedIn 检查会员是否已登录
func IsMemberLoggedIn(r *http.Request) bool {
	return GetMemberID(r) > 0
}

// GetAdminRank 获取管理员等级
func GetAdminRank(r *http.Request) int {
	session, err := sessionStore.Get(r, "admin-session")
	if err != nil {
		return 0
	}

	if rank, ok := session.Values["admin_rank"].(int); ok {
		return rank
	}

	return 0
}

// CreateAdminSession 创建管理员会话
func CreateAdminSession(w http.ResponseWriter, r *http.Request, adminID int64, adminName string, adminRank int) error {
	logger.Info("创建管理员会话", "adminID", adminID, "adminName", adminName, "adminRank", adminRank)

	session, err := sessionStore.Get(r, "admin-session")
	if err != nil {
		logger.Error("获取会话失败", "error", err)
		return err
	}

	session.Values["admin_id"] = adminID
	session.Values["admin_name"] = adminName
	session.Values["admin_rank"] = adminRank
	session.Values["admin_last_login"] = time.Now().Unix()

	err = session.Save(r, w)
	if err != nil {
		logger.Error("保存会话失败", "error", err)
		return err
	}

	logger.Info("管理员会话创建成功", "adminID", adminID)
	return nil
}

// DestroyAdminSession 销毁管理员会话
func DestroyAdminSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, "admin-session")
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	return session.Save(r, w)
}

// GetAdminSession 获取管理员会话
func GetAdminSession(r *http.Request) (*sessions.Session, error) {
	return sessionStore.Get(r, "admin-session")
}
