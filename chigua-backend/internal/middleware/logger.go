package middleware

import (
	"chigua-backend/utils/ip2region"
	"chigua-backend/utils/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 自定义日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		method := c.Request.Method

		// 请求路由
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()
		logger.Infof("clientIP: %s", clientIP)
		area := ip2region.SearchArea(clientIP)
		if area != "" && strings.HasPrefix(area, "Reserved") == false {
			clientIP = area
		}

		// 记录日志
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
