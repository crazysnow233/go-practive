# Kanban API - Go 语言看板系统

> 🎓 这是一个**初学者友好**的 Go 语言项目，包含详细的代码注释，帮助你理解 Go Web 开发的核心概念。

## 📖 项目简介

这是一个基于 Go 语言开发的看板（Kanban）管理系统的后端 API，实现了用户认证和看板管理的核心功能。

### 核心功能

- ✅ 用户注册和登录
- ✅ JWT 令牌认证
- ✅ 看板的增删改查（CRUD）
- ✅ RESTful API 设计
- ✅ SQLite 数据持久化

## 🛠 技术栈

| 技术 | 说明 | 用途 |
|------|------|------|
| **Go 1.25** | 编程语言 | 项目主语言 |
| **Gin** | Web 框架 | HTTP 路由和中间件 |
| **GORM** | ORM 框架 | 数据库操作 |
| **SQLite** | 数据库 | 数据持久化 |
| **JWT** | 认证方案 | 用户身份验证 |
| **bcrypt** | 加密算法 | 密码哈希 |

## 📁 项目结构（分层架构）

```
kanban_api/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口，应用启动
├── internal/                     # 内部代码（不能被外部导入）
│   ├── model/                   # 【数据模型层】
│   │   ├── user.go              # 用户数据结构
│   │   └── board.go             # 看板数据结构
│   ├── repository/              # 【数据访问层】
│   │   ├── id.go                # ID 生成工具
│   │   ├── user.go              # 用户数据访问（内存）
│   │   ├── board.go             # 看板数据访问（内存）
│   │   └── board_sqlite.go      # 看板数据访问（SQLite）
│   ├── service/                 # 【业务逻辑层】
│   │   ├── auth.go              # 认证业务逻辑
│   │   └── board.go             # 看板业务逻辑
│   ├── middleware/              # 【中间件层】
│   │   ├── requestid.go         # 请求 ID 追踪
│   │   ├── logger.go            # 日志记录
│   │   ├── error.go             # 错误恢复
│   │   └── auth.go              # JWT 认证
│   └── http/                    # 【HTTP 处理层】
│       ├── auth_handler.go      # 认证接口处理
│       └── board_handler.go     # 看板接口处理
├── go.mod                        # Go 模块定义
├── go.sum                        # 依赖版本锁定
├── kanban.db                     # SQLite 数据库文件
└── README.md                     # 项目文档（本文件）
```

### 🏗 分层架构详解

这个项目采用经典的**分层架构**（Layered Architecture），每一层有明确的职责：

```
请求流向：
HTTP 请求 → Middleware → Handler → Service → Repository → Database
                ↓           ↓         ↓          ↓
              中间件      控制器    业务逻辑   数据访问
```

#### 1. **Model 层（数据模型）**
- **职责**：定义数据结构
- **特点**：纯数据，无业务逻辑
- **示例**：`User`、`Board` 结构体

#### 2. **Repository 层（数据访问）**
- **职责**：封装所有数据库操作
- **特点**：提供接口，隐藏实现细节
- **好处**：可以轻松切换存储方式（内存 ↔ 数据库）
- **模式**：Repository Pattern（仓储模式）

#### 3. **Service 层（业务逻辑）**
- **职责**：实现业务规则和流程
- **特点**：调用 Repository，处理复杂逻辑
- **示例**：密码加密、数据验证、JWT 生成

#### 4. **Middleware 层（中间件）**
- **职责**：请求拦截和预处理
- **特点**：可复用、可组合
- **示例**：认证、日志、错误处理

#### 5. **Handler 层（HTTP 处理）**
- **职责**：处理 HTTP 请求和响应
- **特点**：解析参数、调用 Service、返回 JSON
- **对应**：MVC 中的 Controller

## 🚀 快速开始

### 前置要求

- Go 1.25 或更高版本
- Git（可选）

### 安装依赖

```bash
# 进入项目目录
cd kanban_api

# 下载依赖
go mod download
```

### 运行项目

```bash
# 运行服务器
go run cmd/server/main.go

# 或者先编译再运行
go build -o kanban-server cmd/server/main.go
./kanban-server
```

服务器将在 `http://localhost:8080` 启动。

### 环境变量（可选）

```bash
# 设置 JWT 密钥（生产环境必须设置）
export JWT_SECRET="your-secret-key-here"

# 然后运行
go run cmd/server/main.go
```

## 📡 API 接口文档

### 基础 URL

```
http://localhost:8080/api/v1
```

### 认证接口（公共，无需登录）

#### 1. 用户注册

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

**响应示例：**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "createdAt": "2024-01-01T12:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### 2. 用户登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

### 看板接口（需要认证）

> ⚠️ 所有看板接口都需要在请求头中携带 JWT 令牌

```http
Authorization: Bearer <your_jwt_token>
```

#### 3. 获取所有看板

```http
GET /api/v1/boards
Authorization: Bearer <token>
```

#### 4. 获取单个看板

```http
GET /api/v1/boards/:id
Authorization: Bearer <token>
```

#### 5. 创建看板

```http
POST /api/v1/boards
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "我的第一个看板"
}
```

#### 6. 更新看板

```http
PUT /api/v1/boards/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "更新后的标题"
}
```

#### 7. 删除看板

```http
DELETE /api/v1/boards/:id
Authorization: Bearer <token>
```

**响应：** 204 No Content

## 🧪 测试接口（使用 curl）

### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'
```

### 2. 登录并保存 Token

```bash
# 登录并提取 token（需要 jq 工具）
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }' | jq -r '.data.token')

echo $TOKEN
```

### 3. 创建看板

```bash
curl -X POST http://localhost:8080/api/v1/boards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "我的第一个看板"
  }'
```

### 4. 获取所有看板

```bash
curl -X GET http://localhost:8080/api/v1/boards \
  -H "Authorization: Bearer $TOKEN"
```

## 🎓 Go 语言关键概念（本项目涉及）

### 1. 接口（Interface）

**定义：** 一组方法的集合，任何类型只要实现了这些方法就满足该接口。

**示例：**
```go
// 定义接口
type UserRepository interface {
    Create(email, password string) (User, error)
    GetByEmail(email string) (User, error)
}

// 实现接口
type memUserRepo struct { /* ... */ }
func (r *memUserRepo) Create(...) (User, error) { /* ... */ }
func (r *memUserRepo) GetByEmail(...) (User, error) { /* ... */ }
```

**好处：**
- ✅ 依赖倒置：高层依赖接口而非实现
- ✅ 易于测试：可以创建 Mock 实现
- ✅ 灵活切换：可以随时更换实现

### 2. 结构体标签（Struct Tags）

**用途：** 为结构体字段添加元数据，供其他库使用。

```go
type User struct {
    ID    string `json:"id" gorm:"primaryKey"`
    Email string `json:"email" gorm:"unique"`
}
```

**常见标签：**
- `json:"name"` - JSON 序列化/反序列化
- `gorm:"primaryKey"` - GORM 数据库映射
- `binding:"required"` - Gin 参数验证

### 3. defer 延迟执行

**作用：** 延迟函数执行，常用于资源清理。

```go
func example() {
    mu.Lock()
    defer mu.Unlock()  // 函数返回前自动释放锁
    
    // 中间无论如何返回，都会执行 Unlock()
    if err != nil {
        return  // Unlock() 会在这里自动执行
    }
}
```

**执行时机：** 函数返回前（包括 panic 时）

### 4. 错误处理

**Go 的哲学：** 显式处理错误，不使用异常（try-catch）。

```go
// 多返回值：(结果, 错误)
user, err := userRepo.GetByEmail(email)
if err != nil {
    return nil, err  // 向上传递错误
}
```

**panic/recover：** 仅用于不可恢复的严重错误。

### 5. 并发控制

**sync.RWMutex：** 读写锁，保护共享数据。

```go
mu.RLock()           // 获取读锁（允许多个读）
defer mu.RUnlock()   // 释放读锁

mu.Lock()            // 获取写锁（独占）
defer mu.Unlock()    // 释放写锁
```

**为什么需要？** Go 的 map 不是并发安全的，多个 goroutine 同时访问会崩溃。

### 6. 中间件模式

**概念：** 在请求到达处理器前/后执行的函数。

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求前
        start := time.Now()
        
        c.Next()  // 执行下一个中间件/处理器
        
        // 请求后
        latency := time.Since(start)
        log.Printf("耗时: %v", latency)
    }
}
```

**执行顺序：** 洋葱模型（先进后出）

## ⚠️ 重要注意事项

### 安全相关

#### 1. 密码安全

```go
// ✅ 正确：使用 bcrypt 哈希
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ❌ 错误：明文存储密码
user.Password = password  // 绝对不要这样做！
```

**要点：**
- 永远不要存储明文密码
- 使用 bcrypt、argon2 等慢速哈希算法
- 不要使用 MD5、SHA1（太快，易被破解）

#### 2. JWT 密钥管理

```bash
# ✅ 生产环境：必须设置环境变量
export JWT_SECRET="your-very-long-random-secret-key"

# ❌ 开发环境：可以使用默认值（仅限开发！）
# 代码中的 "dev-secret" 只适合学习测试
```

**要点：**
- 密钥至少 32 字符
- 使用随机生成的字符串
- 不要提交到版本控制系统

#### 3. SQL 注入防护

```go
// ✅ 正确：使用参数化查询
db.First(&user, "id=?", id)

// ❌ 错误：字符串拼接
db.First(&user, "id=" + id)  // 容易被 SQL 注入攻击
```

**要点：**
- 使用 GORM 的占位符 `?`
- 不要自己拼接 SQL 字符串

### 并发安全

#### 1. Map 并发访问

```go
// ❌ 不安全：Go 的 map 不是并发安全的
var users = make(map[string]User)
// 多个 goroutine 同时读写会崩溃

// ✅ 安全：使用 sync.RWMutex 保护
type UserRepo struct {
    mu    sync.RWMutex
    users map[string]User
}
```

#### 2. 锁的使用

```go
// ✅ 正确：使用 defer 确保解锁
mu.Lock()
defer mu.Unlock()

// ❌ 错误：可能忘记解锁
mu.Lock()
// ... 一堆代码
mu.Unlock()  // 如果中间 return 了，就不会执行
```

### 错误处理

#### 1. 不要忽略错误

```go
// ❌ 错误：忽略错误
user, _ := repo.GetByEmail(email)

// ✅ 正确：检查错误
user, err := repo.GetByEmail(email)
if err != nil {
    return err
}
```

#### 2. 错误信息安全

```go
// ❌ 错误：暴露敏感信息
return errors.New("用户 test@example.com 不存在")  // 泄露邮箱是否注册

// ✅ 正确：使用通用错误消息
return errors.New("invalid credentials")  // 不泄露具体信息
```

### HTTP 处理

#### 1. 记得 return

```go
func handler(c *gin.Context) {
    if err != nil {
        c.JSON(400, gin.H{"error": "bad request"})
        return  // ⚠️ 必须 return，否则会继续执行
    }
    c.JSON(200, gin.H{"data": "success"})
}
```

#### 2. 使用正确的 HTTP 状态码

```go
// 200 OK - 成功
// 201 Created - 创建成功
// 204 No Content - 删除成功
// 400 Bad Request - 请求错误
// 401 Unauthorized - 未认证
// 403 Forbidden - 无权限
// 404 Not Found - 未找到
// 500 Internal Server Error - 服务器错误
```

## 🎯 学习路径建议

### 第一步：理解项目结构（1-2天）

1. 从 `main.go` 开始，理解程序启动流程
2. 按照分层顺序阅读代码：Model → Repository → Service → Handler
3. 重点理解每一层的职责

### 第二步：运行和测试（1天）

1. 运行项目
2. 使用 curl 或 Postman 测试所有接口
3. 观察日志输出
4. 查看数据库文件 `kanban.db`

### 第三步：修改代码（2-3天）

尝试添加新功能：

**简单功能：**
- [ ] 为看板添加描述字段
- [ ] 添加"获取用户信息"接口
- [ ] 修改 JWT 过期时间

**中等功能：**
- [ ] 添加卡片（Card）功能
- [ ] 为看板添加颜色标签
- [ ] 实现分页查询

**进阶功能：**
- [ ] 添加用户权限控制
- [ ] 实现看板共享功能
- [ ] 切换到 PostgreSQL 数据库

### 第四步：深入学习（持续）

1. **Go 并发：** 学习 goroutine、channel
2. **测试：** 编写单元测试和集成测试
3. **部署：** 学习 Docker、CI/CD
4. **监控：** 添加性能监控和日志分析

## 🔧 常见问题

### Q1: 为什么用 `internal` 目录？

**A:** `internal` 是 Go 的特殊目录，表示这些包只能被本项目导入，不能被外部项目使用。这是一种封装机制。

### Q2: 接口和结构体有什么区别？

**A:** 
- **结构体（struct）**：定义数据结构，类似于类
- **接口（interface）**：定义行为规范，只声明方法不实现

### Q3: 为什么要用三层架构？

**A:** 
- ✅ **职责分离**：每层专注自己的事情
- ✅ **易于测试**：可以单独测试每一层
- ✅ **易于维护**：修改一层不影响其他层
- ✅ **可扩展**：容易添加新功能

### Q4: JWT 是如何工作的？

**A:** 
1. 用户登录成功，服务器生成 JWT 令牌
2. 客户端保存令牌（通常在 localStorage）
3. 后续请求在 Header 中携带令牌
4. 服务器验证令牌的签名和有效期
5. 如果有效，提取用户信息（userID）

### Q5: bcrypt 比 MD5 好在哪里？

**A:** 
- ⏱ **慢速算法**：故意设计得很慢，防止暴力破解
- 🧂 **自动加盐**：每次生成的哈希都不同
- 🔧 **可调强度**：可以通过 cost 参数增加计算复杂度

### Q6: 为什么数据库操作返回接口类型？

**A:** 依赖倒置原则。调用者依赖接口，不依赖具体实现，这样可以轻松切换实现（内存 → SQLite → MySQL）。

## 📚 推荐学习资源

### 官方文档

- [Go 官方教程](https://go.dev/tour/)
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)

### 书籍推荐

- 《Go 程序设计语言》（The Go Programming Language）
- 《Go Web 编程》

### 视频教程

- B站搜索"Go 语言入门"
- YouTube 搜索"Go Web Development"

## 🤝 贡献

这是一个学习项目，欢迎提出改进建议！

## 📄 许可证

MIT License

---

**🎉 祝你学习愉快！遇到问题随时查看代码注释或搜索相关文档。**

记住：编程是实践的艺术，多写多练才能真正掌握！💪

