package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
    // Menggunakan string karena id di database kamu adalah UUID
	ID           string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username     string         `json:"username"`
	PasswordHash string         `json:"password_hash"`
	Role         string         `json:"role"`
	TenantName   string         `json:"tenant_name"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}