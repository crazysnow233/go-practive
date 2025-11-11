// Package middleware 日志中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Logger 日志记录中间件
// 记录每个 HTTP 请求的详细信息
// 这对于调试、监控、审计都非常重要
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		// 用于后续计算请求处理耗时
		start := time.Now()

		// 收集请求信息
		path := c.Request.URL.Path           // 请求路径，如 /api/v1/boards
		raw := c.Request.URL.RawQuery        // 查询参数，如 page=1&size=10
		method := c.Request.Method           // HTTP 方法，如 GET、POST
		ip := c.ClientIP()                   // 客户端 IP 地址
		ua := c.Request.UserAgent()          // User-Agent（浏览器/客户端信息）

		// 执行下一个中间件/处理器
		// 注意：这里是分界线！
		// 上面的代码在处理器之前执行
		// 下面的代码在处理器之后执行
		c.Next()

		// 处理器执行完毕，收集响应信息
		status := c.Writer.Status()          // HTTP 状态码，如 200、404、500
		latency := time.Since(start)         // 请求处理耗时
		size := c.Writer.Size()              // 响应体大小（字节）

		// 获取用户 ID（如果已登录）
		// 从上下文中获取，由认证中间件设置
		userID := c.GetString("userID")
		if userID == "" {
			userID = "-"  // 未登录用 - 表示
		}

		// 获取请求 ID
		reqID := c.GetString("requestID")
		if reqID == "" {
			reqID = "-"
		}

		// 获取错误信息（如果有）
		// c.Errors 是 Gin 收集的错误列表
		errMsg := ""
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		// 构建完整的查询字符串
		q := ""
		if raw != "" {
			q = "?" + raw
		}

		// 打印结构化的日志
		// 使用 key=value 格式，方便日志分析工具解析
		// 生产环境建议使用专业的日志库（如 zap、logrus）
		log.Printf(
			"req_id=%s status=%d method=%s path=%s%s ip=%s user=%s size=%dB latency=%s ua=%q err=%q",
			reqID,   // 请求 ID
			status,  // 状态码
			method,  // HTTP 方法
			path, q, // 路径和查询参数
			ip,      // 客户端 IP
			userID,  // 用户 ID
			size,    // 响应大小
			latency, // 耗时
			ua,      // User-Agent
			errMsg,  // 错误信息
		)
	}
}
