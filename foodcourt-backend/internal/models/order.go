package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerName  string         `gorm:"type:varchar(100);not null" json:"customer_name"`
	TotalAmount   float64        `gorm:"type:numeric;not null" json:"total_amount"`
	PaymentMethod string         `gorm:"type:varchar(50)" json:"payment_method"` // cash, qris
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, cancelled
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}