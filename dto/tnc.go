package dto

type CreateTnCRequest struct {
	Title     string `json:"title" binding:"required,min=3,max=255"`
	Content   string `json:"content" binding:"required"`
	Version   string `json:"version" binding:"required"`
	CreatedBy *uint  `json:"created_by,omitempty"`
}

type UpdateTnCRequest struct {
	Title     string `json:"title" binding:"omitempty,min=3,max=255"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	IsActive  *bool  `json:"is_active"`
	UpdatedBy *uint  `json:"updated_by,omitempty"`
}

type TnCResponse struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	Version       string  `json:"version"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}
