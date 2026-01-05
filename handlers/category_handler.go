package handlers

import (
	"fmt"
	"myposcore/dto"
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

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category for the tenant
// @Tags categories
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string false "Category description"
// @Param image formData file false "Category image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	var req dto.CreateCategoryRequest
	contentType := c.GetHeader("Content-Type")

	// Check if request is multipart form-data
	if strings.Contains(contentType, "multipart/form-data") {
		req.Name = c.PostForm("name")
		req.Description = c.PostForm("description")

		if req.Name == "" {
			utils.BadRequest(c, "Name is required")
			return
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")
	req.CreatedBy = &currentUserID

	// Handle image upload if provided
	var imageURL string
	if strings.Contains(contentType, "multipart/form-data") {
		file, err := c.FormFile("image")
		if err == nil {
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

			// Create uploads directory
			uploadDir := "uploads/categories"
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				utils.InternalError(c, "Failed to create upload directory")
				return
			}

			// Generate unique filename
			filename := fmt.Sprintf("category_%s_%d%s", strings.ReplaceAll(req.Name, " ", "_"), time.Now().Unix(), ext)
			filePath := filepath.Join(uploadDir, filename)

			// Save file
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				utils.InternalError(c, "Failed to save image")
				return
			}

			imageURL = fmt.Sprintf("/uploads/categories/%s", filename)
		}
	}

	category, err := h.categoryService.CreateCategory(tenantID.(uint), req.Name, req.Description, imageURL, req.CreatedBy)
	if err != nil {
		// Delete uploaded image if creation fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/categories", filepath.Base(imageURL)))
		}
		utils.BadRequest(c, err.Error())
		return
	}

	// Reload to get audit info
	category, _ = h.categoryService.GetCategory(category.ID, tenantID.(uint))

	var createdByName *string
	if category.Creator != nil {
		name := category.Creator.FullName
		createdByName = &name
	}

	utils.Success(c, "Category created successfully", dto.CategoryResponse{
		ID:            category.ID,
		TenantID:      category.TenantID,
		Name:          category.Name,
		Description:   category.Description,
		Image:         category.Image,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     category.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     category.CreatedBy,
		CreatedByName: createdByName,
	})
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Get a specific category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategory(uint(categoryID), tenantID.(uint))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	var createdByName, updatedByName *string
	if category.Creator != nil {
		name := category.Creator.FullName
		createdByName = &name
	}
	if category.Updater != nil {
		name := category.Updater.FullName
		updatedByName = &name
	}

	utils.Success(c, "Category retrieved successfully", dto.CategoryResponse{
		ID:            category.ID,
		TenantID:      category.TenantID,
		Name:          category.Name,
		Description:   category.Description,
		Image:         category.Image,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     category.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     category.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     category.UpdatedBy,
		UpdatedByName: updatedByName,
	})
}

// ListCategories godoc
// @Summary Get all categories
// @Description Get list of all categories for the tenant
// @Tags categories
// @Produce json
// @Param active_only query bool false "Show only active categories"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	activeOnly := c.Query("active_only") == "true"
	// Get search parameter (optional)
	search := c.Query("search")

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	categories, total, err := h.categoryService.ListCategories(tenantID.(uint), search, activeOnly, pagination.Page, pagination.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		var createdByName, updatedByName *string
		if category.Creator != nil {
			name := category.Creator.FullName
			createdByName = &name
		}
		if category.Updater != nil {
			name := category.Updater.FullName
			updatedByName = &name
		}

		responses[i] = dto.CategoryResponse{
			ID:            category.ID,
			TenantID:      category.TenantID,
			Name:          category.Name,
			Description:   category.Description,
			Image:         category.Image,
			IsActive:      category.IsActive,
			CreatedAt:     category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     category.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     category.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     category.UpdatedBy,
			UpdatedByName: updatedByName,
		}
	}

	paginatedResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total, responses)
	c.JSON(http.StatusOK, paginatedResponse)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update an existing category (supports multipart/form-data for image upload or application/json)
// @Tags categories
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string false "Category name"
// @Param description formData string false "Category description"
// @Param is_active formData boolean false "Category status"
// @Param image formData file false "Category image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")

	var name, description *string
	var isActive *bool
	var imageURL *string

	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data
		if nameVal := c.PostForm("name"); nameVal != "" {
			name = &nameVal
		}
		if descVal := c.PostForm("description"); descVal != "" {
			description = &descVal
		}
		if isActiveVal := c.PostForm("is_active"); isActiveVal != "" {
			isActiveBool := isActiveVal == "true" || isActiveVal == "1"
			isActive = &isActiveBool
		}

		// Handle image upload
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			// Validate file type
			allowedExtensions := map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
				".gif":  true,
				".webp": true,
			}
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if !allowedExtensions[ext] {
				utils.BadRequest(c, "Invalid file type. Allowed: jpg, jpeg, png, gif, webp")
				return
			}

			// Validate file size (max 5MB)
			if file.Size > 5*1024*1024 {
				utils.BadRequest(c, "File size exceeds 5MB limit")
				return
			}

			// Get existing category to retrieve old image path
			existingCategory, err := h.categoryService.GetCategory(uint(categoryID), tenantID.(uint))
			if err != nil {
				utils.BadRequest(c, err.Error())
				return
			}

			// Create upload directory
			uploadDir := "uploads/categories"
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				utils.InternalError(c, "Failed to create upload directory")
				return
			}

			// Generate unique filename
			timestamp := time.Now().UnixNano() / 1e6
			filename := fmt.Sprintf("category_%s_%d%s",
				strings.ReplaceAll(strings.ToLower(existingCategory.Name), " ", "_"),
				timestamp,
				ext,
			)
			filepath := filepath.Join(uploadDir, filename)

			// Save the file
			if err := c.SaveUploadedFile(file, filepath); err != nil {
				utils.InternalError(c, "Failed to save image")
				return
			}

			imageURL = &filepath

			// Delete old image if exists
			if existingCategory.Image != "" {
				if err := os.Remove(existingCategory.Image); err != nil && !os.IsNotExist(err) {
					// Log error but continue (file might not exist)
					fmt.Printf("Warning: Failed to delete old image: %v\n", err)
				}
			}
		}
	} else {
		// Handle JSON request (backward compatibility)
		var req dto.UpdateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}

		if req.Name != "" {
			name = &req.Name
		}
		if req.Description != "" {
			description = &req.Description
		}
		isActive = req.IsActive
	}

	category, err := h.categoryService.UpdateCategory(uint(categoryID), tenantID.(uint), name, description, imageURL, isActive, &currentUserID)
	if err != nil {
		// Rollback: delete uploaded image if database update fails
		if imageURL != nil && *imageURL != "" {
			os.Remove(*imageURL)
		}
		utils.BadRequest(c, err.Error())
		return
	}

	var createdByName, updatedByName *string
	if category.Creator != nil {
		name := category.Creator.FullName
		createdByName = &name
	}
	if category.Updater != nil {
		name := category.Updater.FullName
		updatedByName = &name
	}

	utils.Success(c, "Category updated successfully", dto.CategoryResponse{
		ID:            category.ID,
		TenantID:      category.TenantID,
		Name:          category.Name,
		Description:   category.Description,
		Image:         category.Image,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     category.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     category.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     category.UpdatedBy,
		UpdatedByName: updatedByName,
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category (only if not used by any products). Also deletes associated image file.
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant not found")
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")

	// Get category to retrieve image path before deletion
	category, err := h.categoryService.GetCategory(uint(categoryID), tenantID.(uint))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Delete category from database
	if err := h.categoryService.DeleteCategory(uint(categoryID), tenantID.(uint), &currentUserID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Delete image file if exists
	if category.Image != "" {
		if err := os.Remove(category.Image); err != nil && !os.IsNotExist(err) {
			// Log error but don't fail the request (image might not exist)
			fmt.Printf("Warning: Failed to delete category image: %v\n", err)
		}
	}

	utils.SuccessWithoutData(c, "Category deleted successfully")
}
