// Package model 定义看板应用的数据模型
package model

import "time"

// Board 看板结构体，代表一个看板（Kanban Board）
// 每个看板有标题、创建时间和更新时间
type Board struct {
	// ID 看板的唯一标识符，使用 UUID 格式
	ID string `json:"id"`

	// Title 看板的标题，例如："我的待办事项"、"项目A任务板"
	Title string `json:"title"`

	// CreatedAt 看板的创建时间
	// 创建时设置一次，之后不再修改
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt 看板的最后更新时间
	// 每次修改看板信息时都要更新这个字段
	UpdatedAt time.Time `json:"updatedAt"`
}
