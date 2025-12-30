package dto

type CreateOrderRequest struct {
	Items     []OrderItemRequest `json:"items" binding:"required,min=1"`
	Notes     string             `json:"notes"`
	CreatedBy *uint              `json:"-"` // Set internally, not from request
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	ID            uint                `json:"id"`
	TenantID      uint                `json:"tenant_id"`
	BranchID      uint                `json:"branch_id"`
	UserID        uint                `json:"user_id"`
	OrderNumber   string              `json:"order_number"`
	TotalAmount   float64             `json:"total_amount"`
	Status        string              `json:"status"`
	Notes         string              `json:"notes"`
	OrderItems    []OrderItemResponse `json:"order_items"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
	CreatedBy     *uint               `json:"created_by,omitempty"`
	CreatedByName *string             `json:"created_by_name,omitempty"`
	UpdatedBy     *uint               `json:"updated_by,omitempty"`
	UpdatedByName *string             `json:"updated_by_name,omitempty"`
}

type OrderItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductSKU  string  `json:"product_sku"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}
