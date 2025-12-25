package handlers

import (
	"myposcore/config"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	config *config.Config
}

func NewBaseHandler(cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		config: cfg,
	}
}

// SuccessResponse sends a standardized success response
func (h *BaseHandler) SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := gin.H{
		"message": message,
	}
	if data != nil {
		response["data"] = data
	}
	c.JSON(statusCode, response)
}

// ErrorResponse sends a standardized error response
func (h *BaseHandler) ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}
