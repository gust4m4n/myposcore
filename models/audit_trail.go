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
	Action     string         `gorm:"size:100;not null;index" json:"action"`      // create, update, delete
	EntityType string         `gorm:"size:100;not null;index" json:"entity_type"` // users, products, orders, etc
	EntityID   *uint          `gorm:"index" json:"entity_id,omitempty"`
	Changes    string         `gorm:"type:jsonb" json:"changes,omitempty"` // JSON string of changes
	IPAddress  string         `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"size:255" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Offline Sync Fields
	SyncStatus     string     `gorm:"size:20;default:'synced';index" json:"sync_status"`
	ClientID       string     `gorm:"size:255;index" json:"client_id,omitempty"`
	LocalTimestamp *time.Time `json:"local_timestamp,omitempty"`
	Version        int        `gorm:"default:1" json:"version"`
	ConflictData   string     `gorm:"type:jsonb" json:"conflict_data,omitempty"`

	// Relations (constraint:- prevents GORM auto-FK creation, we manage FK via SQL migrations)
	User   *User   `gorm:"foreignKey:UserID;references:ID;constraint:-" json:"-"`
	Tenant *Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:-" json:"-"`
	Branch *Branch `gorm:"foreignKey:BranchID;references:ID;constraint:-" json:"-"`
}

func (AuditTrail) TableName() string {
	return "audit_trails"
}
