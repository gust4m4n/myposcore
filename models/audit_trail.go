package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditTrail struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	TenantID   *uint          `gorm:"index" json:"tenant_id,omitempty"`
	BranchID   *uint          `gorm:"index" json:"branch_id,omitempty"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	EntityType string         `gorm:"size:50;not null;index" json:"entity_type"` // user, product, order, payment, category, faq, tnc
	EntityID   uint           `gorm:"not null;index" json:"entity_id"`
	Action     string         `gorm:"size:20;not null;index" json:"action"` // create, update, delete
	Changes    string         `gorm:"type:jsonb" json:"changes,omitempty"`  // JSON string of changes
	IPAddress  string         `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"size:255" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (constraint:- prevents GORM auto-FK creation, we manage FK via SQL migrations)
	User   *User   `gorm:"foreignKey:UserID;references:ID;constraint:-" json:"-"`
	Tenant *Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:-" json:"-"`
	Branch *Branch `gorm:"foreignKey:BranchID;references:ID;constraint:-" json:"-"`
}

func (AuditTrail) TableName() string {
	return "audit_trails"
}
