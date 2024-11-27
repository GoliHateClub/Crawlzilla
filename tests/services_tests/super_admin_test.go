package services_tests

import (
	"Crawlzilla/models"
	"Crawlzilla/services/super_admin"
	"errors"
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
			name: "Valid Ad Data",
			result: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				ContactNumber: "09121111111",
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
			},
			expectErr: false,
		},
		{
			name: "Empty Title",
			result: &models.Ads{
				Title:         "",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				ContactNumber: "09121111111",
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
			},
			expectErr: true,
		},
		{
			name: "Invalid Category Type",
			result: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				ContactNumber: "0",
				CategoryType:  "invalid",
				PropertyType:  "آپارتمانی",
			},
			expectErr: true,
		},
		{
			name: "Invalid Latitude",
			result: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      100.0, // Invalid latitude
				Longitude:     90.0,
				ContactNumber: "09121111111",
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := super_admin.ValidateAdData(tt.result)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateAdData() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestCreateAd(t *testing.T) {
	db := SetupTestDB()

	tests := []struct {
		name      string
		inputAd   *models.Ads
		expectErr bool
		validate  func(t *testing.T, ad *models.Ads)
	}{
		{
			name: "Valid Ad Data",
			inputAd: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				CategoryType:  "فروش",      // Persian for "sell"
				PropertyType:  "آپارتمانی", // Persian for "house"
				ContactNumber: "09123456789",
			},
			expectErr: false,
			validate: func(t *testing.T, ad *models.Ads) {
				assert.Equal(t, "Valid Title", ad.Title)
				assert.Equal(t, "sell", ad.CategoryType)
				assert.Equal(t, "house", ad.PropertyType)
				assert.Contains(t, ad.URL, "super-admin-")
				assert.NotEmpty(t, ad.LocationURL)
			},
		},
		{
			name: "Invalid Ad Data - Empty Title",
			inputAd: &models.Ads{
				Title:         "",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
				ContactNumber: "09123456789",
			},
			expectErr: true,
		},
		{
			name: "Invalid Ad Data - Invalid Phone Number",
			inputAd: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
				ContactNumber: "invalid-phone",
			},
			expectErr: true,
		},
		{
			name: "Invalid Ad Data - Invalid Coordinates",
			inputAd: &models.Ads{
				Title:         "Valid Title",
				Price:         100,
				Latitude:      100.0, // Out of range
				Longitude:     200.0, // Out of range
				CategoryType:  "فروش",
				PropertyType:  "آپارتمانی",
				ContactNumber: "09123456789",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean the database before each test
			db.Exec("DELETE FROM ads")

			// Call CreateAd
			err := super_admin.CreateAd(db, tt.inputAd)

			// Assert error presence
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Validate the inserted ad
				var savedAd models.Ads
				err := db.First(&savedAd, "title = ?", tt.inputAd.Title).Error
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, &savedAd)
				}
			}
		})
	}
}

func TestRemoveAdByID(t *testing.T) {
	db := SetupTestDB() // Initialize the test database

	// Seed the database with test data
	testAd := models.Ads{
		ID:          "91a91cd0-4d06-4ddc-8deb-b4522ff1e1db",
		Title:       "Test Ad",
		LocationURL: "https://example.com",
		Price:       150,
		Latitude:    45.0,
		Longitude:   90.0,
		VisitCount:  0,
	}
	if err := db.Create(&testAd).Error; err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	tests := []struct {
		name      string
		id        string
		expectErr bool
	}{
		{
			name:      "Valid ID - Ad exists",
			id:        testAd.ID,
			expectErr: false,
		},
		{
			name:      "Invalid ID - Ad does not exist",
			id:        "invalid-id-12345",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := super_admin.RemoveAdByID(db, tt.id)

			if (err != nil) != tt.expectErr {
				t.Errorf("RemoveAdByID() error = %v, wantErr %v", err, tt.expectErr)
			}

			// If no error is expected, check the ad was actually deleted
			if !tt.expectErr {
				var ad models.Ads
				result := db.First(&ad, "id = ?", tt.id)

				if result.Error == nil || !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					t.Errorf("Expected ad to be deleted, but found: %+v", ad)
				}
			}
		})
	}
}
