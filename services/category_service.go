package services

import (
	"encoding/json"
	"errors"
	"myposcore/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db                *gorm.DB
	auditTrailService *AuditTrailService
}

func NewCategoryService(db *gorm.DB, auditTrailService *AuditTrailService) *CategoryService {
	return &CategoryService{
		db:                db,
		auditTrailService: auditTrailService,
	}
}

func (s *CategoryService) CreateCategory(tenantID uint, name, description string, imageURL string, createdBy *uint) (*models.Category, error) {
	// Check if category name already exists for this tenant
	var existing models.Category
	if err := s.db.Where("tenant_id = ? AND name = ?", tenantID, name).First(&existing).Error; err == nil {
		return nil, errors.New("category name already exists")
	}

	category := &models.Category{
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Image:       imageURL,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"image":       category.Image,
		"is_active":   category.IsActive,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if createdBy != nil {
		userID = *createdBy
	}
	_ = s.auditTrailService.CreateAuditTrail(&tenantID, nil, userID, "category", category.ID, "create", changesMap, "", "")

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

func (s *CategoryService) ListCategories(tenantID uint, activeOnly bool, page, pageSize int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := s.db.Model(&models.Category{}).Where("tenant_id = ?", tenantID)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query2 := s.db.Preload("Creator").Preload("Updater").Where("tenant_id = ?", tenantID)
	if activeOnly {
		query2 = query2.Where("is_active = ?", true)
	}
	if err := query2.Order("name ASC").Limit(pageSize).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (s *CategoryService) UpdateCategory(categoryID, tenantID uint, name, description *string, imageURL *string, isActive *bool, updatedBy *uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("id = ? AND tenant_id = ?", categoryID, tenantID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Save old values for audit trail
	oldValues := map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"image":       category.Image,
		"is_active":   category.IsActive,
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

	if imageURL != nil && *imageURL != "" {
		updates["image"] = *imageURL
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

		// Create audit trail with old/new values
		changes := make(map[string]interface{})
		for key, newVal := range updates {
			if key != "updated_by" {
				if oldVal, exists := oldValues[key]; exists {
					changes[key] = map[string]interface{}{
						"old": oldVal,
						"new": newVal,
					}
				}
			}
		}
		changesJSON, _ := json.Marshal(changes)
		var changesMap map[string]interface{}
		_ = json.Unmarshal(changesJSON, &changesMap)
		var userID uint
		if updatedBy != nil {
			userID = *updatedBy
		}
		_ = s.auditTrailService.CreateAuditTrail(&tenantID, nil, userID, "category", categoryID, "update", changesMap, "", "")
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

	if err := s.db.Delete(&category).Error; err != nil {
		return err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"image":       category.Image,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if deletedBy != nil {
		userID = *deletedBy
	}
	_ = s.auditTrailService.CreateAuditTrail(&tenantID, nil, userID, "category", categoryID, "delete", changesMap, "", "")

	return nil
}
