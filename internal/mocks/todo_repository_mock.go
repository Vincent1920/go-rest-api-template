package mocks

import "todo-api/internal/domain"

// TodoRepositoryMock adalah mock berbasis function field untuk unit test.
type TodoRepositoryMock struct {
	CreateFunc  func(todo *domain.Todo) error
	GetByIDFunc func(id uint) (*domain.Todo, error)
	GetAllFunc  func(userID uint, page, limit int, search string) ([]domain.Todo, int64, error)
	UpdateFunc  func(todo *domain.Todo) error
	DeleteFunc  func(id uint) error
}

var _ domain.TodoRepository = (*TodoRepositoryMock)(nil)

func (m *TodoRepositoryMock) Create(todo *domain.Todo) error {
	if m.CreateFunc == nil {
		panic("TodoRepositoryMock.CreateFunc is not configured")
	}
	return m.CreateFunc(todo)
}

func (m *TodoRepositoryMock) GetByID(id uint) (*domain.Todo, error) {
	if m.GetByIDFunc == nil {
		panic("TodoRepositoryMock.GetByIDFunc is not configured")
	}
	return m.GetByIDFunc(id)
}

func (m *TodoRepositoryMock) GetAll(userID uint, page, limit int, search string) ([]domain.Todo, int64, error) {
	if m.GetAllFunc == nil {
		panic("TodoRepositoryMock.GetAllFunc is not configured")
	}
	return m.GetAllFunc(userID, page, limit, search)
}

func (m *TodoRepositoryMock) Update(todo *domain.Todo) error {
	if m.UpdateFunc == nil {
		panic("TodoRepositoryMock.UpdateFunc is not configured")
	}
	return m.UpdateFunc(todo)
}

func (m *TodoRepositoryMock) Delete(id uint) error {
	if m.DeleteFunc == nil {
		panic("TodoRepositoryMock.DeleteFunc is not configured")
	}
	return m.DeleteFunc(id)
}
