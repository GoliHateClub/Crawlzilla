package repositories

import (
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// CreateUser creates a new user
func CreateUser(db *gorm.DB, telegramId int64, chatID int64) (models.Users, error) {
	var user models.Users

	// Check if the user already exists
	if err := db.Where("telegram_id = ?", telegramId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new user if not found
			user.Role = models.RoleUser // Set the default role
			user.Telegram_ID = telegramId
			user.ChatID = chatID

			if err := db.Create(&user).Error; err != nil { // Pass &user instead of user
				return user, err
			}
			return user, nil
		}
		return user, err // Handle unexpected errors
	}
	return user, nil // Return the existing user's role
}

// CreateUser creates a new user
func CreateAdmin(db *gorm.DB, telegram_id int64) (models.Users, error) {
	var user models.Users

	// Check if the user already exists
	if err := db.Where("telegram_id = ?", telegram_id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new user if not found
			user.Role = models.RoleAdmin // Set the default role
			user.Telegram_ID = telegram_id
			if err := db.Save(&user).Error; err != nil { // Pass &user instead of user
				return user, err
			}
			return user, nil
		}
		return user, err // Handle unexpected errors
	}

	// for not existed admins
	user.Role = models.RoleAdmin
	err := db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil // Return the existing user's role
}

// GetUserByID retrieves a user by their Telegram ID
func GetUserByID(db *gorm.DB, userID string) (models.Users, error) {
	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
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

// GetUserByID retrieves a user by their Telegram ID
func GetUserID(db *gorm.DB, telegramID string) (string, error) {
	var user models.Users
	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return "", err
	}
	return user.ID, nil
}

// SetChatID updates the ChatID for a user identified by their Telegram_ID
func SetChatID(db *gorm.DB, telegramID int64, chatID int64) error {
	var user models.Users

	// Check if the user exists
	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound // Return if user is not found
		}
		return err // Return any other error
	}

	// Update the ChatID
	user.ChatID = chatID
	if err := db.Save(&user).Error; err != nil {
		return err // Return error if save operation fails
	}

	return nil // Success
}
func GetUserByTelegramID(db *gorm.DB, telegramID int64) (models.Users, error) {
	var user models.Users
	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return models.Users{}, err
	}
	return user, nil
}
