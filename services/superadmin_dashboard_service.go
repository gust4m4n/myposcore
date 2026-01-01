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
	var totalTenants, totalBranches, totalUsers, totalProducts, totalCategories int64

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

	// Count total categories
	if err := s.db.Model(&models.Category{}).Count(&totalCategories).Error; err != nil {
		return nil, err
	}

	// Get time boundaries
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start7Days := startOfDay.AddDate(0, 0, -7)
	start30Days := startOfDay.AddDate(0, 0, -30)
	start90Days := startOfDay.AddDate(0, 0, -90)
	start180Days := startOfDay.AddDate(0, 0, -180)
	start360Days := startOfDay.AddDate(0, 0, -360)

	// Calculate order statistics
	orderStats := dto.OrderStats{}

	// Orders all time
	if err := s.db.Model(&models.Order{}).Count(&orderStats.AllTime).Error; err != nil {
		return nil, err
	}

	// Orders today
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", startOfDay).Count(&orderStats.Today).Error; err != nil {
		return nil, err
	}

	// Orders last 7 days
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", start7Days).Count(&orderStats.Last7Days).Error; err != nil {
		return nil, err
	}

	// Orders last 30 days
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", start30Days).Count(&orderStats.Last30Days).Error; err != nil {
		return nil, err
	}

	// Orders last 90 days
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", start90Days).Count(&orderStats.Last90Days).Error; err != nil {
		return nil, err
	}

	// Orders last 180 days
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", start180Days).Count(&orderStats.Last180Days).Error; err != nil {
		return nil, err
	}

	// Orders last 360 days
	if err := s.db.Model(&models.Order{}).Where("created_at >= ?", start360Days).Count(&orderStats.Last360Days).Error; err != nil {
		return nil, err
	}

	// Calculate payment statistics
	paymentStats := dto.PaymentStats{}

	// Payments all time
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Scan(&paymentStats.AllTime).Error; err != nil {
		return nil, err
	}

	// Payments today
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", startOfDay).Scan(&paymentStats.Today).Error; err != nil {
		return nil, err
	}

	// Payments last 7 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start7Days).Scan(&paymentStats.Last7Days).Error; err != nil {
		return nil, err
	}

	// Payments last 30 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start30Days).Scan(&paymentStats.Last30Days).Error; err != nil {
		return nil, err
	}

	// Payments last 90 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start90Days).Scan(&paymentStats.Last90Days).Error; err != nil {
		return nil, err
	}

	// Payments last 180 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start180Days).Scan(&paymentStats.Last180Days).Error; err != nil {
		return nil, err
	}

	// Payments last 360 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start360Days).Scan(&paymentStats.Last360Days).Error; err != nil {
		return nil, err
	}

	// Calculate transaction statistics (count of payment records)
	transactionStats := dto.TransactionStats{}

	// Transactions all time
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Scan(&transactionStats.AllTime).Error; err != nil {
		return nil, err
	}

	// Transactions today
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", startOfDay).Scan(&transactionStats.Today).Error; err != nil {
		return nil, err
	}

	// Transactions last 7 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start7Days).Scan(&transactionStats.Last7Days).Error; err != nil {
		return nil, err
	}

	// Transactions last 30 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start30Days).Scan(&transactionStats.Last30Days).Error; err != nil {
		return nil, err
	}

	// Transactions last 90 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start90Days).Scan(&transactionStats.Last90Days).Error; err != nil {
		return nil, err
	}

	// Transactions last 180 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start180Days).Scan(&transactionStats.Last180Days).Error; err != nil {
		return nil, err
	}

	// Transactions last 360 days
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", start360Days).Scan(&transactionStats.Last360Days).Error; err != nil {
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
			IsActive:  tenant.IsActive,
			CreatedAt: tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.DashboardResponse{
		TotalTenants:    totalTenants,
		TotalBranches:   totalBranches,
		TotalUsers:      totalUsers,
		TotalProducts:   totalProducts,
		TotalCategories: totalCategories,
		Orders:          orderStats,
		Payments:        paymentStats,
		Transactions:    transactionStats,
		Tenants:         tenantResponses,
	}, nil
}
