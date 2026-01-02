package handlers

import (
	"myposcore/config"
	"myposcore/services"
	"myposcore/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DevHandler struct {
	*BaseHandler
	service *services.DevService
}

func NewDevHandler(cfg *config.Config) *DevHandler {
	return &DevHandler{
		BaseHandler: NewBaseHandler(cfg),
		service:     services.NewDevService(),
	}
}

func (h *DevHandler) ListTenants(c *gin.Context) {
	tenants, err := h.service.ListTenants()
	if err != nil {
		utils.InternalError(c, "Failed to fetch tenants")
		return
	}

	utils.Success(c, "Tenants retrieved successfully", tenants)
}

func (h *DevHandler) ListBranchesByTenant(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tenant ID")
		return
	}

	branches, err := h.service.ListBranchesByTenant(uint(tenantID))
	if err != nil {
		utils.InternalError(c, "Failed to fetch branches")
		return
	}

	utils.Success(c, "Branches retrieved successfully", branches)
}
