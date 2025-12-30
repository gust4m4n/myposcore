package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuditTrailHandler struct {
	BaseHandler
	auditTrailService *services.AuditTrailService
}

func NewAuditTrailHandler(auditTrailService *services.AuditTrailService) *AuditTrailHandler {
	return &AuditTrailHandler{
		auditTrailService: auditTrailService,
	}
}

// ListAuditTrails godoc
// @Summary List audit trails
// @Description Get list of audit trails with filters
// @Tags audit-trails
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param user_id query int false "Filter by user ID"
// @Param entity_type query string false "Filter by entity type (user, product, order, payment, category, faq, tnc)"
// @Param entity_id query int false "Filter by entity ID"
// @Param action query string false "Filter by action (create, update, delete)"
// @Param date_from query string false "Filter from date (YYYY-MM-DD)"
// @Param date_to query string false "Filter to date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/audit-trails [get]
func (h *AuditTrailHandler) ListAuditTrails(c *gin.Context) {
	// Get tenant ID from context
	var tenantID *uint
	if tid, exists := c.Get("tenant_id"); exists {
		tidUint := tid.(uint)
		tenantID = &tidUint
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var userID, entityID *uint
	var entityType, action *string
	var dateFrom, dateTo *time.Time

	// Parse user_id
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uidUint := uint(uid)
			userID = &uidUint
		}
	}

	// Parse entity_type
	if et := c.Query("entity_type"); et != "" {
		entityType = &et
	}

	// Parse entity_id
	if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
		if eid, err := strconv.ParseUint(entityIDStr, 10, 32); err == nil {
			eidUint := uint(eid)
			entityID = &eidUint
		}
	}

	// Parse action
	if act := c.Query("action"); act != "" {
		action = &act
	}

	// Parse date_from
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if df, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &df
		}
	}

	// Parse date_to
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dt, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &dt
		}
	}

	// Get audit trails
	auditTrails, total, err := h.auditTrailService.ListAuditTrails(
		tenantID, userID, entityID, entityType, action, dateFrom, dateTo, page, limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit trails"})
		return
	}

	// Build responses
	responses := make([]dto.AuditTrailResponse, len(auditTrails))
	for i, trail := range auditTrails {
		userName := ""
		if trail.User != nil {
			userName = trail.User.FullName
		}

		var changes, ipAddress, userAgent *string
		if trail.Changes != "" {
			changes = &trail.Changes
		}
		if trail.IPAddress != "" {
			ipAddress = &trail.IPAddress
		}
		if trail.UserAgent != "" {
			userAgent = &trail.UserAgent
		}

		responses[i] = dto.AuditTrailResponse{
			ID:         trail.ID,
			TenantID:   trail.TenantID,
			BranchID:   trail.BranchID,
			UserID:     trail.UserID,
			UserName:   userName,
			EntityType: trail.EntityType,
			EntityID:   trail.EntityID,
			Action:     trail.Action,
			Changes:    changes,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CreatedAt:  trail.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetAuditTrailByID godoc
// @Summary Get audit trail by ID
// @Description Get specific audit trail by ID
// @Tags audit-trails
// @Produce json
// @Param id path int true "Audit Trail ID"
// @Success 200 {object} dto.AuditTrailResponse
// @Router /api/v1/audit-trails/{id} [get]
func (h *AuditTrailHandler) GetAuditTrailByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit trail ID"})
		return
	}

	// Get tenant ID from context
	var tenantID *uint
	if tid, exists := c.Get("tenant_id"); exists {
		tidUint := tid.(uint)
		tenantID = &tidUint
	}

	auditTrail, err := h.auditTrailService.GetAuditTrailByID(uint(id), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userName := ""
	if auditTrail.User != nil {
		userName = auditTrail.User.FullName
	}

	var changes, ipAddress, userAgent *string
	if auditTrail.Changes != "" {
		changes = &auditTrail.Changes
	}
	if auditTrail.IPAddress != "" {
		ipAddress = &auditTrail.IPAddress
	}
	if auditTrail.UserAgent != "" {
		userAgent = &auditTrail.UserAgent
	}

	response := dto.AuditTrailResponse{
		ID:         auditTrail.ID,
		TenantID:   auditTrail.TenantID,
		BranchID:   auditTrail.BranchID,
		UserID:     auditTrail.UserID,
		UserName:   userName,
		EntityType: auditTrail.EntityType,
		EntityID:   auditTrail.EntityID,
		Action:     auditTrail.Action,
		Changes:    changes,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  auditTrail.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetEntityAuditHistory godoc
// @Summary Get entity audit history
// @Description Get audit history for a specific entity
// @Tags audit-trails
// @Produce json
// @Param entity_type path string true "Entity type (user, product, order, payment, category, faq, tnc)"
// @Param entity_id path int true "Entity ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/audit-trails/entity/{entity_type}/{entity_id} [get]
func (h *AuditTrailHandler) GetEntityAuditHistory(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityIDStr := c.Param("entity_id")

	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get tenant ID from context
	var tenantID *uint
	if tid, exists := c.Get("tenant_id"); exists {
		tidUint := tid.(uint)
		tenantID = &tidUint
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	auditTrails, total, err := h.auditTrailService.GetEntityAuditHistory(
		tenantID, entityType, uint(entityID), page, limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit history"})
		return
	}

	// Build responses
	responses := make([]dto.AuditTrailResponse, len(auditTrails))
	for i, trail := range auditTrails {
		userName := ""
		if trail.User != nil {
			userName = trail.User.FullName
		}

		var changes, ipAddress, userAgent *string
		if trail.Changes != "" {
			changes = &trail.Changes
		}
		if trail.IPAddress != "" {
			ipAddress = &trail.IPAddress
		}
		if trail.UserAgent != "" {
			userAgent = &trail.UserAgent
		}

		responses[i] = dto.AuditTrailResponse{
			ID:         trail.ID,
			TenantID:   trail.TenantID,
			BranchID:   trail.BranchID,
			UserID:     trail.UserID,
			UserName:   userName,
			EntityType: trail.EntityType,
			EntityID:   trail.EntityID,
			Action:     trail.Action,
			Changes:    changes,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CreatedAt:  trail.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUserActivityLog godoc
// @Summary Get user activity log
// @Description Get all activities performed by a specific user
// @Tags audit-trails
// @Produce json
// @Param user_id path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/audit-trails/user/{user_id} [get]
func (h *AuditTrailHandler) GetUserActivityLog(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get tenant ID from context
	var tenantID *uint
	if tid, exists := c.Get("tenant_id"); exists {
		tidUint := tid.(uint)
		tenantID = &tidUint
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	auditTrails, total, err := h.auditTrailService.GetUserActivityLog(
		tenantID, uint(userID), page, limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user activity log"})
		return
	}

	// Build responses
	responses := make([]dto.AuditTrailResponse, len(auditTrails))
	for i, trail := range auditTrails {
		userName := ""
		if trail.User != nil {
			userName = trail.User.FullName
		}

		var changes, ipAddress, userAgent *string
		if trail.Changes != "" {
			changes = &trail.Changes
		}
		if trail.IPAddress != "" {
			ipAddress = &trail.IPAddress
		}
		if trail.UserAgent != "" {
			userAgent = &trail.UserAgent
		}

		responses[i] = dto.AuditTrailResponse{
			ID:         trail.ID,
			TenantID:   trail.TenantID,
			BranchID:   trail.BranchID,
			UserID:     trail.UserID,
			UserName:   userName,
			EntityType: trail.EntityType,
			EntityID:   trail.EntityID,
			Action:     trail.Action,
			Changes:    changes,
			IPAddress:  ipAddress,
			UserAgent:  userAgent,
			CreatedAt:  trail.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
