package services

import (
	"myposcore/database"
	"myposcore/models"

	"gorm.io/gorm"
)

type SuperAdminUserService struct {
	db *gorm.DB
}

func NewSuperAdminUserService() *SuperAdminUserService {
	return &SuperAdminUserService{
		db: database.GetDB(),
	}
}

func (s *SuperAdminUserService) ListUsersByBranch(branchID uint) ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("branch_id = ?", branchID).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
