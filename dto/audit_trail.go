package dto

type AuditTrailResponse struct {
	ID         uint    `json:"id"`
	TenantID   *uint   `json:"tenant_id,omitempty"`
	BranchID   *uint   `json:"branch_id,omitempty"`
	UserID     uint    `json:"user_id"`
	UserName   string  `json:"user_name"`
	EntityType string  `json:"entity_type"`
	EntityID   uint    `json:"entity_id"`
	Action     string  `json:"action"`
	Changes    *string `json:"changes,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	UserAgent  *string `json:"user_agent,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

type AuditTrailListRequest struct {
	Page       int     `json:"page" form:"page"`
	Limit      int     `json:"limit" form:"limit"`
	UserID     *uint   `json:"user_id" form:"user_id"`
	EntityType *string `json:"entity_type" form:"entity_type"`
	EntityID   *uint   `json:"entity_id" form:"entity_id"`
	Action     *string `json:"action" form:"action"`
	DateFrom   *string `json:"date_from" form:"date_from"` // Format: 2006-01-02
	DateTo     *string `json:"date_to" form:"date_to"`     // Format: 2006-01-02
}
