package users

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"

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
