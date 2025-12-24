package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	BaseHandler
	paymentService *services.PaymentService
}

func NewPaymentHandler(cfg *config.Config, paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		BaseHandler:    BaseHandler{config: cfg},
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")

	payment, err := h.paymentService.CreatePayment(req.OrderID, req.Amount, req.PaymentMethod, req.Notes, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		Notes:         payment.Notes,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	tenantID := c.GetUint("tenant_id")

	payment, err := h.paymentService.GetPayment(uint(paymentID), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	response := dto.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		Notes:         payment.Notes,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *PaymentHandler) GetPaymentsByOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	tenantID := c.GetUint("tenant_id")

	payments, err := h.paymentService.GetPaymentsByOrder(uint(orderID), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = dto.PaymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			Notes:         payment.Notes,
			CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     payment.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	branchID := c.GetUint("branch_id")

	payments, err := h.paymentService.ListPayments(tenantID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = dto.PaymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			Notes:         payment.Notes,
			CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     payment.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
