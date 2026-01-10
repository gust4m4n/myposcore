package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configService *services.ConfigService
}

func NewConfigHandler(configService *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// SetConfig handles POST /api/config/set
func (h *ConfigHandler) SetConfig(c *gin.Context) {
	var req dto.SetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.configService.SetConfig(req.Key, req.Value); err != nil {
		utils.InternalError(c, "Failed to set config: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"key":   req.Key,
		"value": req.Value,
	}

	utils.Success(c, "Config set successfully", response)
}

// GetConfig handles GET /api/config/get/:key
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		utils.BadRequest(c, "Key parameter is required")
		return
	}

	value, err := h.configService.GetConfig(key)
	if err != nil {
		if err.Error() == "config key not found" {
			utils.NotFound(c, "Config key not found")
			return
		}
		utils.InternalError(c, "Failed to get config: "+err.Error())
		return
	}

	response := dto.GetConfigResponse{
		Key:   key,
		Value: value,
	}

	utils.Success(c, "Config retrieved successfully", response)
}
