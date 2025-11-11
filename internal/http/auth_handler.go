// Package http 包含 HTTP 处理器（Handler）
// Handler 是 MVC 模式中的 Controller 层
// 职责：
// 1. 接收 HTTP 请求
// 2. 解析请求参数
// 3. 调用 Service 层处理业务逻辑
// 4. 返回 HTTP 响应
package http

import (
	"github.com/gin-gonic/gin"
	"kanban_api/internal/service"
	"net/http"
	"strings"
)

// AuthHandler 认证处理器
// 处理用户注册、登录等认证相关的 HTTP 请求
type AuthHandler struct {
	// svc 认证服务
	// Handler 通过接口依赖 Service，不知道具体实现
	svc service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// RegisterRoutes 注册路由
// 将 HTTP 路径和处理方法关联起来
// rg *gin.RouterGroup 是 Gin 的路由组，可以给一组路由添加统一的前缀或中间件
func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	// POST 方法表示要创建资源
	// "/auth/register" 是完整路径（会加上路由组的前缀）
	// h.register 是处理函数
	rg.POST("/auth/register", h.register)
	rg.POST("/auth/login", h.login)
}

// register 处理用户注册请求
// HTTP 方法：POST
// 路径：/api/v1/auth/register
// 请求体：{"email": "user@example.com", "password": "123456"}
func (h *AuthHandler) register(c *gin.Context) {
	// 定义请求体的结构
	// 使用匿名结构体，只在这个函数内使用
	// `json:"email"` 标签：JSON 中的字段名映射到结构体字段
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// ShouldBindJSON 解析 JSON 请求体
	// 它会：
	// 1. 读取 HTTP 请求体
	// 2. 将 JSON 解析到 req 结构体
	// 3. 进行基本的数据验证（如果有验证标签的话）
	// 如果解析失败（JSON 格式错误、字段类型不匹配等），返回错误
	if err := c.ShouldBindJSON(&req); err != nil {
		// http.StatusBadRequest = 400（错误的请求）
		// gin.H 是 map[string]interface{} 的别名，用于构建 JSON 响应
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// 调用 Service 层处理注册逻辑
	u, token, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		// 注册失败，根据错误类型返回不同的 HTTP 状态码
		msg := err.Error()

		// 如果是邮箱已存在的错误
		if strings.Contains(msg, "exists") {
			// http.StatusConflict = 409（冲突）
			// 表示请求与当前资源状态冲突（邮箱已注册）
			c.JSON(http.StatusConflict, gin.H{"error": msg})
			return
		}

		// 其他错误（邮箱格式错误、密码为空等）
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// 注册成功！
	// http.StatusCreated = 201（已创建）
	// 201 是创建资源成功的标准状态码
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			// 返回用户信息（不包含密码哈希！）
			"user": gin.H{
				"id":        u.ID,
				"email":     u.Email,
				"createdAt": u.CreatedAt,
			},
			// 返回 JWT 令牌，客户端保存后用于后续请求的认证
			"token": token,
		},
	})
}

// login 处理用户登录请求
// HTTP 方法：POST
// 路径：/api/v1/auth/login
// 请求体：{"email": "user@example.com", "password": "123456"}
func (h *AuthHandler) login(c *gin.Context) {
	// 请求体结构（与注册相同）
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	// 调用 Service 层验证登录
	u, token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		// 登录失败（用户不存在或密码错误）
		// http.StatusUnauthorized = 401（未授权）
		// 注意：无论是邮箱不存在还是密码错误，都返回相同的错误信息
		// 这是安全最佳实践，防止攻击者枚举有效邮箱
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 登录成功！
	// http.StatusOK = 200（成功）
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			// 返回用户信息
			"user": gin.H{"id": u.ID, "email": u.Email, "createdAt": u.CreatedAt},
			// 返回 JWT 令牌
			"token": token,
		},
	})
}
