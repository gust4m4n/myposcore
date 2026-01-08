package handlers

import (
	"fmt"
	"mime/multipart"
	"myposcore/config"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/services"
	"myposcore/utils"
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

func NewProductHandler(cfg *config.Config, productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		BaseHandler: NewBaseHandler(cfg),
		service:     productService,
	}
}

// Helper function to map category to DTO
func mapCategoryToDTO(category *models.Category) *dto.CategorySummary {
	if category == nil {
		return nil
	}
	return &dto.CategorySummary{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Image:       utils.GetFullImageURL(category.Image),
	}
}

// ListProducts godoc
// @Summary List all products
// @Description Get list of all products for the tenant with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Param search query string false "Search by name, description, or SKU"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	// Get query parameters
	search := c.Query("search")

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	products, total, err := h.service.ListProducts(tenantID.(uint), search, pagination.Page, pagination.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	var response []dto.ProductResponse
	for _, product := range products {
		var createdByName, updatedByName *string
		if product.Creator != nil {
			name := product.Creator.FullName
			createdByName = &name
		}
		if product.Updater != nil {
			name := product.Updater.FullName
			updatedByName = &name
		}

		response = append(response, dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,

			CategoryID:     product.CategoryID,
			CategoryDetail: mapCategoryToDTO(product.CategoryDetail),
			SKU:            product.SKU,
			Price:          product.Price,
			Stock:          product.Stock,
			Image:          utils.GetFullImageURL(product.Image),
			IsActive:       product.IsActive,
			CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:      product.CreatedBy,
			CreatedByName:  createdByName,
			UpdatedBy:      product.UpdatedBy,
			UpdatedByName:  updatedByName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"message":     "Products retrieved successfully",
		"page":        pagination.Page,
		"page_size":   pagination.PageSize,
		"total_items": total,
		"total_pages": (int(total) + pagination.PageSize - 1) / pagination.PageSize,
		"data":        response,
	})
}

// ListProductsByCategoryID godoc
// @Summary List products by category ID
// @Description Get list of products filtered by category ID with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param search query string false "Search by name or SKU"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.PaginationResponse
// @Router /api/v1/products/by-category/{category_id} [get]
func (h *ProductHandler) ListProductsByCategoryID(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	// Get category_id from path
	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	// Get query parameters
	search := c.Query("search")

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err == nil {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	products, total, err := h.service.ListProductsByCategoryID(tenantID.(uint), uint(categoryID), search, pagination.Page, pagination.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	var response []dto.ProductResponse
	for _, product := range products {
		var createdByName, updatedByName *string
		if product.Creator != nil {
			name := product.Creator.FullName
			createdByName = &name
		}
		if product.Updater != nil {
			name := product.Updater.FullName
			updatedByName = &name
		}

		response = append(response, dto.ProductResponse{
			ID:          product.ID,
			TenantID:    product.TenantID,
			Name:        product.Name,
			Description: product.Description,

			CategoryID:     product.CategoryID,
			CategoryDetail: mapCategoryToDTO(product.CategoryDetail),
			SKU:            product.SKU,
			Price:          product.Price,
			Stock:          product.Stock,
			Image:          utils.GetFullImageURL(product.Image),
			IsActive:       product.IsActive,
			CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:      product.CreatedBy,
			CreatedByName:  createdByName,
			UpdatedBy:      product.UpdatedBy,
			UpdatedByName:  updatedByName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"message":     "Products retrieved successfully",
		"page":        pagination.Page,
		"page_size":   pagination.PageSize,
		"total_items": total,
		"total_pages": (int(total) + pagination.PageSize - 1) / pagination.PageSize,
		"data":        response,
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
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID")
		return
	}

	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	var createdByName, updatedByName *string
	if product.Creator != nil {
		name := product.Creator.FullName
		createdByName = &name
	}
	if product.Updater != nil {
		name := product.Updater.FullName
		updatedByName = &name
	}

	utils.Success(c, "Product retrieved successfully", dto.ProductResponse{
		ID:             product.ID,
		TenantID:       product.TenantID,
		Name:           product.Name,
		Description:    product.Description,
		CategoryID:     product.CategoryID,
		CategoryDetail: mapCategoryToDTO(product.CategoryDetail),
		SKU:            product.SKU,
		Price:          product.Price,
		Stock:          product.Stock,
		Image:          utils.GetFullImageURL(product.Image),
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      product.CreatedBy,
		CreatedByName:  createdByName,
		UpdatedBy:      product.UpdatedBy,
		UpdatedByName:  updatedByName,
	})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for the tenant (supports both JSON and multipart/form-data with optional image)
// @Tags products
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param request body dto.CreateProductRequest true "Product data (JSON)" (when using application/json)
// @Param name formData string true "Product name" (when using multipart/form-data)
// @Param description formData string false "Product description" (when using multipart/form-data)
// @Param category_id formData integer false "Product category ID" (when using multipart/form-data)
// @Param sku formData string false "Product SKU" (when using multipart/form-data)
// @Param price formData number true "Product price" (when using multipart/form-data)
// @Param stock formData integer false "Product stock" (when using multipart/form-data)
// @Param is_active formData boolean false "Is product active" (when using multipart/form-data)
// @Param image formData file false "Product image file (optional)" (when using multipart/form-data)
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	var req dto.CreateProductRequest
	contentType := c.GetHeader("Content-Type")

	// Check if request is multipart form-data
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Name = c.PostForm("name")
		req.Description = c.PostForm("description")
		req.SKU = c.PostForm("sku")

		// Parse category_id
		if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
			categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err == nil {
				catID := uint(categoryID)
				req.CategoryID = &catID
			}
		}

		// Parse price
		if priceStr := c.PostForm("price"); priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				utils.BadRequest(c, "Invalid price format")
				return
			}
			req.Price = price
		}

		// Parse stock
		if stockStr := c.PostForm("stock"); stockStr != "" {
			stock, err := strconv.Atoi(stockStr)
			if err != nil {
				utils.BadRequest(c, "Invalid stock format")
				return
			}
			req.Stock = stock
		}

		// Parse is_active
		if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
			req.IsActive = isActiveStr == "true" || isActiveStr == "1"
		}

		// Validate required fields
		if req.Name == "" {
			utils.BadRequest(c, "Name is required")
			return
		}
		if req.Price <= 0 {
			utils.BadRequest(c, "Price must be greater than 0")
			return
		}
	} else {
		// Parse JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}

	// Set created_by to current user
	currentUserID := c.GetUint("user_id")
	req.CreatedBy = &currentUserID

	product, err := h.service.CreateProduct(tenantID.(uint), req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Handle image upload if provided (only for multipart form-data)
	if strings.Contains(contentType, "multipart/form-data") {
		if file, err := c.FormFile("image"); err == nil {
			imageURL, uploadErr := h.handleImageUpload(file, product.ID)
			if uploadErr != nil {
				utils.BadRequest(c, fmt.Sprintf("Product created but image upload failed: %s", uploadErr.Error()))
				return
			}
			// Update product with image URL
			product, _ = h.service.UpdateProductImage(product.ID, tenantID.(uint), imageURL)
		}
	}

	utils.Success(c, "Product created successfully", dto.ProductResponse{
		ID:             product.ID,
		TenantID:       product.TenantID,
		Name:           product.Name,
		Description:    product.Description,
		CategoryID:     product.CategoryID,
		CategoryDetail: mapCategoryToDTO(product.CategoryDetail),
		SKU:            product.SKU,
		Price:          product.Price,
		Stock:          product.Stock,
		Image:          utils.GetFullImageURL(product.Image),
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      product.CreatedBy,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product details (supports both JSON and multipart/form-data with optional image)
// @Tags products
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product data (JSON)" (when using application/json)
// @Param name formData string false "Product name" (when using multipart/form-data)
// @Param description formData string false "Product description" (when using multipart/form-data)
// @Param category_id formData integer false "Product category ID" (when using multipart/form-data)
// @Param sku formData string false "Product SKU" (when using multipart/form-data)
// @Param price formData number false "Product price" (when using multipart/form-data)
// @Param stock formData integer false "Product stock" (when using multipart/form-data)
// @Param is_active formData boolean false "Is product active" (when using multipart/form-data)
// @Param image formData file false "Product image file (optional)" (when using multipart/form-data)
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID")
		return
	}

	var req dto.UpdateProductRequest
	contentType := c.GetHeader("Content-Type")

	// Check if request is multipart form-data
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if name := c.PostForm("name"); name != "" {
			req.Name = name
		}
		if description := c.PostForm("description"); description != "" {
			req.Description = description
		}
		if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
			categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
			if err == nil {
				catID := uint(categoryID)
				req.CategoryID = &catID
			}
		}
		if sku := c.PostForm("sku"); sku != "" {
			req.SKU = sku
		}

		// Parse price
		if priceStr := c.PostForm("price"); priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				utils.BadRequest(c, "Invalid price format")
				return
			}
			req.Price = price
		}

		// Parse stock
		if stockStr := c.PostForm("stock"); stockStr != "" {
			stock, err := strconv.Atoi(stockStr)
			if err != nil {
				utils.BadRequest(c, "Invalid stock format")
				return
			}
			req.Stock = stock
		}

		// Parse is_active
		if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true" || isActiveStr == "1"
			req.IsActive = &isActive
		}
	} else {
		// Parse JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}

	// Set updated_by to current user
	currentUserID := c.GetUint("user_id")
	req.UpdatedBy = &currentUserID

	product, err := h.service.UpdateProduct(uint(id), tenantID.(uint), req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Handle image upload if provided (only for multipart form-data)
	if strings.Contains(contentType, "multipart/form-data") {
		if file, err := c.FormFile("image"); err == nil {
			// Delete old image if exists
			if product.Image != "" {
				uploadDir := "uploads/products"
				oldPath := filepath.Join(uploadDir, filepath.Base(product.Image))
				os.Remove(oldPath) // Ignore error if file doesn't exist
			}

			imageURL, uploadErr := h.handleImageUpload(file, product.ID)
			if uploadErr != nil {
				utils.BadRequest(c, fmt.Sprintf("Product updated but image upload failed: %s", uploadErr.Error()))
				return
			}
			// Update product with image URL
			product, _ = h.service.UpdateProductImage(product.ID, tenantID.(uint), imageURL)
		}
	}

	utils.Success(c, "Product updated successfully", dto.ProductResponse{
		ID:             product.ID,
		TenantID:       product.TenantID,
		Name:           product.Name,
		Description:    product.Description,
		CategoryID:     product.CategoryID,
		CategoryDetail: mapCategoryToDTO(product.CategoryDetail),
		SKU:            product.SKU,
		Price:          product.Price,
		Stock:          product.Stock,
		Image:          utils.GetFullImageURL(product.Image),
		IsActive:       product.IsActive,
		CreatedAt:      product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      product.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      product.CreatedBy,
		UpdatedBy:      product.UpdatedBy,
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
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID")
		return
	}

	// Get current user ID
	currentUserID := c.GetUint("user_id")

	if err := h.service.DeleteProduct(uint(id), tenantID.(uint), &currentUserID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "Product deleted successfully")
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
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	categories, err := h.service.GetCategories(tenantID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, "Categories retrieved successfully", categories)
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
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID")
		return
	}

	// Get product first to verify it exists and belongs to tenant
	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		utils.NotFound(c, "Product not found")
		return
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "Image file is required")
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
		utils.BadRequest(c, "Invalid file type. Allowed: jpg, jpeg, png, gif, webp")
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		utils.BadRequest(c, "File size too large. Maximum 5MB")
		return
	}

	// Create uploads directory if not exists
	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.InternalError(c, "Failed to create upload directory")
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
		utils.InternalError(c, "Failed to save photo")
		return
	}

	// Update product image URL
	imageURL := fmt.Sprintf("/uploads/products/%s", filename)
	updatedProduct, err := h.service.UpdateProductImage(uint(id), tenantID.(uint), imageURL)
	if err != nil {
		// Delete uploaded file if database update fails
		os.Remove(filePath)
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, "Image uploaded successfully", dto.ProductResponse{
		ID:             updatedProduct.ID,
		TenantID:       updatedProduct.TenantID,
		Name:           updatedProduct.Name,
		Description:    updatedProduct.Description,
		CategoryID:     updatedProduct.CategoryID,
		CategoryDetail: mapCategoryToDTO(updatedProduct.CategoryDetail),
		SKU:            updatedProduct.SKU,
		Price:          updatedProduct.Price,
		Stock:          updatedProduct.Stock,
		Image:          utils.GetFullImageURL(updatedProduct.Image),
		IsActive:       updatedProduct.IsActive,
		CreatedAt:      updatedProduct.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      updatedProduct.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      updatedProduct.CreatedBy,
	})
}

// handleImageUpload is a helper method to handle image upload for products
func (h *ProductHandler) handleImageUpload(file *multipart.FileHeader, productID uint) (string, error) {
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
		return "", fmt.Errorf("invalid file type. Allowed: jpg, jpeg, png, gif, webp")
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file size too large. Maximum 5MB")
	}

	// Create uploads directory if not exists
	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	// Generate unique filename
	filename := fmt.Sprintf("product_%d_%d%s", productID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file")
	}
	defer dst.Close()

	// Copy file content
	if _, err = dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save file")
	}

	// Return image URL
	imageURL := fmt.Sprintf("/uploads/products/%s", filename)
	return imageURL, nil
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
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID")
		return
	}

	// Get product first
	product, err := h.service.GetProduct(uint(id), tenantID.(uint))
	if err != nil {
		utils.NotFound(c, "Product not found")
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
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "Image deleted successfully")
}
