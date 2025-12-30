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
	if err := s.db.Preload("Creator").Preload("Updater").Order("name ASC").Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s *SuperAdminTenantService) CreateTenant(req dto.CreateTenantRequest, imageURL string, createdBy *uint) (*models.Tenant, error) {
	// Create tenant
	tenant := models.Tenant{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Website:     req.Website,
		Email:       req.Email,
		Phone:       req.Phone,
		Image:       imageURL,
		IsActive:    req.Active,
		CreatedBy:   createdBy,
	}

	if err := s.db.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *SuperAdminTenantService) UpdateTenant(id uint, req dto.UpdateTenantRequest, imageURL string, updatedBy *uint) (*models.Tenant, error) {
	// Check if tenant exists
	var tenant models.Tenant
	if err := s.db.First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	// Update tenant
	tenant.Name = req.Name
	tenant.Description = req.Description
	tenant.Address = req.Address
	tenant.Website = req.Website
	tenant.Email = req.Email
	tenant.Phone = req.Phone
	tenant.IsActive = req.Active
	tenant.UpdatedBy = updatedBy

	// Update image if provided
	if imageURL != "" {
		tenant.Image = imageURL
	}

	if err := s.db.Save(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *SuperAdminTenantService) GetTenantByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := s.db.First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return &tenant, nil
}

func (s *SuperAdminTenantService) DeleteTenant(id uint, deletedBy *uint) error {
	// Check if tenant exists
	var tenant models.Tenant
	if err := s.db.First(&tenant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tenant not found")
		}
		return err
	}

	// Set deleted_by before soft delete
	if deletedBy != nil {
		if err := s.db.Model(&tenant).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}

	// Soft delete tenant
	if err := s.db.Delete(&tenant).Error; err != nil {
		return err
	}

	return nil
}
