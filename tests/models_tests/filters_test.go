package tests

import (
	"Crawlzilla/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBeforeCreateSetsUUID(t *testing.T) {
	f := models.Filters{}
	db := &gorm.DB{} // Mocked *gorm.DB
	err := f.BeforeCreate(db)
	assert.Nil(t, err)
	assert.NotEmpty(t, f.ID, "ID should not be empty after BeforeCreate")
}
