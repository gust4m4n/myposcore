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

func (s *PaymentService) CreatePayment(orderID uint, amount float64, paymentMethod, notes string, tenantID uint) (*models.Payment, error) {
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
	if err := s.db.Joins("JOIN orders ON orders.id = payments.order_id").
		Where("payments.id = ? AND orders.tenant_id = ?", paymentID, tenantID).
		First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) ListPayments(tenantID, branchID uint) ([]models.Payment, error) {
	var payments []models.Payment
	query := s.db.Joins("JOIN orders ON orders.id = payments.order_id").
		Where("orders.tenant_id = ?", tenantID)

	if branchID > 0 {
		query = query.Where("orders.branch_id = ?", branchID)
	}

	if err := query.Order("payments.created_at DESC").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
