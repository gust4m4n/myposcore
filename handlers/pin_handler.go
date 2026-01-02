package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type PINHandler struct {
	*BaseHandler
	service           *services.PINService
	auditTrailService *services.AuditTrailService
}

func NewPINHandler(cfg *config.Config, auditTrailService *services.AuditTrailService) *PINHandler {
	return &PINHandler{
		BaseHandler:       NewBaseHandler(cfg),
		service:           services.NewPINService(),
		auditTrailService: auditTrailService,
	}
}

// CreatePIN godoc
// @Summary Create a new PIN
// @Description Create a 6-digit PIN for the authenticated user
// @Tags pin
// @Accept json
// @Produce json
// @Param request body dto.CreatePINRequest true "PIN data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pin/create [post]
func (h *PINHandler) CreatePIN(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreatePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.CreatePIN(userID, req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "PIN created successfully")
}

// ChangePIN godoc
// @Summary Change existing PIN
// @Description Change the user's existing 6-digit PIN
// @Tags pin
// @Accept json
// @Produce json
// @Param request body dto.ChangePINRequest true "PIN data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pin/change [put]
func (h *PINHandler) ChangePIN(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ChangePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.service.ChangePIN(userID, req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithoutData(c, "PIN created successfully")
}

// CheckPIN godoc
// @Summary Check if user has PIN
// @Description Check whether the authenticated user has set a PIN
// @Tags pin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pin/check [get]
func (h *PINHandler) CheckPIN(c *gin.Context) {
	userID := c.GetUint("user_id")

	hasPIN, err := h.service.HasPIN(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, "PIN check completed", gin.H{"has_pin": hasPIN})
}
