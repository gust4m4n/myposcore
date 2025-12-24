package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterHandler struct {
	*BaseHandler
	registerService *services.RegisterService
}

func NewRegisterHandler(cfg *config.Config) *RegisterHandler {
	return &RegisterHandler{
		BaseHandler:     NewBaseHandler(cfg),
		registerService: services.NewRegisterService(),
	}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, branch, err := h.registerService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Username, h.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:         user.ID,
			TenantID:   user.TenantID,
			BranchID:   user.BranchID,
			BranchName: branch.Name,
			Username:   user.Username,
			Email:      user.Email,
			FullName:   user.FullName,
			Role:       user.Role,
			IsActive:   user.IsActive,
		},
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}
