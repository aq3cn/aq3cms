package middleware

import (
	"net/http"
	"time"
)

// Security 安全中间件
func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全相关的HTTP头
		
		// 防止点击劫持
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		
		// 启用XSS过滤
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// 防止MIME类型嗅探
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// 内容安全策略
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'")
		
		// 引用策略
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// 特性策略
		w.Header().Set("Feature-Policy", "camera 'none'; microphone 'none'; geolocation 'none'")
		
		// 处理请求
		next.ServeHTTP(w, r)
	})
}

// RateLimit 速率限制中间件
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	// 创建IP限制映射
	type client struct {
		count    int
		lastSeen time.Time
	}
	clients := make(map[string]*client)
	
	// 清理过期记录
	go func() {
		for {
			time.Sleep(time.Minute)
			
			// 删除过期记录
			now := time.Now()
			for ip, c := range clients {
				if now.Sub(c.lastSeen) > window {
					delete(clients, ip)
				}
			}
		}
	}()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取客户端IP
			ip := getClientIP(r)
			
			// 检查是否超过限制
			c, exists := clients[ip]
			if !exists {
				clients[ip] = &client{
					count:    1,
					lastSeen: time.Now(),
				}
			} else {
				// 检查时间窗口
				if time.Since(c.lastSeen) > window {
					// 重置计数
					c.count = 1
					c.lastSeen = time.Now()
				} else {
					// 增加计数
					c.count++
					c.lastSeen = time.Now()
					
					// 检查是否超过限制
					if c.count > limit {
						w.Header().Set("Retry-After", window.String())
						http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
						return
					}
				}
			}
			
			// 处理请求
			next.ServeHTTP(w, r)
		})
	}
}

// CORS CORS中间件
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取请求的Origin
			origin := r.Header.Get("Origin")
			
			// 检查Origin是否在允许列表中
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
			
			// 如果Origin被允许，设置CORS头
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}
			
			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// 处理请求
			next.ServeHTTP(w, r)
		})
	}
}
