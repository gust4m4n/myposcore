package dto

type CreateTenantRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Active bool   `json:"is_active"`
}

type TenantResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type BranchResponse struct {
	ID        uint   `json:"id"`
	TenantID  uint   `json:"tenant_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	TenantID  uint   `json:"tenant_id"`
	BranchID  uint   `json:"branch_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
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
