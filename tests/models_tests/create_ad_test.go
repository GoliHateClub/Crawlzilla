package tests

import (
	"Crawlzilla/models"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB initializes an in-memory SQLite database for testing with models.Ads
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Run AutoMigrate to create the Ads table
	if err := db.AutoMigrate(&models.Ads{}); err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

// TestCreateAdWithHash tests adding unique and duplicate Ads records and the Hash field
func TestCreateAdWithHash(t *testing.T) {
	// Initialize test database
	database := SetupTestDB()

	// First Ads entry
	entry1 := &models.Ads{
		Title:       "Sample Name",
		URL:         "http://example.com/item1",
		Description: "Sample Description",
	}
	err := database.Create(entry1).Error
	assert.NoError(t, err)

	// Check if the record was created
	var count int64
	database.Model(&models.Ads{}).Count(&count)
	assert.Equal(t, int64(1), count, "Record should be inserted")

	// Verify the Hash for the first entry
	var firstRecord models.Ads
	database.First(&firstRecord)
	assert.NotEmpty(t, firstRecord.Hash, "Hash should be set for the first record")

	// Second Ads entry with the same URL, Title, and Description (should trigger duplicate hash)
	entry2 := &models.Ads{
		Title:       "Sample Name",              // Same as entry1
		URL:         "http://example.com/item1", // Same URL as entry1
		Description: "Sample Description",       // Same Description as entry1
	}
	err = database.Create(entry2).Error
	assert.Error(t, err, "Duplicate record should not be inserted due to unique constraint on Hash")

	// Verify that the duplicate was not inserted based on the hash
	database.Model(&models.Ads{}).Count(&count)
	assert.Equal(t, int64(1), count, "Duplicate record should not be inserted based on hash")

	// Verify that the hash of the second entry is the same as the first entry
	var secondRecord models.Ads
	database.Last(&secondRecord)
	assert.Equal(t, firstRecord.Hash, secondRecord.Hash, "Hashes of the two records should be the same")

	// Clean up test database
	database.Exec("DROP TABLE ads")
}
