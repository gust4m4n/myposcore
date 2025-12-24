package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChangePasswordHandler struct {
	*BaseHandler
	service *services.ChangePasswordService
}

func NewChangePasswordHandler(cfg *config.Config) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		BaseHandler: NewBaseHandler(cfg),
		service:     services.NewChangePasswordService(),
	}
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password data"
// @Success 200 {object} map[string]string
// @Router /api/v1/change-password [put]
func (h *ChangePasswordHandler) Handle(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangePassword(userID.(uint), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
