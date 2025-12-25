package database

import (
	"fmt"
	"log"
	"myposcore/config"
	"myposcore/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error

	// Configure logger based on environment
	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Warn
	}

	dsn := cfg.GetDSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.Tenant{},
		&models.Branch{},
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.TermsAndConditions{},
		&models.FAQ{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed")

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
