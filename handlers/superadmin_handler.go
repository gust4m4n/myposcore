package handlers

import (
	"myposcore/config"
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuperAdminHandler struct {
	*BaseHandler
	tenantService    *services.SuperAdminTenantService
	branchService    *services.SuperAdminBranchService
	userService      *services.SuperAdminUserService
	dashboardService *services.SuperAdminDashboardService
}

func NewSuperAdminHandler(cfg *config.Config) *SuperAdminHandler {
	return &SuperAdminHandler{
		BaseHandler:      NewBaseHandler(cfg),
		tenantService:    services.NewSuperAdminTenantService(),
		branchService:    services.NewSuperAdminBranchService(),
		userService:      services.NewSuperAdminUserService(),
		dashboardService: services.NewSuperAdminDashboardService(),
	}
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get list of all tenants (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Success 200 {object} []dto.TenantResponse
// @Router /superadmin/tenants [get]
func (h *SuperAdminHandler) ListTenants(c *gin.Context) {
	tenants, err := h.tenantService.ListTenants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.TenantResponse
	for _, tenant := range tenants {
		response = append(response, dto.TenantResponse{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Code:      tenant.Code,
			IsActive:  tenant.IsActive,
			CreatedAt: tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param request body dto.CreateTenantRequest true "Tenant data"
// @Success 201 {object} dto.TenantResponse
// @Router /superadmin/tenants [post]
func (h *SuperAdminHandler) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.tenantService.CreateTenant(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tenant created successfully",
		"data": dto.TenantResponse{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Code:      tenant.Code,
			IsActive:  tenant.IsActive,
			CreatedAt: tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// ListBranches godoc
// @Summary List branches by tenant
// @Description Get list of branches for a specific tenant (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Success 200 {object} []dto.BranchResponse
// @Router /superadmin/tenants/{tenant_id}/branches [get]
func (h *SuperAdminHandler) ListBranches(c *gin.Context) {
	tenantID, err := strconv.ParseUint(c.Param("tenant_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	branches, err := h.branchService.ListBranchesByTenant(uint(tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.BranchResponse
	for _, branch := range branches {
		response = append(response, dto.BranchResponse{
			ID:        branch.ID,
			TenantID:  branch.TenantID,
			Name:      branch.Name,
			Code:      branch.Code,
			Address:   branch.Address,
			Phone:     branch.Phone,
			IsActive:  branch.IsActive,
			CreatedAt: branch.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ListUsers godoc
// @Summary List users by branch
// @Description Get list of users for a specific branch (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {object} []dto.UserResponse
// @Router /superadmin/branches/{branch_id}/users [get]
func (h *SuperAdminHandler) ListUsers(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	users, err := h.userService.ListUsersByBranch(uint(branchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:        user.ID,
			TenantID:  user.TenantID,
			BranchID:  user.BranchID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// Dashboard godoc
// @Summary Get dashboard statistics
// @Description Get total counts of tenants, branches, and users (superadmin only)
// @Tags superadmin
// @Accept json
// @Produce json
// @Success 200 {object} dto.DashboardResponse
// @Router /superadmin/dashboard [get]
func (h *SuperAdminHandler) Dashboard(c *gin.Context) {
	dashboard, err := h.dashboardService.GetDashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dashboard,
	})
}
