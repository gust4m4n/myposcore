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
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Address:     req.Address,
		Website:     req.Website,
		Email:       req.Email,
		Phone:       req.Phone,
		IsActive:    req.Active,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *SuperAdminTenantService) UpdateTenant(id uint, req dto.UpdateTenantRequest) (*models.Tenant, error) {
	// Check if tenant exists
	var tenant models.Tenant
	if err := s.db.First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	// Check if code is being changed and if new code already exists
	if tenant.Code != req.Code {
		var existing models.Tenant
		err := s.db.Where("code = ? AND id != ?", req.Code, id).First(&existing).Error
		if err == nil {
			return nil, errors.New("tenant code already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Update tenant
	tenant.Name = req.Name
	tenant.Code = req.Code
	tenant.Description = req.Description
	tenant.Address = req.Address
	tenant.Website = req.Website
	tenant.Email = req.Email
	tenant.Phone = req.Phone
	tenant.IsActive = req.Active

	if err := s.db.Save(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *SuperAdminTenantService) DeleteTenant(id uint) error {
	// Check if tenant exists
	var tenant models.Tenant
	if err := s.db.First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tenant not found")
		}
		return err
	}

	// Soft delete tenant
	if err := s.db.Delete(&tenant).Error; err != nil {
		return err
	}

	return nil
}
