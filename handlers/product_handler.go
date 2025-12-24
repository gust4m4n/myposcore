package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	*BaseHandler
	service *services.ProductService
}

func NewProductHandler(cfg *config.Config) *ProductHandler {
	return &ProductHandler{
		BaseHandler: NewBaseHandler(cfg),
		service:     services.NewProductService(),
	}
}

// ListProducts godoc
// @Summary List all products
// @Description Get list of all products for the tenant
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ProductResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	products, err := h.service.ListProducts(tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for the tenant
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(tenantID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data": dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product data"
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(uint(id), tenantID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data": dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			IsActive:    product.IsActive,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.service.DeleteProduct(uint(id), tenantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
