package utils

import (
	"myposcore/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFullImageURL converts a relative image path to a full URL
// Returns empty string if imagePath is empty
func GetFullImageURL(imagePath string) string {
	if imagePath == "" {
		return ""
	}
	return config.GetBaseURL() + imagePath
}

// APIResponse represents the standard API response structure
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a successful response with code 0
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// SuccessWithoutData sends a successful response without data
func SuccessWithoutData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: message,
	})
}

// Error sends an error response with a non-zero code
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, APIResponse{
		Code:    code,
		Message: message,
	})
}

// BadRequest sends a 400 error with code 1
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 1, message)
}

// Unauthorized sends a 401 error with code 2
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 2, message)
}

// Forbidden sends a 403 error with code 3
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 3, message)
}

// NotFound sends a 404 error with code 4
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 4, message)
}

// InternalError sends a 500 error with code 5
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 5, message)
}

// Conflict sends a 409 error with code 6
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, 6, message)
}

// UnprocessableEntity sends a 422 error with code 7
func UnprocessableEntity(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, 7, message)
}
