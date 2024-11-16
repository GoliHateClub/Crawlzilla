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
func LoginUser(db *gorm.DB, telegram_id string) (string, error) {
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

// CreateUserService validates and creates a new user
func CreateUserService(db *gorm.DB, user *models.Users) error {
	// Validate Role
	if user.Role != "admin" && user.Role != "user" && user.Role != "super-admin" {
		return errors.New("invalid role, must be 'admin', 'user', or 'super-admin'")
	}

	// Validate Telegram ID (e.g., must be non-empty and match a pattern)
	if user.Telegram_ID == "" || !isValidTelegramID(user.Telegram_ID) {
		return errors.New("invalid Telegram ID")
	}

	// Call repository function to create the user
	return repositories.CreateUser(db, user)
}

// GetUserByIDService retrieves a user by their ID with validation
func GetUserByIDService(db *gorm.DB, userID string) (*models.Users, error) {
	// Validate user ID (e.g., must not be empty)
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Call repository function to retrieve the user
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Helper function to validate Telegram ID (example regex, customize as needed)
func isValidTelegramID(telegramID string) bool {
	// Example: Telegram IDs must be alphanumeric and between 5-32 characters
	var telegramIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{5,32}$`)
	return telegramIDRegex.MatchString(telegramID)
}
