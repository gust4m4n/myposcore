package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

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
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(tenantID.(uint), req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category created successfully",
		"data": dto.CategoryResponse{
			ID:          category.ID,
			TenantID:    category.TenantID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategory(uint(categoryID), tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.CategoryResponse{
			ID:          category.ID,
			TenantID:    category.TenantID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// ListCategories godoc
// @Summary Get all categories
// @Description Get list of all categories for the tenant
// @Tags categories
// @Produce json
// @Param active_only query bool false "Show only active categories"
// @Success 200 {array} dto.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	activeOnly := c.Query("active_only") == "true"

	categories, err := h.categoryService.ListCategories(tenantID.(uint), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = dto.CategoryResponse{
			ID:          category.ID,
			TenantID:    category.TenantID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} dto.CategoryResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var name, description *string
	if req.Name != "" {
		name = &req.Name
	}
	if req.Description != "" {
		description = &req.Description
	}

	category, err := h.categoryService.UpdateCategory(uint(categoryID), tenantID.(uint), name, description, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data": dto.CategoryResponse{
			ID:          category.ID,
			TenantID:    category.TenantID,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category (only if not used by any products)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant not found"})
		return
	}

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(uint(categoryID), tenantID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
