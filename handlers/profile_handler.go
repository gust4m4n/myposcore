package handlers

import (
	"myposcore/config"
	"myposcore/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	*BaseHandler
	authService *services.AuthService
}

func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler: NewBaseHandler(cfg),
		authService: services.NewAuthService(),
	}
}

func (h *ProfileHandler) Handle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
