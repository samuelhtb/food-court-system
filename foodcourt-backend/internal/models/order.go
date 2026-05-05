package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerName   string     `gorm:"type:varchar(100);not null" json:"customer_name"`
	CustomerUserID *uuid.UUID `gorm:"type:uuid" json:"customer_user_id,omitempty"`
	PaymentMethod  string     `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentStatus  string     `gorm:"type:varchar(50);default:'menunggu'" json:"payment_status"`
	OrderStatus    string     `gorm:"type:varchar(50);default:'diproses'" json:"order_status"` // TAMBAHAN BARU
	TotalAmount    float64    `gorm:"type:numeric;not null" json:"total_amount"`
	CreatedAt      time.Time  `json:"created_at"`

	SubOrders []SubOrder `gorm:"foreignKey:ParentOrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_orders,omitempty"`
}