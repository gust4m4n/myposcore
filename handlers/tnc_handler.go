package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TnCHandler struct {
	BaseHandler
	tncService *services.TnCService
}

func NewTnCHandler(tncService *services.TnCService) *TnCHandler {
	return &TnCHandler{
		tncService: tncService,
	}
}

// CreateTnC godoc
// @Summary Create new terms and conditions
// @Description Create a new terms and conditions entry
// @Tags TnC
// @Accept json
// @Produce json
// @Param request body dto.CreateTnCRequest true "TnC request"
// @Success 200 {object} dto.TnCResponse
// @Router /api/tnc [post]
func (h *TnCHandler) CreateTnC(c *gin.Context) {
	var req dto.CreateTnCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")
	req.CreatedBy = &currentUserID

	tnc, err := h.tncService.CreateTnC(req.Title, req.Content, req.Version, req.CreatedBy)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to create terms and conditions")
		return
	}

	// Reload to get audit info
	tnc, _ = h.tncService.GetTnCByID(tnc.ID)

	var createdByName *string
	if tnc.Creator != nil {
		name := tnc.Creator.FullName
		createdByName = &name
	}

	response := dto.TnCResponse{
		ID:            tnc.ID,
		Title:         tnc.Title,
		Content:       tnc.Content,
		Version:       tnc.Version,
		IsActive:      tnc.IsActive,
		CreatedAt:     tnc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tnc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     tnc.CreatedBy,
		CreatedByName: createdByName,
	}

	h.SuccessResponse(c, http.StatusOK, "Terms and conditions created successfully", response)
}

// GetActiveTnC godoc
// @Summary Get active terms and conditions
// @Description Get the currently active terms and conditions
// @Tags TnC
// @Produce json
// @Success 200 {object} dto.TnCResponse
// @Router /api/tnc/active [get]
func (h *TnCHandler) GetActiveTnC(c *gin.Context) {
	tnc, err := h.tncService.GetActiveTnC()
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "No active terms and conditions found")
		return
	}

	var createdByName, updatedByName *string
	if tnc.Creator != nil {
		name := tnc.Creator.FullName
		createdByName = &name
	}
	if tnc.Updater != nil {
		name := tnc.Updater.FullName
		updatedByName = &name
	}

	response := dto.TnCResponse{
		ID:            tnc.ID,
		Title:         tnc.Title,
		Content:       tnc.Content,
		Version:       tnc.Version,
		IsActive:      tnc.IsActive,
		CreatedAt:     tnc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tnc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     tnc.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     tnc.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	h.SuccessResponse(c, http.StatusOK, "Active terms and conditions retrieved successfully", response)
}

// GetAllTnC godoc
// @Summary Get all terms and conditions
// @Description Get all terms and conditions entries
// @Tags TnC
// @Produce json
// @Success 200 {array} dto.TnCResponse
// @Router /api/tnc [get]
func (h *TnCHandler) GetAllTnC(c *gin.Context) {
	tncs, err := h.tncService.GetAllTnC()
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve terms and conditions")
		return
	}

	var responses []dto.TnCResponse
	for _, tnc := range tncs {
		var createdByName, updatedByName *string
		if tnc.Creator != nil {
			name := tnc.Creator.FullName
			createdByName = &name
		}
		if tnc.Updater != nil {
			name := tnc.Updater.FullName
			updatedByName = &name
		}

		responses = append(responses, dto.TnCResponse{
			ID:            tnc.ID,
			Title:         tnc.Title,
			Content:       tnc.Content,
			Version:       tnc.Version,
			IsActive:      tnc.IsActive,
			CreatedAt:     tnc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     tnc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedBy:     tnc.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     tnc.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	h.SuccessResponse(c, http.StatusOK, "Terms and conditions retrieved successfully", responses)
}

// GetTnCByID godoc
// @Summary Get terms and conditions by ID
// @Description Get specific terms and conditions by ID
// @Tags TnC
// @Produce json
// @Param id path int true "TnC ID"
// @Success 200 {object} dto.TnCResponse
// @Router /api/tnc/{id} [get]
func (h *TnCHandler) GetTnCByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid TnC ID")
		return
	}

	tnc, err := h.tncService.GetTnCByID(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "Terms and conditions not found")
		return
	}

	var createdByName, updatedByName *string
	if tnc.Creator != nil {
		name := tnc.Creator.FullName
		createdByName = &name
	}
	if tnc.Updater != nil {
		name := tnc.Updater.FullName
		updatedByName = &name
	}

	response := dto.TnCResponse{
		ID:            tnc.ID,
		Title:         tnc.Title,
		Content:       tnc.Content,
		Version:       tnc.Version,
		IsActive:      tnc.IsActive,
		CreatedAt:     tnc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tnc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     tnc.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     tnc.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	h.SuccessResponse(c, http.StatusOK, "Terms and conditions retrieved successfully", response)
}

// UpdateTnC godoc
// @Summary Update terms and conditions
// @Description Update an existing terms and conditions entry
// @Tags TnC
// @Accept json
// @Produce json
// @Param id path int true "TnC ID"
// @Param request body dto.UpdateTnCRequest true "Update TnC request"
// @Success 200 {object} dto.TnCResponse
// @Router /api/tnc/{id} [put]
func (h *TnCHandler) UpdateTnC(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid TnC ID")
		return
	}

	var req dto.UpdateTnCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")
	req.UpdatedBy = &currentUserID

	var title, content, version *string
	if req.Title != "" {
		title = &req.Title
	}
	if req.Content != "" {
		content = &req.Content
	}
	if req.Version != "" {
		version = &req.Version
	}

	tnc, err := h.tncService.UpdateTnC(uint(id), title, content, version, req.IsActive, req.UpdatedBy)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to update terms and conditions")
		return
	}

	// Reload to get audit info
	tnc, _ = h.tncService.GetTnCByID(uint(id))

	var createdByName, updatedByName *string
	if tnc.Creator != nil {
		name := tnc.Creator.FullName
		createdByName = &name
	}
	if tnc.Updater != nil {
		name := tnc.Updater.FullName
		updatedByName = &name
	}

	response := dto.TnCResponse{
		ID:            tnc.ID,
		Title:         tnc.Title,
		Content:       tnc.Content,
		Version:       tnc.Version,
		IsActive:      tnc.IsActive,
		CreatedAt:     tnc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     tnc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     tnc.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     tnc.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	h.SuccessResponse(c, http.StatusOK, "Terms and conditions updated successfully", response)
}

// DeleteTnC godoc
// @Summary Delete terms and conditions
// @Description Delete a terms and conditions entry
// @Tags TnC
// @Produce json
// @Param id path int true "TnC ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/tnc/{id} [delete]
func (h *TnCHandler) DeleteTnC(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid TnC ID")
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")

	if err := h.tncService.DeleteTnC(uint(id), &currentUserID); err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete terms and conditions")
		return
	}

	h.SuccessResponse(c, http.StatusOK, "Terms and conditions deleted successfully", nil)
}
