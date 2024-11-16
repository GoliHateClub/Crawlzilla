package services_tests

import (
	"Crawlzilla/models"
	"Crawlzilla/services/filters"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB sets up an in-memory database and migrates the models
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Users{}, &models.Filters{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestFilterService_CreateOrUpdateFilter(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Define test cases
	tests := []struct {
		name    string
		filter  models.Filters
		wantID  string // Expected filter ID
		wantErr bool
	}{
		{
			name: "Valid filter with existing user and all fields provided",
			filter: models.Filters{
				USER_ID:      "existing-user-id",
				City:         "Sample City",
				Neighborhood: "Sample Neighborhood",
				Reference:    "Sample Ref",
				CategoryType: "Category",
				PropertyType: "Property",
				Sort:         "price",
				Order:        "asc",
				MinArea:      50,
				MaxArea:      100,
				MinPrice:     1000,
				MaxPrice:     5000,
			},
			wantID:  "", // We don’t expect a specific ID here, since it will be generated
			wantErr: false,
		},
		{
			name: "Valid filter with new user and optional fields missing",
			filter: models.Filters{
				USER_ID:      "new-user-id",
				CategoryType: "Category",
				MinArea:      100,
				MaxArea:      200,
			},
			wantID:  "", // Expect a new ID to be generated
			wantErr: false,
		},
		{
			name: "Invalid filter with min area greater than max area",
			filter: models.Filters{
				USER_ID: "user-id-4",
				MinArea: 150,
				MaxArea: 100,
			},
			wantID:  "", // Expect no ID since it’s invalid
			wantErr: true,
		},
		{
			name: "Invalid filter with negative price values",
			filter: models.Filters{
				USER_ID:  "user-id-5",
				MinPrice: -500,
				MaxPrice: 1000,
			},
			wantID:  "", // Expect no ID since it’s invalid
			wantErr: true,
		},
		{
			name: "Invalid filter with min price greater than max price",
			filter: models.Filters{
				USER_ID:  "user-id-6",
				MinPrice: 1000,
				MaxPrice: 500,
			},
			wantID:  "", // Expect no ID since it’s invalid
			wantErr: true,
		},
		{
			name: "Invalid filter with no user id",
			filter: models.Filters{
				MinPrice: 1000,
				MaxPrice: 500,
			},
			wantID:  "", // Expect no ID since it’s invalid
			wantErr: true,
		},
		{
			name: "Valid filter with minimum values for integers",
			filter: models.Filters{
				USER_ID:  "user-id-7",
				MinArea:  0,
				MaxArea:  50,
				MinPrice: 0,
				MaxPrice: 1000,
			},
			wantID:  "", // Expect a new ID to be generated
			wantErr: false,
		},
		{
			name: "Valid filter with no optional fields provided",
			filter: models.Filters{
				USER_ID: "user-id-8",
			},
			wantID:  "", // Expect a new ID to be generated
			wantErr: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := filters.CreateOrUpdateFilter(db, tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterService.CreateOrUpdateFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotID == "" {
				t.Error("FilterService.CreateOrUpdateFilter() returned an empty ID, expected a generated ID")
			}
			if tt.wantID != "" && gotID != tt.wantID {
				t.Errorf("FilterService.CreateOrUpdateFilter() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestFilterService_GetFiltersByUserID(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Seed data
	userID := "example-user-id"

	// Add sample filters to the database
	sampleFilters := []models.Filters{
		{USER_ID: userID, City: "City1"},
		{USER_ID: userID, City: "City2"},
		{USER_ID: userID, City: "City3"},
		{USER_ID: userID, City: "City4"},
		{USER_ID: userID, City: "City5"},
		{USER_ID: userID, City: "City6"},
	}
	for _, filter := range sampleFilters {
		db.Create(&filter)
	}

	// Define test cases
	tests := []struct {
		name       string
		userID     string
		pageIndex  int
		pageSize   int
		wantCount  int
		wantErr    bool
		totalPages int
	}{
		{
			name:       "Request with invalid page index",
			userID:     userID,
			pageIndex:  -1,
			pageSize:   2,
			wantCount:  0,
			wantErr:    true,
			totalPages: 0,
		},
		{
			name:       "Valid request with page size 2 and page 1",
			userID:     userID,
			pageIndex:  1,
			pageSize:   2,
			wantCount:  2,
			wantErr:    false,
			totalPages: 3,
		},
		{
			name:       "Valid request with page size 3 and page 2",
			userID:     userID,
			pageIndex:  2,
			pageSize:   3,
			wantCount:  3,
			wantErr:    false,
			totalPages: 2,
		},
		{
			name:       "Valid request with page size larger than total filters",
			userID:     userID,
			pageIndex:  1,
			pageSize:   10,
			wantCount:  6,
			wantErr:    false,
			totalPages: 1,
		},
		{
			name:       "Request with no filters",
			userID:     "nonexistent-user-id",
			pageIndex:  1,
			pageSize:   2,
			wantCount:  0,
			wantErr:    false,
			totalPages: 0,
		},
		{
			name:       "Request with invalid page index",
			userID:     userID,
			pageIndex:  -1,
			pageSize:   2,
			wantCount:  0,
			wantErr:    true,
			totalPages: 0,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filters.GetFiltersByUserID(db, tt.userID, tt.pageIndex, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterService.GetFiltersByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Only validate these fields if no error is expected
				if len(result.Data) != tt.wantCount {
					t.Errorf("FilterService.GetFiltersByUserID() count = %v, want %v", len(result.Data), tt.wantCount)
				}
				if result.TotalPages != tt.totalPages {
					t.Errorf("FilterService.GetFiltersByUserID() totalPages = %v, want %v", result.TotalPages, tt.totalPages)
				}
				if result.PageIndex != tt.pageIndex {
					t.Errorf("FilterService.GetFiltersByUserID() pageIndex = %v, want %v", result.PageIndex, tt.pageIndex)
				}
			}
		})
	}

}
func TestFilterService_GetAllFilters(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Seed users and filters
	superAdmin := models.Users{ID: "super-admin-id", Role: "super-admin"}
	admin := models.Users{ID: "admin-id", Role: "admin"}
	user := models.Users{ID: "user-id", Role: "user"}

	db.Create(&superAdmin)
	db.Create(&admin)
	db.Create(&user)

	filtersData := []models.Filters{
		{USER_ID: superAdmin.ID, City: "City1"},
		{USER_ID: superAdmin.ID, City: "City2"},
		{USER_ID: admin.ID, City: "City3"},
		{USER_ID: admin.ID, City: "City4"},
		{USER_ID: user.ID, City: "City5"},
	}
	for _, filter := range filtersData {
		db.Create(&filter)
	}

	// Define test cases
	tests := []struct {
		name       string
		userID     string
		pageIndex  int
		pageSize   int
		wantCount  int
		wantErr    bool
		totalPages int
		hideUserID bool
	}{
		{
			name:       "Super admin retrieving all filters",
			userID:     superAdmin.ID,
			pageIndex:  1,
			pageSize:   2,
			wantCount:  2,
			wantErr:    false,
			totalPages: 3,
			hideUserID: false,
		},
		{
			name:       "Admin retrieving all filters with hidden USER_ID",
			userID:     admin.ID,
			pageIndex:  1,
			pageSize:   3,
			wantCount:  3,
			wantErr:    false,
			totalPages: 2,
			hideUserID: true,
		},
		{
			name:       "User trying to retrieve all filters (unauthorized)",
			userID:     user.ID,
			pageIndex:  1,
			pageSize:   2,
			wantCount:  0,
			wantErr:    true,
			totalPages: 0,
			hideUserID: false,
		},
		{
			name:       "Super admin retrieving all filters with page size larger than total filters",
			userID:     superAdmin.ID,
			pageIndex:  1,
			pageSize:   10,
			wantCount:  5,
			wantErr:    false,
			totalPages: 1,
			hideUserID: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filters.GetAllFilters(db, tt.userID, tt.pageIndex, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterService.GetAllFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result.Data) != tt.wantCount {
					t.Errorf("FilterService.GetAllFilters() count = %v, want %v", len(result.Data), tt.wantCount)
				}
				if result.TotalPages != tt.totalPages {
					t.Errorf("FilterService.GetAllFilters() totalPages = %v, want %v", result.TotalPages, tt.totalPages)
				}
				if tt.hideUserID {
					for _, filter := range result.Data {
						if filter.USER_ID != "" {
							t.Errorf("FilterService.GetAllFilters() USER_ID should be hidden, got %v", filter.USER_ID)
						}
					}
				}
			}
		})
	}
}
func TestFilterService_RemoveFilter(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Seed users and filters
	superAdmin := models.Users{ID: "super-admin-id", Role: "super-admin"}
	admin := models.Users{ID: "admin-id", Role: "admin"}
	user := models.Users{ID: "user-id", Role: "user"}

	db.Save(&superAdmin)
	db.Save(&admin)
	db.Save(&user)

	filtersData := []models.Filters{
		{ID: "filter1", USER_ID: superAdmin.ID, City: "City1"},
		{ID: "filter2", USER_ID: admin.ID, City: "City2"},
		{ID: "filter3", USER_ID: user.ID, City: "City3"},
	}

	// Utility to reset test data
	resetTestDB := func(db *gorm.DB, filters []models.Filters) error {
		if err := db.Exec("DELETE FROM filters").Error; err != nil {
			return err
		}
		for _, filter := range filters {
			if err := db.Save(&filter).Error; err != nil {
				return err
			}
		}
		return nil
	}

	// Define test cases
	tests := []struct {
		name        string
		userID      string
		filterID    string
		wantErr     bool
		expectedErr string
	}{
		{"Super-admin deletes any filter", superAdmin.ID, "filter1", false, ""},
		{"Admin deletes their own filter", admin.ID, "filter2", false, ""},
		{"Admin tries to delete another user's filter", admin.ID, "filter3", true, "unauthorized to delete this filter"},
		{"User deletes their own filter", user.ID, "filter3", false, ""},
		{"User tries to delete another user's filter", user.ID, "filter1", true, "unauthorized to delete this filter"},
		{"Super-admin tries to delete non-existent filter", superAdmin.ID, "non-existent-filter", true, "filter not found"},
		{"Regular user tries to delete non-existent filter", user.ID, "non-existent-filter", true, "filter not found"},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset database before each test
			if err := resetTestDB(db, filtersData); err != nil {
				t.Fatalf("Failed to reset test database: %v", err)
			}

			err := filters.RemoveFilter(db, tt.userID, tt.filterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterService.RemoveFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("FilterService.RemoveFilter() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			// Verify the filter is deleted when no error is expected
			if !tt.wantErr {
				var filter models.Filters
				if err := db.Where("id = ?", tt.filterID).First(&filter).Error; err == nil {
					t.Errorf("FilterService.RemoveFilter() failed to delete filter: %v", tt.filterID)
				}
			}
		})
	}
}
