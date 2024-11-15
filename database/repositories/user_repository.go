package repositories

import (
	"Crawlzilla/models"
	"errors"

	"gorm.io/gorm"
)

// CreateUser creates a new user with a given role or updates the role of an existing user.
func CreateUser(db *gorm.DB, telegramID string, role models.Role) (models.Role, error) {
	var user models.Users

	// Check if the user already exists
	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new user if not found
			user.Role = role // Set the role provided (could be default or custom role)
			user.Telegram_ID = telegramID
			if err := db.Create(&user).Error; err != nil {
				return "", err
			}
			return user.Role, nil
		}
		// Handle unexpected errors
		return "", err
	}

	// If the user already exists, update the role
	user.Role = role // Update the user's role
	if err := db.Save(&user).Error; err != nil {
		return "", err
	}

	// Return the updated user's role
	return user.Role, nil
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
