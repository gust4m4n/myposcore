package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	*BaseHandler
	loginService      *services.LoginService
	auditTrailService *services.AuditTrailService
}

func NewLoginHandler(cfg *config.Config, auditTrailService *services.AuditTrailService) *LoginHandler {
	return &LoginHandler{
		BaseHandler:       NewBaseHandler(cfg),
		loginService:      services.NewLoginService(),
		auditTrailService: auditTrailService,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, tenant, branch, err := h.loginService.Login(req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Email, h.config.JWTSecret)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	response := dto.AuthResponse{
		Token: token,
		User: dto.UserProfile{
			ID:         user.ID,
			TenantID:   user.TenantID,
			BranchID:   user.BranchID,
			BranchName: branch.Name,
			Email:      user.Email,
			FullName:   user.FullName,
			Role:       user.Role,
			IsActive:   user.IsActive,
		},
		Tenant: dto.TenantInfo{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Description: tenant.Description,
			Address:     tenant.Address,
			Website:     tenant.Website,
			Email:       tenant.Email,
			Phone:       tenant.Phone,
			Image:       tenant.Image,
			IsActive:    tenant.IsActive,
		},
		Branch: dto.BranchInfo{
			ID:          branch.ID,
			Name:        branch.Name,
			Description: branch.Description,
			Address:     branch.Address,
			Website:     branch.Website,
			Email:       branch.Email,
			Phone:       branch.Phone,
			Image:       branch.Image,
			IsActive:    branch.IsActive,
		},
	}

	utils.Success(c, "Login successful", response)
}
