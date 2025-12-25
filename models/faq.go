package models

import (
	"time"

	"gorm.io/gorm"
)

type FAQ struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Question  string         `gorm:"size:500;not null" json:"question"`
	Answer    string         `gorm:"type:text;not null" json:"answer"`
	Category  string         `gorm:"size:100" json:"category"`
	Order     int            `gorm:"default:0" json:"order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
