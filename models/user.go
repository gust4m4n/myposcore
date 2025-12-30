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
	Image     string         `gorm:"type:varchar(500)" json:"image"`
	Role      string         `gorm:"size:50;default:'user'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedBy *uint          `gorm:"index" json:"created_by"`
	UpdatedBy *uint          `gorm:"index" json:"updated_by"`
	DeletedBy *uint          `gorm:"index" json:"deleted_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant  Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Branch  Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Creator *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Updater *User  `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Deleter *User  `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
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
