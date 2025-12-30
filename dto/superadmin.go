package dto

type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type UpdateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
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
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Active      bool   `json:"is_active"`
}

type UpdateBranchRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
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
	Code          string  `json:"code"`
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
	Code          string  `json:"code"`
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
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	FullName      string  `json:"full_name"`
	Role          string  `json:"role"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	CreatedBy     *uint   `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name,omitempty"`
	UpdatedBy     *uint   `json:"updated_by,omitempty"`
	UpdatedByName *string `json:"updated_by_name,omitempty"`
}

type DashboardResponse struct {
	TotalTenants         int64            `json:"total_tenants"`
	TotalBranches        int64            `json:"total_branches"`
	TotalUsers           int64            `json:"total_users"`
	TotalProducts        int64            `json:"total_products"`
	TotalOrders          int64            `json:"total_orders"`
	TotalOrdersToday     int64            `json:"total_orders_today"`
	TotalOrdersThisWeek  int64            `json:"total_orders_this_week"`
	TotalOrdersThisMonth int64            `json:"total_orders_this_month"`
	Tenants              []TenantResponse `json:"tenants"`
}
