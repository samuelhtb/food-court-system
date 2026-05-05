package repositories

import (
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderTransaction(order *models.Order, menusToUpdate []models.Menu) error
	
	// Untuk Tenant
	GetSubOrdersByTenantID(tenantID uuid.UUID) ([]models.SubOrder, error)
	UpdateSubOrderStatus(subOrderID, tenantID uuid.UUID, status string) error
	
	// Untuk Admin / Kasir
	GetAllOrders() ([]models.Order, error)

	MarkOrderAsPaid(orderID uuid.UUID) error
	GetOrderByID(orderID uuid.UUID) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) CreateOrderTransaction(order *models.Order, menusToUpdate []models.Menu) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Order Utama (GORM akan otomatis menyimpan SubOrder dan OrderItem)
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Potong stok menu
		for _, menu := range menusToUpdate {
			if err := tx.Model(&menu).Update("stock", menu.Stock).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// 1. Upgrade Laporan Pemasukan Tenant (Hanya melihat sub-order miliknya + rincian menu)
func (r *orderRepository) GetSubOrdersByTenantID(tenantID uuid.UUID) ([]models.SubOrder, error) {
	var subOrders []models.SubOrder
	// Preload wajib agar Tenant tahu harus masak apa
	err := r.db.Preload("OrderItems").Where("tenant_id = ?", tenantID).Find(&subOrders).Error
	return subOrders, err
}

// 2. Implementasi Update Status Tenant
func (r *orderRepository) UpdateSubOrderStatus(subOrderID, tenantID uuid.UUID, status string) error {
	result := r.db.Model(&models.SubOrder{}).
		Where("id = ? AND tenant_id = ?", subOrderID, tenantID).
		Update("tenant_status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 3. Laporan Pemasukan Admin (Melihat semua orderan utama)
func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	// Disarankan menambah Preload SubOrders agar admin bisa lihat rincian pecahannya juga
	err := r.db.Preload("SubOrders").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) MarkOrderAsPaid(orderID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Ubah status transaksi utama jadi lunas
		if err := tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Update("payment_status", "paid").Error; err != nil {
			return err
		}

		// 2. Trigger semua dapur (SubOrder) untuk mulai memproses makanan
		if err := tx.Model(&models.SubOrder{}).
			Where("parent_order_id = ?", orderID).
			Update("tenant_status", "diproses").Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) GetOrderByID(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", orderID).First(&order).Error
	return &order, err
}