package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	OrderID       uint           `gorm:"not null;index" json:"order_id"`
	Amount        float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	PaymentMethod string         `gorm:"size:50;not null" json:"payment_method"`        // cash, card, transfer, qris
	Status        string         `gorm:"size:20;default:'pending';index" json:"status"` // pending, completed, failed
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedBy     *uint          `gorm:"index" json:"created_by"`
	UpdatedBy     *uint          `gorm:"index" json:"updated_by"`
	DeletedBy     *uint          `gorm:"index" json:"deleted_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Order   Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Updater *User `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Deleter *User `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
}
