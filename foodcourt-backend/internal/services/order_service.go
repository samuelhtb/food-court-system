package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/dto"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest) (uuid.UUID, error)
	GetTenantOrders(tenantID uuid.UUID) ([]models.SubOrder, error)
	UpdateTenantOrderStatus(subOrderID, tenantID uuid.UUID, req dto.UpdateSubOrderStatusRequest) error
	GetAllOrders() ([]models.Order, error)
	MarkOrderAsPaid(orderID uuid.UUID) error
	GetTenantEarnings(tenantID uuid.UUID) (float64, int64, error)
	GetOrderWithDetails(orderID uuid.UUID) (*models.Order, error)
	GetAdminIncomeReport(startDate, endDate string) (*dto.AdminIncomeReportResponse, error)
	GetTenantIncomeReport(tenantID uuid.UUID, startDate, endDate string) (*dto.TenantIncomeReportResponse, error)
}

type orderService struct {
	orderRepo repositories.OrderRepository
	menuRepo  repositories.MenuRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, menuRepo repositories.MenuRepository) OrderService {
	return &orderService{orderRepo, menuRepo}
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest) (uuid.UUID, error) {
	var totalOrderAmount float64
	var menusToUpdate []models.Menu

	// Map untuk memecah pesanan berdasarkan TenantID
	subOrdersMap := make(map[uuid.UUID]*models.SubOrder)

	for _, itemReq := range req.Items {
		menu, err := s.menuRepo.GetMenuByID(itemReq.MenuID)
		if err != nil {
			// Jika error, return uuid.Nil sebagai penanda kosong
			return uuid.Nil, errors.New("menu tidak ditemukan")
		}

		if !menu.IsAvailable || menu.Stock < itemReq.Quantity {
			return uuid.Nil, errors.New("stok tidak cukup untuk menu: " + menu.Name)
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

	// Add order_id to JSON response
	err := s.orderRepo.CreateOrderTransaction(&newOrder, menusToUpdate)
	if err != nil {
		return uuid.Nil, err
	}

	return newOrder.ID, nil
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
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	// Validasi metode pembayaran
	if order.PaymentMethod != "cash" {
		return errors.New("ditolak: pesanan QRIS/Non-Tunai hanya bisa dikonfirmasi secara otomatis oleh sistem gateway")
	}

	// Jika cash, lanjutkan proses lunas
	return s.orderRepo.MarkOrderAsPaid(orderID)
}

func (s *orderService) GetTenantEarnings(tenantID uuid.UUID) (float64, int64, error) {
	return s.orderRepo.GetTenantEarnings(tenantID)
}

func (s *orderService) GetOrderWithDetails(orderID uuid.UUID) (*models.Order, error) {
	return s.orderRepo.GetOrderWithDetails(orderID)
}

func (s *orderService) GetAdminIncomeReport(startDate, endDate string) (*dto.AdminIncomeReportResponse, error) {
	totalRev, totalOrders, subOrders, tenants, err := s.orderRepo.GetAdminIncomeReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Map to aggregate by tenant
	tenantMap := make(map[string]*dto.TenantIncomeBreakdown)

	// Initialize all tenants with 0 revenue
	for _, t := range tenants {
		tenantMap[t.ID] = &dto.TenantIncomeBreakdown{
			TenantID:    t.ID,
			TenantName:  t.TenantName,
			Revenue:     0,
			TotalOrders: 0,
		}
	}

	for _, sub := range subOrders {
		tenantID := sub.TenantID.String()
		if _, exists := tenantMap[tenantID]; !exists {
			// Fallback if somehow a suborder has a tenant not in the user table (shouldn't happen but just in case)
			tenantMap[tenantID] = &dto.TenantIncomeBreakdown{
				TenantID:    tenantID,
				TenantName:  sub.Tenant.TenantName,
				Revenue:     0,
				TotalOrders: 0,
			}
		}

		tenantMap[tenantID].TotalOrders++
		for _, item := range sub.OrderItems {
			tenantMap[tenantID].Revenue += (float64(item.Quantity) * item.PriceAtOrder)
		}
	}

	// Convert map to slice
	var breakdown []dto.TenantIncomeBreakdown
	for _, b := range tenantMap {
		breakdown = append(breakdown, *b)
	}

	return &dto.AdminIncomeReportResponse{
		TotalRevenue: totalRev,
		TotalOrders:  totalOrders,
		Breakdown:    breakdown,
	}, nil
}

func (s *orderService) GetTenantIncomeReport(tenantID uuid.UUID, startDate, endDate string) (*dto.TenantIncomeReportResponse, error) {
	subOrders, err := s.orderRepo.GetTenantIncomeReport(tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalRev float64
	var totalOrders int64 = int64(len(subOrders))

	var history []dto.TenantSalesHistory

	for _, sub := range subOrders {
		for _, item := range sub.OrderItems {
			revenue := float64(item.Quantity) * item.PriceAtOrder
			totalRev += revenue

			history = append(history, dto.TenantSalesHistory{
				MenuName:     item.Menu.Name,
				QuantitySold: item.Quantity,
				Revenue:      revenue,
				PurchaseDate: sub.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	return &dto.TenantIncomeReportResponse{
		TotalRevenue: totalRev,
		TotalOrders:  totalOrders,
		History:      history,
	}, nil
}



