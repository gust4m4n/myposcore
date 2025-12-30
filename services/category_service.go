package services

import (
	"errors"
	"myposcore/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) CreateCategory(tenantID uint, name, description string, createdBy *uint) (*models.Category, error) {
	// Check if category name already exists for this tenant
	var existing models.Category
	if err := s.db.Where("tenant_id = ? AND name = ?", tenantID, name).First(&existing).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	category := &models.Category{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategory(categoryID, tenantID uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Creator").Preload("Updater").Where("id = ? AND tenant_id = ?", categoryID, tenantID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) ListCategories(tenantID uint, activeOnly bool) ([]models.Category, error) {
	var categories []models.Category
	query := s.db.Preload("Creator").Preload("Updater").Where("tenant_id = ?", tenantID)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) UpdateCategory(categoryID, tenantID uint, name, description *string, isActive *bool, updatedBy *uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("id = ? AND tenant_id = ?", categoryID, tenantID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})

	if name != nil {
		// Check if new name already exists for another category
		var existing models.Category
		if err := s.db.Where("tenant_id = ? AND name = ? AND id != ?", tenantID, *name, categoryID).First(&existing).Error; err == nil {
			return nil, errors.New("category name already exists")
		}
		updates["name"] = *name
	}

	if description != nil {
		updates["description"] = *description
	}

	if isActive != nil {
		updates["is_active"] = *isActive
	}

	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}

	if len(updates) > 0 {
		if err := s.db.Model(&category).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload to get updated values
	if err := s.db.Preload("Creator").Preload("Updater").First(&category, categoryID).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) DeleteCategory(categoryID, tenantID uint, deletedBy *uint) error {
	var category models.Category
	if err := s.db.Where("id = ? AND tenant_id = ?", categoryID, tenantID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category is used by any products
	var productCount int64
	if err := s.db.Model(&models.Product{}).Where("tenant_id = ? AND category = ?", tenantID, category.Name).Count(&productCount).Error; err != nil {
		return err
	}

	if productCount > 0 {
		return errors.New("cannot delete category that is used by products")
	}

	if deletedBy != nil {
		if err := s.db.Model(&category).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}

	return s.db.Delete(&category).Error
}
