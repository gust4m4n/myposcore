package handlers

import (
	"myposcore/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	*BaseHandler
}

func NewProfileHandler(cfg *config.Config) *ProfileHandler {
	return &ProfileHandler{
		BaseHandler: NewBaseHandler(cfg),
	}
}

func (h *ProfileHandler) Handle(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"tenant_id": tenantID,
		"username":  username,
	})
}
