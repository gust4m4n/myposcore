package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Address     string         `gorm:"type:text" json:"address"`
	City        string         `gorm:"size:100" json:"city"`
	Country     string         `gorm:"size:100" json:"country"`
	PostalCode  string         `gorm:"size:20" json:"postal_code"`
	Website     string         `gorm:"size:255" json:"website"`
	Email       string         `gorm:"size:255" json:"email"`
	Phone       string         `gorm:"size:50" json:"phone"`
	Image       string         `gorm:"size:500" json:"image"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Audit tracking
	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint `gorm:"index" json:"deleted_by,omitempty"`

	// Offline Sync Fields
	SyncStatus     string     `gorm:"size:20;default:'synced';index" json:"sync_status"`
	ClientID       string     `gorm:"size:255;index" json:"client_id,omitempty"`
	LocalTimestamp *time.Time `json:"local_timestamp,omitempty"`
	Version        int        `gorm:"default:1" json:"version"`
	ConflictData   *string     `gorm:"type:jsonb" json:"conflict_data,omitempty"`

	// Relations
	Users   []User `gorm:"foreignKey:TenantID" json:"-"`
	Creator *User  `gorm:"foreignKey:CreatedBy;references:ID;constraint:-" json:"-"`
	Updater *User  `gorm:"foreignKey:UpdatedBy;references:ID;constraint:-" json:"-"`
	Deleter *User  `gorm:"foreignKey:DeletedBy;references:ID;constraint:-" json:"-"`
}
