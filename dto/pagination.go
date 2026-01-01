package dto

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func NewPaginationRequest(page, pageSize int) *PaginationRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 32 // default
	}
	return &PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

func NewPaginationResponse(page, pageSize int, totalItems int64, data interface{}) *PaginationResponse {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Data:       data,
	}
}
