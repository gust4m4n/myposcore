package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	CreatedBy   *uint  `json:"created_by,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	UpdatedBy   *uint  `json:"updated_by,omitempty"`
}

type CategoryResponse struct {
	ID            uint    `json:"id"`
	TenantID      uint    `json:"tenant_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}
