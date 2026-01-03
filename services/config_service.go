package services

import (
	"errors"
	"myposcore/models"

	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

// SetConfig creates or updates a config key-value pair
func (s *ConfigService) SetConfig(key, value string) error {
	config := models.Config{
		Key:   key,
		Value: value,
	}

	// Use UPSERT pattern: try to create, if key exists, update
	result := s.db.Where("key = ?", key).Assign(models.Config{Value: value}).FirstOrCreate(&config)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetConfig retrieves a config value by key
func (s *ConfigService) GetConfig(key string) (string, error) {
	var config models.Config
	result := s.db.Where("key = ?", key).First(&config)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("config key not found")
		}
		return "", result.Error
	}

	return config.Value, nil
}
