package database

import (
	"fmt"
	"log"
	"myposcore/config"
	"myposcore/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	var err error

	dsn := cfg.GetDSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
