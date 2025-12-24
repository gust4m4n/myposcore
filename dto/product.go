package dto

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	IsActive    bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" binding:"omitempty,min=0"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	TenantID    uint    `json:"tenant_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
