package service_test

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"todo-api/internal/domain"
	"todo-api/internal/dto"
	"todo-api/internal/mocks"
	"todo-api/internal/service"
)

func stringPtr(value string) *string { return &value }

func TestTodoServiceCreateSuccess(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		CreateFunc: func(todo *domain.Todo) error {
			todo.ID = 1
			return nil
		},
	}

	result, err := service.NewTodoService(repo).Create(10, dto.CreateTodoRequest{
		Title: "  Belajar Golang  ", Description: "  Unit testing  ", Status: "",
	})
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}
	if result.ID != 1 || result.UserID != 10 || result.Title != "Belajar Golang" || result.Status != "Pending" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestTodoServiceCreateEmptyTitle(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{}
	_, err := service.NewTodoService(repo).Create(1, dto.CreateTodoRequest{Title: "   "})
	if !errors.Is(err, service.ErrEmptyTodoTitle) {
		t.Fatalf("expected ErrEmptyTodoTitle, got %v", err)
	}
}

func TestTodoServiceCreateInvalidStatus(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{}
	_, err := service.NewTodoService(repo).Create(1, dto.CreateTodoRequest{
		Title: "Belajar Golang", Status: "Unknown",
	})
	if !errors.Is(err, service.ErrInvalidTodoStatus) {
		t.Fatalf("expected ErrInvalidTodoStatus, got %v", err)
	}
}

func TestTodoServiceGetAllSuccess(t *testing.T) {
	now := time.Now()
	repo := &mocks.TodoRepositoryMock{
		GetAllFunc: func(userID uint, page, limit int, search string) ([]domain.Todo, int64, error) {
			if userID != 5 || page != 1 || limit != 10 || search != "go" {
				t.Fatalf("unexpected repository args: userID=%d page=%d limit=%d search=%q", userID, page, limit, search)
			}
			return []domain.Todo{{ID: 1, UserID: 5, Title: "Go", Status: "Pending", CreatedAt: now}}, 1, nil
		},
	}

	result, err := service.NewTodoService(repo).GetAll(5, 0, 0, "  go  ")
	if err != nil {
		t.Fatalf("GetAll() unexpected error: %v", err)
	}
	if result.Total != 1 || len(result.Data) != 1 || result.Page != 1 || result.Limit != 10 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestTodoServiceGetByIDSuccess(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return &domain.Todo{ID: id, UserID: 2, Title: "Todo", Status: "Pending"}, nil
		},
	}

	result, err := service.NewTodoService(repo).GetByID(2, 3)
	if err != nil {
		t.Fatalf("GetByID() unexpected error: %v", err)
	}
	if result.ID != 3 || result.UserID != 2 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestTodoServiceGetByIDNotFound(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	_, err := service.NewTodoService(repo).GetByID(1, 99)
	if !errors.Is(err, service.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoServiceRejectsOtherUsersTodo(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return &domain.Todo{ID: id, UserID: 999}, nil
		},
	}

	_, err := service.NewTodoService(repo).GetByID(1, 5)
	if !errors.Is(err, service.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

func TestTodoServiceUpdateSuccess(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return &domain.Todo{ID: id, UserID: 4, Title: "Old", Description: "Old", Status: "Pending"}, nil
		},
		UpdateFunc: func(todo *domain.Todo) error { return nil },
	}

	result, err := service.NewTodoService(repo).Update(4, 1, dto.UpdateTodoRequest{
		Title: stringPtr(" New Title "), Description: stringPtr(" New Description "), Status: stringPtr("Done"),
	})
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}
	if result.Title != "New Title" || result.Description != "New Description" || result.Status != "Done" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestTodoServiceUpdateInvalidStatus(t *testing.T) {
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return &domain.Todo{ID: id, UserID: 4, Title: "Todo", Status: "Pending"}, nil
		},
	}

	_, err := service.NewTodoService(repo).Update(4, 1, dto.UpdateTodoRequest{Status: stringPtr("Invalid")})
	if !errors.Is(err, service.ErrInvalidTodoStatus) {
		t.Fatalf("expected ErrInvalidTodoStatus, got %v", err)
	}
}

func TestTodoServiceDeleteSuccess(t *testing.T) {
	deletedID := uint(0)
	repo := &mocks.TodoRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.Todo, error) {
			return &domain.Todo{ID: id, UserID: 8}, nil
		},
		DeleteFunc: func(id uint) error {
			deletedID = id
			return nil
		},
	}

	if err := service.NewTodoService(repo).Delete(8, 12); err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}
	if deletedID != 12 {
		t.Fatalf("expected deleted ID 12, got %d", deletedID)
	}
}
