package dto

type CreatePaymentRequest struct {
	OrderID       uint    `json:"order_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash card transfer qris"`
	Notes         string  `json:"notes"`
	CreatedBy     *uint   `json:"-"` // Set internally, not from request
}

type PaymentResponse struct {
	ID            uint    `json:"id"`
	OrderID       uint    `json:"order_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
	Notes         string  `json:"notes"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}

type PaymentDetailResponse struct {
	ID            uint               `json:"id"`
	OrderID       uint               `json:"order_id"`
	OrderNumber   string             `json:"order_number"`
	Amount        float64            `json:"amount"`
	PaymentMethod string             `json:"payment_method"`
	Status        string             `json:"status"`
	Notes         string             `json:"notes"`
	Change        float64            `json:"change"`
	CreatedAt     string             `json:"created_at"`
	Order         PaymentOrderDetail `json:"order"`
}

type PaymentOrderDetail struct {
	TotalAmount float64                  `json:"total_amount"`
	Status      string                   `json:"status"`
	Notes       string                   `json:"notes"`
	OrderItems  []PaymentOrderItemDetail `json:"order_items"`
	CashierName string                   `json:"cashier_name"`
	BranchName  string                   `json:"branch_name"`
}

type PaymentOrderItemDetail struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}
