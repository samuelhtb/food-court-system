package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
    Username     string         `gorm:"type:varchar(255);unique;not null" json:"username"`
    PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
    Role         string         `gorm:"type:role_enum;not null" json:"role"`
    TenantName   string         `gorm:"type:varchar(255)" json:"tenant_name,omitempty"`
    Name         string         `gorm:"type:varchar(100);not null" json:"name"`
    Email        string         `gorm:"type:varchar(100);unique;not null" json:"email"`
    CreatedAt    time.Time      `gorm:"default:now()" json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}