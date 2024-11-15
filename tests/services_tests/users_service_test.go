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
		user := models.Users{Telegram_ID: "user_" + strconv.Itoa(i), Role: "user"}
		repositories.CreateUser(db, &user)
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
