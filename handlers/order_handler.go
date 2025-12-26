package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	BaseHandler
	orderService *services.OrderService
}

func NewOrderHandler(cfg *config.Config, orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		BaseHandler:  BaseHandler{config: cfg},
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")
	branchID := c.GetUint("branch_id")
	userID := c.GetUint("user_id")

	// Convert items to service format
	items := make([]struct {
		ProductID uint
		Quantity  int
	}, len(req.Items))

	for i, item := range req.Items {
		items[i].ProductID = item.ProductID
		items[i].Quantity = item.Quantity
	}

	order, err := h.orderService.CreateOrder(tenantID, branchID, userID, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response
	response := dto.OrderResponse{
		ID:          order.ID,
		TenantID:    order.TenantID,
		BranchID:    order.BranchID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Notes:       order.Notes,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Add order items
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			ProductSKU:  item.Product.SKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		}
	}
	response.OrderItems = orderItems

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	tenantID := c.GetUint("tenant_id")

	order, err := h.orderService.GetOrder(uint(orderID), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Build response
	response := dto.OrderResponse{
		ID:          order.ID,
		TenantID:    order.TenantID,
		BranchID:    order.BranchID,
		UserID:      order.UserID,
		OrderNumber: order.OrderNumber,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Notes:       order.Notes,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Add order items
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			ProductSKU:  item.Product.SKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		}
	}
	response.OrderItems = orderItems

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
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

	orders, total, err := h.orderService.ListOrders(tenantID, branchID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response
	responses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = dto.OrderResponse{
			ID:          order.ID,
			TenantID:    order.TenantID,
			BranchID:    order.BranchID,
			UserID:      order.UserID,
			OrderNumber: order.OrderNumber,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Notes:       order.Notes,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Add order items
		orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))
		for j, item := range order.OrderItems {
			orderItems[j] = dto.OrderItemResponse{
				ID:          item.ID,
				ProductID:   item.ProductID,
				ProductName: item.Product.Name,
				ProductSKU:  item.Product.SKU,
				Quantity:    item.Quantity,
				Price:       item.Price,
				Subtotal:    item.Subtotal,
			}
		}
		responses[i].OrderItems = orderItems
	}

	totalPages := (int(total) + perPage - 1) / perPage
	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
