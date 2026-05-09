package dto

import "github.com/google/uuid"

type OrderItemRequest struct {
	MenuID   uuid.UUID `json:"menu_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	CustomerName   string             `json:"customer_name" binding:"required"`
	CustomerUserID *string            `json:"customer_user_id,omitempty"` // Opsional jika belum login
	PaymentMethod  string             `json:"payment_method" binding:"required,oneof=cash qris"`
	Items          []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type UpdateSubOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=menunggu_pembayaran diproses selesai"`
}