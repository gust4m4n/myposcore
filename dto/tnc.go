package dto

type CreateTnCRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=255"`
	Content string `json:"content" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type UpdateTnCRequest struct {
	Title    string `json:"title" binding:"omitempty,min=3,max=255"`
	Content  string `json:"content"`
	Version  string `json:"version"`
	IsActive *bool  `json:"is_active"`
}

type TnCResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Version   string `json:"version"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
