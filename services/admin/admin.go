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
