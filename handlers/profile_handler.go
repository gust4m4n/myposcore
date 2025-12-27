package handlers

import (
	"fmt"
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	*BaseHandler
	authService *services.AuthService
}

func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler: NewBaseHandler(cfg),
		authService: services.NewAuthService(),
	}
}

func (h *ProfileHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile information (email, full name, PIN)
// @Tags profile
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.authService.UpdateProfile(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    profile,
	})
}

// UploadProfileImage godoc
// @Summary Upload profile image
// @Description Upload or update user profile image
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile image file"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profile/photo [post]
func (h *ProfileHandler) UploadProfileImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user profile
	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
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
	uploadDir := "uploads/profiles"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Delete old image if exists
	if profile.User.Image != "" {
		oldPath := filepath.Join(uploadDir, filepath.Base(profile.User.Image))
		os.Remove(oldPath) // Ignore error if file doesn't exist
	}

	// Generate unique filename
	filename := fmt.Sprintf("user_%d_%d%s", userID.(uint), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}

	// Update user image URL
	imageURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	updatedProfile, err := h.authService.UpdateProfileImage(userID.(uint), imageURL)
	if err != nil {
		// Delete uploaded file if database update fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"data":    updatedProfile,
	})
}

// DeleteProfileImage godoc
// @Summary Delete profile image
// @Description Delete user profile image
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/profile/photo [delete]
func (h *ProfileHandler) DeleteProfileImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user profile
	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Delete image file if exists
	if profile.User.Image != "" {
		uploadDir := "uploads/profiles"
		filePath := filepath.Join(uploadDir, filepath.Base(profile.User.Image))
		os.Remove(filePath) // Ignore error if file doesn't exist
	}

	// Update database
	_, err = h.authService.UpdateProfileImage(userID.(uint), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
