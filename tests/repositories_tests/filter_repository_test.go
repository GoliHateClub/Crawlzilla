package repositories_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilterByID(t *testing.T) {
	db := SetupTestDB()

	// Seed test data
	filter := models.Filters{Title: "Test", City: "Test City", HasParking: true}
	db.Create(&filter)
	log.Println("HEREEEEEEEEEEEEEE ", filter)
	// Test repository function
	retrievedFilter, err := repositories.GetFilterByID(db, filter.ID)
	assert.NoError(t, err)
	assert.Equal(t, filter.City, retrievedFilter.City)
	assert.True(t, retrievedFilter.HasParking)
}

func TestCountFilteredAds(t *testing.T) {
	db := SetupTestDB()

	// Seed test data
	ads := []models.Ads{
		{City: "City1", Price: 1000},
		{City: "City2", Price: 2000},
	}
	db.Create(&ads)

	query := db.Model(&models.Ads{}).Where("price >= ?", 1000)
	total, err := repositories.CountFilteredAds(db, query)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
}

func TestGetFilteredAds(t *testing.T) {
	db := SetupTestDB()

	// Seed test data
	ads := []models.Ads{
		{City: "City1", Price: 1000},
		{City: "City2", Price: 2000},
	}
	db.Create(&ads)

	query := db.Model(&models.Ads{}).Where("price >= ?", 1000)
	results, err := repositories.GetFilteredAds(db, query, 1, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "City1", results[0].City)
}
