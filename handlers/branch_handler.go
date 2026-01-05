package handlers

import (
	"fmt"
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	BaseHandler
	cfg           *config.Config
	branchService *services.SuperAdminBranchService
}

func NewBranchHandler(cfg *config.Config, branchService *services.SuperAdminBranchService) *BranchHandler {
	return &BranchHandler{
		cfg:           cfg,
		branchService: branchService,
	}
}

// GetBranches godoc
// @Summary Get branches for current user's tenant
// @Description Get list of branches from the tenant where the logged-in user is registered
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches [get]
func (h *BranchHandler) GetBranches(c *gin.Context) {
	// Get tenant_id from JWT context
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	// Get search parameter (optional)
	search := c.Query("search")

	// Get branches for this tenant (without pagination, get all)
	branches, _, err := h.branchService.ListBranches(tenantID.(uint), search, 1, 9999)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve branches: "+err.Error())
		return
	}

	utils.Success(c, "Branches retrieved successfully", branches)
}

// GetBranch godoc
// @Summary Get branch detail
// @Description Get detail of a specific branch in user's tenant
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches/{id} [get]
func (h *BranchHandler) GetBranch(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	branchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid branch ID")
		return
	}

	branch, err := h.branchService.GetBranchByID(uint(branchID))
	if err != nil {
		utils.NotFound(c, "Branch not found")
		return
	}

	// Verify branch belongs to user's tenant
	if branch.TenantID != tenantID.(uint) {
		utils.Forbidden(c, "Access denied to this branch")
		return
	}

	utils.Success(c, "Branch retrieved successfully", branch)
}

// CreateBranch godoc
// @Summary Create new branch
// @Description Create a new branch in user's tenant
// @Tags Branch
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Branch name"
// @Param description formData string false "Branch description"
// @Param address formData string false "Branch address"
// @Param website formData string false "Branch website"
// @Param email formData string false "Branch email"
// @Param phone formData string false "Branch phone"
// @Param is_active formData boolean false "Is active"
// @Param image formData file false "Branch image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches [post]
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User ID not found in token")
		return
	}

	// Parse form data
	req := dto.CreateBranchRequest{
		TenantID:    tenantID.(uint),
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
		uploadDir := "uploads/branches"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.InternalError(c, "Failed to create upload directory")
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("branch_%d_%d%s", time.Now().Unix(), time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			utils.InternalError(c, "Failed to save image")
			return
		}

		imageURL = fmt.Sprintf("/uploads/branches/%s", filename)
	}

	// Convert userID to *uint
	uid := userID.(uint)

	// Create branch using service
	branch, err := h.branchService.CreateBranch(req, imageURL, &uid)
	if err != nil {
		// Delete uploaded image if creation fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/branches", filepath.Base(imageURL)))
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, "Branch created successfully", branch)
}

// UpdateBranch godoc
// @Summary Update branch
// @Description Update branch data in user's tenant
// @Tags Branch
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Branch ID"
// @Param name formData string false "Branch name"
// @Param description formData string false "Branch description"
// @Param address formData string false "Branch address"
// @Param website formData string false "Branch website"
// @Param email formData string false "Branch email"
// @Param phone formData string false "Branch phone"
// @Param is_active formData boolean false "Is active"
// @Param image formData file false "Branch image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches/{id} [put]
func (h *BranchHandler) UpdateBranch(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User ID not found in token")
		return
	}

	branchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid branch ID")
		return
	}

	// Get existing branch
	branch, err := h.branchService.GetBranchByID(uint(branchID))
	if err != nil {
		utils.NotFound(c, "Branch not found")
		return
	}

	// Verify branch belongs to user's tenant
	if branch.TenantID != tenantID.(uint) {
		utils.Forbidden(c, "Access denied to this branch")
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
		uploadDir := "uploads/branches"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.InternalError(c, "Failed to create upload directory")
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("branch_%d_%d%s", uint(branchID), time.Now().Unix(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			utils.InternalError(c, "Failed to save image")
			return
		}

		// Delete old image if exists
		if branch.Image != "" {
			oldPath := filepath.Join("uploads/branches", filepath.Base(branch.Image))
			os.Remove(oldPath)
		}

		imageURL = fmt.Sprintf("/uploads/branches/%s", filename)
	}

	// Convert userID to *uint
	uid := userID.(uint)

	// Update branch using service
	updatedBranch, err := h.branchService.UpdateBranch(uint(branchID), req, imageURL, &uid)
	if err != nil {
		// Delete uploaded image if update fails
		if imageURL != "" {
			os.Remove(filepath.Join("uploads/branches", filepath.Base(imageURL)))
		}
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, "Branch updated successfully", updatedBranch)
}

// DeleteBranch godoc
// @Summary Delete branch
// @Description Soft delete a branch in user's tenant
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches/{id} [delete]
func (h *BranchHandler) DeleteBranch(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User ID not found in token")
		return
	}

	branchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid branch ID")
		return
	}

	// Get existing branch
	branch, err := h.branchService.GetBranchByID(uint(branchID))
	if err != nil {
		utils.NotFound(c, "Branch not found")
		return
	}

	// Verify branch belongs to user's tenant
	if branch.TenantID != tenantID.(uint) {
		utils.Forbidden(c, "Access denied to this branch")
		return
	}

	// Convert userID to *uint
	uid := userID.(uint)

	// Delete branch using service
	if err := h.branchService.DeleteBranch(uint(branchID), &uid); err != nil {
		utils.InternalError(c, "Failed to delete branch: "+err.Error())
		return
	}

	utils.Success(c, "Branch deleted successfully", nil)
}

// GetBranchUsers godoc
// @Summary Get users in branch
// @Description Get list of users in a specific branch
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Branch ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/branches/{id}/users [get]
func (h *BranchHandler) GetBranchUsers(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		utils.Unauthorized(c, "Tenant ID not found in token")
		return
	}

	branchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid branch ID")
		return
	}

	// Get existing branch
	branch, err := h.branchService.GetBranchByID(uint(branchID))
	if err != nil {
		utils.NotFound(c, "Branch not found")
		return
	}

	// Verify branch belongs to user's tenant
	if branch.TenantID != tenantID.(uint) {
		utils.Forbidden(c, "Access denied to this branch")
		return
	}

	// For now, return empty array since service doesn't have this method yet
	// TODO: Implement GetBranchUsers in service
	utils.Success(c, "Users retrieved successfully", gin.H{
		"users": []interface{}{},
	})
}
