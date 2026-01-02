package handlers

import (
	"myposcore/config"
	"myposcore/utils"

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
	code := 1 // Default error code
	switch statusCode {
	case 400:
		code = 1
	case 401:
		code = 2
	case 403:
		code = 3
	case 404:
		code = 4
	case 500:
		code = 5
	case 409:
		code = 6
	case 422:
		code = 7
	}
	utils.Error(c, statusCode, code, message)
}
