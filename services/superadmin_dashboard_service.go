package services

import (
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"time"

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
	var totalTenants, totalBranches, totalUsers, totalProducts int64
	var totalOrders, totalOrdersToday, totalOrdersThisWeek, totalOrdersThisMonth int64

	// Count total tenants
	if err := s.db.Model(&models.Tenant{}).Count(&totalTenants).Error; err != nil {
		return nil, err
	}

	// Count total branches
	if err := s.db.Model(&models.Branch{}).Count(&totalBranches).Error; err != nil {
		return nil, err
	}

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	// Count total products
	if err := s.db.Model(&models.Product{}).Count(&totalProducts).Error; err != nil {
		return nil, err
	}

	// Count total orders (all time)
	if err := s.db.Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}

	// Get time boundaries
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Count orders today
	if err := s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfDay).
		Count(&totalOrdersToday).Error; err != nil {
		return nil, err
	}

	// Count orders this week
	if err := s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfWeek).
		Count(&totalOrdersThisWeek).Error; err != nil {
		return nil, err
	}

	// Count orders this month
	if err := s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfMonth).
		Count(&totalOrdersThisMonth).Error; err != nil {
		return nil, err
	}

	// Get all tenants
	var tenants []models.Tenant
	if err := s.db.Find(&tenants).Error; err != nil {
		return nil, err
	}

	var tenantResponses []dto.TenantResponse
	for _, tenant := range tenants {
		tenantResponses = append(tenantResponses, dto.TenantResponse{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Code:      tenant.Code,
			IsActive:  tenant.IsActive,
			CreatedAt: tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.DashboardResponse{
		TotalTenants:         totalTenants,
		TotalBranches:        totalBranches,
		TotalUsers:           totalUsers,
		TotalProducts:        totalProducts,
		TotalOrders:          totalOrders,
		TotalOrdersToday:     totalOrdersToday,
		TotalOrdersThisWeek:  totalOrdersThisWeek,
		TotalOrdersThisMonth: totalOrdersThisMonth,
		Tenants:              tenantResponses,
	}, nil
}
