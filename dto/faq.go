package dto

type CreateFAQRequest struct {
	Question string `json:"question" binding:"required,min=5,max=500"`
	Answer   string `json:"answer" binding:"required"`
	Category string `json:"category" binding:"omitempty,max=100"`
	Order    int    `json:"order"`
}

type UpdateFAQRequest struct {
	Question string `json:"question" binding:"omitempty,min=5,max=500"`
	Answer   string `json:"answer"`
	Category string `json:"category" binding:"omitempty,max=100"`
	Order    *int   `json:"order"`
	IsActive *bool  `json:"is_active"`
}

type FAQResponse struct {
	ID        uint   `json:"id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Category  string `json:"category"`
	Order     int    `json:"order"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
