package repositories

import (
	"Crawlzilla/models"
	"errors"

	"gorm.io/gorm"
)

// CreateUser creates a new user
func CreateUser(db *gorm.DB, user *models.Users) error {
	return db.Create(user).Error
}

// GetUserByID retrieves a user by their Telegram ID
func GetUserByID(db *gorm.DB, userID string) (*models.Users, error) {
	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetAllUsersPaginated retrieves users with pagination
func GetAllUsersPaginated(db *gorm.DB, page int, pageSize int) ([]models.Users, int64, error) {
	var users []models.Users
	var totalRecords int64

	// Count total records for pagination info
	if err := db.Model(&models.Users{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records
	err := db.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, totalRecords, err
}
