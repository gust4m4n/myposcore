package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type AdminChangePINService struct {
	db *gorm.DB
}

func NewAdminChangePINService() *AdminChangePINService {
	return &AdminChangePINService{
		db: database.GetDB(),
	}
}

// AdminChangePIN allows higher role users to change PIN of lower role users
func (s *AdminChangePINService) AdminChangePIN(adminID uint, req dto.AdminChangePINRequest) error {
	// Validate PIN match
	if req.PIN != req.ConfirmPIN {
		return errors.New("PIN and confirm PIN do not match")
	}

	// Validate PIN format (6 digits)
	if len(req.PIN) != 6 {
		return errors.New("PIN must be exactly 6 digits")
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

	// Tenant validation - admins can only change PINs within their tenant
	// Superadmin can change any user's PIN
	if adminUser.Role != "superadmin" {
		if adminUser.TenantID != targetUser.TenantID {
			return errors.New("cannot change PIN for users from different tenant")
		}

		// Additional validation for owner - can only change PINs in their branch
		if adminUser.Role == "owner" && adminUser.BranchID != targetUser.BranchID {
			return errors.New("owner can only change PIN for users in their branch")
		}
	}

	// Hash new PIN
	hashedPIN, err := utils.HashPassword(req.PIN)
	if err != nil {
		return err
	}

	// Update PIN
	if err := s.db.Model(&targetUser).Update("pin", hashedPIN).Error; err != nil {
		return err
	}

	return nil
}

// validateRoleHierarchy checks if admin role has permission to change target user's PIN
func (s *AdminChangePINService) validateRoleHierarchy(adminRole, targetRole string) error {
	// Role hierarchy: superadmin > owner > admin > staff
	roleHierarchy := map[string]int{
		"superadmin": 4,
		"owner":      3,
		"admin":      2,
		"staff":      1,
	}

	adminLevel, adminExists := roleHierarchy[adminRole]
	targetLevel, targetExists := roleHierarchy[targetRole]

	if !adminExists || !targetExists {
		return errors.New("invalid role")
	}

	// Admin must have higher or equal role level
	if adminLevel <= targetLevel {
		return errors.New("insufficient permission: can only change PIN for lower role users")
	}

	return nil
}
