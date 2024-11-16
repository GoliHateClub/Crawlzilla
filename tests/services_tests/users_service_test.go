package services_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/services/users"
	"strconv"
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
		user := models.Users{Telegram_ID: "user_" + strconv.Itoa(i)}
		repositories.CreateUser(db, user.Telegram_ID)
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
func TestCreateUserService(t *testing.T) {
	db := setupServiceTestDB()
	defer db.Exec("DROP TABLE users")

	// Test valid user creation
	user := models.Users{
		Telegram_ID: "validTelegramID",
		Role:        "admin",
	}

	err := users.CreateUserService(db, &user)
	assert.NoError(t, err, "Creating a valid user should not return an error")

	// Verify the user was saved to the database
	retrievedUser, err := repositories.GetUserByID(db, user.ID)
	assert.NoError(t, err, "Retrieving the created user should not return an error")
	assert.Equal(t, user.ID, retrievedUser.ID, "The created and retrieved user IDs should match")
	assert.Equal(t, user.Role, retrievedUser.Role, "The created and retrieved user roles should match")

	// Test invalid role
	invalidUser := models.Users{
		Telegram_ID: "validTelegramID2",
		Role:        "invalid-role",
	}

	err = users.CreateUserService(db, &invalidUser)
	assert.Error(t, err, "Creating a user with an invalid role should return an error")
	assert.Equal(t, "invalid role, must be 'admin', 'user', or 'super-admin'", err.Error())

	// Test invalid Telegram ID
	invalidUser = models.Users{
		Telegram_ID: "",
		Role:        "user",
	}

	err = users.CreateUserService(db, &invalidUser)
	assert.Error(t, err, "Creating a user with an invalid Telegram ID should return an error")
	assert.Equal(t, "invalid Telegram ID", err.Error())
}

func TestGetUserByIDService(t *testing.T) {
	db := setupServiceTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert a test user
	user := models.Users{
		Telegram_ID: "testTelegramID",
		Role:        "super-admin",
	}
	repositories.CreateUser(db, &user)

	// Test retrieving the user by ID
	retrievedUser, err := users.GetUserByIDService(db, user.ID)
	assert.NoError(t, err, "Retrieving a valid user should not return an error")
	assert.NotNil(t, retrievedUser, "Retrieved user should not be nil")
	assert.Equal(t, user.ID, retrievedUser.ID, "The retrieved user ID should match the created user ID")
	assert.Equal(t, user.Role, retrievedUser.Role, "The retrieved user role should match the created user role")

	// Test retrieving a user with a non-existent ID
	nonExistentID := "non-existent-id"
	retrievedUser, err = users.GetUserByIDService(db, nonExistentID)
	assert.Error(t, err, "Retrieving a user with a non-existent ID should return an error")
	assert.Equal(t, "user not found", err.Error())
}
