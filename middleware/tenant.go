package middleware

import (
	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tenant ID is set by auth middleware from JWT token
		// This middleware is kept for future extensions if needed
		c.Next()
	}
}
