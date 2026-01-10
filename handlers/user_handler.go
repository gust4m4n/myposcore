package handlers

import (
	"fmt"
	"mime/multipart"
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

type UserHandler struct {
	*BaseHandler
	userService *services.UserService
}

func NewUserHandler(cfg *config.Config, userService *services.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(cfg),
		userService: userService,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Get list of all users for the authenticated tenant
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(32)
// @Success 200 {object} dto.PaginationResponse
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")

	// Get search parameter (optional)
	search := c.Query("search")

	// Parse pagination parameters
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = *dto.NewPaginationRequest(1, 32)
	} else {
		pagination = *dto.NewPaginationRequest(pagination.Page, pagination.PageSize)
	}

	users, total, err := h.userService.ListUsers(tenantID, search, pagination.Page, pagination.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	var response []dto.UserResponse
	for _, user := range users {
		var createdByName, updatedByName *string
		if user.Creator != nil {
			name := user.Creator.FullName
			createdByName = &name
		}
		if user.Updater != nil {
			name := user.Updater.FullName
			updatedByName = &name
		}

		response = append(response, dto.UserResponse{
			ID:            user.ID,
			TenantID:      user.TenantID,
			BranchID:      user.BranchID,
			Email:         user.Email,
			Password:      user.Password,
			PIN:           user.PIN,
			FullName:      user.FullName,
			Image:         utils.GetFullImageURL(user.Image),
			Role:          user.Role,
			IsActive:      user.IsActive,
			CreatedAt:     user.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy:     user.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     user.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        0,
		"message":     "Users retrieved successfully",
		"page":        pagination.Page,
		"page_size":   pagination.PageSize,
		"total_items": total,
		"total_pages": (int(total) + pagination.PageSize - 1) / pagination.PageSize,
		"data":        response,
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by ID for the authenticated tenant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(tenantID, uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	var createdByName, updatedByName *string
	if user.Creator != nil {
		name := user.Creator.FullName
		createdByName = &name
	}
	if user.Updater != nil {
		name := user.Updater.FullName
		updatedByName = &name
	}

	utils.Success(c, "User retrieved successfully", dto.UserResponse{
		ID:            user.ID,
		TenantID:      user.TenantID,
		BranchID:      user.BranchID,
		Email:         user.Email,
		Password:      user.Password,
		PIN:           user.PIN,
		FullName:      user.FullName,
		Image:         utils.GetFullImageURL(user.Image),
		Role:          user.Role,
		IsActive:      user.IsActive,
		CreatedAt:     user.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:     user.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     user.UpdatedBy,
		UpdatedByName: updatedByName,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user for the authenticated tenant (supports both JSON and multipart/form-data with optional image)
// @Tags users
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param request body dto.CreateUserRequest true "User data (JSON)" (when using application/json)
// @Param email formData string true "Email" (when using multipart/form-data)
// @Param password formData string true "Password" (when using multipart/form-data)
// @Param full_name formData string true "Full name" (when using multipart/form-data)
// @Param role formData string true "Role (user/branchadmin/tenantadmin)" (when using multipart/form-data)
// @Param branch_id formData integer true "Branch ID" (when using multipart/form-data)
// @Param is_active formData boolean false "Is active" (when using multipart/form-data)
// @Param image formData file false "User image file (optional)" (when using multipart/form-data)
// @Success 200 {object} map[string]interface{}
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")

	var req dto.CreateUserRequest
	contentType := c.GetHeader("Content-Type")

	// Check if request is multipart form-data
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		req.Email = c.PostForm("email")
		req.Password = c.PostForm("password")
		req.FullName = c.PostForm("full_name")
		req.Role = c.PostForm("role")

		// Parse branch_id
		if branchIDStr := c.PostForm("branch_id"); branchIDStr != "" {
			branchID, err := strconv.ParseUint(branchIDStr, 10, 32)
			if err != nil {
				utils.BadRequest(c, "Invalid branch_id format")
				return
			}
			req.BranchID = uint(branchID)
		}

		// Parse is_active
		if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
			isActive := isActiveStr == "true" || isActiveStr == "1"
			req.IsActive = &isActive
		}

		// Validate required fields
		if req.Email == "" {
			utils.BadRequest(c, "Email is required")
			return
		}
		if req.Password == "" {
			utils.BadRequest(c, "Password is required")
			return
		}
		if req.FullName == "" {
			utils.BadRequest(c, "Full name is required")
			return
		}
		if req.Role == "" {
			utils.BadRequest(c, "Role is required")
			return
		}
		if req.BranchID == 0 {
			utils.BadRequest(c, "Branch ID is required")
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

	user, err := h.userService.CreateUser(tenantID, req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Handle image upload if provided (only for multipart form-data)
	if strings.Contains(contentType, "multipart/form-data") {
		if file, err := c.FormFile("image"); err == nil {
			imageURL, uploadErr := h.handleUserImageUpload(file, user.ID)
			if uploadErr != nil {
				utils.BadRequest(c, fmt.Sprintf("User created but image upload failed: %s", uploadErr.Error()))
				return
			}
			// Update user with image URL
			user, _ = h.userService.UpdateUserImage(user.ID, tenantID, imageURL)
		}
	}

	utils.Success(c, "User created successfully", dto.UserResponse{
		ID:        user.ID,
		TenantID:  user.TenantID,
		BranchID:  user.BranchID,
		Email:     user.Email,
		Password:  user.Password,
		PIN:       user.PIN,
		FullName:  user.FullName,
		Image:     utils.GetFullImageURL(user.Image),
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: user.CreatedBy,
	})
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update user details for the authenticated tenant (supports both JSON and multipart/form-data with optional image)
// @Tags users
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "User data to update (JSON)" (when using application/json)
// @Param email formData string false "Email" (when using multipart/form-data)
// @Param password formData string false "Password" (when using multipart/form-data)
// @Param full_name formData string false "Full name" (when using multipart/form-data)
// @Param role formData string false "Role (user/branchadmin/tenantadmin)" (when using multipart/form-data)
// @Param branch_id formData integer false "Branch ID" (when using multipart/form-data)
// @Param is_active formData boolean false "Is active" (when using multipart/form-data)
// @Param image formData file false "User image file (optional)" (when using multipart/form-data)
// @Success 200 {object} map[string]interface{}
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
	contentType := c.GetHeader("Content-Type")

	// Check if request is multipart form-data
	if strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		if email := c.PostForm("email"); email != "" {
			req.Email = &email
		}
		if password := c.PostForm("password"); password != "" {
			req.Password = &password
		}
		if fullName := c.PostForm("full_name"); fullName != "" {
			req.FullName = &fullName
		}
		if role := c.PostForm("role"); role != "" {
			req.Role = &role
		}

		// Parse branch_id
		if branchIDStr := c.PostForm("branch_id"); branchIDStr != "" {
			branchID, err := strconv.ParseUint(branchIDStr, 10, 32)
			if err != nil {
				utils.BadRequest(c, "Invalid branch_id format")
				return
			}
			branchIDUint := uint(branchID)
			req.BranchID = &branchIDUint
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

	user, err := h.userService.UpdateUser(tenantID, uint(userID), req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.BadRequest(c, err.Error())
		return
	}

	// Handle image upload if provided (only for multipart form-data)
	if strings.Contains(contentType, "multipart/form-data") {
		if file, err := c.FormFile("image"); err == nil {
			// Delete old image if exists
			if user.Image != "" {
				uploadDir := "uploads/profiles"
				oldPath := filepath.Join(uploadDir, filepath.Base(user.Image))
				os.Remove(oldPath) // Ignore error if file doesn't exist
			}

			imageURL, uploadErr := h.handleUserImageUpload(file, user.ID)
			if uploadErr != nil {
				utils.BadRequest(c, fmt.Sprintf("User updated but image upload failed: %s", uploadErr.Error()))
				return
			}
			// Update user with image URL
			user, _ = h.userService.UpdateUserImage(user.ID, tenantID, imageURL)
		}
	}

	utils.Success(c, "User updated successfully", dto.UserResponse{
		ID:        user.ID,
		TenantID:  user.TenantID,
		BranchID:  user.BranchID,
		Email:     user.Email,
		Password:  user.Password,
		PIN:       user.PIN,
		FullName:  user.FullName,
		Image:     utils.GetFullImageURL(user.Image),
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	})
}

// handleUserImageUpload is a helper method to handle image upload for users
func (h *UserHandler) handleUserImageUpload(file *multipart.FileHeader, userID uint) (string, error) {
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
	uploadDir := "uploads/profiles"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory")
	}

	// Generate unique filename
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
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
	imageURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	return imageURL, nil
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user (soft delete) for the authenticated tenant
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	// Get current user ID
	currentUserID := c.GetUint("user_id")

	if err := h.userService.DeleteUser(tenantID, uint(userID), &currentUserID); err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "User deleted successfully")
}
