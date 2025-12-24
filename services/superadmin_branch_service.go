package services

import (
	"myposcore/database"
	"myposcore/models"

	"gorm.io/gorm"
)

type SuperAdminBranchService struct {
	db *gorm.DB
}

func NewSuperAdminBranchService() *SuperAdminBranchService {
	return &SuperAdminBranchService{
		db: database.GetDB(),
	}
}

func (s *SuperAdminBranchService) ListBranchesByTenant(tenantID uint) ([]models.Branch, error) {
	var branches []models.Branch
	if err := s.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}
