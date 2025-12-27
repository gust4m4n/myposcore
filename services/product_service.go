package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"

	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService() *ProductService {
	return &ProductService{
		db: database.GetDB(),
	}
}

func (s *ProductService) ListProducts(tenantID uint, category, search string) ([]models.Product, error) {
	var products []models.Product
	query := s.db.Where("tenant_id = ?", tenantID)

	// Filter by category if provided
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Search by name or SKU if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR sku ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetCategories returns list of unique categories for a tenant
func (s *ProductService) GetCategories(tenantID uint) ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.Product{}).
		Where("tenant_id = ? AND category IS NOT NULL AND category != ''", tenantID).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *ProductService) GetProduct(id, tenantID uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) CreateProduct(tenantID uint, req dto.CreateProductRequest) (*models.Product, error) {
	product := models.Product{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductService) UpdateProduct(id, tenantID uint, req dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetProduct(id, tenantID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.SKU != "" {
		updates["sku"] = req.SKU
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.db.Model(product).Updates(updates).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(id, tenantID uint) error {
	product, err := s.GetProduct(id, tenantID)
	if err != nil {
		return err
	}

	if err := s.db.Delete(product).Error; err != nil {
		return err
	}

	return nil
}

func (s *ProductService) UpdateProductImage(id, tenantID uint, imageURL string) (*models.Product, error) {
	product, err := s.GetProduct(id, tenantID)
	if err != nil {
		return nil, err
	}

	product.Image = imageURL
	if err := s.db.Save(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}
