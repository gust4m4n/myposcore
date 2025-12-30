package dto

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FullName  string `json:"full_name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=user branchadmin tenantadmin"`
	BranchID  uint   `json:"branch_id" binding:"required"`
	IsActive  *bool  `json:"is_active"`
	CreatedBy *uint  `json:"-"` // Set internally, not from request
}

type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Password  *string `json:"password,omitempty" binding:"omitempty,min=6"`
	FullName  *string `json:"full_name,omitempty"`
	Role      *string `json:"role,omitempty" binding:"omitempty,oneof=user branchadmin tenantadmin"`
	BranchID  *uint   `json:"branch_id,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	UpdatedBy *uint   `json:"-"` // Set internally, not from request
}
