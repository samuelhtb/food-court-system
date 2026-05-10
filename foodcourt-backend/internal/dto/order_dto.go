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

type CreateOrderResponse struct {
	Message       string    `json:"message"`
	OrderID       uuid.UUID `json:"order_id"`
	MidtransToken string    `json:"midtrans_token,omitempty"`
}


type UpdateSubOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=menunggu_pembayaran diproses selesai"`
}

type TrackOrderItemResponse struct {
	MenuName string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type TrackOrderResponse struct {
	ID            uuid.UUID                `json:"id"`
	CustomerName  string                   `json:"customer_name"`
	PaymentMethod string                   `json:"payment_method"`
	PaymentStatus string                   `json:"payment_status"`
	OrderStatus   string                   `json:"order_status"`
	TotalAmount   float64                  `json:"total_amount"`
	Items         []TrackOrderItemResponse `json:"items"`
}

// Laporan Pemasukan (Income Reports)

type TenantIncomeBreakdown struct {
	TenantID    string  `json:"tenant_id"`
	TenantName  string  `json:"tenant_name"`
	Revenue     float64 `json:"revenue"`
	TotalOrders int64   `json:"total_orders"`
}

type AdminIncomeReportResponse struct {
	TotalRevenue float64                 `json:"total_revenue"`
	TotalOrders  int64                   `json:"total_orders"`
	Breakdown    []TenantIncomeBreakdown `json:"breakdown"`
}

type MenuIncomeBreakdown struct {
	MenuID       string  `json:"menu_id"`
	MenuName     string  `json:"menu_name"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}

type TenantIncomeReportResponse struct {
	TotalRevenue float64               `json:"total_revenue"`
	TotalOrders  int64                 `json:"total_orders"`
	Items        []MenuIncomeBreakdown `json:"items"`
}