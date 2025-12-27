package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID  uint           `gorm:"not null;index" json:"branch_id"`
	Username  string         `gorm:"size:100;not null" json:"username"`
	Email     string         `gorm:"size:255;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	PIN       string         `gorm:"size:255" json:"-"`
	FullName  string         `gorm:"size:255" json:"full_name"`
	Role      string         `gorm:"size:50;default:'user'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Branch Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

// Unique constraint on username per tenant
func (User) TableName() string {
	return "users"
}

// Add composite unique index
type UserIndex struct{}

func (UserIndex) TableName() string {
	return "users"
}
