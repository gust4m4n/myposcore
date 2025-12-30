package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	*BaseHandler
	loginService *services.LoginService
}

func NewLoginHandler(cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		BaseHandler:  NewBaseHandler(cfg),
		loginService: services.NewLoginService(),
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tenant, branch, err := h.loginService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		Tenant: dto.TenantInfo{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Code:        tenant.Code,
			Description: tenant.Description,
			Address:     tenant.Address,
			Website:     tenant.Website,
			Email:       tenant.Email,
			Phone:       tenant.Phone,
			IsActive:    tenant.IsActive,
		},
		Branch: dto.BranchInfo{
			ID:          branch.ID,
			Name:        branch.Name,
			Code:        branch.Code,
			Description: branch.Description,
			Address:     branch.Address,
			Website:     branch.Website,
			Email:       branch.Email,
			Phone:       branch.Phone,
			IsActive:    branch.IsActive,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}
