package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
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

func (s *SuperAdminBranchService) ListBranches(tenantID uint, search string, page, pageSize int) ([]models.Branch, int64, error) {
	var branches []models.Branch
	var total int64

	query := s.db.Model(&models.Branch{}).Where("tenant_id = ?", tenantID)

	// Search by name or id if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR CAST(id AS TEXT) = ?", searchPattern, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query2 := s.db.Preload("Creator").Preload("Updater").Where("tenant_id = ?", tenantID)

	if search != "" {
		searchPattern := "%" + search + "%"
		query2 = query2.Where("name ILIKE ? OR CAST(id AS TEXT) = ?", searchPattern, search)
	}

	if err := query2.Order("name ASC").Limit(pageSize).Offset(offset).Find(&branches).Error; err != nil {
		return nil, 0, err
	}
	return branches, total, nil
}

func (s *SuperAdminBranchService) CreateBranch(req dto.CreateBranchRequest, imageURL string, createdBy *uint) (*models.Branch, error) {
	// Check if tenant exists
	var tenant models.Tenant
	if err := s.db.First(&tenant, req.TenantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	// Create branch
	branch := models.Branch{
		TenantID:    req.TenantID,
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

	if err := s.db.Create(&branch).Error; err != nil {
		return nil, err
	}

	return &branch, nil
}

func (s *SuperAdminBranchService) UpdateBranch(id uint, req dto.UpdateBranchRequest, imageURL string, updatedBy *uint) (*models.Branch, error) {
	// Check if branch exists
	var branch models.Branch
	if err := s.db.First(&branch, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}

	// Update branch
	branch.Name = req.Name
	branch.Description = req.Description
	branch.Address = req.Address
	branch.Website = req.Website
	branch.Email = req.Email
	branch.Phone = req.Phone
	branch.IsActive = req.Active
	branch.UpdatedBy = updatedBy

	// Update image if provided
	if imageURL != "" {
		branch.Image = imageURL
	}

	if err := s.db.Save(&branch).Error; err != nil {
		return nil, err
	}

	return &branch, nil
}

func (s *SuperAdminBranchService) GetBranchByID(id uint) (*models.Branch, error) {
	var branch models.Branch
	if err := s.db.First(&branch, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}
	return &branch, nil
}

func (s *SuperAdminBranchService) DeleteBranch(id uint, deletedBy *uint) error {
	// Check if branch exists
	var branch models.Branch
	if err := s.db.First(&branch, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("branch not found")
		}
		return err
	}

	// Set deleted_by before soft delete
	if deletedBy != nil {
		if err := s.db.Model(&branch).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}

	// Soft delete branch
	if err := s.db.Delete(&branch).Error; err != nil {
		return err
	}

	return nil
}
