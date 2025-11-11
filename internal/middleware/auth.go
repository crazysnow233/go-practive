// Package middleware 认证中间件
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// CustomClaims JWT 声明结构
// 必须与 service/auth.go 中的 customClaims 保持一致
type CustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthRequired 认证中间件
// 要求请求必须携带有效的 JWT 令牌
// 用于保护需要登录才能访问的接口
func AuthRequired(secret []byte) gin.HandlerFunc {
	// 返回一个闭包（closure），捕获了 secret 变量
	// 这样每次请求都可以使用同一个密钥来验证令牌
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 字段
		// 标准格式是：Authorization: Bearer <token>
		// Bearer 是一种认证类型，表示"持有者令牌"
		authz := c.GetHeader("Authorization")

		// 检查是否以 "Bearer " 开头
		if !strings.HasPrefix(authz, "Bearer ") {
			// 没有令牌或格式错误
			// AbortWithStatusJSON 会终止请求处理，不再调用后续的处理器
			// http.StatusUnauthorized = 401（未授权）
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		// 提取令牌字符串（去掉 "Bearer " 前缀）
		raw := strings.TrimPrefix(authz, "Bearer ")

		// jwt.ParseWithClaims 解析并验证 JWT
		// 参数说明：
		// 1. raw: JWT 字符串
		// 2. &CustomClaims{}: 用于存储解析结果的结构体
		// 3. 回调函数：返回用于验证签名的密钥
		tok, err := jwt.ParseWithClaims(raw, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			// 这个函数会被 JWT 库调用，用于获取验证密钥
			// 返回签名时使用的同一个密钥
			return secret, nil
		})

		// 检查解析和验证结果
		// err != nil: 解析失败（格式错误、签名不匹配等）
		// !tok.Valid: 令牌无效（过期、未生效等）
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 类型断言：将 interface{} 转换为 *CustomClaims
		// tok.Claims 的类型是 interface{}，需要转换为具体类型才能使用
		// ok 表示转换是否成功
		claims, ok := tok.Claims.(*CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 验证通过！将用户信息存入上下文
		// 后续的处理器可以通过 c.GetString("userID") 获取当前用户的 ID
		c.Set("userID", claims.Subject) // Subject 存储的是用户 ID
		c.Set("email", claims.Email)

		// 继续执行后续的处理器
		// 此时请求已经通过认证，可以访问受保护的资源
		c.Next()
	}
}
