package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest) error
	GetTenantOrders(tenantID uuid.UUID) ([]models.SubOrder, error)
	UpdateTenantOrderStatus(subOrderID, tenantID uuid.UUID, req dto.UpdateSubOrderStatusRequest) error
	GetAllOrders() ([]models.Order, error)
	MarkOrderAsPaid(orderID uuid.UUID) error
}

type orderService struct {
	orderRepo repositories.OrderRepository
	menuRepo  repositories.MenuRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, menuRepo repositories.MenuRepository) OrderService {
	return &orderService{orderRepo, menuRepo}
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest) error {
	var totalOrderAmount float64
	var menusToUpdate []models.Menu

	// Map untuk memecah pesanan berdasarkan TenantID
	subOrdersMap := make(map[uuid.UUID]*models.SubOrder)

	for _, itemReq := range req.Items {
		menu, err := s.menuRepo.GetMenuByID(itemReq.MenuID)
		if err != nil {
			return errors.New("menu tidak ditemukan")
		}

		if !menu.IsAvailable || menu.Stock < itemReq.Quantity {
			return errors.New("stok tidak cukup untuk menu: " + menu.Name)
		}

		// Kurangi stok di memori (akan di-commit saat transaksi DB)
		menu.Stock -= itemReq.Quantity
		menusToUpdate = append(menusToUpdate, *menu)

		itemTotalPrice := menu.Price * float64(itemReq.Quantity)
		totalOrderAmount += itemTotalPrice

		// Buat SubOrder untuk Tenant ini jika belum ada di map
		if subOrdersMap[menu.TenantID] == nil {
			subOrdersMap[menu.TenantID] = &models.SubOrder{
				TenantID:     menu.TenantID,
				TenantStatus: "menunggu_pembayaran", // Sesuai DBML
				OrderItems:   []models.OrderItem{},
			}
		}

		// Masukkan rincian item dengan mengunci harga saat ini (PriceAtOrder)
		subOrdersMap[menu.TenantID].OrderItems = append(subOrdersMap[menu.TenantID].OrderItems, models.OrderItem{
			MenuID:       menu.ID,
			Quantity:     itemReq.Quantity,
			PriceAtOrder: menu.Price,
		})
	}

	// Konversi Map ke Slice
	var subOrdersList []models.SubOrder
	for _, so := range subOrdersMap {
		subOrdersList = append(subOrdersList, *so)
	}

	// Konversi CustomerUserID string ke pointer UUID jika ada
	var customerUserIDPtr *uuid.UUID
	if req.CustomerUserID != nil && *req.CustomerUserID != "" {
		parsedID, err := uuid.Parse(*req.CustomerUserID)
		if err == nil {
			customerUserIDPtr = &parsedID
		}
	}

	newOrder := models.Order{
		CustomerName:   req.CustomerName,
		CustomerUserID: customerUserIDPtr,
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  "menunggu", // Sesuai DBML
		TotalAmount:    totalOrderAmount,
		SubOrders:      subOrdersList,
	}

	return s.orderRepo.CreateOrderTransaction(&newOrder, menusToUpdate)
}

func (s *orderService) GetTenantOrders(tenantID uuid.UUID) ([]models.SubOrder, error) {
	return s.orderRepo.GetSubOrdersByTenantID(tenantID)
}

func (s *orderService) UpdateTenantOrderStatus(subOrderID, tenantID uuid.UUID, req dto.UpdateSubOrderStatusRequest) error {
	err := s.orderRepo.UpdateSubOrderStatus(subOrderID, tenantID, req.Status)
	if err != nil {
		return errors.New("gagal memperbarui status pesanan: pesanan tidak ditemukan atau akses ditolak")
	}
	return nil
}

func (s *orderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *orderService) MarkOrderAsPaid(orderID uuid.UUID) error {
	// (Opsional: Di masa depan kamu bisa tambahkan validasi di sini untuk 
	// mengecek apakah order tersebut benar-benar menggunakan metode 'cash')
	return s.orderRepo.MarkOrderAsPaid(orderID)
}