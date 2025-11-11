// Package http 看板处理器
package http

import (
	"github.com/gin-gonic/gin"
	"kanban_api/internal/service"
	"net/http"
)

// BoardHandler 看板处理器
// 处理看板相关的 HTTP 请求
type BoardHandler struct {
	svc service.BoardService
}

// NewBoardHandler 创建看板处理器实例
func NewBoardHandler(svc service.BoardService) *BoardHandler {
	return &BoardHandler{svc: svc}
}

// Register 注册路由
// 这里展示了 RESTful API 的标准设计：
// - GET /boards: 列出所有资源
// - POST /boards: 创建新资源
// - GET /boards/:id: 获取单个资源
// - PUT /boards/:id: 更新资源
// - DELETE /boards/:id: 删除资源
func (h *BoardHandler) Register(rg *gin.RouterGroup) {
	// GET 用于查询数据
	rg.GET("/boards", h.list)

	// POST 用于创建新资源
	rg.POST("/boards", h.create)

	// :id 是路径参数，会匹配任意值
	// 例如：/boards/123 中的 123 就是 id
	rg.GET("/boards/:id", h.get)

	// PUT 用于完整更新资源（替换整个资源）
	// PATCH 用于部分更新（只更新部分字段）
	// 这里用 PUT，虽然实际上只更新了 title 字段
	rg.PUT("/boards/:id", h.update)

	// DELETE 用于删除资源
	rg.DELETE("/boards/:id", h.delete)
}

// list 列出所有看板
// GET /api/v1/boards
func (h *BoardHandler) list(c *gin.Context) {
	// 调用 Service 层获取所有看板
	items, err := h.svc.ListBoards()
	if err != nil {
		// http.StatusInternalServerError = 500（服务器内部错误）
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回看板列表
	// items 是 []model.Board，会被自动序列化为 JSON 数组
	c.JSON(http.StatusOK, gin.H{"data": items})
}

// create 创建新看板
// POST /api/v1/boards
// 请求体：{"title": "我的看板"}
func (h *BoardHandler) create(c *gin.Context) {
	// 定义请求体结构
	var req struct {
		Title string `json:"title"`
	}

	// 解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// 调用 Service 层创建看板
	b, err := h.svc.CreateBoard(req.Title)
	if err != nil {
		// 注意：这里缺少 return
		// 如果不加 return，会继续执行下面的代码，导致返回两个响应（会报错）
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // 应该加上 return
	}

	// 创建成功，返回 201
	c.JSON(http.StatusCreated, gin.H{"data": b})
}

// get 获取单个看板
// GET /api/v1/boards/:id
func (h *BoardHandler) get(c *gin.Context) {
	// c.Param 获取路径参数
	// 例如：GET /boards/123 中，c.Param("id") 返回 "123"
	// 注意：路径参数总是字符串类型
	//
	// 对比：
	// - c.Param("id"): 路径参数，如 /boards/:id
	// - c.Query("id"): 查询参数，如 /boards?id=123
	// - c.GetHeader("id"): 请求头
	id := c.Param("id")

	// 调用 Service 层获取看板
	b, err := h.svc.GetBoard(id)
	if err != nil {
		// http.StatusNotFound = 404（未找到）
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// 返回看板数据
	c.JSON(http.StatusOK, gin.H{"data": b})
}

// update 更新看板
// PUT /api/v1/boards/:id
// 请求体：{"title": "新标题"}
func (h *BoardHandler) update(c *gin.Context) {
	// 获取路径参数（看板 ID）
	id := c.Param("id")

	// 定义请求体结构
	var req struct {
		Title string `json:"title"`
	}

	// 解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return // 应该加上 return
	}

	// 调用 Service 层更新看板
	b, err := h.svc.UpdateBoard(id, req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return // 应该加上 return
	}

	// 更新成功，返回更新后的看板
	c.JSON(http.StatusOK, gin.H{"data": b})
}

// delete 删除看板
// DELETE /api/v1/boards/:id
func (h *BoardHandler) delete(c *gin.Context) {
	// 获取要删除的看板 ID
	id := c.Param("id")

	// 调用 Service 层删除看板
	if err := h.svc.DeleteBoard(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return // 应该加上 return
	}

	// http.StatusNoContent = 204（无内容）
	// 204 表示请求成功，但没有内容返回
	// 这是删除操作的标准响应
	// c.Status 只设置状态码，不返回响应体
	c.Status(http.StatusNoContent)
}
