package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint    `gorm:"primarykey" json:"id"`
	TenantID      uint    `gorm:"not null;index" json:"tenant_id"`
	BranchID      uint    `gorm:"not null;index" json:"branch_id"`
	OrderID       uint    `gorm:"not null;index" json:"order_id"`
	Amount        float64 `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentMethod string  `gorm:"size:50;not null" json:"payment_method"`        // cash, card, transfer, qris
	Status        string  `gorm:"size:20;default:'pending';index" json:"status"` // pending, completed, failed
	Notes         string  `gorm:"type:text" json:"notes"`

	// Offline sync fields
	SyncStatus     string     `gorm:"size:20;default:'synced';index" json:"sync_status"`
	ClientID       string     `gorm:"size:100;index" json:"client_id"`
	LocalTimestamp *time.Time `json:"local_timestamp"`
	Version        int        `gorm:"default:1" json:"version"`
	ConflictData   string     `gorm:"type:jsonb" json:"conflict_data,omitempty"`

	CreatedBy *uint          `gorm:"index" json:"created_by"`
	UpdatedBy *uint          `gorm:"index" json:"updated_by"`
	DeletedBy *uint          `gorm:"index" json:"deleted_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Order   Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Creator *User `gorm:"foreignKey:CreatedBy;constraint:-" json:"creator,omitempty"`
	Updater *User `gorm:"foreignKey:UpdatedBy;constraint:-" json:"updater,omitempty"`
	Deleter *User `gorm:"foreignKey:DeletedBy;constraint:-" json:"deleter,omitempty"`
}
