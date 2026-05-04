package services

import (
	"github.com/google/uuid"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/models"
	"github.com/samuelhtb/food-court-system/foodcourt-backend/internal/repositories"
)

type OrderService interface {
	CreateOrder(customerName string, items []models.OrderItem) (*models.Order, error)
}

type orderService struct {
	repo     repositories.OrderRepository
	menuRepo repositories.MenuRepository
}

func NewOrderService(repo repositories.OrderRepository, menuRepo repositories.MenuRepository) OrderService {
	return &orderService{repo, menuRepo}
}

func (s *orderService) CreateOrder(customerName string, items []models.OrderItem) (*models.Order, error) {
	orderID := uuid.New()
	var totalAmount float64

	// Logika Hitung total harga dari semua item yang dipesan
	for i := range items {
		totalAmount += items[i].Price * float64(items[i].Quantity)
		items[i].ID = uuid.New()
	}

	order := &models.Order{
		ID:           orderID,
		CustomerName: customerName,
		TotalAmount:  totalAmount,
		Status:       "pending",
	}

	err := s.repo.CreateFullOrder(order, nil, items) 
	return order, err
}