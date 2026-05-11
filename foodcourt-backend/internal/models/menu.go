package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Menu merepresentasikan tabel menus di database
type Menu struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null" json:"tenant_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:numeric;not null" json:"price"`
	Stock       int            `gorm:"type:int;not null;default:0" json:"stock"` // Kita pertahankan stock
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	Tenant      User           `gorm:"foreignKey:TenantID" json:"tenant"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}