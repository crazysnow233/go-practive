// Package model 包含了应用程序的核心数据模型（实体）
// 这一层定义了我们要操作的数据结构，不包含任何业务逻辑
package model

import "time"

// User 用户结构体，代表系统中的一个用户
// 在 Go 中，结构体（struct）类似于其他语言中的类（class）
type User struct {
	// ID 用户的唯一标识符，使用 UUID 格式
	// `json:"id"` 是结构体标签（struct tag），告诉 Go 在序列化为 JSON 时使用 "id" 作为字段名
	ID string `json:"id"`

	// Email 用户的邮箱地址，用于登录
	// `json:"email"` 表示 JSON 中的字段名为 "email"
	Email string `json:"email"`

	// PasswordHash 用户密码的哈希值（不是明文密码！）
	// 重要：永远不要存储明文密码，我们使用 bcrypt 算法生成的哈希值
	// `json:"-"` 这个特殊标签表示：在序列化为 JSON 时忽略这个字段，保护用户密码安全
	PasswordHash string `json:"-"`

	// CreatedAt 用户创建时间
	// time.Time 是 Go 内置的时间类型
	// `json:"createdAt"` 表示 JSON 中使用驼峰命名
	CreatedAt time.Time `json:"createdAt"`
}
