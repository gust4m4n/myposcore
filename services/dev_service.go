package services

import (
	"myposcore/database"
	"myposcore/models"
)

type DevService struct{}

func NewDevService() *DevService {
	return &DevService{}
}

func (s *DevService) ListTenants() ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := database.DB.Where("is_active = ?", true).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s *DevService) ListBranchesByTenant(tenantID uint) ([]models.Branch, error) {
	var branches []models.Branch
	if err := database.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true).Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}
