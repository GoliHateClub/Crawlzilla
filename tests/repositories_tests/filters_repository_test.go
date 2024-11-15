package repositories_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDatabase() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Filters{})
	return db
}

func TestCreateOrUpdateFilterCreatesNewWithoutID(t *testing.T) {
	db := setupDatabase()
	filter := models.Filters{
		City: "TestCity",
	}

	created, err := repositories.CreateOrUpdateFilter(db, &filter)
	assert.True(t, created)
	assert.Nil(t, err)
	assert.NotEmpty(t, filter.ID, "ID should be set by GORM after creation")
}

func TestRemoveFilter(t *testing.T) {
	db := setupDatabase()
	filter := models.Filters{City: "TestCity"}
	db.Create(&filter)

	removed, err := repositories.RemoveFilter(db, filter.ID)
	assert.True(t, removed)
	assert.Nil(t, err)
}
