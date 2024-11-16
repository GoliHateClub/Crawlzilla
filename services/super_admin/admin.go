package super_admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// CreateAdminUser creates a new user with the admin role
func CreateAdminUser(db *gorm.DB, telegramID int64) (models.Role, error) {

	user, err := repositories.CreateAdmin(db, telegramID)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

// IsAdmin checks if the user with the given ID is an admin
func IsAdmin(db *gorm.DB, userID string) (bool, error) {
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return false, err
	}
	if user.Role == models.RoleAdmin {
		return true, nil
	}
	return false, nil
}
