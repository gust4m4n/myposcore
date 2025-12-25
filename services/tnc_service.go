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

func (s *TnCService) CreateTnC(title, content, version string) (*models.TermsAndConditions, error) {
	tnc := &models.TermsAndConditions{
		Title:    title,
		Content:  content,
		Version:  version,
		IsActive: true,
	}

	if err := s.db.Create(tnc).Error; err != nil {
		return nil, err
	}

	return tnc, nil
}

func (s *TnCService) GetTnCByID(id uint) (*models.TermsAndConditions, error) {
	var tnc models.TermsAndConditions
	if err := s.db.First(&tnc, id).Error; err != nil {
		return nil, err
	}
	return &tnc, nil
}

func (s *TnCService) GetActiveTnC() (*models.TermsAndConditions, error) {
	var tnc models.TermsAndConditions
	if err := s.db.Where("is_active = ?", true).Order("created_at DESC").First(&tnc).Error; err != nil {
		return nil, err
	}
	return &tnc, nil
}

func (s *TnCService) GetAllTnC() ([]models.TermsAndConditions, error) {
	var tncs []models.TermsAndConditions
	if err := s.db.Order("created_at DESC").Find(&tncs).Error; err != nil {
		return nil, err
	}
	return tncs, nil
}

func (s *TnCService) UpdateTnC(id uint, title, content, version *string, isActive *bool) (*models.TermsAndConditions, error) {
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

	if err := s.db.Model(&tnc).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &tnc, nil
}

func (s *TnCService) DeleteTnC(id uint) error {
	return s.db.Delete(&models.TermsAndConditions{}, id).Error
}
