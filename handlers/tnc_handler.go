package handlers

import (
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type TnCHandler struct {
	configService *services.ConfigService
}

func NewTnCHandler(configService *services.ConfigService) *TnCHandler {
	return &TnCHandler{
		configService: configService,
	}
}

// GetTnC godoc
// @Summary Get terms and conditions
// @Description Get the terms and conditions content from config table
// @Tags TnC
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tnc [get]
func (h *TnCHandler) GetTnC(c *gin.Context) {
	// Get TnC content from config table with key "tnc"
	content, err := h.configService.GetConfig("tnc")
	if err != nil {
		utils.NotFound(c, "Terms and conditions not found")
		return
	}

	response := map[string]interface{}{
		"title":   "Terms and Conditions - MyPOS Core System",
		"content": content,
	}

	utils.Success(c, "Terms and conditions retrieved successfully", response)
}
