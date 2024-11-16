package repositories_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Users{})
	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE users")

	user := models.Users{
		Telegram_ID: int64(12345),
	}

	user, err := repositories.CreateUser(db, user.Telegram_ID)
	assert.NoError(t, err, "Creating user should not return an error")
	assert.Equal(t, user.Role, models.RoleUser, "Role should match")
}

func TestGetUserByID(t *testing.T) {
	// Setup the test database
	db := setupTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert test user
	user := models.Users{
		Role: "admin",
	}
	err := db.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	// Retrieve the generated ID from the user
	userID := user.ID // This is populated by GORM after creation

	retrievedUser, err := repositories.GetUserByID(db, userID)
	assert.NoError(t, err, "Retrieving user should not return an error")
	assert.Equal(t, user.ID, retrievedUser.ID, "User ID should match")
	assert.Equal(t, user.Role, retrievedUser.Role, "Role should match")
}

func TestGetAllUsersPaginated(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert test users
	for i := 0; i < 15; i++ {
		user := models.Users{Telegram_ID: int64(i), Role: models.RoleUser}
		db.Create(&user)
	}

	users, totalRecords, err := repositories.GetAllUsersPaginated(db, 1, 10)
	assert.NoError(t, err, "Paginated retrieval should not return an error")
	assert.Equal(t, int64(15), totalRecords, "Total records should match")
	assert.Equal(t, 10, len(users), "Page size should be 10")
}
