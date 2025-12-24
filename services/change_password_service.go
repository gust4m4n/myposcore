package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type ChangePasswordService struct {
	db *gorm.DB
}

func NewChangePasswordService() *ChangePasswordService {
	return &ChangePasswordService{
		db: database.GetDB(),
	}
}

func (s *ChangePasswordService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	// Get user from database
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify old password
	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := s.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
