package services

import (
	"encoding/json"
	"myposcore/models"

	"gorm.io/gorm"
)

type TnCService struct {
	db                *gorm.DB
	auditTrailService *AuditTrailService
}

func NewTnCService(db *gorm.DB, auditTrailService *AuditTrailService) *TnCService {
	return &TnCService{
		db:                db,
		auditTrailService: auditTrailService,
	}
}

func (s *TnCService) CreateTnC(title, content, version string, createdBy *uint) (*models.TermsAndConditions, error) {
	tnc := &models.TermsAndConditions{
		Title:     title,
		Content:   content,
		Version:   version,
		IsActive:  true,
		CreatedBy: createdBy,
	}

	if err := s.db.Create(tnc).Error; err != nil {
		return nil, err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"title":     tnc.Title,
		"content":   tnc.Content,
		"version":   tnc.Version,
		"is_active": tnc.IsActive,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if createdBy != nil {
		userID = *createdBy
	}
	_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "tnc", tnc.ID, "create", changesMap, "", "")

	return tnc, nil
}

func (s *TnCService) GetTnCByID(id uint) (*models.TermsAndConditions, error) {
	var tnc models.TermsAndConditions
	if err := s.db.Preload("Creator").Preload("Updater").First(&tnc, id).Error; err != nil {
		return nil, err
	}
	return &tnc, nil
}

func (s *TnCService) GetActiveTnC() (*models.TermsAndConditions, error) {
	var tnc models.TermsAndConditions
	if err := s.db.Preload("Creator").Preload("Updater").Where("is_active = ?", true).Order("created_at DESC").First(&tnc).Error; err != nil {
		return nil, err
	}
	return &tnc, nil
}

func (s *TnCService) GetAllTnC() ([]models.TermsAndConditions, error) {
	var tncs []models.TermsAndConditions
	if err := s.db.Preload("Creator").Preload("Updater").Order("created_at DESC").Find(&tncs).Error; err != nil {
		return nil, err
	}
	return tncs, nil
}

func (s *TnCService) UpdateTnC(id uint, title, content, version *string, isActive *bool, updatedBy *uint) (*models.TermsAndConditions, error) {
	var tnc models.TermsAndConditions
	if err := s.db.First(&tnc, id).Error; err != nil {
		return nil, err
	}

	// Save old values for audit trail
	oldValues := map[string]interface{}{
		"title":     tnc.Title,
		"content":   tnc.Content,
		"version":   tnc.Version,
		"is_active": tnc.IsActive,
	}

	updates := make(map[string]interface{})
	if title != nil {
		updates["title"] = *title
	}
	if content != nil {
		updates["content"] = *content
	}
	if version != nil {
		updates["version"] = *version
	}
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}

	if err := s.db.Model(&tnc).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Create audit trail with old/new values
	if len(updates) > 0 {
		changes := make(map[string]interface{})
		for key, newVal := range updates {
			if key != "updated_by" {
				if oldVal, exists := oldValues[key]; exists {
					changes[key] = map[string]interface{}{
						"old": oldVal,
						"new": newVal,
					}
				}
			}
		}
		changesJSON, _ := json.Marshal(changes)
		var changesMap map[string]interface{}
		_ = json.Unmarshal(changesJSON, &changesMap)
		var userID uint
		if updatedBy != nil {
			userID = *updatedBy
		}
		_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "tnc", id, "update", changesMap, "", "")
	}

	return &tnc, nil
}

func (s *TnCService) DeleteTnC(id uint, deletedBy *uint) error {
	// Get TnC for audit trail
	var tnc models.TermsAndConditions
	if err := s.db.First(&tnc, id).Error; err != nil {
		return err
	}

	if deletedBy != nil {
		if err := s.db.Model(&models.TermsAndConditions{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}

	if err := s.db.Delete(&models.TermsAndConditions{}, id).Error; err != nil {
		return err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"title":   tnc.Title,
		"version": tnc.Version,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if deletedBy != nil {
		userID = *deletedBy
	}
	_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "tnc", id, "delete", changesMap, "", "")

	return nil
}
