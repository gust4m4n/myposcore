package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db                *gorm.DB
	auditTrailService *AuditTrailService
}

func NewUserService(auditTrailService *AuditTrailService) *UserService {
	return &UserService{
		db:                database.GetDB(),
		auditTrailService: auditTrailService,
	}
}

func (s *UserService) ListUsers(tenantID uint, search string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("tenant_id = ?", tenantID)

	// Search by full_name, email, or id if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR email ILIKE ? OR CAST(id AS TEXT) = ?", searchPattern, searchPattern, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query2 := s.db.Preload("Creator").Preload("Updater").Where("tenant_id = ?", tenantID)

	if search != "" {
		searchPattern := "%" + search + "%"
		query2 = query2.Where("full_name ILIKE ? OR email ILIKE ? OR CAST(id AS TEXT) = ?", searchPattern, searchPattern, search)
	}

	if err := query2.Order("full_name ASC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) GetUser(tenantID, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Creator").Preload("Updater").Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(tenantID uint, req dto.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Verify branch belongs to tenant
	var branch models.Branch
	if err := s.db.Where("id = ? AND tenant_id = ?", req.BranchID, tenantID).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found or doesn't belong to this tenant")
		}
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user := models.User{
		TenantID:  tenantID,
		BranchID:  &req.BranchID,
		Email:     req.Email,
		Password:  hashedPassword,
		FullName:  req.FullName,
		Role:      req.Role,
		IsActive:  isActive,
		CreatedBy: req.CreatedBy,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"email":     user.Email,
		"full_name": user.FullName,
		"role":      user.Role,
		"branch_id": user.BranchID,
		"is_active": user.IsActive,
	}
	var auditUserID uint
	if req.CreatedBy != nil {
		auditUserID = *req.CreatedBy
	}
	_ = s.auditTrailService.CreateAuditTrail(&tenantID, user.BranchID, auditUserID, "user", user.ID, "create", changes, "", "")

	return &user, nil
}

func (s *UserService) UpdateUser(tenantID, userID uint, req dto.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.GetUser(tenantID, userID)
	if err != nil {
		return nil, err
	}

	// Prepare update map
	updates := make(map[string]interface{})

	if req.Email != nil {
		// Check email uniqueness
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", *req.Email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		updates["email"] = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		updates["password"] = hashedPassword
	}

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.BranchID != nil {
		// Verify branch belongs to tenant
		var branch models.Branch
		if err := s.db.Where("id = ? AND tenant_id = ?", *req.BranchID, tenantID).First(&branch).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("branch not found or doesn't belong to this tenant")
			}
			return nil, err
		}
		updates["branch_id"] = *req.BranchID
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Set updated_by
	if req.UpdatedBy != nil {
		updates["updated_by"] = *req.UpdatedBy
	}

	if len(updates) == 0 {
		return user, nil
	}

	// Save old values for audit
	oldValues := make(map[string]interface{})
	if req.Email != nil {
		oldValues["email"] = user.Email
	}
	if req.FullName != nil {
		oldValues["full_name"] = user.FullName
	}
	if req.Role != nil {
		oldValues["role"] = user.Role
	}
	if req.BranchID != nil {
		oldValues["branch_id"] = user.BranchID
	}
	if req.IsActive != nil {
		oldValues["is_active"] = user.IsActive
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload user
	if err := s.db.First(user, userID).Error; err != nil {
		return nil, err
	}

	// Create audit trail with changes
	if len(oldValues) > 0 {
		changes := make(map[string]interface{})
		for key := range oldValues {
			changes[key] = map[string]interface{}{
				"old": oldValues[key],
				"new": updates[key],
			}
		}
		auditorID := user.ID
		if req.UpdatedBy != nil {
			auditorID = *req.UpdatedBy
		}
		_ = s.auditTrailService.CreateAuditTrail(&tenantID, user.BranchID, auditorID, "user", user.ID, "update", changes, "", "")
	}

	return user, nil
}

func (s *UserService) DeleteUser(tenantID, userID uint, deletedBy *uint) error {
	user, err := s.GetUser(tenantID, userID)
	if err != nil {
		return err
	}

	// Set deleted_by before soft delete
	if deletedBy != nil {
		user.DeletedBy = deletedBy
		if err := s.db.Save(user).Error; err != nil {
			return err
		}
	}

	if err := s.db.Delete(user).Error; err != nil {
		return err
	}

	// Create audit trail
	auditorID := userID
	if deletedBy != nil {
		auditorID = *deletedBy
	}
	changes := map[string]interface{}{
		"email":     user.Email,
		"full_name": user.FullName,
		"role":      user.Role,
	}
	_ = s.auditTrailService.CreateAuditTrail(&tenantID, user.BranchID, auditorID, "user", user.ID, "delete", changes, "", "")

	return nil
}

func (s *UserService) UpdateUserImage(userID, tenantID uint, imageURL string) (*models.User, error) {
	user, err := s.GetUser(tenantID, userID)
	if err != nil {
		return nil, err
	}

	user.Image = imageURL
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
