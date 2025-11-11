// Package repository 提供用户数据的存储和访问接口
package repository

import (
	"errors"
	"kanban_api/internal/model"
	"sync"
	"time"
)

// ErrUserExists 当尝试创建已存在的用户时返回的错误
// 在 Go 中，习惯用 Err 开头命名错误变量
var ErrUserExists = errors.New("user already exists")

// UserRepository 用户仓储接口
// 接口（interface）定义了一组方法，任何实现了这些方法的类型都满足这个接口
// 使用接口的好处：
// 1. 依赖倒置：高层代码依赖接口而不是具体实现
// 2. 易于测试：可以创建 mock 实现来测试
// 3. 灵活切换：可以轻松切换不同的存储方式（内存、MySQL、MongoDB 等）
type UserRepository interface {
	// Create 创建新用户
	// 参数：email(邮箱), passwordHash(密码哈希值)
	// 返回：创建的用户对象, 错误信息
	Create(email, passwordHash string) (model.User, error)

	// GetByEmail 通过邮箱查询用户
	// 用于登录时验证用户
	GetByEmail(email string) (model.User, error)

	// GetByID 通过 ID 查询用户
	// 用于鉴权后获取用户信息
	GetByID(id string) (model.User, error)
}

// memUserRepo 是 UserRepository 接口的内存实现
// 数据存储在内存中，程序重启后数据会丢失
// 适用于：开发测试、原型演示
// 小写字母开头表示这是包内私有的（不能被其他包访问）
type memUserRepo struct {
	// mu 读写锁，用于保护并发访问
	// sync.RWMutex 允许多个读操作同时进行，但写操作是独占的
	// 为什么需要锁？因为多个 HTTP 请求可能同时访问这个数据结构
	mu sync.RWMutex

	// users 存储所有用户，key 是用户 ID，value 是用户对象
	// map 是 Go 的哈希表数据结构，查询速度是 O(1)
	users map[string]model.User

	// emailIdx 邮箱索引，用于快速通过邮箱查找用户 ID
	// 这是一个常见的数据库优化技巧：建立索引加速查询
	emailIdx map[string]string // email -> userID
}

// NewMemUserRepo 创建一个新的内存用户仓储
// 这是一个构造函数（Go 中的惯例是用 New 开头）
// 返回类型是接口 UserRepository，而不是具体的 *memUserRepo
// 这样调用者只知道接口，不知道具体实现，实现了依赖倒置
func NewMemUserRepo() UserRepository {
	return &memUserRepo{
		// make 函数用于创建 map、slice、channel 等类型
		// 必须用 make 初始化 map，否则是 nil，无法使用
		users:    make(map[string]model.User),
		emailIdx: make(map[string]string),
	}
}

// Create 实现 UserRepository 接口的 Create 方法
// (r *memUserRepo) 是接收者（receiver），表示这个方法属于 memUserRepo 类型
// 类似于其他语言中的 this 或 self
func (r *memUserRepo) Create(email, password string) (model.User, error) {
	// Lock() 获取写锁，确保同一时刻只有一个 goroutine 可以修改数据
	r.mu.Lock()

	// defer 关键字：延迟执行，函数返回前自动调用
	// 这里确保无论函数如何结束（正常返回或 panic），都会释放锁
	// defer 是 Go 中非常重要的特性，常用于资源清理
	defer r.mu.Unlock()

	// 检查邮箱是否已存在
	// _, ok := map[key] 是 Go 的惯用法，检查 key 是否存在
	// _ 表示忽略返回的值，我们只关心 key 存不存在
	if _, ok := r.emailIdx[email]; ok {
		// 如果邮箱已存在，返回空用户对象和错误
		// Go 支持多返回值，常用模式是 (结果, 错误)
		return model.User{}, ErrUserExists
	}

	// 创建新用户对象
	u := model.User{
		ID:           generateID(), // 生成唯一 ID
		Email:        email,        // 保存邮箱
		PasswordHash: password,     // 保存密码哈希（不是明文！）
		CreatedAt:    time.Now(),   // 记录创建时间
	}

	// 保存到主存储
	r.users[u.ID] = u
	// 同时更新邮箱索引，方便通过邮箱查询
	r.emailIdx[email] = u.ID

	// 返回创建的用户和 nil（表示没有错误）
	return u, nil
}

// GetByEmail 通过邮箱查询用户
func (r *memUserRepo) GetByEmail(email string) (model.User, error) {
	// RLock() 获取读锁
	// 读锁的特点：多个 goroutine 可以同时持有读锁
	// 但如果有写锁，读锁会等待
	// 这样可以提高并发读取的性能
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 先从邮箱索引中查找用户 ID
	id, ok := r.emailIdx[email]
	if !ok {
		// 如果找不到，返回错误
		return model.User{}, ErrNotFound
	}

	// 根据 ID 从主存储中获取完整的用户对象
	// nil 表示没有错误
	return r.users[id], nil
}

// GetByID 通过用户 ID 查询用户
func (r *memUserRepo) GetByID(id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 直接从主存储中查询
	u, ok := r.users[id]
	if !ok {
		return model.User{}, ErrNotFound
	}

	return u, nil
}
