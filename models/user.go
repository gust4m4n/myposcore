package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID  *uint          `gorm:"index" json:"branch_id,omitempty"`
	Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	PIN       string         `gorm:"size:255" json:"-"`
	FullName  string         `gorm:"size:255;index" json:"full_name"`
	Phone     string         `gorm:"size:50" json:"phone"`
	Image     string         `gorm:"type:varchar(500)" json:"image"`
	Role      string         `gorm:"size:50;default:'staff'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedBy *uint          `gorm:"index" json:"created_by"`
	UpdatedBy *uint          `gorm:"index" json:"updated_by"`
	DeletedBy *uint          `gorm:"index" json:"deleted_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Offline Sync Fields
	SyncStatus     string     `gorm:"size:20;default:'synced';index" json:"sync_status"`
	ClientID       string     `gorm:"size:255;index" json:"client_id,omitempty"`
	LocalTimestamp *time.Time `json:"local_timestamp,omitempty"`
	Version        int        `gorm:"default:1" json:"version"`
	ConflictData   string     `gorm:"type:jsonb" json:"conflict_data,omitempty"`

	// Relations
	Tenant  Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Branch  Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Creator *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Updater *User  `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Deleter *User  `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// Add composite unique index
type UserIndex struct{}

func (UserIndex) TableName() string {
	return "users"
}
