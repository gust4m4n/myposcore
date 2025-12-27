package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	Name        string         `gorm:"size:255;not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"size:100;index" json:"category"`
	SKU         string         `gorm:"size:100;index" json:"sku"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	Image       string         `gorm:"type:varchar(500)" json:"image"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
