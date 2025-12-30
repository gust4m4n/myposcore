package handlers

import (
	"myposcore/dto"
	"myposcore/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FAQHandler struct {
	BaseHandler
	faqService *services.FAQService
}

func NewFAQHandler(faqService *services.FAQService) *FAQHandler {
	return &FAQHandler{
		faqService: faqService,
	}
}

// CreateFAQ godoc
// @Summary Create new FAQ
// @Description Create a new FAQ entry
// @Tags FAQ
// @Accept json
// @Produce json
// @Param request body dto.CreateFAQRequest true "FAQ request"
// @Success 200 {object} dto.FAQResponse
// @Router /api/faq [post]
func (h *FAQHandler) CreateFAQ(c *gin.Context) {
	var req dto.CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")
	req.CreatedBy = &currentUserID

	faq, err := h.faqService.CreateFAQ(req.Question, req.Answer, req.Category, req.Order, req.CreatedBy)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to create FAQ")
		return
	}

	var createdByName *string
	if faq.Creator != nil {
		name := faq.Creator.FullName
		createdByName = &name
	}

	response := dto.FAQResponse{
		ID:            faq.ID,
		Question:      faq.Question,
		Answer:        faq.Answer,
		Category:      faq.Category,
		Order:         faq.Order,
		IsActive:      faq.IsActive,
		CreatedAt:     faq.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     faq.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     faq.CreatedBy,
		CreatedByName: createdByName,
	}

	h.SuccessResponse(c, http.StatusOK, "FAQ created successfully", response)
}

// GetAllFAQ godoc
// @Summary Get all FAQs
// @Description Get all FAQ entries, optionally filtered by category and active status
// @Tags FAQ
// @Produce json
// @Param category query string false "Category filter"
// @Param active query bool false "Show only active FAQs"
// @Success 200 {array} dto.FAQResponse
// @Router /api/faq [get]
func (h *FAQHandler) GetAllFAQ(c *gin.Context) {
	category := c.Query("category")
	activeOnly := c.Query("active_only") == "true"

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	faqs, err := h.faqService.GetAllFAQ(categoryPtr, activeOnly)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve FAQs")
		return
	}

	var responses []dto.FAQResponse
	for _, faq := range faqs {
		var createdByName, updatedByName *string
		if faq.Creator != nil {
			name := faq.Creator.FullName
			createdByName = &name
		}
		if faq.Updater != nil {
			name := faq.Updater.FullName
			updatedByName = &name
		}

		responses = append(responses, dto.FAQResponse{
			ID:            faq.ID,
			Question:      faq.Question,
			Answer:        faq.Answer,
			Category:      faq.Category,
			Order:         faq.Order,
			IsActive:      faq.IsActive,
			CreatedAt:     faq.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     faq.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedBy:     faq.CreatedBy,
			CreatedByName: createdByName,
			UpdatedBy:     faq.UpdatedBy,
			UpdatedByName: updatedByName,
		})
	}

	h.SuccessResponse(c, http.StatusOK, "FAQs retrieved successfully", responses)
}

// GetFAQByID godoc
// @Summary Get FAQ by ID
// @Description Get specific FAQ by ID
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} dto.FAQResponse
// @Router /api/faq/{id} [get]
func (h *FAQHandler) GetFAQByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid FAQ ID")
		return
	}

	faq, err := h.faqService.GetFAQByID(uint(id))
	if err != nil {
		h.ErrorResponse(c, http.StatusNotFound, "FAQ not found")
		return
	}

	var createdByName, updatedByName *string
	if faq.Creator != nil {
		name := faq.Creator.FullName
		createdByName = &name
	}
	if faq.Updater != nil {
		name := faq.Updater.FullName
		updatedByName = &name
	}

	response := dto.FAQResponse{
		ID:            faq.ID,
		Question:      faq.Question,
		Answer:        faq.Answer,
		Category:      faq.Category,
		Order:         faq.Order,
		IsActive:      faq.IsActive,
		CreatedAt:     faq.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     faq.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     faq.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     faq.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	h.SuccessResponse(c, http.StatusOK, "FAQ retrieved successfully", response)
}

// UpdateFAQ godoc
// @Summary Update FAQ
// @Description Update an existing FAQ entry
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path int true "FAQ ID"
// @Param request body dto.UpdateFAQRequest true "Update FAQ request"
// @Success 200 {object} dto.FAQResponse
// @Router /api/faq/{id} [put]
func (h *FAQHandler) UpdateFAQ(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid FAQ ID")
		return
	}

	var req dto.UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")
	req.UpdatedBy = &currentUserID

	var question, answer, category *string
	if req.Question != "" {
		question = &req.Question
	}
	if req.Answer != "" {
		answer = &req.Answer
	}
	if req.Category != "" {
		category = &req.Category
	}

	faq, err := h.faqService.UpdateFAQ(uint(id), question, answer, category, req.Order, req.IsActive, req.UpdatedBy)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to update FAQ")
		return
	}

	// Reload to get audit info
	faq, _ = h.faqService.GetFAQByID(uint(id))

	var createdByName, updatedByName *string
	if faq.Creator != nil {
		name := faq.Creator.FullName
		createdByName = &name
	}
	if faq.Updater != nil {
		name := faq.Updater.FullName
		updatedByName = &name
	}

	response := dto.FAQResponse{
		ID:            faq.ID,
		Question:      faq.Question,
		Answer:        faq.Answer,
		Category:      faq.Category,
		Order:         faq.Order,
		IsActive:      faq.IsActive,
		CreatedAt:     faq.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     faq.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:     faq.CreatedBy,
		CreatedByName: createdByName,
		UpdatedBy:     faq.UpdatedBy,
		UpdatedByName: updatedByName,
	}

	h.SuccessResponse(c, http.StatusOK, "FAQ updated successfully", response)
}

// DeleteFAQ godoc
// @Summary Delete FAQ
// @Description Delete a FAQ entry
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/faq/{id} [delete]
func (h *FAQHandler) DeleteFAQ(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid FAQ ID")
		return
	}

	// Get current user ID from context
	currentUserID := c.GetUint("user_id")

	if err := h.faqService.DeleteFAQ(uint(id), &currentUserID); err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete FAQ")
		return
	}

	h.SuccessResponse(c, http.StatusOK, "FAQ deleted successfully", nil)
}
