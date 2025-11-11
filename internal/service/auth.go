// Package service 包含业务逻辑层
// 这一层处理业务规则、数据验证、事务协调等
// Service 层调用 Repository 层，被 Handler 层调用
package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5" // JWT（JSON Web Token）库，用于生成和验证令牌
	"golang.org/x/crypto/bcrypt"   // bcrypt 加密库，用于密码哈希
	"kanban_api/internal/model"
	"kanban_api/internal/repository"
	"os"
	"strings"
	"time"
)

// AuthService 认证服务接口
// 负责用户注册、登录、JWT 令牌生成等认证相关的业务逻辑
type AuthService interface {
	// Register 用户注册
	// 返回：用户对象、JWT令牌、错误
	Register(email, password string) (model.User, string, error)

	// Login 用户登录
	// 返回：用户对象、JWT令牌、错误
	Login(email, password string) (model.User, string, error)
}

// authService 认证服务的具体实现
// 小写字母开头，包外不可见
type authService struct {
	// users 用户仓储，用于访问用户数据
	users repository.UserRepository

	// jwtSecret JWT 签名密钥
	// 用于生成和验证 JWT 令牌的安全性
	// 必须保密！泄露会导致他人可以伪造令牌
	jwtSecret []byte

	// tokenTTL JWT 令牌的有效期（Time To Live）
	// 例如 24*time.Hour 表示令牌 24 小时后过期
	tokenTTL time.Duration
}

// NewAuthService 创建认证服务实例
// 这是构造函数，返回接口类型
func NewAuthService(users repository.UserRepository, jwtSecret []byte, tokenTTL time.Duration) AuthService {
	return &authService{
		users:     users,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// Register 实现用户注册逻辑
func (s *authService) Register(email, password string) (model.User, string, error) {
	// 数据清理和标准化
	// TrimSpace: 去除首尾空格，防止 "user@example.com " 和 "user@example.com" 被当作不同邮箱
	// ToLower: 转为小写，确保邮箱不区分大小写（User@Example.com 和 user@example.com 是同一个）
	email = strings.TrimSpace(strings.ToLower(email))

	// 数据验证：邮箱和密码不能为空
	if email == "" || password == "" {
		return model.User{}, "", errors.New("email and password required")
	}

	// 验证邮箱是否注册过

	// bcrypt.GenerateFromPassword 生成密码哈希
	// 参数说明：
	// - []byte(password): 将字符串转为字节数组
	// - bcrypt.DefaultCost: 加密强度（默认是 10，越大越安全但越慢）
	// bcrypt 的特点：
	// 1. 单向加密：无法从哈希值还原密码
	// 2. 加盐（salt）：即使相同密码，每次生成的哈希值也不同
	// 3. 慢速算法：故意设计得很慢，防止暴力破解
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, "", err
	}

	// 调用仓储层创建用户
	// 注意：存储的是哈希值 string(hash)，不是明文密码！
	u, err := s.users.Create(email, string(hash))
	if err != nil {
		return model.User{}, "", err
	}

	// 注册成功后，立即颁发 JWT 令牌
	// 这样用户注册后就自动登录了，提供更好的用户体验
	tok, err := s.issueToken(u)
	return u, tok, err
}

// Login 实现用户登录逻辑
func (s *authService) Login(email, password string) (model.User, string, error) {
	// 同样对邮箱进行标准化处理
	email = strings.TrimSpace(strings.ToLower(email))

	// 根据邮箱查询用户
	u, err := s.users.GetByEmail(email)
	if err != nil {
		// 注意：不管是用户不存在还是其他错误，都返回相同的错误信息
		// 这是安全最佳实践：不要泄露"用户是否存在"的信息
		// 否则攻击者可以枚举有效的邮箱地址
		return model.User{}, "", errors.New("invalid credentials")
	}

	// bcrypt.CompareHashAndPassword 验证密码
	// 参数1：数据库中存储的哈希值
	// 参数2：用户输入的明文密码
	// 如果密码正确返回 nil，否则返回错误
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		// 密码错误，返回相同的错误信息（同样是安全考虑）
		return model.User{}, "", errors.New("invalid credentials")
	}

	// 验证通过，颁发 JWT 令牌
	tok, err := s.issueToken(u)
	return u, tok, err
}

// customClaims JWT 令牌中存储的自定义声明（Claims）
// JWT 由三部分组成：Header（头部）、Payload（载荷）、Signature（签名）
// Claims 就是 Payload 中存储的信息
type customClaims struct {
	// Email 用户邮箱（自定义字段）
	Email string `json:"email"`

	// jwt.RegisteredClaims 嵌入标准声明
	// Go 的嵌入（embedding）特性：customClaims 自动拥有 RegisteredClaims 的所有字段
	// RegisteredClaims 包含：
	// - Subject (sub): 主题，通常是用户 ID
	// - IssuedAt (iat): 签发时间
	// - ExpiresAt (exp): 过期时间
	// - Issuer (iss): 签发者
	jwt.RegisteredClaims
}

// issueToken 颁发 JWT 令牌
// 这是一个私有方法（小写字母开头），只在 service 内部使用
func (s *authService) issueToken(u model.User) (string, error) {
	now := time.Now()

	// 构建 JWT Claims（声明）
	claims := customClaims{
		Email: u.Email, // 自定义字段：存储用户邮箱
		RegisteredClaims: jwt.RegisteredClaims{
			// Subject（主题）：通常存储用户 ID
			// 后续请求时可以从 JWT 中提取用户 ID，知道是哪个用户在访问
			Subject: u.ID,

			// IssuedAt（签发时间）：令牌的创建时间
			IssuedAt: jwt.NewNumericDate(now),

			// ExpiresAt（过期时间）：令牌的有效期
			// now.Add(s.tokenTTL): 当前时间 + 有效期（如 24 小时）
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),

			// Issuer（签发者）：标识是哪个应用签发的令牌
			Issuer: "kanban_api",
		},
	}

	// jwt.NewWithClaims 创建令牌
	// jwt.SigningMethodHS256: 使用 HMAC-SHA256 算法签名
	// HS256 是对称加密：签名和验证使用同一个密钥
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// SignedString 生成最终的 JWT 字符串
	// 使用密钥对令牌进行签名，确保令牌不被篡改
	// 返回的字符串格式：header.payload.signature（三部分用 . 分隔）
	return token.SignedString(s.jwtSecret)
}

// MustJWTSecret 获取 JWT 密钥
// Must 前缀是 Go 的命名惯例，表示"必须成功，失败就 panic"
// 不过这个函数实际上总是返回一个值（使用默认值兜底）
func MustJWTSecret() []byte {
	// os.Getenv 从环境变量中读取 JWT_SECRET
	// 环境变量是配置应用的常用方式，适合存储敏感信息
	sec := os.Getenv("JWT_SECRET")

	// 如果环境变量未设置，使用默认值
	// 注意：默认值只适合开发环境，生产环境必须设置真实的密钥！
	if sec == "" {
		sec = "dev-secret"
	}

	// 转换为字节数组返回
	return []byte(sec)
}
