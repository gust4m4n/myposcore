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

	// Get time boundaries
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Calculate start of week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	startOfWeek := startOfDay.AddDate(0, 0, -(weekday - 1))

	// Calculate start of month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	start7Days := startOfDay.AddDate(0, 0, -7)
	start30Days := startOfDay.AddDate(0, 0, -30)
	start90Days := startOfDay.AddDate(0, 0, -90)
	start180Days := startOfDay.AddDate(0, 0, -180)
	start360Days := startOfDay.AddDate(0, 0, -360)

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

	// Transactions this week
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", startOfWeek).Scan(&transactionStats.ThisWeek).Error; err != nil {
		return nil, err
	}

	// Transactions this month
	if err := s.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Where("created_at >= ?", startOfMonth).Scan(&transactionStats.ThisMonth).Error; err != nil {
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

	return &dto.DashboardResponse{
		TotalTenants:  totalTenants,
		TotalBranches: totalBranches,
		TotalUsers:    totalUsers,
		TotalProducts: totalProducts,
		Transactions:  transactionStats,
	}, nil
}
