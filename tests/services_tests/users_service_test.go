package services_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/services/users"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupServiceTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Users{})
	return db
}

func TestGetAllUsersPaginatedService(t *testing.T) {
	db := setupServiceTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert test users
	for i := 0; i < 25; i++ {
		telegramID := int64(i)
		chatID := int64(1000 + i) // Generate a unique chatID for each user
		repositories.CreateUser(db, telegramID, chatID)
	}

	// Test with page 1, page size 10
	result, err := users.GetAllUsersPaginatedService(db, 1, 10)
	assert.NoError(t, err, "Paginated retrieval should not return an error")
	assert.Equal(t, 10, len(result.Data), "Page size should be 10")
	assert.Equal(t, 3, result.Pages, "Total pages should be 3 for 25 records with page size 10")
	assert.Equal(t, 1, result.Page, "Current page should be 1")

	// Test with page 2, page size 10
	result, err = users.GetAllUsersPaginatedService(db, 2, 10)
	assert.NoError(t, err, "Paginated retrieval should not return an error")
	assert.Equal(t, 10, len(result.Data), "Page size should be 10")
	assert.Equal(t, 3, result.Pages, "Total pages should be 3")
	assert.Equal(t, 2, result.Page, "Current page should be 2")
}

func TestGetUserByIDService(t *testing.T) {
	db := setupServiceTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert a test user
	user := models.Users{
		Telegram_ID: int64(12345),
		ChatID:      int64(67890),
	}
	createdUser, err := repositories.CreateUser(db, user.Telegram_ID, user.ChatID)
	assert.NoError(t, err, "Creating user should not return an error")

	// Test retrieving the user by ID
	retrievedUser, err := users.GetUserByIDService(db, createdUser.ID)
	assert.NoError(t, err, "Retrieving a valid user should not return an error")
	assert.NotNil(t, retrievedUser, "Retrieved user should not be nil")
	assert.Equal(t, createdUser.ID, retrievedUser.ID, "The retrieved user ID should match the created user ID")
	assert.Equal(t, createdUser.Role, retrievedUser.Role, "The retrieved user role should match the created user role")
	assert.Equal(t, createdUser.ChatID, retrievedUser.ChatID, "The retrieved user ChatID should match the created user ChatID")

	// Test retrieving a user with a non-existent ID
	nonExistentID := "00000000-0000-0000-0000-000000000000"
	retrievedUser, err = users.GetUserByIDService(db, nonExistentID)
	assert.Error(t, err, "Retrieving a user with a non-existent ID should return an error")
	assert.Equal(t, "user not found", err.Error(), "Error message should be 'user not found'")
	assert.Equal(t, models.Users{}, retrievedUser, "Retrieved user should be an empty struct for a non-existent ID")
}
