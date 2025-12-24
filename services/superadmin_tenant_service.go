package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"

	"gorm.io/gorm"
)

type SuperAdminTenantService struct {
	db *gorm.DB
}

func NewSuperAdminTenantService() *SuperAdminTenantService {
	return &SuperAdminTenantService{
		db: database.GetDB(),
	}
}

func (s *SuperAdminTenantService) ListTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := s.db.Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s *SuperAdminTenantService) CreateTenant(req dto.CreateTenantRequest) (*models.Tenant, error) {
	// Check if tenant code already exists
	var existing models.Tenant
	err := s.db.Where("code = ?", req.Code).First(&existing).Error
	if err == nil {
		return nil, errors.New("tenant code already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create tenant
	tenant := models.Tenant{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: req.Active,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}
