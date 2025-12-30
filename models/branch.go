package models

import (
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	TenantID    uint   `gorm:"not null;index:idx_branch_tenant" json:"tenant_id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Address     string `gorm:"type:text" json:"address"`
	Website     string `gorm:"type:varchar(255)" json:"website"`
	Email       string `gorm:"type:varchar(255)" json:"email"`
	Phone       string `gorm:"type:varchar(50)" json:"phone"`
	Image       string `gorm:"type:varchar(500)" json:"image"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Audit tracking
	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint `gorm:"index" json:"deleted_by,omitempty"`

	// Relations
	Tenant  *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users   []User  `gorm:"foreignKey:BranchID" json:"users,omitempty"`
	Creator *User   `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
	Updater *User   `gorm:"foreignKey:UpdatedBy;references:ID" json:"-"`
	Deleter *User   `gorm:"foreignKey:DeletedBy;references:ID" json:"-"`
}
