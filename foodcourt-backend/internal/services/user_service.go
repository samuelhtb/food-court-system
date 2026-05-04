package services

import (
	"errors"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterTenant(user *models.User) error
	Login(email, password string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// Register Tenant oleh Admin
func (s *userService) RegisterTenant(user *models.User) error {
	// 1. Password encrypt: cost 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	
	user.Password = string(hashedPassword)
	user.PasswordHash = string(hashedPassword)
	user.Role = "tenant"
	
	return s.repo.CreateUser(user)
}

// Login Tenant/Admin
func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password salah")
	}

	return user, nil
}