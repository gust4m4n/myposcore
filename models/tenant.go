package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex:idx_tenant_code;not null" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Address     string         `gorm:"type:text" json:"address"`
	Website     string         `gorm:"size:255" json:"website"`
	Email       string         `gorm:"size:255" json:"email"`
	Phone       string         `gorm:"size:50" json:"phone"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Users []User `gorm:"foreignKey:TenantID" json:"-"`
}
