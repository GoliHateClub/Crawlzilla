package repositories

import (
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// CreateUser creates a new user
func CreateUser(db *gorm.DB, user *models.Users) error {
	return db.Create(user).Error
}

// GetUserByTelegramID retrieves a user by their Telegram ID
func GetUserByTelegramID(db *gorm.DB, telegramID string) (models.Users, error) {
	var user models.Users
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
	return user, err
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
