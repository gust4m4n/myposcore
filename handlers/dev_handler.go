package handlers

import (
	"myposcore/config"
	"myposcore/services"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tenants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenants retrieved successfully",
		"data":    tenants,
	})
}

func (h *DevHandler) ListBranchesByTenant(c *gin.Context) {
	tenantIDStr := c.Param("tenant_id")
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tenant ID",
		})
		return
	}

	branches, err := h.service.ListBranchesByTenant(uint(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch branches",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Branches retrieved successfully",
		"data":    branches,
	})
}
