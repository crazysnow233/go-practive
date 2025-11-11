// Package middleware 错误处理中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RecoverJSON panic 恢复中间件
// 作用：捕获处理器中的 panic，防止程序崩溃
// 什么是 panic？
// - panic 是 Go 中的严重错误，类似于其他语言的 exception
// - 如果不处理，panic 会导致整个程序崩溃
// - 常见原因：空指针访问、数组越界等
func RecoverJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		// defer + recover 是 Go 处理 panic 的标准模式
		// defer: 延迟执行，函数返回前（或 panic 时）执行
		// recover: 捕获 panic，类似于 try-catch 中的 catch
		defer func() {
			// recover() 捕获 panic
			// 如果没有 panic，返回 nil
			// 如果有 panic，返回 panic 的值
			if rec := recover(); rec != nil {
				// 捕获到 panic！
				// 不要让程序崩溃，而是返回一个友好的 JSON 错误响应

				// AbortWithStatusJSON 终止请求处理并返回 JSON
				// http.StatusInternalServerError = 500（服务器内部错误）
				// gin.H 是 map[string]interface{} 的简写，用于构建 JSON
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})

				// 注意：实际生产环境中，应该：
				// 1. 记录详细的错误日志（包括堆栈信息）
				// 2. 不要向用户暴露敏感的错误细节
				// 3. 可以考虑发送告警通知
			}
		}()

		// 继续执行下一个中间件/处理器
		// 如果它们发生 panic，会被上面的 defer 捕获
		c.Next()
	}
}
