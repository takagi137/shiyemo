package log

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware 日志中间件
func Middleware(logService *LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 获取用户ID
		var userID int64
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(int64); ok {
				userID = id
			}
		}

		// 设置请求ID到上下文
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 创建请求日志记录器
		reqLogger := logService.WithRequest(requestID, userID)

		// 记录请求开始
		reqLogger.Info("请求开始", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 记录请求结束
		fields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     statusCode,
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			reqLogger.Error("请求结束", fields)
		} else if statusCode >= 400 {
			reqLogger.Warn("请求结束", fields)
		} else {
			reqLogger.Info("请求结束", fields)
		}
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}