package handlers

import (
	"fmt"
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
// @Description Get list of all products for the tenant with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param search query string false "Search by name, description, or SKU"
// @Success 200 {object} []dto.ProductResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	// Get query parameters
	category := c.Query("category")
	search := c.Query("search")

	products, err := h.service.ListProducts(tenantID.(uint), category, search)
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
			Category:    product.Category,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
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
			Category:    product.Category,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
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
// @Success 200 {object} dto.ProductResponse
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfully",
		"data": dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
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
			Category:    product.Category,
			SKU:         product.SKU,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
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

// GetCategories godoc
// @Summary Get product categories
// @Description Get list of unique product categories for the tenant
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /api/v1/products/categories [get]
func (h *ProductHandler) GetCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categories, err := h.service.GetCategories(tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// UploadProductImage godoc
// @Summary Upload product image
// @Description Upload or update product image
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param image formData file true "Product image file"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/products/{id}/photo [post]
func (h *ProductHandler) UploadProductImage(c *gin.Context) {
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

	// Get product first to verify it exists and belongs to tenant
	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB"})
		return
	}

	// Create uploads directory if not exists
	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Delete old image if exists
	if product.Image != "" {
		oldPath := filepath.Join(uploadDir, filepath.Base(product.Image))
		os.Remove(oldPath) // Ignore error if file doesn't exist
	}

	// Generate unique filename
	filename := fmt.Sprintf("product_%d_%d%s", product.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}

	// Update product image URL
	imageURL := fmt.Sprintf("/uploads/products/%s", filename)
	updatedProduct, err := h.service.UpdateProductImage(uint(id), tenantID.(uint), imageURL)
	if err != nil {
		// Delete uploaded file if database update fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"data": dto.ProductResponse{
			ID:          updatedProduct.ID,
			TenantID:    updatedProduct.TenantID,
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Category:    updatedProduct.Category,
			SKU:         updatedProduct.SKU,
			Price:       updatedProduct.Price,
			Stock:       updatedProduct.Stock,
			Image:       updatedProduct.Image,
			IsActive:    updatedProduct.IsActive,
			CreatedAt:   updatedProduct.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   updatedProduct.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteProductImage godoc
// @Summary Delete product image
// @Description Delete product image
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/products/{id}/photo [delete]
func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
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

	// Get product first
	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Delete photo file if exists
	if product.Image != "" {
		uploadDir := "uploads/products"
		filePath := filepath.Join(uploadDir, filepath.Base(product.Image))
		os.Remove(filePath) // Ignore error if file doesn't exist
	}

	// Update database
	_, err = h.service.UpdateProductImage(uint(id), tenantID.(uint), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
