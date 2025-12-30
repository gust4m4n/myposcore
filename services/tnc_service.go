package services

import (
	"myposcore/models"

	"gorm.io/gorm"
)

type TnCService struct {
	db *gorm.DB
}

func NewTnCService(db *gorm.DB) *TnCService {
	return &TnCService{db: db}
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

	return &tnc, nil
}

func (s *TnCService) DeleteTnC(id uint, deletedBy *uint) error {
	if deletedBy != nil {
		if err := s.db.Model(&models.TermsAndConditions{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}
	return s.db.Delete(&models.TermsAndConditions{}, id).Error
}
