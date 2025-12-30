package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID    uint           `gorm:"not null;index" json:"branch_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	OrderNumber string         `gorm:"size:50;uniqueIndex;not null" json:"order_number"`
	TotalAmount float64        `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Status      string         `gorm:"size:20;default:'pending';index" json:"status"` // pending, confirmed, completed, cancelled
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedBy   *uint          `gorm:"index" json:"created_by"`
	UpdatedBy   *uint          `gorm:"index" json:"updated_by"`
	DeletedBy   *uint          `gorm:"index" json:"deleted_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant     Tenant      `gorm:"foreignKey:TenantID" json:"-"`
	Branch     Branch      `gorm:"foreignKey:BranchID" json:"-"`
	User       User        `gorm:"foreignKey:UserID" json:"-"`
	Creator    *User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Updater    *User       `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Deleter    *User       `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Payments   []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Subtotal  float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Order   Order   `gorm:"foreignKey:OrderID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
