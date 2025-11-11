// Package service 看板业务逻辑层
package service

import (
	"errors"
	"kanban_api/internal/model"
	"kanban_api/internal/repository"
	"strings"
)

// BoardService 看板服务接口
// 定义看板相关的业务操作
type BoardService interface {
	// ListBoards 列出所有看板
	ListBoards() ([]model.Board, error)

	// GetBoard 获取单个看板
	GetBoard(id string) (model.Board, error)

	// CreateBoard 创建新看板
	CreateBoard(title string) (model.Board, error)

	// UpdateBoard 更新看板
	UpdateBoard(id, title string) (model.Board, error)

	// DeleteBoard 删除看板
	DeleteBoard(id string) error
}

// boardService 看板服务的具体实现
type boardService struct {
	// repo 看板仓储，用于数据访问
	repo repository.BoardRepository
}

// NewBoardService 创建看板服务实例
func NewBoardService(repo repository.BoardRepository) BoardService {
	return &boardService{repo: repo}
}

// ListBoards 列出所有看板
// 这个方法比较简单，直接调用仓储层
func (s *boardService) ListBoards() ([]model.Board, error) {
	return s.repo.List()
}

// GetBoard 获取单个看板
// 同样直接调用仓储层
func (s *boardService) GetBoard(id string) (model.Board, error) {
	return s.repo.Get(id)
}

// CreateBoard 创建新看板
// Service 层负责业务验证
func (s *boardService) CreateBoard(title string) (model.Board, error) {
	// 清理标题：去除首尾空格
	title = strings.TrimSpace(title)

	// 业务规则验证：标题不能为空
	// 这是 Service 层的职责：确保数据符合业务规则
	if title == "" {
		return model.Board{}, errors.New("title required")
	}

	// 验证通过，调用仓储层创建
	return s.repo.Create(title)
}

// UpdateBoard 更新看板
func (s *boardService) UpdateBoard(id, title string) (model.Board, error) {
	// 同样进行数据清理和验证
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Board{}, errors.New("title required")
	}

	return s.repo.Update(id, title)
}

// DeleteBoard 删除看板
func (s *boardService) DeleteBoard(id string) error {
	// 直接调用仓储层删除
	// 如果需要更复杂的业务逻辑（例如：删除看板前要先删除所有任务），
	// 就在这里添加
	return s.repo.Delete(id)
}
