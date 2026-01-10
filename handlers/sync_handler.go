package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"myposcore/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	syncService *services.SyncService
}

func NewSyncHandler(syncService *services.SyncService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

// UploadFromClient godoc
// @Summary Upload data from mobile client to server
// @Description Upload orders and payments from mobile app (offline mode) to server. Supports conflict detection and resolution.
// @Tags Sync
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.SyncUploadRequest true "Upload sync data"
// @Success 200 {object} dto.SuccessResponse{data=dto.SyncUploadResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sync/upload [post]
func (h *SyncHandler) UploadFromClient(c *gin.Context) {
	var req dto.SyncUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 1, err.Error())
		return
	}

	// Get user info from JWT
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	branchID, _ := c.Get("branch_id")

	result, err := h.syncService.UploadFromClient(&req, tenantID.(uint), branchID.(uint), userID.(uint))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 1, err.Error())
		return
	}

	// Check if there are any errors or conflicts
	if len(result.Errors) > 0 || len(result.Conflicts) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Sync completed with warnings. Check errors and conflicts.",
			"data":    result,
		})
		return
	}

	utils.Success(c, "Data synced successfully", result)
}

// DownloadToClient godoc
// @Summary Download master data from server to mobile client
// @Description Download products, categories and other master data to mobile app. Supports delta sync (only updated data).
// @Tags Sync
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.SyncDownloadRequest true "Download sync request"
// @Success 200 {object} dto.SuccessResponse{data=dto.SyncDownloadResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sync/download [post]
func (h *SyncHandler) DownloadToClient(c *gin.Context) {
	var req dto.SyncDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 1, err.Error())
		return
	}

	tenantID, _ := c.Get("tenant_id")

	result, err := h.syncService.DownloadToClient(&req, tenantID.(uint))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 1, err.Error())
		return
	}

	utils.Success(c, "Data downloaded successfully", result)
}

// GetSyncStatus godoc
// @Summary Get sync status for mobile client
// @Description Get current sync status including last sync time, pending uploads, and unresolved conflicts
// @Tags Sync
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param client_id query string true "Client device ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.SyncStatusResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sync/status [get]
func (h *SyncHandler) GetSyncStatus(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		utils.Error(c, http.StatusBadRequest, 1, "client_id is required")
		return
	}

	tenantID, _ := c.Get("tenant_id")

	result, err := h.syncService.GetSyncStatus(clientID, tenantID.(uint))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 1, err.Error())
		return
	}

	utils.Success(c, "Sync status retrieved successfully", result)
}

// GetSyncLogs godoc
// @Summary Get sync history logs
// @Description Get paginated list of sync operations history for debugging and monitoring
// @Tags Sync
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param client_id query string true "Client device ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Success 200 {object} dto.SuccessResponse{data=dto.PaginatedResponse{data=[]dto.SyncLogResponse}}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sync/logs [get]
func (h *SyncHandler) GetSyncLogs(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		utils.Error(c, http.StatusBadRequest, 1, "client_id is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tenantID, _ := c.Get("tenant_id")

	logs, total, err := h.syncService.GetSyncLogs(clientID, tenantID.(uint), page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 1, err.Error())
		return
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	utils.Success(c, "Sync logs retrieved successfully", gin.H{
		"items": logs,
		"pagination": gin.H{
			"page":        page,
			"per_page":    pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// ResolveConflict godoc
// @Summary Manually resolve sync conflict
// @Description Resolve detected sync conflict with chosen strategy: server_wins, client_wins, or manual merge
// @Tags Sync
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.ResolveConflictRequest true "Conflict resolution request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sync/conflicts/resolve [post]
func (h *SyncHandler) ResolveConflict(c *gin.Context) {
	var req dto.ResolveConflictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 1, err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.syncService.ResolveConflict(&req, userID.(uint)); err != nil {
		if err.Error() == "record not found" {
			utils.Error(c, http.StatusNotFound, 1, "Conflict not found")
			return
		}
		utils.Error(c, http.StatusInternalServerError, 1, err.Error())
		return
	}

	utils.Success(c, "Conflict resolved successfully", nil)
}

// GetServerTime godoc
// @Summary Get current server time
// @Description Get server timestamp for time synchronization between client and server
// @Tags Sync
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse{data=map[string]string}
// @Router /api/sync/time [get]
func (h *SyncHandler) GetServerTime(c *gin.Context) {
	now := time.Now()
	utils.Success(c, "Server time retrieved successfully", gin.H{
		"server_time": now.Format("2006-01-02T15:04:05.000Z07:00"),
		"unix_time":   now.Unix(),
	})
}
