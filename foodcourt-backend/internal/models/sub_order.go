package models

import (
	"time"

	"github.com/google/uuid"
)

type SubOrder struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ParentOrderID uuid.UUID  `gorm:"type:uuid;not null" json:"parent_order_id"`
	TenantID      uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	TenantStatus  string     `gorm:"type:varchar(50);default:'menunggu_pembayaran'" json:"tenant_status"`
	CreatedAt     time.Time  `json:"created_at"`

	OrderItems []OrderItem `gorm:"foreignKey:SubOrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items,omitempty"`
}