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

type SuperAdminHandler struct {
	*BaseHandler
	tenantService    *services.SuperAdminTenantService
	branchService    *services.SuperAdminBranchService
	userService      *services.SuperAdminUserService
	dashboardService *services.SuperAdminDashboardService
}

func NewSuperAdminHandler(cfg *config.Config) *SuperAdminHandler {
	return &SuperAdminHandler{
		BaseHandler:      NewBaseHandler(cfg),
		tenantService:    services.NewSuperAdminTenantService(),
		branchService:    services.NewSuperAdminBranchService(),
		userService:      services.NewSuperAdminUserService(),
		dashboardService: services.NewSuperAdminDashboardService(),
	}
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get list of all tenants (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /superadmin/tenants [get]
func (h *SuperAdminHandler) ListTenants(c *gin.Context) {
	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	tenants, total, err := h.tenantService.ListTenants(pagination.Page, pagination.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			Image:         tenant.Image,
			IsActive:      tenant.IsActive,
			CreatedAt:     tenant.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     tenant.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     tenant.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	paginatedResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total, response)
	c.JSON(http.StatusOK, paginatedResponse)
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant with optional image upload (superadmin only)
// @Tags superadmin
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
// @Router /superadmin/tenants [post]
func (h *SuperAdminHandler) CreateTenant(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB"})
			return
		}

		// Create uploads directory
		uploadDir := "uploads/tenants"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("tenant_%d_%d%s", time.Now().Unix(), time.Now().UnixNano(), ext)
		imageURL = fmt.Sprintf("/uploads/tenants/%s", filename)
	}

	// Get current user ID from context (superadmin)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant created successfully",
		"data": dto.TenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Description: tenant.Description,
			Address:     tenant.Address,
			Website:     tenant.Website,
			Email:       tenant.Email,
			Phone:       tenant.Phone,
			Image:       tenant.Image,
			IsActive:    tenant.IsActive,
			CreatedAt:   tenant.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:   tenant.CreatedBy,
		},
	})
}

// ListBranches godoc
// @Summary List branches by tenant
// @Description Get list of branches for a specific tenant (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /superadmin/tenants/{tenant_id}/branches [get]
func (h *SuperAdminHandler) ListBranches(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("tenant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	branches, total, err := h.branchService.ListBranches(uint(tenantID), pagination.Page, pagination.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.BranchResponse
	for _, branch := range branches {
		var createdByName, updatedByName *string
		if branch.Creator != nil {
			name := branch.Creator.FullName
			createdByName = &name
		}
		if branch.Updater != nil {
			name := branch.Updater.FullName
			updatedByName = &name
		}

		response = append(response, dto.BranchResponse{
			ID:            branch.ID,
			TenantID:      branch.TenantID,
			Name:          branch.Name,
			Description:   branch.Description,
			Address:       branch.Address,
			Website:       branch.Website,
			Email:         branch.Email,
			Phone:         branch.Phone,
			Image:         branch.Image,
			IsActive:      branch.IsActive,
			CreatedAt:     branch.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     branch.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     branch.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     branch.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	paginatedResponse := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, total, response)
	c.JSON(http.StatusOK, paginatedResponse)
}

// UpdateTenant godoc
// @Summary Update a tenant
// @Description Update an existing tenant with optional image upload (superadmin only)
// @Tags superadmin
// @Accept multipart/form-data
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param name formData string true "Tenant name"

// @Param description formData string false "Description"
// @Param address formData string false "Address"
// @Param website formData string false "Website URL"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param is_active formData boolean true "Active status"
// @Param image formData file false "Tenant image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.TenantResponse
// @Router /superadmin/tenants/{tenant_id} [put]
func (h *SuperAdminHandler) UpdateTenant(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("tenant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB"})
			return
		}

		// Get existing tenant to delete old image
		existingTenant, err := h.tenantService.GetTenantByID(uint(tenantID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
			return
		}

		// Create uploads directory
		uploadDir := "uploads/tenants"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant updated successfully",
		"data": dto.TenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Description: tenant.Description,
			Address:     tenant.Address,
			Website:     tenant.Website,
			Email:       tenant.Email,
			Phone:       tenant.Phone,
			Image:       tenant.Image,
			IsActive:    tenant.IsActive,
			CreatedAt:   tenant.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:   tenant.UpdatedBy,
		},
	})
}

// DeleteTenant godoc
// @Summary Delete a tenant
// @Description Soft delete a tenant (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Success 200 {object} map[string]interface{}
// @Router /superadmin/tenants/{tenant_id} [delete]
func (h *SuperAdminHandler) DeleteTenant(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("tenant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	// Get existing tenant to delete image file
	existingTenant, err := h.tenantService.GetTenantByID(uint(tenantID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete image file if exists
	if existingTenant.Image != "" {
		uploadDir := "uploads/tenants"
		imagePath := filepath.Join(uploadDir, filepath.Base(existingTenant.Image))
		os.Remove(imagePath) // Ignore error if file doesn't exist
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant deleted successfully",
	})
}

// CreateBranch godoc
// @Summary Create a new branch
// @Description Create a new branch for a tenant (superadmin only)
// @Tags superadmin
// @Accept multipart/form-data
// @Produce json
// @Param tenant_id formData integer true "Tenant ID"
// @Param name formData string true "Branch name"

// @Param description formData string false "Description"
// @Param address formData string false "Address"
// @Param website formData string false "Website URL"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param is_active formData boolean true "Active status"
// @Param image formData file false "Branch image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.BranchResponse
// @Router /superadmin/branches [post]
func (h *SuperAdminHandler) CreateBranch(c *gin.Context) {
	// Parse form data
	tenantID, err := strconv.ParseUint(c.PostForm("tenant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	req := dto.CreateBranchRequest{
		TenantID:    uint(tenantID),
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB"})
			return
		}

		// Create uploads directory
		uploadDir := "uploads/branches"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("branch_%d_%d%s", time.Now().Unix(), time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL = fmt.Sprintf("/uploads/branches/%s", filename)
	}

	// Get current user ID from context
	var createdBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		createdBy = &uid
	}

	branch, err := h.branchService.CreateBranch(req, imageURL, createdBy)
	if err != nil {
		// Delete uploaded image if creation fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/branches", filepath.Base(imageURL)))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branch created successfully",
		"data": dto.BranchResponse{
			ID:          branch.ID,
			TenantID:    branch.TenantID,
			Name:        branch.Name,
			Description: branch.Description,
			Address:     branch.Address,
			Website:     branch.Website,
			Email:       branch.Email,
			Phone:       branch.Phone,
			Image:       branch.Image,
			IsActive:    branch.IsActive,
			CreatedAt:   branch.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   branch.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:   branch.CreatedBy,
		},
	})
}

// UpdateBranch godoc
// @Summary Update a branch
// @Description Update an existing branch (superadmin only)
// @Tags superadmin
// @Accept multipart/form-data
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param name formData string true "Branch name"

// @Param description formData string false "Description"
// @Param address formData string false "Address"
// @Param website formData string false "Website URL"
// @Param email formData string false "Email"
// @Param phone formData string false "Phone"
// @Param is_active formData boolean true "Active status"
// @Param image formData file false "Branch image (jpg, jpeg, png, gif, webp, max 5MB)"
// @Success 200 {object} dto.BranchResponse
// @Router /superadmin/branches/{branch_id} [put]
func (h *SuperAdminHandler) UpdateBranch(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	// Parse form data
	req := dto.UpdateBranchRequest{
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp"})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB"})
			return
		}

		// Get existing branch to delete old image
		existingBranch, err := h.branchService.GetBranchByID(uint(branchID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}

		// Create uploads directory
		uploadDir := "uploads/branches"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Delete old image if exists
		if existingBranch.Image != "" {
			oldPath := filepath.Join(uploadDir, filepath.Base(existingBranch.Image))
			os.Remove(oldPath) // Ignore error if file doesn't exist
		}

		// Generate unique filename
		filename := fmt.Sprintf("branch_%d_%d%s", branchID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL = fmt.Sprintf("/uploads/branches/%s", filename)
	}

	// Get current user ID from context
	var updatedBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		updatedBy = &uid
	}

	branch, err := h.branchService.UpdateBranch(uint(branchID), req, imageURL, updatedBy)
	if err != nil {
		// Delete uploaded image if update fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/branches", filepath.Base(imageURL)))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branch updated successfully",
		"data": dto.BranchResponse{
			ID:          branch.ID,
			TenantID:    branch.TenantID,
			Name:        branch.Name,
			Description: branch.Description,
			Address:     branch.Address,
			Website:     branch.Website,
			Email:       branch.Email,
			Phone:       branch.Phone,
			Image:       branch.Image,
			IsActive:    branch.IsActive,
			CreatedAt:   branch.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   branch.UpdatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:   branch.UpdatedBy,
		},
	})
}

// DeleteBranch godoc
// @Summary Delete a branch
// @Description Soft delete a branch (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Router /superadmin/branches/{branch_id} [delete]
func (h *SuperAdminHandler) DeleteBranch(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	// Get existing branch to delete image file
	existingBranch, err := h.branchService.GetBranchByID(uint(branchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
		return
	}

	// Get current user ID from context
	var deletedBy *uint
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		deletedBy = &uid
	}

	// Delete branch from database
	if err := h.branchService.DeleteBranch(uint(branchID), deletedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete image file if exists
	if existingBranch.Image != "" {
		uploadDir := "uploads/branches"
		imagePath := filepath.Join(uploadDir, filepath.Base(existingBranch.Image))
		os.Remove(imagePath) // Ignore error if file doesn't exist
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branch deleted successfully",
	})
}

// ListUsers godoc
// @Summary List users by branch
// @Description Get list of users for a specific branch (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {object} []dto.UserResponse
// @Router /superadmin/branches/{branch_id}/users [get]
func (h *SuperAdminHandler) ListUsers(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	users, err := h.userService.ListUsersByBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:        user.ID,
			TenantID:  user.TenantID,
			BranchID:  user.BranchID,
			Email:     user.Email,
			FullName:  user.FullName,
			Image:     user.Image,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// Dashboard godoc
// Dashboard godoc
// @Summary Get dashboard statistics
// @Description Get comprehensive dashboard with tenant list, counts of branches, users, products, and orders (all time, today, this week, this month)
// @Tags superadmin
// @Accept json
// @Produce json
// @Success 200 {object} dto.DashboardResponse
// @Router /superadmin/dashboard [get]
func (h *SuperAdminHandler) Dashboard(c *gin.Context) {
	dashboard, err := h.dashboardService.GetDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dashboard,
	})
}
