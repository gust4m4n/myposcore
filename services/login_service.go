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

func (s *LoginService) Login(req dto.LoginRequest) (*models.User, *models.Branch, error) {
	// Get tenant
	var tenant models.Tenant
	if err := s.db.Where("code = ? AND is_active = ?", req.TenantCode, true).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, err
	}

	// Check if branch exists for this tenant
	var branch models.Branch
	if err := s.db.Where("code = ? AND tenant_id = ? AND is_active = ?", req.BranchCode, tenant.ID, true).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, err
	}

	// Get user
	var user models.User
	if err := s.db.Where("branch_id = ? AND username = ? AND is_active = ?", branch.ID, req.Username, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, nil, errors.New("invalid credentials")
	}

	return &user, &branch, nil
}
