package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SubOrderID uuid.UUID      `gorm:"type:uuid;not null" json:"sub_order_id"`
	MenuID     uuid.UUID      `gorm:"type:uuid;not null" json:"menu_id"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	Price      float64        `gorm:"type:numeric;not null" json:"price"` // Harga saat dibeli (jaga-jaga jika harga menu berubah nanti)
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}