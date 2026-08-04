package mocks

import (
	"todo-api/internal/domain"
)

// UserRepositoryMock adalah mock berbasis function field untuk unit test.
// Set hanya function yang dibutuhkan oleh setiap test case.
type UserRepositoryMock struct {
	CreateFunc     func(user *domain.User) error
	GetByIDFunc    func(id uint) (*domain.User, error)
	GetByEmailFunc func(email string) (*domain.User, error)
	UpdateFunc     func(user *domain.User) error
	DeleteFunc     func(id uint) error
}

var _ domain.UserRepository = (*UserRepositoryMock)(nil)

func (m *UserRepositoryMock) Create(user *domain.User) error {
	if m.CreateFunc == nil {
		panic("UserRepositoryMock.CreateFunc is not configured")
	}
	return m.CreateFunc(user)
}

func (m *UserRepositoryMock) GetByID(id uint) (*domain.User, error) {
	if m.GetByIDFunc == nil {
		panic("UserRepositoryMock.GetByIDFunc is not configured")
	}
	return m.GetByIDFunc(id)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*domain.User, error) {
	if m.GetByEmailFunc == nil {
		panic("UserRepositoryMock.GetByEmailFunc is not configured")
	}
	return m.GetByEmailFunc(email)
}

func (m *UserRepositoryMock) Update(user *domain.User) error {
	if m.UpdateFunc == nil {
		panic("UserRepositoryMock.UpdateFunc is not configured")
	}
	return m.UpdateFunc(user)
}

func (m *UserRepositoryMock) Delete(id uint) error {
	if m.DeleteFunc == nil {
		panic("UserRepositoryMock.DeleteFunc is not configured")
	}
	return m.DeleteFunc(id)
}
