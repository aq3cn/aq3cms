package middleware

import (
	"net/http"
	"time"

	"aq3cms/pkg/logger"
)

// Logger 日志中间件
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 创建自定义响应写入器
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// 处理请求
		next.ServeHTTP(lrw, r)
		
		// 计算请求处理时间
		duration := time.Since(start)
		
		// 记录请求日志
		logger.Info("HTTP请求",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", lrw.statusCode,
			"duration", duration,
			"ip", getClientIP(r),
			"user-agent", r.UserAgent(),
			"referer", r.Referer(),
		)
	})
}

// loggingResponseWriter 自定义响应写入器
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写WriteHeader方法
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// getClientIP 获取客户端IP
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For获取
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	
	// 尝试从X-Real-IP获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	
	// 使用RemoteAddr
	return r.RemoteAddr
}
