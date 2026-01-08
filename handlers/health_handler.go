package handlers

import (
	"fmt"
	"myposcore/config"
	"myposcore/database"
	"myposcore/utils"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config: cfg,
	}
}

func (h *HealthHandler) Handle(c *gin.Context) {
	// Check database connection
	dbStatus := "connected"
	var dbVersion string
	var dbStats gin.H

	db, err := database.DB.DB()
	if err != nil || db.Ping() != nil {
		dbStatus = "disconnected"
	} else {
		// Get database version
		database.DB.Raw("SELECT version()").Scan(&dbVersion)

		// Get connection pool stats
		stats := db.Stats()
		dbStats = gin.H{
			"max_open_connections": stats.MaxOpenConnections,
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
			"wait_count":           stats.WaitCount,
			"wait_duration":        stats.WaitDuration.String(),
		}
	}

	// Calculate uptime
	uptime := time.Since(h.config.StartupTime)
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	uptimeStr := fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)

	// Get memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Get hostname
	hostname, _ := os.Hostname()

	// Get environment
	environment := os.Getenv("GIN_MODE")
	if environment == "" {
		environment = "development"
	}

	response := gin.H{
		"status":       "ok",
		"version":      config.AppVersion,
		"startup_time": h.config.StartupTime.Format("2006-01-02 15:04:05"),
		"uptime":       uptimeStr,
		"uptime_details": gin.H{
			"days":    days,
			"hours":   hours,
			"minutes": minutes,
			"seconds": seconds,
		},
		"system": gin.H{
			"hostname":     hostname,
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
			"go_version":   runtime.Version(),
			"cpu_cores":    runtime.NumCPU(),
		},
		"resources": gin.H{
			"goroutines": runtime.NumGoroutine(),
			"memory": gin.H{
				"allocated_mb":   fmt.Sprintf("%.2f", float64(memStats.Alloc)/1024/1024),
				"total_alloc_mb": fmt.Sprintf("%.2f", float64(memStats.TotalAlloc)/1024/1024),
				"system_mb":      fmt.Sprintf("%.2f", float64(memStats.Sys)/1024/1024),
				"gc_cycles":      memStats.NumGC,
			},
		},
		"database": gin.H{
			"status":  dbStatus,
			"version": dbVersion,
			"pool":    dbStats,
		},
		"server": gin.H{
			"port":        h.config.ServerPort,
			"environment": environment,
		},
	}

	utils.Success(c, "Health check successful", response)
}
