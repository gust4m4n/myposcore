package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type AdminChangePINHandler struct {
	*BaseHandler
	service           *services.AdminChangePINService
	auditTrailService *services.AuditTrailService
}

func NewAdminChangePINHandler(cfg *config.Config, auditTrailService *services.AuditTrailService) *AdminChangePINHandler {
	return &AdminChangePINHandler{
		BaseHandler:       NewBaseHandler(cfg),
		service:           services.NewAdminChangePINService(),
		auditTrailService: auditTrailService,
	}
}

// AdminChangePIN godoc
// @Summary Change user PIN by admin
// @Description Allows higher role users (owner, admin, superadmin) to change PIN of lower role users
// @Tags admin
// @Accept json
// @Produce json
// @Param request body dto.AdminChangePINRequest true "Change PIN data"
// @Success 200 {object} map[string]string
// @Router /api/v1/admin/change-pin [put]
func (h *AdminChangePINHandler) Handle(c *gin.Context) {
	// Get admin user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.AdminChangePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.AdminChangePIN(userID.(uint), req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "PIN changed successfully")
}
