package repositories

import (
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"gorm.io/gorm"
)

// 1. Interface: Daftar kontrak/tugas yang harus bisa dilakukan oleh Menu Repository
type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
	GetAllMenus() ([]models.Menu, error)
}

// 2. Struct: Penyimpan koneksi database
type menuRepository struct {
	db *gorm.DB
}

// 3. Constructor: Fungsi untuk menciptakan Menu Repository baru
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

// 4. Implementasi: Fungsi untuk menyimpan menu ke database
func (r *menuRepository) CreateMenu(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// 5. Implementasi: Fungsi untuk mengambil semua menu dari database
func (r *menuRepository) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	// GORM akan otomatis membuat query "SELECT * FROM menus"
	err := r.db.Find(&menus).Error
	return menus, err
}