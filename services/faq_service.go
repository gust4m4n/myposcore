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

func (s *FAQService) CreateFAQ(question, answer, category string, order int) (*models.FAQ, error) {
	faq := &models.FAQ{
		Question: question,
		Answer:   answer,
		Category: category,
		Order:    order,
		IsActive: true,
	}

	if err := s.db.Create(faq).Error; err != nil {
		return nil, err
	}

	return faq, nil
}

func (s *FAQService) GetFAQByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.First(&faq, id).Error; err != nil {
		return nil, err
	}
	return &faq, nil
}

func (s *FAQService) GetAllFAQ(category *string, activeOnly bool) ([]models.FAQ, error) {
	var faqs []models.FAQ
	query := s.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	if err := query.Order("`order` ASC, created_at DESC").Find(&faqs).Error; err != nil {
		return nil, err
	}

	return faqs, nil
}

func (s *FAQService) UpdateFAQ(id uint, question, answer, category *string, order *int, isActive *bool) (*models.FAQ, error) {
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

	if err := s.db.Model(&faq).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &faq, nil
}

func (s *FAQService) DeleteFAQ(id uint) error {
	return s.db.Delete(&models.FAQ{}, id).Error
}
