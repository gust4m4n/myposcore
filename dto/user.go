package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=user branchadmin tenantadmin"`
	BranchID uint   `json:"branch_id" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	FullName *string `json:"full_name,omitempty"`
	Role     *string `json:"role,omitempty" binding:"omitempty,oneof=user branchadmin tenantadmin"`
	BranchID *uint   `json:"branch_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
