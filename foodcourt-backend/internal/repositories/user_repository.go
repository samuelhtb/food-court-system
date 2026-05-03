package repositories

import (
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Digunakan oleh Admin untuk mendaftarkan Tenant baru
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// Digunakan saat Login untuk mencari data berdasarkan email
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Query: SELECT * FROM users WHERE email = ? LIMIT 1
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err // Mengembalikan error jika email tidak ditemukan
	}
	return &user, nil
}