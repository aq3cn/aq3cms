package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"aq3cms/pkg/logger"
)

// CSRFToken CSRF令牌
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFProtection CSRF防护
type CSRFProtection struct {
	tokens   map[string]*CSRFToken
	mutex    sync.RWMutex
	lifetime time.Duration
}

// NewCSRFProtection 创建CSRF防护
func NewCSRFProtection(lifetime time.Duration) *CSRFProtection {
	protection := &CSRFProtection{
		tokens:   make(map[string]*CSRFToken),
		lifetime: lifetime,
	}
	
	// 启动清理过期令牌的协程
	go protection.cleanExpiredTokens()
	
	return protection
}

// GenerateToken 生成CSRF令牌
func (p *CSRFProtection) GenerateToken(sessionID string) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	// 生成随机令牌
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logger.Error("生成随机令牌失败", "error", err)
		return ""
	}
	
	// 编码令牌
	token := base64.StdEncoding.EncodeToString(b)
	
	// 保存令牌
	p.tokens[sessionID] = &CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(p.lifetime),
	}
	
	return token
}

// ValidateToken 验证CSRF令牌
func (p *CSRFProtection) ValidateToken(sessionID, token string) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	// 获取令牌
	csrfToken, ok := p.tokens[sessionID]
	if !ok {
		return false
	}
	
	// 检查令牌是否过期
	if csrfToken.ExpiresAt.Before(time.Now()) {
		return false
	}
	
	// 验证令牌
	return csrfToken.Token == token
}

// RemoveToken 移除CSRF令牌
func (p *CSRFProtection) RemoveToken(sessionID string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	delete(p.tokens, sessionID)
}

// 清理过期令牌
func (p *CSRFProtection) cleanExpiredTokens() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		p.mutex.Lock()
		
		now := time.Now()
		for sessionID, token := range p.tokens {
			if token.ExpiresAt.Before(now) {
				delete(p.tokens, sessionID)
			}
		}
		
		p.mutex.Unlock()
	}
}

// CSRFMiddleware CSRF中间件
func CSRFMiddleware(protection *CSRFProtection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 跳过GET、HEAD、OPTIONS请求
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}
			
			// 获取会话ID
			session, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			
			// 获取CSRF令牌
			token := r.FormValue("csrf_token")
			if token == "" {
				token = r.Header.Get("X-CSRF-Token")
			}
			
			// 验证令牌
			if !protection.ValidateToken(session.Value, token) {
				http.Error(w, "CSRF token validation failed", http.StatusForbidden)
				return
			}
			
			// 继续处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// CSRFField 生成CSRF字段
func CSRFField(protection *CSRFProtection, sessionID string) string {
	token := protection.GenerateToken(sessionID)
	return fmt.Sprintf("<input type=\"hidden\" name=\"csrf_token\" value=\"%s\">", token)
}
