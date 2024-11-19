package services_tests

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/services/filters"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
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
				Title:        "Test",
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
				Title:        "Test",
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
				Title:   "Test",
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
				Title:    "Test",
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
				Title:    "Test",
				MinPrice: 1000,
				MaxPrice: 500,
			},
			wantID:  "", // Expect no ID since it’s invalid
			wantErr: true,
		},
		{
			name: "Invalid filter with no user id",
			filter: models.Filters{
				Title:    "Test",
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
				Title:    "Test",
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
				Title:   "Test",
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
		{USER_ID: userID, City: "City1", Title: "Title1"},
		{USER_ID: userID, City: "City2", Title: "Title2"},
		{USER_ID: userID, City: "City3", Title: "Title3"},
		{USER_ID: userID, City: "City4", Title: "Title4"},
		{USER_ID: userID, City: "City5", Title: "Title5"},
		{USER_ID: userID, City: "City6", Title: "Title6"},
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
	user := models.Users{ID: "user-id", Role: models.RoleUser}

	db.Create(&superAdmin)
	db.Create(&admin)
	db.Create(&user)

	filtersData := []models.Filters{
		{USER_ID: superAdmin.ID, City: "City1", Title: "Title1"},
		{USER_ID: superAdmin.ID, City: "City2", Title: "Title2"},
		{USER_ID: admin.ID, City: "City3", Title: "Title3"},
		{USER_ID: admin.ID, City: "City4", Title: "Title4"},
		{USER_ID: user.ID, City: "City5", Title: "Title5"},
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

// TestRemoveFilter tests the RemoveFilter function
func TestRemoveFilter(t *testing.T) {
	// Set up in-memory database
	db, _ := setupTestDB()

	// Create test users
	superAdmin := models.Users{Role: models.RoleSuperAdmin}
	admin := models.Users{Role: models.RoleAdmin}
	user := models.Users{Role: models.RoleUser}
	otherUser := models.Users{Role: models.RoleUser}

	// Insert users into the database
	err := db.Create(&superAdmin).Error
	assert.NoError(t, err)
	err = db.Create(&admin).Error
	assert.NoError(t, err)
	err = db.Create(&user).Error
	assert.NoError(t, err)
	err = db.Create(&otherUser).Error
	assert.NoError(t, err)

	// Create test filter entries
	filter1 := models.Filters{ID: "filter1", USER_ID: otherUser.ID, City: "Test City", CategoryType: "Apartment", Title: "Title1"}
	filter2 := models.Filters{ID: "filter2", USER_ID: admin.ID, City: "Test City", CategoryType: "Apartment", Title: "Title2"}
	filter3 := models.Filters{ID: "filter3", USER_ID: user.ID, City: "Test City", CategoryType: "Apartment", Title: "Title3"}
	filter4 := models.Filters{ID: "filter4", USER_ID: otherUser.ID, City: "Test City", CategoryType: "Apartment", Title: "Title4"}

	// Insert filters into the database

	err = db.Create(&filter1).Error
	assert.NoError(t, err)
	err = db.Create(&filter2).Error
	assert.NoError(t, err)
	err = db.Create(&filter3).Error
	assert.NoError(t, err)
	err = db.Create(&filter4).Error
	assert.NoError(t, err)

	// Test: User cannot delete someone else's filter
	err = filters.RemoveFilter(db, user.ID, filter2.ID)
	assert.Error(t, err) // Unauthorized to delete another user's filter

	// Test: Unauthorized role (e.g., guest) should fail
	unauthorizedUser := models.Users{ID: "5", Role: "guest"}
	err = db.Create(&unauthorizedUser).Error
	assert.NoError(t, err)

	err = filters.RemoveFilter(db, unauthorizedUser.ID, filter1.ID)
	assert.Error(t, err) // Unauthorized role should fail

	// Test: Super Admin can delete any filter
	err = filters.RemoveFilter(db, superAdmin.ID, filter4.ID)
	assert.NoError(t, err)

	// Test: Admin can delete their own filter
	err = filters.RemoveFilter(db, admin.ID, filter1.ID)
	assert.Error(t, err) // Admin should not be able to delete another user's filter

	// Test: User can delete their own filter
	err = filters.RemoveFilter(db, user.ID, filter3.ID)
	assert.NoError(t, err)

	// Test: Non-existent filter ID
	err = filters.RemoveFilter(db, user.ID, "non-existent-filter")
	assert.Error(t, err)

}

func TestRemoveAllFilters(t *testing.T) {
	// Set up in-memory database
	db := SetupSearchTestDB()

	// Prepare mock user data
	userID := "user-12345" // Use a UUID or string for the user ID

	// Create some filters for the user
	testFilters := []models.Filters{
		{City: "City1", Sort: "price", Order: "asc", MinArea: 50, MaxArea: 100, MinPrice: 1000, MaxPrice: 5000, USER_ID: userID},
		{City: "City2", Sort: "area", Order: "desc", MinArea: 70, MaxArea: 120, MinPrice: 2000, MaxPrice: 6000, USER_ID: userID},
	}

	// Insert filters into the database
	for _, filter := range testFilters {
		err := repositories.CreateOrUpdateFilter(db, &filter)
		assert.NoError(t, err, "Creating filter should not return an error")
	}

	// Verify that filters are created
	var createdFilters []models.Filters
	err := db.Where("user_id = ?", userID).Find(&createdFilters).Error
	assert.NoError(t, err)
	assert.Equal(t, len(testFilters), len(createdFilters), "The number of filters should match")

	// Call RemoveAllFilters to remove the filters
	err = filters.RemoveAllFilters(db, userID)
	assert.NoError(t, err, "Removing filters should not return an error")

	// Verify that all filters are removed
	var remainingFilters []models.Filters
	err = db.Where("user_id = ?", userID).Find(&remainingFilters).Error
	assert.NoError(t, err)
	assert.Equal(t, 0, len(remainingFilters), "There should be no filters left for the user")
}
