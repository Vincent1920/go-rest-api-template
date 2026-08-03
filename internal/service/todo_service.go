package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"todo-api/internal/domain"
	"todo-api/internal/dto"
)

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrInvalidTodoStatus = errors.New("invalid todo status")
	ErrEmptyTodoTitle    = errors.New("title cannot be empty")
)

type TodoService interface {
	Create(
		userID uint,
		req dto.CreateTodoRequest,
	) (*dto.TodoResponse, error)

	GetAll(
		userID uint,
		page int,
		limit int,
		search string,
	) (*dto.TodoListResponse, error)

	GetByID(
		userID uint,
		todoID uint,
	) (*dto.TodoResponse, error)

	Update(
		userID uint,
		todoID uint,
		req dto.UpdateTodoRequest,
	) (*dto.TodoResponse, error)

	Delete(
		userID uint,
		todoID uint,
	) error
}

type todoService struct {
	todoRepo domain.TodoRepository
}

func NewTodoService(todoRepo domain.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (s *todoService) Create(
	userID uint,
	req dto.CreateTodoRequest,
) (*dto.TodoResponse, error) {
	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	status := strings.TrimSpace(req.Status)

	if title == "" {
		return nil, ErrEmptyTodoTitle
	}

	if status == "" {
		status = "Pending"
	}

	if !isValidTodoStatus(status) {
		return nil, ErrInvalidTodoStatus
	}

	todo := &domain.Todo{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      status,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	response := mapTodoToResponse(todo)

	return &response, nil
}

func (s *todoService) GetAll(
	userID uint,
	page int,
	limit int,
	search string,
) (*dto.TodoListResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)

	todos, total, err := s.todoRepo.GetAll(
		userID,
		page,
		limit,
		search,
	)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TodoResponse, 0, len(todos))

	for i := range todos {
		responses = append(
			responses,
			mapTodoToResponse(&todos[i]),
		)
	}

	return &dto.TodoListResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *todoService) GetByID(
	userID uint,
	todoID uint,
) (*dto.TodoResponse, error) {
	todo, err := s.findOwnedTodo(userID, todoID)
	if err != nil {
		return nil, err
	}

	response := mapTodoToResponse(todo)

	return &response, nil
}

func (s *todoService) Update(
	userID uint,
	todoID uint,
	req dto.UpdateTodoRequest,
) (*dto.TodoResponse, error) {
	todo, err := s.findOwnedTodo(userID, todoID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)

		if title == "" {
			return nil, ErrEmptyTodoTitle
		}

		todo.Title = title
	}

	if req.Description != nil {
		todo.Description = strings.TrimSpace(*req.Description)
	}

	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)

		if !isValidTodoStatus(status) {
			return nil, ErrInvalidTodoStatus
		}

		todo.Status = status
	}

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	response := mapTodoToResponse(todo)

	return &response, nil
}

func (s *todoService) Delete(
	userID uint,
	todoID uint,
) error {
	todo, err := s.findOwnedTodo(userID, todoID)
	if err != nil {
		return err
	}

	if err := s.todoRepo.Delete(todo.ID); err != nil {
		return err
	}

	return nil
}

func (s *todoService) findOwnedTodo(
	userID uint,
	todoID uint,
) (*domain.Todo, error) {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoNotFound
		}

		return nil, err
	}

	if todo.UserID != userID {
		// Jangan bocorkan keberadaan Todo milik user lain.
		return nil, ErrTodoNotFound
	}

	return todo, nil
}

func isValidTodoStatus(status string) bool {
	switch status {
	case "Pending", "In Progress", "Done":
		return true
	default:
		return false
	}
}

func mapTodoToResponse(todo *domain.Todo) dto.TodoResponse {
	return dto.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		UserID:      todo.UserID,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
