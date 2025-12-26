package services

import (
	"errors"
	"fmt"
	"myposcore/models"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) CreateOrder(tenantID, branchID, userID uint, items []struct {
	ProductID uint
	Quantity  int
}) (*models.Order, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate all products exist and belong to tenant
	var products []models.Product
	productIDs := make([]uint, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}

	if err := tx.Where("id IN ? AND tenant_id = ? AND is_active = ?", productIDs, tenantID, true).
		Find(&products).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(products) != len(items) {
		tx.Rollback()
		return nil, errors.New("some products not found or inactive")
	}

	// Create product map for easy lookup
	productMap := make(map[uint]*models.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", tenantID, time.Now().Unix())

	// Create order
	order := &models.Order{
		TenantID:    tenantID,
		BranchID:    branchID,
		UserID:      userID,
		OrderNumber: orderNumber,
		Status:      "pending",
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order items and calculate total
	var totalAmount float64
	orderItems := make([]models.OrderItem, len(items))
	for i, item := range items {
		product := productMap[item.ProductID]
		if product == nil {
			tx.Rollback()
			return nil, fmt.Errorf("product ID %d not found", item.ProductID)
		}

		// Check stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		subtotal := product.Price * float64(item.Quantity)
		orderItems[i] = models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		}
		totalAmount += subtotal

		// Update product stock
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update order total
	order.TotalAmount = totalAmount
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load order items with products
	s.db.Preload("OrderItems.Product").First(order, order.ID)

	return order, nil
}

func (s *OrderService) GetOrder(orderID, tenantID uint) (*models.Order, error) {
	var order models.Order
	if err := s.db.Preload("OrderItems.Product").
		Where("id = ? AND tenant_id = ?", orderID, tenantID).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) ListOrders(tenantID, branchID uint, page, perPage int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := s.db.Model(&models.Order{}).Where("tenant_id = ?", tenantID)

	if branchID > 0 {
		query = query.Where("branch_id = ?", branchID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	if err := query.Preload("OrderItems.Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (s *OrderService) UpdateOrderStatus(orderID, tenantID uint, status string) error {
	return s.db.Model(&models.Order{}).
		Where("id = ? AND tenant_id = ?", orderID, tenantID).
		Update("status", status).Error
}
