package dto

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	CategoryID  *uint   `json:"category_id"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	IsActive    bool    `json:"is_active"`
	CreatedBy   *uint   `json:"-"` // Set internally, not from request
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  *uint   `json:"category_id"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" binding:"omitempty,min=0"`
	Stock       int     `json:"stock" binding:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active"`
	UpdatedBy   *uint   `json:"-"` // Set internally, not from request
}

type ProductResponse struct {
	ID             uint             `json:"id"`
	TenantID       uint             `json:"tenant_id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	CategoryID     *uint            `json:"category_id"`
	CategoryDetail *CategorySummary `json:"category_detail,omitempty"`
	SKU            string           `json:"sku"`
	Price          float64          `json:"price"`
	Stock          int              `json:"stock"`
	Image          string           `json:"image"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	CreatedBy      *uint            `json:"created_by,omitempty"`
	CreatedByName  *string          `json:"created_by_name,omitempty"`
	UpdatedBy      *uint            `json:"updated_by,omitempty"`
	UpdatedByName  *string          `json:"updated_by_name,omitempty"`
}

type CategorySummary struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
}
