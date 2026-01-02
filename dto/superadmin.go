package dto

type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type CreateBranchRequest struct {
	TenantID    uint   `json:"tenant_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type UpdateBranchRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type TenantResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Address       string  `json:"address"`
	Website       string  `json:"website"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Image         string  `json:"image"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}

type BranchResponse struct {
	ID            uint    `json:"id"`
	TenantID      uint    `json:"tenant_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Address       string  `json:"address"`
	Website       string  `json:"website"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Image         string  `json:"image"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}

type UserResponse struct {
	ID            uint    `json:"id"`
	TenantID      uint    `json:"tenant_id"`
	BranchID      uint    `json:"branch_id"`
	Email         string  `json:"email"`
	FullName      string  `json:"full_name"`
	Image         string  `json:"image"`
	Role          string  `json:"role"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}

type DashboardResponse struct {
	TotalTenants  int64            `json:"total_tenants"`
	TotalBranches int64            `json:"total_branches"`
	TotalUsers    int64            `json:"total_users"`
	TotalProducts int64            `json:"total_products"`
	Transactions  TransactionStats `json:"transactions"`
}

type OrderStats struct {
	AllTime     int64 `json:"all_time"`
	Today       int64 `json:"today"`
	Last7Days   int64 `json:"last_7_days"`
	Last30Days  int64 `json:"last_30_days"`
	Last90Days  int64 `json:"last_90_days"`
	Last180Days int64 `json:"last_180_days"`
	Last360Days int64 `json:"last_360_days"`
}

type PaymentStats struct {
	AllTime     float64 `json:"all_time"`
	Today       float64 `json:"today"`
	Last7Days   float64 `json:"last_7_days"`
	Last30Days  float64 `json:"last_30_days"`
	Last90Days  float64 `json:"last_90_days"`
	Last180Days float64 `json:"last_180_days"`
	Last360Days float64 `json:"last_360_days"`
}

type TransactionStats struct {
	AllTime     float64 `json:"all_time"`
	Today       float64 `json:"today"`
	ThisWeek    float64 `json:"this_week"`
	ThisMonth   float64 `json:"this_month"`
	Last7Days   float64 `json:"last_7_days"`
	Last30Days  float64 `json:"last_30_days"`
	Last90Days  float64 `json:"last_90_days"`
	Last180Days float64 `json:"last_180_days"`
	Last360Days float64 `json:"last_360_days"`
}
