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
	// Initialize services
	orderService := services.NewOrderService(database.DB)
	paymentService := services.NewPaymentService(database.DB)
	tncService := services.NewTnCService(database.DB)
	faqService := services.NewFAQService(database.DB)
	categoryService := services.NewCategoryService(database.DB)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	registerHandler := handlers.NewRegisterHandler(cfg)
	loginHandler := handlers.NewLoginHandler(cfg)
	profileHandler := handlers.NewProfileHandler(cfg)
	changePasswordHandler := handlers.NewChangePasswordHandler(cfg)
	pinHandler := handlers.NewPINHandler(cfg)
	productHandler := handlers.NewProductHandler(cfg)
	orderHandler := handlers.NewOrderHandler(cfg, orderService)
	paymentHandler := handlers.NewPaymentHandler(cfg, paymentService)
	tncHandler := handlers.NewTnCHandler(tncService)
	faqHandler := handlers.NewFAQHandler(faqService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	userHandler := handlers.NewUserHandler(cfg)
	devHandler := handlers.NewDevHandler(cfg)
	superAdminHandler := handlers.NewSuperAdminHandler(cfg)

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
			auth.POST("/register", registerHandler.Handle)
			auth.POST("/login", loginHandler.Handle)
		}

		// Public routes
		public := v1.Group("")
		{
			// TnC routes (public)
			public.GET("/tnc/active", tncHandler.GetActiveTnC)
			public.GET("/tnc", tncHandler.GetAllTnC)
			public.GET("/tnc/:id", tncHandler.GetTnCByID)

			// FAQ routes (public)
			public.GET("/faq", faqHandler.GetAllFAQ)
			public.GET("/faq/:id", faqHandler.GetFAQByID)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.Use(middleware.TenantMiddleware())
		{
			protected.GET("/profile", profileHandler.Handle)
			protected.PUT("/change-password", changePasswordHandler.Handle)

			// PIN routes
			protected.POST("/pin/create", pinHandler.CreatePIN)
			protected.PUT("/pin/change", pinHandler.ChangePIN)
			protected.GET("/pin/check", pinHandler.CheckPIN)

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

			// Order routes
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.ListOrders)
			protected.GET("/orders/:id", orderHandler.GetOrder)
			protected.GET("/orders/:id/payments", paymentHandler.GetPaymentsByOrder)

			// Payment routes
			protected.POST("/payments", paymentHandler.CreatePayment)
			protected.GET("/payments", paymentHandler.ListPayments)
			protected.GET("/payments/:id", paymentHandler.GetPayment)

			// User routes
			protected.GET("/users", userHandler.ListUsers)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.POST("/users", userHandler.CreateUser)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)
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

			// TnC management routes (superadmin only)
			superadmin.POST("/tnc", tncHandler.CreateTnC)
			superadmin.PUT("/tnc/:id", tncHandler.UpdateTnC)
			superadmin.DELETE("/tnc/:id", tncHandler.DeleteTnC)

			// FAQ management routes (superadmin only)
			superadmin.POST("/faq", faqHandler.CreateFAQ)
			superadmin.PUT("/faq/:id", faqHandler.UpdateFAQ)
			superadmin.DELETE("/faq/:id", faqHandler.DeleteFAQ)
		}
	}
}
