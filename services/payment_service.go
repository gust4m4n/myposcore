package services

import (
	"errors"
	"myposcore/models"

	"gorm.io/gorm"
)

type PaymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (s *PaymentService) CreatePayment(orderID uint, amount float64, paymentMethod, notes string, tenantID uint, createdBy *uint) (*models.Payment, error) {
	// Verify order exists and belongs to tenant
	var order models.Order
	if err := s.db.Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Check if order is already paid
	if order.Status == "completed" {
		return nil, errors.New("order already completed")
	}

	// Check if order is cancelled
	if order.Status == "cancelled" {
		return nil, errors.New("cannot pay cancelled order")
	}

	// Validate payment amount
	if amount < order.TotalAmount {
		return nil, errors.New("payment amount less than order total")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create payment
	payment := &models.Payment{
		OrderID:       orderID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "completed",
		Notes:         notes,
		CreatedBy:     createdBy,
	}

	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update order status
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status": "completed",
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPaymentWithDetails(paymentID, tenantID uint) (*models.Payment, *models.Order, error) {
	var payment models.Payment
	if err := s.db.Joins("JOIN orders ON orders.id = payments.order_id").
		Where("payments.id = ? AND orders.tenant_id = ?", paymentID, tenantID).
		First(&payment).Error; err != nil {
		return nil, nil, err
	}

	// Get order with all related data
	var order models.Order
	if err := s.db.Preload("OrderItems.Product").
		Preload("User").
		Preload("Branch").
		Where("id = ?", payment.OrderID).
		First(&order).Error; err != nil {
		return nil, nil, err
	}

	return &payment, &order, nil
}

func (s *PaymentService) GetPaymentsByOrder(orderID, tenantID uint) ([]models.Payment, error) {
	var payments []models.Payment
	if err := s.db.Joins("JOIN orders ON orders.id = payments.order_id").
		Where("payments.order_id = ? AND orders.tenant_id = ?", orderID, tenantID).
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *PaymentService) GetPayment(paymentID, tenantID uint) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.Preload("Creator").Preload("Updater").Joins("JOIN orders ON orders.id = payments.order_id").
		Where("payments.id = ? AND orders.tenant_id = ?", paymentID, tenantID).
		First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) ListPayments(tenantID, branchID uint, page, perPage int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := s.db.Model(&models.Payment{}).
		Joins("JOIN orders ON orders.id = payments.order_id").
		Where("orders.tenant_id = ?", tenantID)

	if branchID > 0 {
		query = query.Where("orders.branch_id = ?", branchID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	if err := query.Preload("Creator").Preload("Updater").Order("payments.created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}
