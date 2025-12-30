package services

import (
	"myposcore/models"

	"gorm.io/gorm"
)

type FAQService struct {
	db *gorm.DB
}

func NewFAQService(db *gorm.DB) *FAQService {
	return &FAQService{db: db}
}

func (s *FAQService) CreateFAQ(question, answer, category string, order int, createdBy *uint) (*models.FAQ, error) {
	faq := &models.FAQ{
		Question:  question,
		Answer:    answer,
		Category:  category,
		Order:     order,
		IsActive:  true,
		CreatedBy: createdBy,
	}

	if err := s.db.Create(faq).Error; err != nil {
		return nil, err
	}

	return faq, nil
}

func (s *FAQService) GetFAQByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.Preload("Creator").Preload("Updater").First(&faq, id).Error; err != nil {
		return nil, err
	}
	return &faq, nil
}

func (s *FAQService) GetAllFAQ(category *string, activeOnly bool) ([]models.FAQ, error) {
	var faqs []models.FAQ
	query := s.db.Preload("Creator").Preload("Updater")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	if err := query.Order("\"order\" ASC, created_at DESC").Find(&faqs).Error; err != nil {
		return nil, err
	}

	return faqs, nil
}

func (s *FAQService) UpdateFAQ(id uint, question, answer, category *string, order *int, isActive *bool, updatedBy *uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.First(&faq, id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if question != nil {
		updates["question"] = *question
	}
	if answer != nil {
		updates["answer"] = *answer
	}
	if category != nil {
		updates["category"] = *category
	}
	if order != nil {
		updates["order"] = *order
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

	return &faq, nil
}

func (s *FAQService) DeleteFAQ(id uint, deletedBy *uint) error {
	if deletedBy != nil {
		if err := s.db.Model(&models.FAQ{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
	}
	return s.db.Delete(&models.FAQ{}, id).Error
}
