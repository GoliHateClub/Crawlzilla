package repositories_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"strconv"
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
		Telegram_ID: "test_telegram_id",
		Role:        "admin",
	}

	err := repositories.CreateUser(db, &user)
	assert.NoError(t, err, "Creating user should not return an error")
	assert.NotEmpty(t, user.ID, "User ID should be set")
}

func TestGetUserByTelegramID(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert test user
	user := models.Users{
		Telegram_ID: "test_telegram_id",
		Role:        "admin",
	}
	db.Create(&user)

	retrievedUser, err := repositories.GetUserByTelegramID(db, "test_telegram_id")
	assert.NoError(t, err, "Retrieving user should not return an error")
	assert.Equal(t, user.Telegram_ID, retrievedUser.Telegram_ID, "Telegram ID should match")
}

func TestGetAllUsersPaginated(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE users")

	// Insert test users
	for i := 0; i < 15; i++ {
		user := models.Users{Telegram_ID: "user_" + strconv.Itoa(i), Role: "user"}
		db.Create(&user)
	}

	users, totalRecords, err := repositories.GetAllUsersPaginated(db, 1, 10)
	assert.NoError(t, err, "Paginated retrieval should not return an error")
	assert.Equal(t, int64(15), totalRecords, "Total records should match")
	assert.Equal(t, 10, len(users), "Page size should be 10")
}
