package tests

import (
	"Crawlzilla/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUsersBeforeCreate(t *testing.T) {
	// Mock the GORM DB transaction
	db := &gorm.DB{}
	user := models.Users{}

	err := user.BeforeCreate(db)
	assert.NoError(t, err, "BeforeCreate should not return an error")
	assert.NotEmpty(t, user.ID, "User ID should be set with UUID")
}
