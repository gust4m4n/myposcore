package models

import (
	"time"

	"gorm.io/gorm"
)

type TermsAndConditions struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Version   string         `gorm:"size:20;not null" json:"version"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Audit tracking
	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint `gorm:"index" json:"deleted_by,omitempty"`

	// Relations for audit tracking
	Creator *User `gorm:"foreignKey:CreatedBy;references:ID;constraint:-" json:"-"`
	Updater *User `gorm:"foreignKey:UpdatedBy;references:ID;constraint:-" json:"-"`
	Deleter *User `gorm:"foreignKey:DeletedBy;references:ID;constraint:-" json:"-"`
}
