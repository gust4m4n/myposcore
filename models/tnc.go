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
}
