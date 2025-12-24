package routes

import (
	"myposcore/config"
	"myposcore/handlers"
	"myposcore/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	registerHandler := handlers.NewRegisterHandler(cfg)
	loginHandler := handlers.NewLoginHandler(cfg)
	profileHandler := handlers.NewProfileHandler(cfg)
	changePasswordHandler := handlers.NewChangePasswordHandler(cfg)
	superAdminHandler := handlers.NewSuperAdminHandler(cfg)

	// Health check
	router.GET("/health", healthHandler.Handle)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", registerHandler.Handle)
			auth.POST("/login", loginHandler.Handle)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.Use(middleware.TenantMiddleware())
		{
			protected.GET("/profile", profileHandler.Handle)
			protected.PUT("/change-password", changePasswordHandler.Handle)
		}

		// Superadmin routes
		superadmin := v1.Group("/superadmin")
		superadmin.Use(middleware.AuthMiddleware(cfg))
		superadmin.Use(middleware.SuperAdminMiddleware(cfg))
		{
			superadmin.GET("/dashboard", superAdminHandler.Dashboard)
			superadmin.GET("/tenants", superAdminHandler.ListTenants)
			superadmin.POST("/tenants", superAdminHandler.CreateTenant)
			superadmin.GET("/tenants/:tenant_id/branches", superAdminHandler.ListBranches)
			superadmin.GET("/branches/:branch_id/users", superAdminHandler.ListUsers)
		}
	}
}
