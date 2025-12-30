package services

import (
	"encoding/json"
	"errors"
	"myposcore/models"
	"time"

	"gorm.io/gorm"
)

type AuditTrailService struct {
	db *gorm.DB
}

func NewAuditTrailService(db *gorm.DB) *AuditTrailService {
	return &AuditTrailService{db: db}
}

// CreateAuditTrail creates a new audit trail entry
func (s *AuditTrailService) CreateAuditTrail(tenantID, branchID *uint, userID uint, entityType string, entityID uint, action string, changes interface{}, ipAddress, userAgent string) error {
	var changesJSON string
	if changes != nil {
		changesBytes, err := json.Marshal(changes)
		if err != nil {
			return err
		}
		changesJSON = string(changesBytes)
	}

	auditTrail := &models.AuditTrail{
		TenantID:   tenantID,
		BranchID:   branchID,
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Changes:    changesJSON,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	return s.db.Create(auditTrail).Error
}

// ListAuditTrails retrieves audit trails with filters and pagination
func (s *AuditTrailService) ListAuditTrails(tenantID *uint, userID, entityID *uint, entityType, action *string, dateFrom, dateTo *time.Time, page, limit int) ([]models.AuditTrail, int64, error) {
	var auditTrails []models.AuditTrail
	var total int64

	query := s.db.Model(&models.AuditTrail{}).Preload("User")

	// Apply filters
	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if entityType != nil && *entityType != "" {
		query = query.Where("entity_type = ?", *entityType)
	}

	if entityID != nil {
		query = query.Where("entity_id = ?", *entityID)
	}

	if action != nil && *action != "" {
		query = query.Where("action = ?", *action)
	}

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}

	if dateTo != nil {
		// Add one day to include the entire end date
		endDate := dateTo.Add(24 * time.Hour)
		query = query.Where("created_at < ?", endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditTrails).Error; err != nil {
		return nil, 0, err
	}

	return auditTrails, total, nil
}

// GetAuditTrailByID retrieves a single audit trail by ID
func (s *AuditTrailService) GetAuditTrailByID(id uint, tenantID *uint) (*models.AuditTrail, error) {
	var auditTrail models.AuditTrail
	query := s.db.Preload("User").Where("id = ?", id)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.First(&auditTrail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit trail not found")
		}
		return nil, err
	}

	return &auditTrail, nil
}

// GetEntityAuditHistory retrieves audit history for a specific entity
func (s *AuditTrailService) GetEntityAuditHistory(tenantID *uint, entityType string, entityID uint, page, limit int) ([]models.AuditTrail, int64, error) {
	var auditTrails []models.AuditTrail
	var total int64

	query := s.db.Model(&models.AuditTrail{}).Preload("User").
		Where("entity_type = ? AND entity_id = ?", entityType, entityID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditTrails).Error; err != nil {
		return nil, 0, err
	}

	return auditTrails, total, nil
}

// GetUserActivityLog retrieves all activities performed by a specific user
func (s *AuditTrailService) GetUserActivityLog(tenantID *uint, userID uint, page, limit int) ([]models.AuditTrail, int64, error) {
	var auditTrails []models.AuditTrail
	var total int64

	query := s.db.Model(&models.AuditTrail{}).Preload("User").
		Where("user_id = ?", userID)

	if tenantID != nil {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditTrails).Error; err != nil {
		return nil, 0, err
	}

	return auditTrails, total, nil
}

// DeleteOldAuditTrails deletes audit trails older than the specified days (for cleanup/maintenance)
func (s *AuditTrailService) DeleteOldAuditTrails(daysOld int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	result := s.db.Where("created_at < ?", cutoffDate).Delete(&models.AuditTrail{})
	return result.RowsAffected, result.Error
}
