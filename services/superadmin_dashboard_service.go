package services

import (
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"

	"gorm.io/gorm"
)

type SuperAdminDashboardService struct {
	db *gorm.DB
}

func NewSuperAdminDashboardService() *SuperAdminDashboardService {
	return &SuperAdminDashboardService{
		db: database.GetDB(),
	}
}

func (s *SuperAdminDashboardService) GetDashboard() (*dto.DashboardResponse, error) {
	var totalTenants, totalBranches, totalUsers int64

	if err := s.db.Model(&models.Tenant{}).Count(&totalTenants).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Branch{}).Count(&totalBranches).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	return &dto.DashboardResponse{
		TotalTenants:  totalTenants,
		TotalBranches: totalBranches,
		TotalUsers:    totalUsers,
	}, nil
}
