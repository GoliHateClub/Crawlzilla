package repositories_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"errors"
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

func TestCreateAd(t *testing.T) {
	db := SetupTestDB()
	defer db.Exec("DROP TABLE ads;")

	// Create a sample ad
	ad := models.Ads{
		Title: "Sample Ad",
	}

	// Test case 1: Successfully add an ad
	id, err := repositories.CreateAd(&ad, db)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Test case 2: Attempt to add a duplicate ad
	adDuplicate := models.Ads{
		Title: ad.Title,
		// Hash:  ad.Hash, // Manually set to the same hash to simulate duplication
	}
	_, err = repositories.CreateAd(&adDuplicate, db)
	assert.Error(t, err)
	assert.Equal(t, "cant't add to database", err.Error())
}

func TestGetAllAds(t *testing.T) {
	db := SetupTestDB()
	defer db.Exec("DROP TABLE ads;")

	// Add a sample ad to the database
	ad := models.Ads{Title: "Sample Ad"}
	db.Create(&ad)

	// Test case: Retrieve all ads
	results, totalRecords, err := repositories.GetAllAds(db, 0, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.Equal(t, "Sample Ad", results[0].Title)
	assert.Equal(t, int64(1), totalRecords)
}

func TestGetAdByID(t *testing.T) {
	// Set up the test database
	db := SetupTestDB()
	defer db.Exec("DROP TABLE ads;")

	// Add a sample ad to the database
	ad := models.Ads{Title: "Sample Ad"}
	id, err := repositories.CreateAd(&ad, db)
	assert.NoError(t, err, "Expected no error when creating ad")

	// Test case 1: Retrieve an existing ad by ID
	result, err := repositories.GetAdByID(`"`+id+`"`, db) // Wrap ID in quotes
	assert.NoError(t, err, "Expected no error when retrieving an existing ad")
	assert.Equal(t, "Sample Ad", result.Title, "Expected title to match the inserted ad")

	// Test case 2: Retrieve a non-existent ad by ID
	nonExistentID := "00000000-0000-0000-0000-000000000000" // A valid but non-existent UUID
	_, err = repositories.GetAdByID(nonExistentID, db)
	assert.Error(t, err, "Expected an error when retrieving a non-existent ad")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Expected gorm.ErrRecordNotFound error")
}

func TestDeleteAdById(t *testing.T) {
	db := SetupTestDB()
	defer db.Exec("DROP TABLE ads;")

	// Add a sample ad to the database
	ad := models.Ads{Title: "Sample Ad"}
	id, err := repositories.CreateAd(&ad, db)

	// Test case: Delete the ad by ID
	err = repositories.DeleteAdById(`"`+id+`"`, db)
	assert.NoError(t, err)

	// Verify that the ad is deleted
	_, err = repositories.GetAdByID(`"`+ad.ID+`"`, db)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
