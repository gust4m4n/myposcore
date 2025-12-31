package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Image       string         `gorm:"type:varchar(500)" json:"image"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Audit tracking
	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint `gorm:"index" json:"deleted_by,omitempty"`

	// Relations
	Tenant  Tenant `gorm:"foreignKey:TenantID" json:"-"`
	Creator *User  `gorm:"foreignKey:CreatedBy;references:ID;constraint:-" json:"-"`
	Updater *User  `gorm:"foreignKey:UpdatedBy;references:ID;constraint:-" json:"-"`
	Deleter *User  `gorm:"foreignKey:DeletedBy;references:ID;constraint:-" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}
