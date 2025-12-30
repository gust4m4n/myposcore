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
	// Get user by email (email is unique across all tenants)
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("user tidak ditemukan")
		}
		return nil, nil, nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, nil, errors.New("user tidak aktif")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, nil, nil, errors.New("password salah")
	}

	// Get branch
	var branch models.Branch
	if err := s.db.Where("id = ?", user.BranchID).First(&branch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("branch tidak ditemukan")
		}
		return nil, nil, nil, err
	}

	// Check if branch is active
	if !branch.IsActive {
		return nil, nil, nil, errors.New("branch tidak aktif")
	}

	// Get tenant
	var tenant models.Tenant
	if err := s.db.Where("id = ?", user.TenantID).First(&tenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, errors.New("tenant tidak ditemukan")
		}
		return nil, nil, nil, err
	}

	// Check if tenant is active
	if !tenant.IsActive {
		return nil, nil, nil, errors.New("tenant tidak aktif")
	}

	return &user, &tenant, &branch, nil
}
