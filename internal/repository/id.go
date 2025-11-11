// Package repository 负责数据持久化和访问
// 这一层处理所有与数据存储相关的操作（数据库、内存等）
package repository

import "github.com/google/uuid"

// uuidNew 生成一个新的 UUID（通用唯一识别码）
// UUID 是一个 128 位的随机数，格式类似：550e8400-e29b-41d4-a716-446655440000
// 用于生成全局唯一的 ID，避免数据冲突
func uuidNew() string {
	return uuid.NewString()
}

// generateID 是一个辅助函数，用于生成唯一 ID
// 使用这个包装函数的好处：
// 1. 如果以后想换其他 ID 生成方式，只需修改这一处
// 2. 代码更清晰，表明这是在生成 ID
func generateID() string {
	return uuidNew()
}
