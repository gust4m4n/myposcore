package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"
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
		utils.BadRequest(c, err.Error())
		return
	}

	tenantID := c.GetUint("tenant_id")

	// Set created_by to current user
	currentUserID := c.GetUint("user_id")
	req.CreatedBy = &currentUserID

	payment, err := h.paymentService.CreatePayment(req.OrderID, req.Amount, req.PaymentMethod, req.Notes, tenantID, req.CreatedBy)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	// Get payment with full details for receipt
	payment, order, err := h.paymentService.GetPaymentWithDetails(payment.ID, tenantID)
	if err != nil {
		utils.InternalError(c, "Failed to get payment details")
		return
	}

	// Build order items detail
	orderItems := make([]dto.PaymentOrderItemDetail, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dto.PaymentOrderItemDetail{
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		}
	}

	// Calculate change
	change := req.Amount - order.TotalAmount
	if change < 0 {
		change = 0
	}

	response := dto.PaymentDetailResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		OrderNumber:   order.OrderNumber,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
		Notes:         payment.Notes,
		Change:        change,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
		Order: dto.PaymentOrderDetail{
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Notes:       order.Notes,
			OrderItems:  orderItems,
			CashierName: order.User.FullName,
			BranchName:  order.Branch.Name,
		},
	}

	utils.Success(c, "Success", response)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid payment ID")
		return
	}

	tenantID := c.GetUint("tenant_id")

	payment, err := h.paymentService.GetPayment(uint(paymentID), tenantID)
	if err != nil {
		utils.NotFound(c, "Payment not found")
		return
	}

	var createdByName, updatedByName *string
	if payment.Creator != nil {
		name := payment.Creator.FullName
		createdByName = &name
	}
	if payment.Updater != nil {
		name := payment.Updater.FullName
		updatedByName = &name
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
		CreatedBy:     payment.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     payment.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	utils.Success(c, "Success", response)
}

func (h *PaymentHandler) GetPaymentsByOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	tenantID := c.GetUint("tenant_id")

	payments, err := h.paymentService.GetPaymentsByOrder(uint(orderID), tenantID)
	if err != nil {
		utils.InternalError(c, err.Error())
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

	utils.Success(c, "Success", responses)
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	branchID := c.GetUint("branch_id")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "32"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 32
	}

	payments, total, err := h.paymentService.ListPayments(tenantID, branchID, page, perPage)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		var createdByName, updatedByName *string
		if payment.Creator != nil {
			name := payment.Creator.FullName
			createdByName = &name
		}
		if payment.Updater != nil {
			name := payment.Updater.FullName
			updatedByName = &name
		}

		responses[i] = dto.PaymentResponse{
			ID:            payment.ID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			Notes:         payment.Notes,
			CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     payment.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     payment.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     payment.UpdatedBy,
			UpdatedByName: updatedByName,
		}
	}

	totalPages := (int(total) + perPage - 1) / perPage
	utils.Success(c, "Payments retrieved successfully", gin.H{
		"items": responses,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetPaymentPerformance godoc
// @Summary Get payment performance statistics
// @Description Get daily payment statistics for the last N days (qty and total amount)
// @Tags payments
// @Produce json
// @Param days query int false "Number of days to look back" default(7)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/payments/performance [get]
func (h *PaymentHandler) GetPaymentPerformance(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	branchID := c.GetUint("branch_id")

	// Parse days parameter, default to 7
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days < 1 {
		days = 7
	}

	performance, err := h.paymentService.GetPaymentPerformance(tenantID, branchID, days)
	if err != nil {
		utils.InternalError(c, "Failed to get payment performance")
		return
	}

	utils.Success(c, "Success", performance)
}
