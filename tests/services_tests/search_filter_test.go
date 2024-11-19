package services_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/services/search"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB sets up a test in-memory database
func SetupSearchTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Run AutoMigrate to create the Ads table
	if err := db.AutoMigrate(&models.Ads{}, &models.Users{}, &models.Filters{}); err != nil {
		panic("failed to migrate database schema")
	}

	return db
}

// TestGetFilteredAdsSuccess tests the successful case for GetFilteredAds.
func TestGetFilteredAdsSuccess(t *testing.T) {
	// Set up in-memory database
	db := SetupSearchTestDB()

	// Prepare mock data
	filter := models.Filters{
		City:     "City",
		Sort:     "price",
		Order:    "asc",
		MinArea:  50,
		MaxArea:  100,
		MinPrice: 1000,
		MaxPrice: 5000,
	}

	// Create a filter entry in the database
	err := repositories.CreateOrUpdateFilter(db, &filter)
	assert.NoError(t, err)

	// Add some ads related to the filter in the database
	ads := []models.Ads{
		{Title: "Test1", City: "City", Price: 1500, Area: 60},
		{Title: "Test2", City: "City", Price: 3000, Area: 80},
		{Title: "Test3", City: "City", Price: 4000, Area: 90},
	}
	for _, ad := range ads {
		_, err := repositories.CreateAd(db, &ad)
		assert.NoError(t, err)
	}

	// Call the search service function with pagination (page 1, pageSize 2)
	result, err := search.GetFilteredAds(db, filter.ID, 1, 2)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Pages) // Total pages = 3 records / 2 pageSize
	assert.Equal(t, 1, result.Page)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "City", result.Data[0].City)
}

func TestGetMostFilteredAdsSuccess(t *testing.T) {
	// Set up in-memory database
	db := SetupSearchTestDB()

	// Prepare mock filter data
	filter := models.Filters{
		City:     "City",
		Sort:     "price",
		Order:    "asc",
		MinArea:  50,
		MaxArea:  100,
		MinPrice: 1000,
		MaxPrice: 5000,
	}

	// Create a filter entry in the database
	err := repositories.CreateOrUpdateFilter(db, &filter)
	assert.NoError(t, err)

	// Add some ads related to the filter in the database
	ads := []models.Ads{
		{Title: "Test1", City: "City", Price: 1500, Area: 60},
		{Title: "Test2", City: "City", Price: 3000, Area: 80},
		{Title: "Test3", City: "City", Price: 4000, Area: 90},
	}

	// Insert ads into the database and associate them with the filter
	for _, ad := range ads {
		_, err := repositories.CreateAd(db, &ad)
		assert.NoError(t, err)

		// Update the filter usage count (or the logic of associating the ad with the filter)
		_, err = search.GetFilteredAds(db, filter.ID, 1, 2)
		assert.NoError(t, err)
	}

	// Call the search service function with pagination (page 1, pageSize 2)
	result, err := search.GetMostFilteredAds(db, 1, 2)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Pages) // Total pages = 3 records / 2 pageSize
	assert.Equal(t, 1, result.Page)
	assert.Len(t, result.Data, 2)

	// Ensure that the ads returned are sorted by their filter usage count
	assert.Equal(t, "City", result.Data[0].City) // Verify the city of the first ad in the result
	assert.Equal(t, 1500, result.Data[0].Price)  // Verify the price of the first ad in the result
}

// TestGetFilteredAdsError tests the error case for GetFilteredAds.
func TestGetFilteredAdsError(t *testing.T) {
	// Set up in-memory database
	db := SetupSearchTestDB()

	// Call the search service function with a non-existent filter ID
	result, err := search.GetFilteredAds(db, "non-existent-id", 1, 2)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, result.Data)
}
