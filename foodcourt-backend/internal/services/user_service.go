package services

import (
	"errors"

	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterTenant(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) RegisterTenant(req dto.RegisterRequest) error {
	// 1. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return errors.New("gagal mengamankan password")
	}

	// Mapping DTO into Model
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		Name:         req.Name,
		TenantName:   req.TenantName,
		PasswordHash: string(hashedPassword),
		Role:         "tenant",
	}

	// Save to DB
	return s.repo.CreateUser(&user)
}

func (s *userService) Login(req dto.LoginRequest) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	return user, nil
}