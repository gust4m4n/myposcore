package middleware

import (
	"github.com/gin-gonic/gin"
)

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant from header or from JWT token
		tenantCode := c.GetHeader("X-Tenant-Code")

		if tenantCode != "" {
			c.Set("tenant_code", tenantCode)
		}

		// If tenant_id is already set by auth middleware, use that
		if tenantID, exists := c.Get("tenant_id"); exists {
			c.Set("tenant_id", tenantID)
		}

		c.Next()
	}
}
