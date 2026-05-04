package repositories

import (
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"gorm.io/gorm"
)

// 1. Interface: Daftar kontrak/tugas yang harus bisa dilakukan oleh Menu Repository
type MenuRepository interface {
	CreateMenu(menu *models.Menu) error
	GetMenusByTenantID(tenantID uuid.UUID) ([]models.Menu, error)
	GetMenuByIDAndTenantID(id uuid.UUID, tenantID uuid.UUID) (*models.Menu, error)
	UpdateMenu(menu *models.Menu) error
	DeleteMenu(id uuid.UUID, tenantID uuid.UUID) error
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
	err := r.db.Where("tenant_id = ?", tenantID).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetMenuByIDAndTenantID(id uuid.UUID, tenantID uuid.UUID) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&menu).Error
	return &menu, err
}

func (r *menuRepository) UpdateMenu(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) DeleteMenu(id uuid.UUID, tenantID uuid.UUID) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Menu{}).Error
}