package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type LoginService struct {
	db *gorm.DB
}

func NewLoginService() *LoginService {
	return &LoginService{
		db: database.GetDB(),
	}
}

func (s *LoginService) Login(req dto.LoginRequest) (*models.User, *models.Tenant, *models.Branch, error) {
	// Get user by username (username is unique across all tenants)
	var user models.User
	if err := s.db.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
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
