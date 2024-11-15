package admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"gorm.io/gorm"
)

// CreateAdminUser creates a new user with the admin role
func CreateAdminUser(db *gorm.DB, telegramID string) (models.Role, error) {
	role := models.RoleAdmin
	userRole, err := repositories.CreateUser(db, telegramID, role)
	if err != nil {
		return "", err
	}

	return userRole, nil
}

// IsAdmin checks if the user with the given ID is an admin
func IsAdmin(db *gorm.DB, userID string) (bool, error) {
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	if user.Role == models.RoleAdmin {
		return true, nil
	}
	return false, nil
}
