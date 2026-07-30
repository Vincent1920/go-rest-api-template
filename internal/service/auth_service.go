package service

import (
	"errors"

	"gorm.io/gorm"

	"todo-api/internal/domain"
	"todo-api/internal/dto"
	"todo-api/internal/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID uint) (*dto.UserResponse, error)
}

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Register user baru
func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {

	// Cek apakah email sudah digunakan
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Login
func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

// GetProfile
func (s *authService) GetProfile(userID uint) (*dto.UserResponse, error) {

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
