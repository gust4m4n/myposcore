package main

import (
	"log"
	"myposcore/config"
	"myposcore/database"
	"myposcore/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Setup Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
