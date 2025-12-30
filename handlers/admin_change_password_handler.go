package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminChangePasswordHandler struct {
	*BaseHandler
	service *services.AdminChangePasswordService
}

func NewAdminChangePasswordHandler(cfg *config.Config) *AdminChangePasswordHandler {
	return &AdminChangePasswordHandler{
		BaseHandler: NewBaseHandler(cfg),
		service:     services.NewAdminChangePasswordService(),
	}
}

// AdminChangePassword godoc
// @Summary Change user password by admin
// @Description Allows higher role users (owner, admin, superadmin) to change password of lower role users
// @Tags admin
// @Accept json
// @Produce json
// @Param request body dto.AdminChangePasswordRequest true "Change password data"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/change-password [put]
func (h *AdminChangePasswordHandler) Handle(c *gin.Context) {
	// Get admin user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AdminChangePassword(userID.(uint), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
