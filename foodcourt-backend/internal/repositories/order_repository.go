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

	GetTenantEarnings(tenantID uuid.UUID) (float64, int64, error)

	GetOrderWithDetails(orderID uuid.UUID) (*models.Order, error)
	FindByMidtransOrderID(midtransOrderID string) (*models.Order, error)
	UpdateMidtransInfo(orderID uuid.UUID, token, midtransOrderID string) error
	UpdatePaymentStatus(orderID uuid.UUID, status string) error

	// Reports
	GetAdminIncomeReport(startDate, endDate string) (float64, int64, []models.SubOrder, []models.User, error)
	GetTenantIncomeReport(tenantID uuid.UUID, startDate, endDate string) ([]models.SubOrder, error)
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
	// Hanya tampilkan jika status pembayaran di tabel utama sudah "paid"
	// Preload ParentOrder (untuk nama pembeli) dan OrderItems.Menu (untuk nama menu)
	err := r.db.Preload("ParentOrder").Preload("OrderItems.Menu").
		Joins("JOIN orders ON orders.id = sub_orders.parent_order_id").
		Where("sub_orders.tenant_id = ? AND orders.payment_status = ?", tenantID, "paid").
		Find(&subOrders).Error
	return subOrders, err
}

// 2. Implementasi Update Status Tenant
func (r *orderRepository) UpdateSubOrderStatus(subOrderID, tenantID uuid.UUID, status string) error {
	// Gunakan Transaction agar pengecekan dan update aman dari bentrok data
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update status SubOrder milik tenant tersebut
		result := tx.Model(&models.SubOrder{}).
			Where("id = ? AND tenant_id = ?", subOrderID, tenantID).
			Update("tenant_status", status)

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// 2. Jika status diubah menjadi "selesai", kita lakukan pengecekan cerdas
		if status == "selesai" {
			// Cari tahu siapa Parent Order (Order Utama) dari sub-order ini
			var subOrder models.SubOrder
			if err := tx.Select("parent_order_id").First(&subOrder, "id = ?", subOrderID).Error; err != nil {
				return err
			}

			// Hitung apakah masih ada SubOrder lain dari Parent ini yang statusnya BUKAN "selesai"
			var incompleteCount int64
			if err := tx.Model(&models.SubOrder{}).
				Where("parent_order_id = ? AND tenant_status != ?", subOrder.ParentOrderID, "selesai").
				Count(&incompleteCount).Error; err != nil {
				return err
			}

			// 3. Jika hasilnya 0 (artinya SEMUA dapur sudah selesai memasak pesanan ini)
			if incompleteCount == 0 {
				// Update OrderStatus di pesanan utama menjadi "selesai"
				if err := tx.Model(&models.Order{}).
					Where("id = ?", subOrder.ParentOrderID).
					Update("order_status", "selesai").Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// 3. Laporan Pemasukan Admin (Melihat semua orderan utama)
func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	// Disarankan menambah Preload SubOrders agar admin bisa lihat rincian pecahannya juga
	err := r.db.Preload("SubOrders.OrderItems.Menu").Order("created_at desc").Find(&orders).Error
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

func (r *orderRepository) GetTenantEarnings(tenantID uuid.UUID) (float64, int64, error) {
	var subOrders []models.SubOrder
	
	// Tarik data SubOrder yang berstatus "selesai" beserta rincian itemnya
	err := r.db.Preload("OrderItems").
		Where("tenant_id = ? AND tenant_status = ?", tenantID, "selesai").
		Find(&subOrders).Error

	if err != nil {
		return 0, 0, err
	}

	var totalEarnings float64 = 0
	var totalOrders int64 = int64(len(subOrders))

	// Kalkulasi total pendapatan
	for _, subOrder := range subOrders {
		for _, item := range subOrder.OrderItems {
			totalEarnings += (float64(item.Quantity) * item.PriceAtOrder)
		}
	}

	return totalEarnings, totalOrders, nil
}

func (r *orderRepository) GetOrderWithDetails(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	
	// Preload bertingkat: Ambil SubOrders sekaligus OrderItems di dalamnya
	err := r.db.Preload("SubOrders.OrderItems.Menu").Where("id = ?", orderID).First(&order).Error
	
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByMidtransOrderID(midtransOrderID string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("midtrans_order_id = ?", midtransOrderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateMidtransInfo(orderID uuid.UUID, token, midtransOrderID string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"midtrans_token":    token,
		"midtrans_order_id": midtransOrderID,
	}).Error
}

func (r *orderRepository) UpdatePaymentStatus(orderID uuid.UUID, status string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Order{}).Where("id = ?", orderID).Update("payment_status", status).Error; err != nil {
			return err
		}

		if status == "paid" {
			// Trigger semua dapur (SubOrder) untuk mulai memproses makanan
			if err := tx.Model(&models.SubOrder{}).
				Where("parent_order_id = ?", orderID).
				Update("tenant_status", "diproses").Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Reports

func (r *orderRepository) GetAdminIncomeReport(startDate, endDate string) (float64, int64, []models.SubOrder, []models.User, error) {
	// First, get total revenue and total orders
	var orders []models.Order
	query := r.db.Where("payment_status = ?", "paid")

	if startDate != "" && endDate != "" {
		// Filter by created_at range. We need to parse or assume format YYYY-MM-DD
		query = query.Where("created_at >= ? AND created_at <= ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := query.Find(&orders).Error
	if err != nil {
		return 0, 0, nil, nil, err
	}

	var totalRevenue float64
	var totalOrders int64 = int64(len(orders))
	for _, o := range orders {
		totalRevenue += o.TotalAmount
	}

	// Next, get breakdown per tenant via SubOrders
	var subOrders []models.SubOrder
	subQuery := r.db.Preload("Tenant").Preload("OrderItems").Joins("JOIN orders ON orders.id = sub_orders.parent_order_id").Where("orders.payment_status = ?", "paid")

	if startDate != "" && endDate != "" {
		subQuery = subQuery.Where("orders.created_at >= ? AND orders.created_at <= ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err = subQuery.Find(&subOrders).Error
	if err != nil {
		return 0, 0, nil, nil, err
	}

	// Fetch all tenants to ensure those with 0 revenue are also included
	var tenants []models.User
	err = r.db.Where("role = ?", "tenant").Find(&tenants).Error
	if err != nil {
		return 0, 0, nil, nil, err
	}

	return totalRevenue, totalOrders, subOrders, tenants, nil
}

func (r *orderRepository) GetTenantIncomeReport(tenantID uuid.UUID, startDate, endDate string) ([]models.SubOrder, error) {
	var subOrders []models.SubOrder
	query := r.db.Preload("OrderItems.Menu").Joins("JOIN orders ON orders.id = sub_orders.parent_order_id").Where("sub_orders.tenant_id = ? AND orders.payment_status = ?", tenantID, "paid")

	if startDate != "" && endDate != "" {
		query = query.Where("orders.created_at >= ? AND orders.created_at <= ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := query.Find(&subOrders).Error
	if err != nil {
		return nil, err
	}

	return subOrders, nil
}