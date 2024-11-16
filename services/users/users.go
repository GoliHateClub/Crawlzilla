package users

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type PaginatedUsers struct {
	Data  []models.Users `json:"data"`
	Pages int            `json:"pages"`
	Page  int            `json:"page"`
}

// GetAllUsersPaginatedService retrieves all users with pagination and structures the output
func LoginUser(db *gorm.DB, telegramId int64) (models.Users, error) {
	return repositories.CreateUser(db, telegram_id)
}

// GetAllUsersPaginatedService retrieves all users with pagination and structures the output
func GetAllUsersPaginatedService(db *gorm.DB, page int, pageSize int) (PaginatedUsers, error) {
	users, totalRecords, err := repositories.GetAllUsersPaginated(db, page, pageSize)
	if err != nil {
		return PaginatedUsers{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the output struct
	result := PaginatedUsers{
		Data:  users,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

// GetUserByIDService retrieves a user by their ID with validation
func GetUserByIDService(db *gorm.DB, userID string) (models.Users, error) {
	// Validate user ID (e.g., must not be empty)
	if userID == "" {
		return models.Users{}, errors.New("user ID cannot be empty")
	}

	// Call repository function to retrieve the user
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return models.Users{}, errors.New("user not found")
	}
	return user, nil
}

// Helper function to validate Telegram ID (example regex, customize as needed)
func isValidTelegramID(telegramID string) bool {
	// Example: Telegram IDs must be alphanumeric and between 5-32 characters
	var telegramIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{5,32}$`)
	return telegramIDRegex.MatchString(telegramID)
}
