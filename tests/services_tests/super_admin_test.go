package services_tests

import (
	"Crawlzilla/models"
	"Crawlzilla/services/super_admin"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Optional logging for debugging
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Automatically migrate the schema (create tables)
	if err := db.AutoMigrate(&models.Ads{}); err != nil {
		panic("failed to migrate database schema")
	}

	return db
}
func TestIsSuperAdmin(t *testing.T) {
	err := os.Setenv("SUPER_ADMIN_ID", "1922802339")
	if err != nil {
		return
	}
	tests := []struct {
		userID   int64
		expected bool
	}{
		{userID: 1922802339, expected: true},
		{userID: 1922706439, expected: false},
	}

	for _, tt := range tests {
		t.Run("TestIsSuperAdmin", func(t *testing.T) {
			result := super_admin.IsSuperAdmin(tt.userID)
			if result != tt.expected {
				t.Errorf("IsSuperAdmin(%d) = %v; want %v", tt.userID, result, tt.expected)
			}
		})
	}
}

// TestValidateAdData tests the ValidateAdData function
func TestValidateAdData(t *testing.T) {
	tests := []struct {
		name      string
		result    *models.Ads
		expectErr bool
	}{
		{
			name: "Valid Crawl Result",
			result: &models.Ads{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: false,
		},
		{
			name: "Empty Title",
			result: &models.Ads{
				Title:       "",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		// Additional tests here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := super_admin.ValidateAdData(tt.result)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateAdData() error = %v, wantErr %v", err, tt.expectErr)
			}
			if err != nil && !tt.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAddAdForSuperAdmin(t *testing.T) {
	// Setup the in-memory database for testing
	db := SetupTestDB()

	tests := []struct {
		name      string
		result    *models.Ads
		expectErr bool
	}{
		{
			name: "Valid Crawl Result",
			result: &models.Ads{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: false,
		},
		{
			name: "Invalid Crawl Result - Empty Title",
			result: &models.Ads{
				Title:       "",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
	}

	// Loop over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with the mock database
			err := super_admin.CreateAd(tt.result, db)
			if (err != nil) != tt.expectErr {
				t.Errorf("AddAdForSuperAdmin() error = %v, wantErr %v", err, tt.expectErr)
			}

			// If no error is expected, check the data was actually inserted
			if !tt.expectErr {
				var ad models.Ads
				db.First(&ad, "title = ?", tt.result.Title)
				assert.Equal(t, tt.result.Title, ad.Title)
				assert.Equal(t, tt.result.Price, ad.Price)
				assert.Equal(t, tt.result.LocationURL, ad.LocationURL)
			}
		})
	}
}
