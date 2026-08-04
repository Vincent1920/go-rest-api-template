package service_test

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm"

	"todo-api/internal/domain"
	"todo-api/internal/dto"
	"todo-api/internal/mocks"
	"todo-api/internal/service"
	"todo-api/internal/utils"
)

func TestAuthServiceRegisterSuccess(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		CreateFunc: func(user *domain.User) error {
			user.ID = 1
			return nil
		},
	}

	svc := service.NewAuthService(repo)
	result, err := svc.Register(dto.RegisterRequest{
		Name: "Martin", Email: "martin@test.com", Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register() unexpected error: %v", err)
	}
	if result.ID != 1 || result.Email != "martin@test.com" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestAuthServiceRegisterEmailAlreadyExists(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: email}, nil
		},
	}

	_, err := service.NewAuthService(repo).Register(dto.RegisterRequest{
		Name: "Martin", Email: "martin@test.com", Password: "123456",
	})
	if err == nil || err.Error() != "email already registered" {
		t.Fatalf("expected email already registered, got %v", err)
	}
}

func TestAuthServiceRegisterRepositoryError(t *testing.T) {
	expected := errors.New("database unavailable")
	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return nil, expected
		},
	}

	_, err := service.NewAuthService(repo).Register(dto.RegisterRequest{
		Name: "Martin", Email: "martin@test.com", Password: "123456",
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	hash, err := utils.HashPassword("123456")
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}
	utils.ConfigureJWT("unit-test-secret-at-least-32-characters", 1)

	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{ID: 7, Name: "Martin", Email: email, Password: hash}, nil
		},
	}

	result, err := service.NewAuthService(repo).Login(dto.LoginRequest{
		Email: "martin@test.com", Password: "123456",
	})
	if err != nil {
		t.Fatalf("Login() unexpected error: %v", err)
	}
	if result.Token == "" || result.User.ID != 7 {
		t.Fatalf("unexpected login result: %+v", result)
	}
	claims, err := utils.ValidateToken(result.Token)
	if err != nil || claims.UserID != 7 {
		t.Fatalf("generated token is invalid: claims=%+v err=%v", claims, err)
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	hash, _ := utils.HashPassword("correct-password")
	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{ID: 1, Email: email, Password: hash}, nil
		},
	}

	_, err := service.NewAuthService(repo).Login(dto.LoginRequest{
		Email: "martin@test.com", Password: "wrong-password",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid email or password") {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByEmailFunc: func(email string) (*domain.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	_, err := service.NewAuthService(repo).Login(dto.LoginRequest{
		Email: "missing@test.com", Password: "123456",
	})
	if err == nil || err.Error() != "invalid email or password" {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestAuthServiceGetProfileSuccess(t *testing.T) {
	repo := &mocks.UserRepositoryMock{
		GetByIDFunc: func(id uint) (*domain.User, error) {
			return &domain.User{ID: id, Name: "Martin", Email: "martin@test.com"}, nil
		},
	}

	result, err := service.NewAuthService(repo).GetProfile(9)
	if err != nil {
		t.Fatalf("GetProfile() unexpected error: %v", err)
	}
	if result.ID != 9 || result.Name != "Martin" {
		t.Fatalf("unexpected profile: %+v", result)
	}
}
