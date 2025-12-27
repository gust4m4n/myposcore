package services

import (
	"errors"
	"myposcore/database"
	"myposcore/dto"
	"myposcore/models"
	"myposcore/utils"

	"gorm.io/gorm"
)

type PINService struct {
	db *gorm.DB
}

func NewPINService() *PINService {
	return &PINService{
		db: database.GetDB(),
	}
}

func (s *PINService) CreatePIN(userID uint, req dto.CreatePINRequest) error {
	// Validate PIN and confirm PIN match
	if req.PIN != req.ConfirmPIN {
		return errors.New("PIN and confirm PIN do not match")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if PIN already exists
	if user.PIN != "" {
		return errors.New("PIN already exists, use change PIN instead")
	}

	// Hash PIN
	hashedPIN, err := utils.HashPassword(req.PIN)
	if err != nil {
		return errors.New("failed to hash PIN")
	}

	// Update user with new PIN
	if err := s.db.Model(&user).Update("pin", hashedPIN).Error; err != nil {
		return err
	}

	return nil
}

func (s *PINService) ChangePIN(userID uint, req dto.ChangePINRequest) error {
	// Validate new PIN and confirm PIN match
	if req.NewPIN != req.ConfirmPIN {
		return errors.New("new PIN and confirm PIN do not match")
	}

	// Validate new PIN is different from old PIN
	if req.OldPIN == req.NewPIN {
		return errors.New("new PIN must be different from old PIN")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if PIN exists
	if user.PIN == "" {
		return errors.New("PIN not set, use create PIN instead")
	}

	// Verify old PIN
	if !utils.CheckPasswordHash(req.OldPIN, user.PIN) {
		return errors.New("old PIN is incorrect")
	}

	// Hash new PIN
	hashedPIN, err := utils.HashPassword(req.NewPIN)
	if err != nil {
		return errors.New("failed to hash new PIN")
	}

	// Update user with new PIN
	if err := s.db.Model(&user).Update("pin", hashedPIN).Error; err != nil {
		return err
	}

	return nil
}

func (s *PINService) HasPIN(userID uint) (bool, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, err
	}

	return user.PIN != "", nil
}

func (s *PINService) VerifyPIN(userID uint, pin string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if user.PIN == "" {
		return errors.New("PIN not set")
	}

	if !utils.CheckPasswordHash(pin, user.PIN) {
		return errors.New("PIN is incorrect")
	}

	return nil
}
