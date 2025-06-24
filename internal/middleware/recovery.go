package middleware

import (
	"net/http"
	"runtime/debug"

	"aq3cms/pkg/logger"
)

// Recovery 恢复中间件
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()

				// 记录错误日志
				logger.Error("HTTP请求处理异常",
					"error", err,
					"stack", string(stack),
					"method", r.Method,
					"path", r.URL.Path,
					"query", r.URL.RawQuery,
					"ip", getClientIP(r),
					"user-agent", r.UserAgent(),
				)

				// 返回500错误
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// 处理请求
		next.ServeHTTP(w, r)
	})
}

// NotFound 404处理中间件
func NotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建自定义响应写入器
		nfw := &notFoundWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			written:        false,
		}

		// 处理请求
		next.ServeHTTP(nfw, r)

		// 如果没有写入响应，返回404
		if !nfw.written {
			logger.Info("页面未找到",
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"ip", getClientIP(r),
				"user-agent", r.UserAgent(),
			)

			// 返回404错误
			http.NotFound(w, r)
		}
	})
}

// notFoundWriter 自定义响应写入器
type notFoundWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader 重写WriteHeader方法
func (nfw *notFoundWriter) WriteHeader(code int) {
	nfw.statusCode = code
	nfw.written = true
	nfw.ResponseWriter.WriteHeader(code)
}

// Write 重写Write方法
func (nfw *notFoundWriter) Write(b []byte) (int, error) {
	nfw.written = true
	return nfw.ResponseWriter.Write(b)
}
