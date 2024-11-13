package tests

import (
	"Crawlzilla/models/ads"
	"Crawlzilla/services/super_admin"
	"os"
	"testing"
)

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

// TestValidateCrawlResultData tests the ValidateCrawlResultData function
func TestValidateCrawlResultData(t *testing.T) {
	tests := []struct {
		name      string
		result    *ads.CrawlResult
		expectErr bool
	}{
		{
			name: "Valid Crawl Result",
			result: &ads.CrawlResult{
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
			result: &ads.CrawlResult{
				Title:       "",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Long Title",
			result: &ads.CrawlResult{
				Title:       "This is a very long title that should fail validation because it exceeds fifty characters",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Empty Location URL",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Long Location URL",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "https://this-is-a-very-long-url-that-should-fail-validation-because-it-exceeds-the-max-lengthsgdgfgfsgfhffsxhfhfhfhfhfhfhfghfghfghfghfhfhfhfhfhfhfhfhfhfhflengthsgdgfgfsgfhffsxhfhfhfhfhfhfhfghfghfghfghfhfhfhfhfhfhfhfhfhfhflengthsgdgfgfsgfhffsxhfhfhfhfhfhfhfghfghfghfghfhfhfhfhfhfhfhfhfhfhflengthsgdgfgfsgfhffsxhfhfhfhfhfhfhfghfghfghfghfhfhfhfhfhfhfhfhfhfhf",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Negative Price",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       -1,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Invalid Latitude",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    100.0, // Out of range
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Invalid Longitude",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   200.0, // Out of range
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := super_admin.ValidateCrawlResultData(tt.result)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateCrawlResultData() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestAddAdForSuperAdmin tests the AddAdForSuperAdmin function
func TestAddAdForSuperAdmin(t *testing.T) {
	tests := []struct {
		name      string
		result    *ads.CrawlResult
		expectErr bool
	}{
		{
			name: "Valid Crawl Result",
			result: &ads.CrawlResult{
				Title:         "Valid Title",
				LocationURL:   "https://valid.url",
				Price:         100,
				Latitude:      45.0,
				Longitude:     90.0,
				Description:   "stringddddddddddddddddd",
				ImageURL:      "string  gorm:type:varchar(255)",
				URL:           "string  gorm:type:varchar(255)",
				Reference:     "string  gorm:type:varchar(10)",
				Area:          2,
				Room:          2,
				FloorNumber:   1,
				TotalFloors:   5,
				ContactNumber: 3,
				VisitCount:    1000,
				HasElevator:   true,
				HasStorage:    true,
				HasParking:    false,
			},
			expectErr: false,
		},
		{
			name: "Invalid Crawl Result - Empty Title",
			result: &ads.CrawlResult{
				Title:       "",
				LocationURL: "https://valid.url",
				Price:       100,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
		{
			name: "Invalid Crawl Result - Negative Price",
			result: &ads.CrawlResult{
				Title:       "Valid Title",
				LocationURL: "https://valid.url",
				Price:       -1,
				Latitude:    45.0,
				Longitude:   90.0,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := super_admin.AddAdForSuperAdmin(tt.result)
			if (err != nil) != tt.expectErr {
				t.Errorf("AddAdForSuperAdmin() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
