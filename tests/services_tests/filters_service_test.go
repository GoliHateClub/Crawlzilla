package services_tests

import (
	"Crawlzilla/models"
	"Crawlzilla/services/filters"
	"strconv"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&models.Filters{})
	return db
}

func TestGetAllFiltersService(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE filters") // Clean up after the test

	// Insert test filters
	for i := 0; i < 30; i++ {
		db.Create(&models.Filters{
			City:         "City_" + strconv.Itoa(i),
			Neighborhood: "Neighborhood_" + strconv.Itoa(i),
			Area:         i * 10,
			Price:        i * 1000,
		})
	}

	// Test with page 1, page size 10
	result, err := filters.GetAllFilters(db, 1, 10, "SUPER_ADMIN")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(result.Data), "Page size should be 10")
	assert.Equal(t, 3, result.Pages, "Total pages should be 3 for 30 records with page size 10")
	assert.Equal(t, 1, result.Page, "Current page should be 1")
}

func TestCreateOrUpdateFilterService(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE filters") // Clean up after the test

	filter := models.Filters{City: "TestCity"}
	created, err := filters.CreateOrUpdateFilte(db, filter)
	assert.True(t, created)
	assert.NoError(t, err)
}

func TestRemoveFilterService(t *testing.T) {
	db := setupTestDB()
	defer db.Exec("DROP TABLE filters") // Clean up after the test

	filter := models.Filters{City: "TestCity"}
	db.Create(&filter)

	removed, err := filters.RemoveFilter(db, filter.ID)
	assert.True(t, removed)
	assert.NoError(t, err)
}
