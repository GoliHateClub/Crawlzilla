package services_tests

import (
	"Crawlzilla/database/repositories"
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

	// Create a mock repository
	mockRepo := repositories.NewFilterRepository()
	s := filters.NewFilterService(mockRepo)

	// Define test cases
	tests := []struct {
		name    string
		filter  models.Filters
		want    bool
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
			want:    true,
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
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid filter with min area greater than max area",
			filter: models.Filters{
				USER_ID: "user-id-4",
				MinArea: 150,
				MaxArea: 100,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid filter with negative price values",
			filter: models.Filters{
				USER_ID:  "user-id-5",
				MinPrice: -500,
				MaxPrice: 1000,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid filter with min price greater than max price",
			filter: models.Filters{
				USER_ID:  "user-id-6",
				MinPrice: 1000,
				MaxPrice: 500,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid filter with no user id",
			filter: models.Filters{
				MinPrice: 1000,
				MaxPrice: 500,
			},
			want:    false,
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
			want:    true,
			wantErr: false,
		},
		{
			name: "Valid filter with no optional fields provided",
			filter: models.Filters{
				USER_ID: "user-id-8",
			},
			want:    true,
			wantErr: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateOrUpdateFilter(db, tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterService.CreateOrUpdateFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilterService.CreateOrUpdateFilter() = %v, want %v", got, tt.want)
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
	mockRepo := repositories.NewFilterRepository()
	s := filters.NewFilterService(mockRepo)

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
			result, err := s.GetFiltersByUserID(db, tt.userID, tt.pageIndex, tt.pageSize)
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

	// Create a mock repository and service
	mockRepo := repositories.NewFilterRepository()
	s := filters.NewFilterService(mockRepo)

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
			result, err := s.GetAllFilters(db, tt.userID, tt.pageIndex, tt.pageSize)
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
