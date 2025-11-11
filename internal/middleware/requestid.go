// Package middleware 包含 HTTP 中间件
// 中间件（Middleware）是什么？
// - 中间件是在请求到达处理器之前/之后执行的函数
// - 可以理解为"拦截器"或"过滤器"
// - 用途：日志记录、鉴权、错误处理、请求ID追踪等
// - 执行顺序：请求 -> 中间件1 -> 中间件2 -> 处理器 -> 中间件2 -> 中间件1 -> 响应
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 请求ID中间件
// 为每个 HTTP 请求生成或使用一个唯一的 ID
// 作用：
// 1. 追踪请求：在日志中可以根据 ID 追踪一个请求的完整生命周期
// 2. 调试：当用户报告问题时，可以通过请求 ID 定位日志
// 3. 分布式追踪：在微服务架构中传递请求 ID
func RequestID() gin.HandlerFunc {
	// gin.HandlerFunc 是 Gin 框架的中间件/处理器类型
	// 类型定义：type HandlerFunc func(*gin.Context)

	// 返回一个函数，这个函数就是中间件的实际逻辑
	return func(c *gin.Context) {
		// c *gin.Context 是 Gin 的上下文对象
		// 包含了请求、响应、参数等所有信息
		// 可以理解为"请求的容器"

		// 从请求头中获取 X-Request-Id
		// 如果客户端已经发送了请求 ID（例如负载均衡器添加的），就使用它
		id := c.GetHeader("X-Request-Id")

		// 如果没有，就生成一个新的 UUID
		if id == "" {
			id = uuid.NewString()
		}

		// 将请求 ID 存储到上下文中
		// 后续的中间件和处理器可以通过 c.GetString("requestID") 获取
		c.Set("requestID", id)

		// 在响应头中也返回请求 ID
		// 这样客户端可以知道这个请求的 ID，方便调试
		c.Writer.Header().Set("X-Request-Id", id)

		// c.Next() 是关键！
		// 调用下一个中间件或处理器
		// 如果不调用 Next()，请求处理就会中断
		c.Next()

		// Next() 之后的代码会在处理器执行完后运行
		// 可以用于清理、记录响应时间等
	}
}
