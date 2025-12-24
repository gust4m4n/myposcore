package dto

type RegisterRequest struct {
	TenantCode string `json:"tenant_code" binding:"required"`
	BranchCode string `json:"branch_code" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FullName   string `json:"full_name"`
}

type LoginRequest struct {
	TenantCode string `json:"tenant_code" binding:"required"`
	BranchCode string `json:"branch_code" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID         uint   `json:"id"`
	TenantID   uint   `json:"tenant_id"`
	BranchID   uint   `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	IsActive   bool   `json:"is_active"`
}
