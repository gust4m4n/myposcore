package models

import (
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	TenantID    uint   `gorm:"not null;index:idx_branch_tenant" json:"tenant_id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Code        string `gorm:"type:varchar(50);uniqueIndex:idx_branch_code;not null" json:"code"`
	Description string `gorm:"type:text" json:"description"`
	Address     string `gorm:"type:text" json:"address"`
	Website     string `gorm:"type:varchar(255)" json:"website"`
	Email       string `gorm:"type:varchar(255)" json:"email"`
	Phone       string `gorm:"type:varchar(50)" json:"phone"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Relations
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users  []User  `gorm:"foreignKey:BranchID" json:"users,omitempty"`
}
