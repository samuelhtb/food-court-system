package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SubOrderID   uuid.UUID `gorm:"type:uuid;not null" json:"sub_order_id"`
	MenuID       uuid.UUID `gorm:"type:uuid;not null" json:"menu_id"`
	Quantity     int       `gorm:"type:int;not null" json:"quantity"`
	PriceAtOrder float64   `gorm:"type:numeric;not null" json:"price_at_order"`
    CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}