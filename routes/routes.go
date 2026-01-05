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
	branchService := services.NewSuperAdminBranchService()

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
	superAdminHandler := handlers.NewSuperAdminHandler(cfg)
	auditTrailHandler := handlers.NewAuditTrailHandler(auditTrailService)
	configHandler := handlers.NewConfigHandler(configService)
	branchHandler := handlers.NewBranchHandler(cfg, branchService)
	tenantHandler := handlers.NewTenantHandler(cfg)

	// Health check
	router.GET("/health", healthHandler.Handle)

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

			// Branch routes
			protected.GET("/branches", branchHandler.GetBranches)
			protected.GET("/branches/:id", branchHandler.GetBranch)
			protected.POST("/branches", branchHandler.CreateBranch)
			protected.PUT("/branches/:id", branchHandler.UpdateBranch)
			protected.DELETE("/branches/:id", branchHandler.DeleteBranch)
			protected.GET("/branches/:id/users", branchHandler.GetBranchUsers)

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

			// Tenant routes
			protected.GET("/tenants", tenantHandler.ListTenants)
			protected.GET("/tenants/:id", tenantHandler.GetTenant)
			protected.POST("/tenants", tenantHandler.CreateTenant)
			protected.PUT("/tenants/:id", tenantHandler.UpdateTenant)
			protected.DELETE("/tenants/:id", tenantHandler.DeleteTenant)

			// Audit trail routes
			protected.GET("/audit-trails", auditTrailHandler.ListAuditTrails)
			protected.GET("/audit-trails/user/:user_id", auditTrailHandler.GetUserActivityLog)
			protected.GET("/audit-trails/entity/:entity_type/:entity_id", auditTrailHandler.GetEntityAuditHistory)
			protected.GET("/audit-trails/:id", auditTrailHandler.GetAuditTrailByID)

			// Dashboard route
			protected.GET("/dashboard", superAdminHandler.Dashboard)

			// FAQ management routes
			protected.POST("/faq", faqHandler.CreateFAQ)
			protected.PUT("/faq/:id", faqHandler.UpdateFAQ)
			protected.DELETE("/faq/:id", faqHandler.DeleteFAQ)
		}
	}
}
