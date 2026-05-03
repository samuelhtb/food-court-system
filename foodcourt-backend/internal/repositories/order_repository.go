package repositories

import (
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateFullOrder(order *models.Order, subOrders []models.SubOrder, items []models.OrderItem) error
	GetSubOrdersByTenant(tenantID string) ([]models.SubOrder, error)
	GetAllOrders() ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

// Transaksi Database: Menyimpan ke 3 tabel sekaligus dengan aman
func (r *orderRepository) CreateFullOrder(order *models.Order, subOrders []models.SubOrder, items []models.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Order Utama
		if err := tx.Create(order).Error; err != nil {
			return err // Jika gagal, batalkan semua (Rollback)
		}

		// 2. Simpan Sub Orders (Pecahan per Tenant)
		if len(subOrders) > 0 {
			if err := tx.Create(&subOrders).Error; err != nil {
				return err
			}
		}

		// 3. Simpan Order Items (Detail Makanan)
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil // Jika semua sukses, simpan permanen (Commit)
	})
}

// Laporan Pemasukan Tenant (Hanya melihat sub-order miliknya)
func (r *orderRepository) GetSubOrdersByTenant(tenantID string) ([]models.SubOrder, error) {
	var subOrders []models.SubOrder
	err := r.db.Where("tenant_id = ?", tenantID).Find(&subOrders).Error
	return subOrders, err
}

// Laporan Pemasukan Admin (Melihat semua orderan utama)
func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Find(&orders).Error
	return orders, err
}