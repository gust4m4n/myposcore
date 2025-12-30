package dto

type RegisterRequest struct {
	TenantCode string `json:"tenant_code" binding:"required"`
	BranchCode string `json:"branch_code" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string      `json:"token"`
	User   UserProfile `json:"user"`
	Tenant TenantInfo  `json:"tenant"`
	Branch BranchInfo  `json:"branch"`
}

type TenantInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
	IsActive    bool   `json:"is_active"`
}

type BranchInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Image       string `json:"image"`
	IsActive    bool   `json:"is_active"`
}

type UserProfile struct {
	ID         uint   `json:"id"`
	TenantID   uint   `json:"tenant_id"`
	BranchID   uint   `json:"branch_id"`
	BranchName string `json:"branch_name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AdminChangePasswordRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	PIN      string `json:"pin"`
}

type ProfileResponse struct {
	User   UserDetailProfile   `json:"user"`
	Tenant TenantDetailProfile `json:"tenant"`
	Branch BranchDetailProfile `json:"branch"`
}

type UserDetailProfile struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Image    string `json:"image"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type TenantDetailProfile struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

type BranchDetailProfile struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}
