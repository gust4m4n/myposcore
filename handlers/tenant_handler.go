package handlers

import (
	"fmt"
	"myposcore/config"
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

type TenantHandler struct {
	*BaseHandler
	tenantService *services.SuperAdminTenantService
}

func NewTenantHandler(cfg *config.Config) *TenantHandler {
	return &TenantHandler{
		BaseHandler:   NewBaseHandler(cfg),
		tenantService: services.NewSuperAdminTenantService(),
	}
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get list of all tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	// Get search parameter
	search := c.Query("search")

	tenants, total, err := h.tenantService.ListTenants(search, pagination.Page, pagination.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	var response []dto.TenantResponse
	for _, tenant := range tenants {
		var createdByName, updatedByName *string
		if tenant.Creator != nil {
			name := tenant.Creator.FullName
			createdByName = &name
		}
		if tenant.Updater != nil {
			name := tenant.Updater.FullName
			updatedByName = &name
		}

		response = append(response, dto.TenantResponse{
			ID:            tenant.ID,
			Name:          tenant.Name,
			Description:   tenant.Description,
			Address:       tenant.Address,
			Website:       tenant.Website,
			Email:         tenant.Email,
			Phone:         tenant.Phone,
			Image:         utils.GetFullImageURL(tenant.Image),
			IsActive:      tenant.IsActive,
			CreatedAt:     tenant.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     tenant.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     tenant.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"message":     "Tenants retrieved successfully",
		"page":        pagination.Page,
		"page_size":   pagination.PageSize,
		"total_items": total,
		"total_pages": (int(total) + pagination.PageSize - 1) / pagination.PageSize,
		"data":        response,
	})
}

// GetTenant godoc
// @Summary Get tenant by ID
// @Description Get tenant details by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 200 {object} dto.TenantResponse
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tenant ID")
		return
	}

	tenant, err := h.tenantService.GetTenantByID(uint(tenantID))
	if err != nil {
		utils.NotFound(c, "Tenant not found")
		return
	}

	var createdByName, updatedByName *string
	if tenant.Creator != nil {
		name := tenant.Creator.FullName
		createdByName = &name
	}
	if tenant.Updater != nil {
		name := tenant.Updater.FullName
		updatedByName = &name
	}

	utils.Success(c, "Operation successful", dto.TenantResponse{
		ID:            tenant.ID,
		Name:          tenant.Name,
		Description:   tenant.Description,
		Address:       tenant.Address,
		Website:       tenant.Website,
		Email:         tenant.Email,
		Phone:         tenant.Phone,
		Image:         utils.GetFullImageURL(tenant.Image),
		IsActive:      tenant.IsActive,
		CreatedAt:     tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     tenant.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     tenant.UpdatedBy,
		UpdatedByName: updatedByName,
	})
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant with optional image upload
// @Tags tenants
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Tenant name"
// @Param description formData string false "Description"
// @Param address formData string false "Address"
// @Param website formData string false "Website URL"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param is_active formData boolean true "Active status"
// @Param image formData file false "Tenant image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.TenantResponse
// @Router /tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	// Parse form data
	req := dto.CreateTenantRequest{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Address:     c.PostForm("address"),
		Website:     c.PostForm("website"),
		Email:       c.PostForm("email"),
		Phone:       c.PostForm("phone"),
		Active:      c.PostForm("is_active") == "true",
	}

	// Validate required fields
	if req.Name == "" {
		utils.BadRequest(c, "name is required")
		return
	}

	// Handle image upload
	var imageURL string
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
		uploadDir := "uploads/tenants"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.InternalError(c, "Failed to create upload directory")
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("tenant_%d_%d%s", time.Now().Unix(), time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			utils.InternalError(c, "Failed to save image")
			return
		}

		imageURL = fmt.Sprintf("/uploads/tenants/%s", filename)
	}

	// Get current user ID from context
	var createdBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		createdBy = &uid
	}

	tenant, err := h.tenantService.CreateTenant(req, imageURL, createdBy)
	if err != nil {
		// Delete uploaded image if tenant creation fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/tenants", filepath.Base(imageURL)))
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, "Tenant created successfully", dto.TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		Address:     tenant.Address,
		Website:     tenant.Website,
		Email:       tenant.Email,
		Phone:       tenant.Phone,
		Image:       utils.GetFullImageURL(tenant.Image),
		IsActive:    tenant.IsActive,
		CreatedAt:   tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   tenant.CreatedBy,
	})
}

// UpdateTenant godoc
// @Summary Update a tenant
// @Description Update an existing tenant with optional image upload
// @Tags tenants
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Tenant ID"
// @Param name formData string true "Tenant name"
// @Param description formData string false "Description"
// @Param address formData string false "Address"
// @Param website formData string false "Website URL"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param is_active formData boolean true "Active status"
// @Param image formData file false "Tenant image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.TenantResponse
// @Router /tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tenant ID")
		return
	}

	// Parse form data
	req := dto.UpdateTenantRequest{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Address:     c.PostForm("address"),
		Website:     c.PostForm("website"),
		Email:       c.PostForm("email"),
		Phone:       c.PostForm("phone"),
		Active:      c.PostForm("is_active") == "true",
	}

	// Validate required fields
	if req.Name == "" {
		utils.BadRequest(c, "name is required")
		return
	}

	// Handle image upload
	var imageURL string
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

		// Get existing tenant to delete old image
		existingTenant, err := h.tenantService.GetTenantByID(uint(tenantID))
		if err != nil {
			utils.NotFound(c, "Tenant not found")
			return
		}

		// Create uploads directory
		uploadDir := "uploads/tenants"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.InternalError(c, "Failed to create upload directory")
			return
		}

		// Delete old image if exists
		if existingTenant.Image != "" {
			oldPath := filepath.Join(uploadDir, filepath.Base(existingTenant.Image))
			os.Remove(oldPath) // Ignore error if file doesn't exist
		}

		// Generate unique filename
		filename := fmt.Sprintf("tenant_%d_%d%s", tenantID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			utils.InternalError(c, "Failed to save image")
			return
		}

		imageURL = fmt.Sprintf("/uploads/tenants/%s", filename)
	}

	// Get current user ID from context
	var updatedBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		updatedBy = &uid
	}

	tenant, err := h.tenantService.UpdateTenant(uint(tenantID), req, imageURL, updatedBy)
	if err != nil {
		// Delete uploaded image if update fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/tenants", filepath.Base(imageURL)))
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, "Tenant updated successfully", dto.TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		Address:     tenant.Address,
		Website:     tenant.Website,
		Email:       tenant.Email,
		Phone:       tenant.Phone,
		Image:       utils.GetFullImageURL(tenant.Image),
		IsActive:    tenant.IsActive,
		CreatedAt:   tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:   tenant.UpdatedBy,
	})
}

// DeleteTenant godoc
// @Summary Delete a tenant
// @Description Soft delete a tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 200 {object} map[string]interface{}
// @Router /tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tenant ID")
		return
	}

	// Get existing tenant to delete image file
	existingTenant, err := h.tenantService.GetTenantByID(uint(tenantID))
	if err != nil {
		utils.NotFound(c, "Tenant not found")
		return
	}

	// Get current user ID from context
	var deletedBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		deletedBy = &uid
	}

	// Delete tenant from database
	if err := h.tenantService.DeleteTenant(uint(tenantID), deletedBy); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Delete image file if exists
	if existingTenant.Image != "" {
		uploadDir := "uploads/tenants"
		imagePath := filepath.Join(uploadDir, filepath.Base(existingTenant.Image))
		os.Remove(imagePath) // Ignore error if file doesn't exist
	}

	utils.SuccessWithoutData(c, "Tenant deleted successfully")
}
