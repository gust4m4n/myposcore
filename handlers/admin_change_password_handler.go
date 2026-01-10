package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type AdminChangePasswordHandler struct {
	*BaseHandler
	service           *services.AdminChangePasswordService
	auditTrailService *services.AuditTrailService
}

func NewAdminChangePasswordHandler(cfg *config.Config, auditTrailService *services.AuditTrailService) *AdminChangePasswordHandler {
	return &AdminChangePasswordHandler{
		BaseHandler:       NewBaseHandler(cfg),
		service:           services.NewAdminChangePasswordService(),
		auditTrailService: auditTrailService,
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
// @Router /api/admin/change-password [put]
func (h *AdminChangePasswordHandler) Handle(c *gin.Context) {
	// Get admin user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.AdminChangePassword(userID.(uint), req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "Password changed successfully")
}
