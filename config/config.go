package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const AppVersion = "1.0.0"

type Config struct {
	ServerPort  string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	StartupTime time.Time
}

func LoadConfig() (*Config, error) {
	// Load .env file if exists
	godotenv.Load()

	config := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "myposcore"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		StartupTime: time.Now(),
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
