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
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

func (s *UserService) ListUsers(tenantID uint) ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUser(tenantID, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(tenantID uint, req dto.CreateUserRequest) (*models.User, error) {
	// Check if username already exists for this tenant
	var existingUser models.User
	if err := s.db.Where("username = ? AND tenant_id = ?", req.Username, tenantID).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists for this tenant")
	}

	// Check if email already exists
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
		TenantID: tenantID,
		BranchID: req.BranchID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     req.Role,
		IsActive: isActive,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

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

	if req.Username != nil {
		// Check username uniqueness for this tenant
		var existingUser models.User
		if err := s.db.Where("username = ? AND tenant_id = ? AND id != ?", *req.Username, tenantID, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("username already exists for this tenant")
		}
		updates["username"] = *req.Username
	}

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

	if len(updates) == 0 {
		return user, nil
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload user
	if err := s.db.First(user, userID).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(tenantID, userID uint) error {
	user, err := s.GetUser(tenantID, userID)
	if err != nil {
		return err
	}

	if err := s.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}
