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

// RemoveAdmin deletes the admin user by their Telegram ID (int64).
func RemoveAdmin(db *gorm.DB, telegramID int64) error {
	err := repositories.RemoveAdminByTelegramID(db, telegramID)
	if err != nil {
		return err
	}

	return nil
}

// GetAdmins fetches paginated admins data from the database.
func GetAdmins(db *gorm.DB, page int, pageSize int) ([]models.Users, int64, error) {
	admins, total, err := repositories.GetAdminsPaginated(db, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}
