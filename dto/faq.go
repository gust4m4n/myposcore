package dto

type CreateFAQRequest struct {
	Question  string `json:"question" binding:"required,min=5,max=500"`
	Answer    string `json:"answer" binding:"required"`
	CreatedBy *uint  `json:"created_by,omitempty"`
}

type UpdateFAQRequest struct {
	Question  string `json:"question" binding:"omitempty,min=5,max=500"`
	Answer    string `json:"answer"`
	IsActive  *bool  `json:"is_active"`
	UpdatedBy *uint  `json:"updated_by,omitempty"`
}

type FAQResponse struct {
	ID            uint    `json:"id"`
	Question      string  `json:"question"`
	Answer        string  `json:"answer"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}
