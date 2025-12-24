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
		}
	}
}
