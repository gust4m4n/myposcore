package handlers

import (
	"myposcore/config"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	*BaseHandler
	auditTrailService *services.AuditTrailService
}

func NewLogoutHandler(cfg *config.Config, auditTrailService *services.AuditTrailService) *LogoutHandler {
	return &LogoutHandler{
		BaseHandler:       NewBaseHandler(cfg),
		auditTrailService: auditTrailService,
	}
}

// Handle godoc
// @Summary Logout user
// @Description Logout user and record audit trail. Since JWT is stateless, this only records the logout event in audit trail. Client should discard the token.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/logout [post]
func (h *LogoutHandler) Handle(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	tenantID, _ := c.Get("tenantID")
	branchID, _ := c.Get("branchID")

	// Record logout audit trail
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	auditChanges := map[string]interface{}{
		"ip_address": ipAddress,
		"user_agent": userAgent,
	}

	uid := userID.(uint)
	tid := tenantID.(uint)
	bid := branchID.(uint)
	_ = h.auditTrailService.CreateAuditTrail(&tid, &bid, uid, "auth", uid, "logout", auditChanges, ipAddress, userAgent)

	utils.SuccessWithoutData(c, "Logout successful. Please discard your token.")
}
