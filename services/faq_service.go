package services

import (
	"encoding/json"
	"myposcore/models"

	"gorm.io/gorm"
)

type FAQService struct {
	db                *gorm.DB
	auditTrailService *AuditTrailService
}

func NewFAQService(db *gorm.DB, auditTrailService *AuditTrailService) *FAQService {
	return &FAQService{
		db:                db,
		auditTrailService: auditTrailService,
	}
}

func (s *FAQService) CreateFAQ(question, answer string, createdBy *uint) (*models.FAQ, error) {
	faq := &models.FAQ{
		Question:  question,
		Answer:    answer,
		IsActive:  true,
		CreatedBy: createdBy,
	}

	if err := s.db.Create(faq).Error; err != nil {
		return nil, err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"question":  faq.Question,
		"answer":    faq.Answer,
		"is_active": faq.IsActive,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if createdBy != nil {
		userID = *createdBy
	}
	_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "faq", faq.ID, "create", changesMap, "", "")

	return faq, nil
}

func (s *FAQService) GetFAQByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.Preload("Creator").Preload("Updater").First(&faq, id).Error; err != nil {
		return nil, err
	}
	return &faq, nil
}

func (s *FAQService) GetAllFAQ(activeOnly bool) ([]models.FAQ, error) {
	var faqs []models.FAQ
	query := s.db.Preload("Creator").Preload("Updater")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("created_at DESC").Find(&faqs).Error; err != nil {
		return nil, err
	}

	return faqs, nil
}

func (s *FAQService) UpdateFAQ(id uint, question, answer *string, isActive *bool, updatedBy *uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.First(&faq, id).Error; err != nil {
		return nil, err
	}

	// Save old values for audit trail
	oldValues := map[string]interface{}{
		"question":  faq.Question,
		"answer":    faq.Answer,
		"is_active": faq.IsActive,
	}

	updates := make(map[string]interface{})
	if question != nil {
		updates["question"] = *question
	}
	if answer != nil {
		updates["answer"] = *answer
	}
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}

	if err := s.db.Model(&faq).Updates(updates).Error; err != nil {
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
		_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "faq", id, "update", changesMap, "", "")
	}

	return &faq, nil
}

func (s *FAQService) DeleteFAQ(id uint, deletedBy *uint) error {
	// Get FAQ for audit trail
	var faq models.FAQ
	if err := s.db.First(&faq, id).Error; err != nil {
		return err
	}

	if deletedBy != nil {
		if err := s.db.Model(&models.FAQ{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}

	if err := s.db.Delete(&models.FAQ{}, id).Error; err != nil {
		return err
	}

	// Create audit trail
	changes := map[string]interface{}{
		"question": faq.Question,
	}
	changesJSON, _ := json.Marshal(changes)
	var changesMap map[string]interface{}
	_ = json.Unmarshal(changesJSON, &changesMap)
	var userID uint
	if deletedBy != nil {
		userID = *deletedBy
	}
	_ = s.auditTrailService.CreateAuditTrail(nil, nil, userID, "faq", id, "delete", changesMap, "", "")

	return nil
}
