// Package main 是程序的入口包
// 每个 Go 程序都必须有一个 main 包和一个 main 函数
package main

import (
	"github.com/gin-gonic/gin"       // Gin Web 框架
	httpx "kanban_api/internal/http" // 导入时使用别名 httpx，避免与标准库 http 冲突
	"kanban_api/internal/middleware"
	"kanban_api/internal/repository"
	"kanban_api/internal/service"
	"log"
	"time"
)

// main 函数是程序的入口点
// 程序启动时会自动执行这个函数
func main() {
	// ========== 第一步：初始化数据访问层（Repository） ==========
	// 采用"依赖注入"的方式，从底层往上层构建
	userRepo, err := repository.NewSQLiteUserRepo("file:kanban.db?cache=shared&_fk=1")
	if err != nil {
		// log.Fatal 会打印错误信息并退出程序（调用 os.Exit(1)）
		// 适用于启动时的致命错误
		log.Fatal(err)
	}
	// 创建看板仓储（SQLite 数据库实现）
	// 连接字符串参数说明：
	// - file:kanban.db: 数据库文件路径
	// - cache=shared: 启用共享缓存，多个连接可以共享缓存
	// - _fk=1: 启用外键约束
	boardRepo, err := repository.NewSQLiteBoardRepo("file:kanban.db?cache=shared&_fk=1")
	if err != nil {
		// log.Fatal 会打印错误信息并退出程序（调用 os.Exit(1)）
		// 适用于启动时的致命错误
		log.Fatal(err)
	}

	// 如果想使用内存实现（不持久化），可以取消下面这行的注释：
	// boardRepo := repository.NewMemBoardRepo()
	// 创建用户仓储（内存实现）
	//userRepo := repository.NewMemUserRepo()

	// ========== 第二步：初始化业务逻辑层（Service） ==========

	// 获取 JWT 密钥（从环境变量读取）
	jwtSecret := service.MustJWTSecret()

	// 创建认证服务
	// 参数：用户仓储、JWT密钥、令牌有效期（24小时）
	authSvc := service.NewAuthService(userRepo, jwtSecret, 24*time.Hour)

	// 创建看板服务
	boardSvc := service.NewBoardService(boardRepo)

	// ========== 第三步：初始化 HTTP 处理器层（Handler） ==========

	// 创建认证处理器
	authH := httpx.NewAuthHandler(authSvc)

	// 创建看板处理器
	boardH := httpx.NewBoardHandler(boardSvc)

	// ========== 第四步：配置路由和中间件 ==========

	// gin.New() 创建一个不带默认中间件的 Gin 引擎
	// 对比：gin.Default() 会自动添加 Logger 和 Recovery 中间件
	r := gin.New()

	// r.Use() 注册全局中间件
	// 中间件按注册顺序执行
	// 执行顺序：RequestID -> Logger -> Recovery -> RecoverJSON -> 处理器
	r.Use(
		middleware.RequestID(),   // 为每个请求生成唯一 ID
		middleware.Logger(),      // 记录请求日志
		gin.Recovery(),           // Gin 自带的 panic 恢复中间件
		middleware.RecoverJSON(), // 自定义的 JSON 格式错误恢复
	)

	// 注意：gin.Recovery() 和 middleware.RecoverJSON() 功能类似
	// gin.Recovery() 会恢复 panic 但返回纯文本错误
	// middleware.RecoverJSON() 返回 JSON 格式错误
	// 实际上只需要一个就够了，这里两个都用是为了演示

	// ========== 第五步：注册路由 ==========

	// r.Group() 创建路由组
	// 路由组的作用：
	// 1. 统一路径前缀（这里是 "api/v1"）
	// 2. 统一应用中间件

	// 公共路由组：不需要认证
	// 包含：注册、登录接口
	public := r.Group("api/v1")
	authH.RegisterRoutes(public)

	// 私有路由组：需要认证
	// middleware.AuthRequired(jwtSecret) 是认证中间件
	// 只有携带有效 JWT 令牌的请求才能访问这组路由
	private := r.Group("api/v1", middleware.AuthRequired(jwtSecret))
	boardH.Register(private)

	// ========== 第六步：启动 HTTP 服务器 ==========

	log.Println("listen on :8080")
	log.Println("公共接口（无需登录）：")
	log.Println("  POST http://localhost:8080/api/v1/auth/register")
	log.Println("  POST http://localhost:8080/api/v1/auth/login")
	log.Println("私有接口（需要登录）：")
	log.Println("  GET    http://localhost:8080/api/v1/boards")
	log.Println("  POST   http://localhost:8080/api/v1/boards")
	log.Println("  GET    http://localhost:8080/api/v1/boards/:id")
	log.Println("  PUT    http://localhost:8080/api/v1/boards/:id")
	log.Println("  DELETE http://localhost:8080/api/v1/boards/:id")

	// r.Run() 启动 HTTP 服务器
	// 参数 ":8080" 表示监听所有网络接口的 8080 端口
	// 等价于 "0.0.0.0:8080"
	// 如果只想本地访问，可以用 "127.0.0.1:8080" 或 "localhost:8080"
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	// 注意：r.Run() 会阻塞，下面的代码不会执行
	// 除非服务器停止（例如按 Ctrl+C）
}
