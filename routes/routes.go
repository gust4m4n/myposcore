package routes

import (
	"myposcore/config"
	"myposcore/database"
	"myposcore/handlers"
	"myposcore/middleware"
	"myposcore/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize audit trail service first (no dependencies)
	auditTrailService := services.NewAuditTrailService(database.DB)

	// Initialize services with audit trail dependency
	orderService := services.NewOrderService(database.DB, auditTrailService)
	paymentService := services.NewPaymentService(database.DB, auditTrailService)
	faqService := services.NewFAQService(database.DB, auditTrailService)
	categoryService := services.NewCategoryService(database.DB, auditTrailService)
	userService := services.NewUserService(auditTrailService)
	productService := services.NewProductService(auditTrailService)
	configService := services.NewConfigService(database.DB)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	loginHandler := handlers.NewLoginHandler(cfg, auditTrailService)
	logoutHandler := handlers.NewLogoutHandler(cfg, auditTrailService)
	profileHandler := handlers.NewProfileHandler(cfg)
	changePasswordHandler := handlers.NewChangePasswordHandler(cfg, auditTrailService)
	adminChangePasswordHandler := handlers.NewAdminChangePasswordHandler(cfg, auditTrailService)
	adminChangePINHandler := handlers.NewAdminChangePINHandler(cfg, auditTrailService)
	pinHandler := handlers.NewPINHandler(cfg, auditTrailService)
	productHandler := handlers.NewProductHandler(cfg, productService)
	orderHandler := handlers.NewOrderHandler(cfg, orderService)
	paymentHandler := handlers.NewPaymentHandler(cfg, paymentService)
	tncHandler := handlers.NewTnCHandler(configService)
	faqHandler := handlers.NewFAQHandler(faqService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	userHandler := handlers.NewUserHandler(cfg, userService)
	devHandler := handlers.NewDevHandler(cfg)
	superAdminHandler := handlers.NewSuperAdminHandler(cfg)
	auditTrailHandler := handlers.NewAuditTrailHandler(auditTrailService)
	configHandler := handlers.NewConfigHandler(configService)

	// Health check
	router.GET("/health", healthHandler.Handle)

	// Dev routes (public - no authentication required)
	dev := router.Group("/dev")
	{
		dev.GET("/tenants", devHandler.ListTenants)
		dev.GET("/tenants/:tenant_id/branches", devHandler.ListBranchesByTenant)
	}

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginHandler.Handle)
		}

		// Public routes
		public := v1.Group("")
		{
			// TnC route (public)
			public.GET("/tnc", tncHandler.GetTnC)

			// FAQ routes (public)
			public.GET("/faq", faqHandler.GetAllFAQ)
			public.GET("/faq/:id", faqHandler.GetFAQByID)

			// Config routes (public)
			public.POST("/config/set", configHandler.SetConfig)
			public.GET("/config/get/:key", configHandler.GetConfig)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.Use(middleware.TenantMiddleware())
		{
			// Auth routes
			protected.POST("/logout", logoutHandler.Handle)
			protected.GET("/profile", profileHandler.Handle)
			protected.PUT("/profile", profileHandler.UpdateProfile)
			protected.PUT("/change-password", changePasswordHandler.Handle)

			// Admin change password (for higher roles to change lower roles password)
			protected.PUT("/admin/change-password", adminChangePasswordHandler.Handle)

			// Profile image routes
			protected.POST("/profile/photo", profileHandler.UploadProfileImage)
			protected.DELETE("/profile/photo", profileHandler.DeleteProfileImage)

			// PIN routes
			protected.POST("/pin/create", pinHandler.CreatePIN)
			protected.PUT("/pin/change", pinHandler.ChangePIN)
			protected.GET("/pin/check", pinHandler.CheckPIN)

			// Admin change PIN (for higher roles to change lower roles PIN)
			protected.PUT("/admin/change-pin", adminChangePINHandler.Handle)

			// Category routes
			protected.GET("/categories", categoryHandler.ListCategories)
			protected.GET("/categories/:id", categoryHandler.GetCategory)
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.PUT("/categories/:id", categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			// Product routes
			protected.GET("/products/categories", productHandler.GetCategories)
			protected.GET("/products", productHandler.ListProducts)
			protected.GET("/products/:id", productHandler.GetProduct)
			protected.POST("/products", productHandler.CreateProduct)
			protected.PUT("/products/:id", productHandler.UpdateProduct)
			protected.DELETE("/products/:id", productHandler.DeleteProduct)
			protected.POST("/products/:id/photo", productHandler.UploadProductImage)
			protected.DELETE("/products/:id/photo", productHandler.DeleteProductImage)

			// Order routes
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.ListOrders)
			protected.GET("/orders/:id", orderHandler.GetOrder)
			protected.GET("/orders/:id/payments", paymentHandler.GetPaymentsByOrder)

			// Payment routes
			protected.POST("/payments", paymentHandler.CreatePayment)
			protected.GET("/payments", paymentHandler.ListPayments)
			protected.GET("/payments/:id", paymentHandler.GetPayment)
			protected.GET("/payments/performance", paymentHandler.GetPaymentPerformance)

			// User routes
			protected.GET("/users", userHandler.ListUsers)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.POST("/users", userHandler.CreateUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)

			// Audit trail routes
			protected.GET("/audit-trails", auditTrailHandler.ListAuditTrails)
			protected.GET("/audit-trails/user/:user_id", auditTrailHandler.GetUserActivityLog)
			protected.GET("/audit-trails/entity/:entity_type/:entity_id", auditTrailHandler.GetEntityAuditHistory)
			protected.GET("/audit-trails/:id", auditTrailHandler.GetAuditTrailByID)
		}

		// Superadmin routes - now accessible by all roles
		superadmin := v1.Group("/superadmin")
		superadmin.Use(middleware.AuthMiddleware(cfg))
		{
			superadmin.GET("/dashboard", superAdminHandler.Dashboard)
			superadmin.GET("/tenants", superAdminHandler.ListTenants)
			superadmin.POST("/tenants", superAdminHandler.CreateTenant)
			superadmin.PUT("/tenants/:tenant_id", superAdminHandler.UpdateTenant)
			superadmin.DELETE("/tenants/:tenant_id", superAdminHandler.DeleteTenant)
			superadmin.GET("/tenants/:tenant_id/branches", superAdminHandler.ListBranches)
			superadmin.POST("/branches", superAdminHandler.CreateBranch)
			superadmin.PUT("/branches/:branch_id", superAdminHandler.UpdateBranch)
			superadmin.DELETE("/branches/:branch_id", superAdminHandler.DeleteBranch)
			superadmin.GET("/branches/:branch_id/users", superAdminHandler.ListUsers)

			// FAQ management routes (superadmin only)
			superadmin.POST("/faq", faqHandler.CreateFAQ)
			superadmin.PUT("/faq/:id", faqHandler.UpdateFAQ)
			superadmin.DELETE("/faq/:id", faqHandler.DeleteFAQ)
		}
	}
}
