package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerName   string     `gorm:"type:varchar(100);not null" json:"customer_name"`
	CustomerUserID *uuid.UUID `gorm:"type:uuid" json:"customer_user_id,omitempty"` // Pointer agar bisa null
	PaymentMethod  string     `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentStatus  string     `gorm:"type:varchar(50);default:'pending'" json:"payment_status"`
	TotalAmount    float64    `gorm:"type:numeric;not null" json:"total_amount"`
	CreatedAt      time.Time  `json:"created_at"`

	// Relasi ke SubOrders
	SubOrders []SubOrder `gorm:"foreignKey:ParentOrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_orders,omitempty"`
}