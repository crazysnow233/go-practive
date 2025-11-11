// Package repository 提供看板仓储的数据库实现
package repository

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban_api/internal/model"
	"time"
)

// sqliteBoardRepo 是 BoardRepository 接口的 SQLite 数据库实现
// 与内存实现不同，数据会持久化到磁盘文件中
type sqliteBoardRepo struct {
	// db 是 GORM 的数据库连接对象
	// GORM 是 Go 语言最流行的 ORM（对象关系映射）库
	// ORM 让我们用面向对象的方式操作数据库，而不用写 SQL
	db *gorm.DB
}

// boardRow 数据库表结构
// 这个结构体对应数据库中的一张表
// GORM 会自动根据这个结构体创建表（auto migration）
type boardRow struct {
	// ID 主键
	// `gorm:"primaryKey"` 是 GORM 的标签，表示这是主键
	ID string `gorm:"primaryKey"`

	// Title 看板标题
	// 没有标签时，GORM 会自动将字段名转为蛇形命名（title）
	Title string

	// CreatedAt 创建时间
	// GORM 会自动识别 CreatedAt 字段，在插入时自动设置
	CreatedAt time.Time

	// UpdatedAt 更新时间
	// GORM 会自动识别 UpdatedAt 字段，在更新时自动刷新
	UpdatedAt time.Time
}

// NewSQLiteBoardRepo 创建一个新的 SQLite 看板仓储
// 参数 path 是数据库文件路径，例如："file:kanban.db?cache=shared&_fk=1"
// 返回 BoardRepository 接口，使用者不需要知道底层是 SQLite
func NewSQLiteBoardRepo(path string) (BoardRepository, error) {
	// gorm.Open 打开数据库连接
	// sqlite.Open(path) 指定使用 SQLite 驱动
	// &gorm.Config{} 是 GORM 的配置选项（这里使用默认配置）
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		// 如果连接失败，返回错误
		return nil, err
	}

	// AutoMigrate 自动迁移数据库表结构
	// 它会根据 boardRow 结构体自动创建表
	// 如果表已存在，会根据结构体更新表结构（增加新字段等）
	// 注意：传入的是指针 &boardRow{}
	if err := db.AutoMigrate(&boardRow{}); err != nil {
		return nil, err
	}

	// 返回仓储实例
	return &sqliteBoardRepo{db: db}, nil
}

// toModel 将数据库行（boardRow）转换为业务模型（model.Board）
// 这是一个辅助方法，用于数据转换
// 为什么需要两个不同的结构体？
// - boardRow: 数据库层的表示，带有 GORM 标签
// - model.Board: 业务层的表示，带有 JSON 标签
// 这种分层设计让各层职责更清晰
func (r *sqliteBoardRepo) toModel(row boardRow) model.Board {
	return model.Board{
		ID:        row.ID,
		Title:     row.Title,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// List 查询所有看板
func (r *sqliteBoardRepo) List() ([]model.Board, error) {
	// 声明一个切片来接收查询结果
	var rows []boardRow

	// GORM 链式调用：
	// Order("created_at desc"): 按创建时间降序排序（最新的在前）
	// Find(&rows): 查询所有记录，结果存入 rows
	// .Error: 获取错误（GORM 用这种方式返回错误）
	if err := r.db.Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}

	// 将数据库行转换为业务模型
	out := make([]model.Board, 0, len(rows))
	for _, rw := range rows {
		out = append(out, r.toModel(rw))
	}

	return out, nil
}

// Get 根据 ID 查询单个看板
func (r *sqliteBoardRepo) Get(id string) (model.Board, error) {
	var rw boardRow

	// First 查询第一条匹配的记录
	// "id=?" 是 SQL 条件，? 是占位符
	// id 是占位符的值，GORM 会自动防止 SQL 注入
	// 相当于 SQL: SELECT * FROM board_rows WHERE id=? LIMIT 1
	if err := r.db.First(&rw, "id=?", id).Error; err != nil {
		// errors.Is 判断错误类型（Go 1.13+ 的标准错误处理方式）
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录不存在，返回我们自定义的 ErrNotFound
			return model.Board{}, ErrNotFound
		}
		// 其他错误直接返回
		return model.Board{}, err
	}

	// 将数据库行转换为业务模型
	return r.toModel(rw), nil
}
// Create 创建新看板
func (r *sqliteBoardRepo) Create(title string) (model.Board, error) {
	now := time.Now()

	// 构建数据库行对象
	rw := boardRow{
		ID:        generateID(), // 生成唯一 ID
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create 插入一条新记录
	// 相当于 SQL: INSERT INTO board_rows (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)
	if err := r.db.Create(&rw).Error; err != nil {
		return model.Board{}, err
	}

	// 返回转换后的业务模型
	return r.toModel(rw), nil
}

// Update 更新看板信息
func (r *sqliteBoardRepo) Update(id, title string) (model.Board, error) {
	var rw boardRow

	// 先查询记录是否存在
	if err := r.db.First(&rw, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Board{}, ErrNotFound
		}
		return model.Board{}, err
	}

	// 修改字段
	rw.Title = title
	rw.UpdatedAt = time.Now()

	// Save 更新记录
	// 相当于 SQL: UPDATE board_rows SET title=?, updated_at=? WHERE id=?
	// Save 会更新所有字段，即使字段值没变
	if err := r.db.Save(&rw).Error; err != nil {
		return model.Board{}, err
	}

	return r.toModel(rw), nil
}

// Delete 删除看板
func (r *sqliteBoardRepo) Delete(id string) error {
	// Delete 删除记录
	// 相当于 SQL: DELETE FROM board_rows WHERE id=?
	// 第一个参数 &boardRow{} 用于指定表名（GORM 会根据类型推断）
	res := r.db.Delete(&boardRow{}, "id=?", id)

	// 检查是否有错误
	if res.Error != nil {
		return res.Error
	}

	// RowsAffected 返回受影响的行数
	// 如果为 0，说明没有找到要删除的记录
	if res.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
