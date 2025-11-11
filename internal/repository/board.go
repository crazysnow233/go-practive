package repository

import (
	"errors"
	"kanban_api/internal/model"
	"sync"
	"time"
)

// ErrNotFound 当查询的资源不存在时返回的错误
var ErrNotFound = errors.New("not found")

// BoardRepository 看板仓储接口
// 定义了对看板数据的 CRUD（增删改查）操作
type BoardRepository interface {
	// List 列出所有看板
	List() ([]model.Board, error)

	// Get 获取单个看板
	Get(id string) (model.Board, error)

	// Create 创建新看板
	Create(title string) (model.Board, error)

	// Update 更新看板信息
	Update(id, title string) (model.Board, error)

	// Delete 删除看板
	Delete(id string) error
}

// memBoardRepo 看板仓储的内存实现
// 与 memUserRepo 类似，数据存在内存中
type memBoardRepo struct {
	mu     sync.RWMutex           // 读写锁，保护并发访问
	boards map[string]model.Board // 存储所有看板，key 是看板 ID
}

// NewMemBoardRepo 创建一个新的内存看板仓储
func NewMemBoardRepo() BoardRepository {
	return &memBoardRepo{
		boards: make(map[string]model.Board),
	}
}

// List 列出所有看板
func (r *memBoardRepo) List() ([]model.Board, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 创建一个切片（slice）来存储结果
	// make([]model.Board, 0, len(r.boards)) 的含义：
	// - []model.Board: 切片类型
	// - 0: 初始长度为 0（当前没有元素）
	// - len(r.boards): 容量（capacity）为 boards 的数量，避免多次扩容
	out := make([]model.Board, 0, len(r.boards))

	// range 用于遍历 map、slice、channel 等
	// for key, value := range map 会遍历所有键值对
	// 这里用 _ 忽略 key（看板ID），只关心 value（看板对象）
	for _, b := range r.boards {
		// append 向切片追加元素
		out = append(out, b)
	}

	return out, nil
}

// Get 根据 ID 获取单个看板
func (r *memBoardRepo) Get(id string) (model.Board, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	b, ok := r.boards[id]
	if !ok {
		// 注意：这里原代码返回 nil 作为错误可能是个 bug
		// 应该返回 ErrNotFound 才对
		return model.Board{}, ErrNotFound
	}
	return b, nil
}

// Create 创建新看板
func (r *memBoardRepo) Create(title string) (model.Board, error) {
	// 获取当前时间，创建时间和更新时间都设置为当前时间
	now := time.Now()

	// 构建看板对象
	b := model.Board{
		ID:        generateID(), // 生成唯一 ID
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 写操作需要获取写锁
	r.mu.Lock()
	r.boards[b.ID] = b
	r.mu.Unlock()

	return b, nil
}

// Update 更新看板信息
func (r *memBoardRepo) Update(id, title string) (model.Board, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 先检查看板是否存在
	b, ok := r.boards[id]
	if !ok {
		return model.Board{}, ErrNotFound
	}

	// 更新标题和更新时间
	b.Title = title
	b.UpdatedAt = time.Now()

	// 注意：在 Go 中，从 map 取出的是值的副本
	// 所以修改 b 后，需要重新放回 map 中
	r.boards[id] = b

	return b, nil
}

// Delete 删除看板
func (r *memBoardRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 先检查是否存在
	if _, ok := r.boards[id]; !ok {
		return ErrNotFound
	}

	// delete 是 Go 内置函数，用于删除 map 中的键值对
	delete(r.boards, id)

	return nil
}
