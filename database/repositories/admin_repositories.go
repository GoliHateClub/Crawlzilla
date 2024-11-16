package repositories

import (
	"Crawlzilla/models"
	"errors"
	"gorm.io/gorm"
)

// CreateAdmin creates an admin if it doesn't exist or make an existing user to be an admin.
func CreateAdmin(db *gorm.DB, telegramID int64) (models.Role, error) {
	var user models.Users

	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAdmin := models.Users{
				Telegram_ID: telegramID,
				Role:        models.RoleAdmin,
			}
			if err := db.Create(&newAdmin).Error; err != nil {
				return "", err
			}
			return newAdmin.Role, nil
		}
		return "", err
	}

	if user.Role != models.RoleAdmin {
		user.Role = models.RoleAdmin
		if err := db.Save(&user).Error; err != nil {
			return "", err
		}
	}

	return user.Role, nil
}

// GetAdminByID retrieves a user with the "admin" role by their Telegram ID
func GetAdminByID(db *gorm.DB, telegramID int64) (*models.Users, error) {
	var user models.Users
	err := db.Where("telegram_id = ? AND role = ?", telegramID, models.RoleAdmin).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
