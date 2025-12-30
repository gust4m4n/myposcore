package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type AdminChangePasswordService struct {
	db *gorm.DB
}

func NewAdminChangePasswordService() *AdminChangePasswordService {
	return &AdminChangePasswordService{
		db: database.GetDB(),
	}
}

// AdminChangePassword allows higher role users to change password of lower role users
func (s *AdminChangePasswordService) AdminChangePassword(adminID uint, req dto.AdminChangePasswordRequest) error {
	// Validate password match
	if req.Password != req.ConfirmPassword {
		return errors.New("password and confirm password do not match")
	}

	// Validate password length
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Get admin user info
	var adminUser models.User
	if err := s.db.First(&adminUser, adminID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("admin user not found")
		}
		return err
	}

	// Get target user by email
	var targetUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&targetUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("target user not found")
		}
		return err
	}

	// Role hierarchy validation
	if err := s.validateRoleHierarchy(adminUser.Role, targetUser.Role); err != nil {
		return err
	}

	// Tenant validation - admins can only change passwords within their tenant
	// Superadmin can change any user's password
	if adminUser.Role != "superadmin" {
		if adminUser.TenantID != targetUser.TenantID {
			return errors.New("cannot change password for users from different tenant")
		}

		// Additional validation for owner - can only change passwords in their branch
		if adminUser.Role == "owner" && adminUser.BranchID != targetUser.BranchID {
			return errors.New("owner can only change password for users in their branch")
		}
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Update password
	if err := s.db.Model(&targetUser).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}

// validateRoleHierarchy checks if admin role has permission to change target user's password
func (s *AdminChangePasswordService) validateRoleHierarchy(adminRole, targetRole string) error {
	// Role hierarchy: superadmin > owner > admin > user
	roleHierarchy := map[string]int{
		"superadmin": 4,
		"owner":      3,
		"admin":      2,
		"user":       1,
	}

	adminLevel, adminExists := roleHierarchy[adminRole]
	targetLevel, targetExists := roleHierarchy[targetRole]

	if !adminExists || !targetExists {
		return errors.New("invalid role")
	}

	// Admin must have higher or equal role level
	if adminLevel <= targetLevel {
		return errors.New("insufficient permission: can only change password for lower role users")
	}

	return nil
}
