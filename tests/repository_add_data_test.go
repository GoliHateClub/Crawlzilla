package tests

import (
	"Crawlzilla/models/ads"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB initializes an in-memory SQLite database for testing with ads.CrawlResult
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Run AutoMigrate to create the CrawlResult table
	if err := db.AutoMigrate(&ads.CrawlResult{}); err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

// TestAddCrawlResult tests adding unique and duplicate CrawlResult records and the Hash field
func TestAddCrawlResultWithHash(t *testing.T) {
	// Initialize test database
	database := SetupTestDB()

	// First CrawlResult entry
	entry1 := &ads.CrawlResult{
		Title:       "Sample Name",
		URL:         "http://example.com/item1",
		Description: "Sample Description",
	}
	err := database.Create(entry1).Error
	assert.NoError(t, err)

	// Check if the record was created
	var count int64
	database.Model(&ads.CrawlResult{}).Count(&count)
	assert.Equal(t, int64(1), count, "Record should be inserted")

	// Verify the Hash for the first entry
	var firstRecord ads.CrawlResult
	database.First(&firstRecord)
	assert.NotEmpty(t, firstRecord.Hash, "Hash should be set for the first record")

	// Second CrawlResult entry with the same URL, Title, and Description (should trigger duplicate hash)
	entry2 := &ads.CrawlResult{
		Title:       "Sample Name",              // Same as entry1
		URL:         "http://example.com/item1", // Same URL as entry1
		Description: "Sample Description",       // Same Description as entry1
	}
	err = database.Create(entry2).Error
	assert.Error(t, err, "Duplicate record should not be inserted due to unique constraint on Hash")

	// Verify that the duplicate was not inserted based on the hash
	database.Model(&ads.CrawlResult{}).Count(&count)
	assert.Equal(t, int64(1), count, "Duplicate record should not be inserted based on hash")

	// Verify that the hash of the second entry is the same as the first entry
	var secondRecord ads.CrawlResult
	database.Last(&secondRecord)
	assert.Equal(t, firstRecord.Hash, secondRecord.Hash, "Hashes of the two records should be the same")

	// Clean up test database
	database.Exec("DROP TABLE crawl_results")
}
