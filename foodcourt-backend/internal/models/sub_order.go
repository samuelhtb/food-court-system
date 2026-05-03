package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubOrder struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID     uuid.UUID      `gorm:"type:uuid;not null" json:"order_id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null" json:"tenant_id"`
	TotalAmount float64        `gorm:"type:numeric;not null" json:"total_amount"`
	Status      string         `gorm:"type:varchar(20);default:'preparing'" json:"status"` // preparing, ready, completed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}