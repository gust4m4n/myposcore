package dto

type CreatePaymentRequest struct {
	OrderID       uint    `json:"order_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash card transfer qris"`
	Notes         string  `json:"notes"`
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
}
