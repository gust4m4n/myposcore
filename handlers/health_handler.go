package handlers

import (
	"myposcore/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

func (h *HealthHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
