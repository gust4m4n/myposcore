package middleware

import (
	"myposcore/config"
	"myposcore/database"
	"myposcore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Get user from database to check role
		var user models.User
		if err := database.GetDB().First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if user has superadmin role
		if user.Role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Superadmin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
