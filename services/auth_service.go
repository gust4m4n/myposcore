package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*models.User, *models.Tenant, *models.Branch, error) {
	// Get user by email (email is unique across all tenants)
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, nil, nil, errors.New("invalid credentials")
	}

	// Get branch
	var branch models.Branch
	if err := s.db.Where("id = ? AND is_active = ?", user.BranchID, true).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("branch not found or inactive")
		}
		return nil, nil, nil, err
	}

	// Get tenant
	var tenant models.Tenant
	if err := s.db.Where("id = ? AND is_active = ?", user.TenantID, true).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("tenant not found or inactive")
		}
		return nil, nil, nil, err
	}

	return &user, &tenant, &branch, nil
}

func (s *AuthService) GetProfile(userID uint) (*dto.ProfileResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var tenant models.Tenant
	if err := s.db.First(&tenant, user.TenantID).Error; err != nil {
		return nil, err
	}

	var branch models.Branch
	if err := s.db.First(&branch, user.BranchID).Error; err != nil {
		return nil, err
	}

	profile := &dto.ProfileResponse{
		User: dto.UserDetailProfile{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Image:    user.Image,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
		Tenant: dto.TenantDetailProfile{
			ID:       tenant.ID,
			Name:     tenant.Name,
			IsActive: tenant.IsActive,
		},
		Branch: dto.BranchDetailProfile{
			ID:       branch.ID,
			Name:     branch.Name,
			Address:  branch.Address,
			Phone:    branch.Phone,
			IsActive: branch.IsActive,
		},
	}

	return profile, nil
}

func (s *AuthService) UpdateProfileImage(userID uint, imageURL string) (*dto.ProfileResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Image = imageURL
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}

func (s *AuthService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	user.Email = req.Email
	user.FullName = req.FullName

	// Update PIN if provided
	if req.PIN != "" {
		hashedPIN, err := utils.HashPassword(req.PIN)
		if err != nil {
			return nil, errors.New("failed to hash PIN")
		}
		user.PIN = hashedPIN
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}
