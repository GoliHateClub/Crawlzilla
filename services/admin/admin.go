package admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"gorm.io/gorm"
)

// CreateAdminUser creates a new user with the admin role
func CreateAdminUser(db *gorm.DB, telegramID int64) (models.Role, error) {
	role, err := repositories.CreateAdmin(db, telegramID)
	if err != nil {
		return "", err
	}

	return role, nil
}

// IsAdmin checks if a user is an admin by their Telegram ID
func IsAdmin(db *gorm.DB, telegramID int64) (bool, error) {

	user, err := repositories.GetAdminByID(db, telegramID)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
