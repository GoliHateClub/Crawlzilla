package filters

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// FilterData represents the paginated filter data returned by the service
type FilterData struct {
	Data  []models.Filters `json:"data"`
	Pages int              `json:"pages"`
	Page  int              `json:"page"`
}

// GetAllFilters retrieves paginated filters with role-based visibility
func GetAllFilters(db *gorm.DB, page int, pageSize int, role string) (FilterData, error) {
	// SuperAdmin can see user information, while Admin and NormalUser cannot
	includeUserData := (role == "super-admin")

	// Fetch filters with pagination
	filters, totalRecords, err := repositories.GetAllFiltersPaginated(db, page, pageSize, includeUserData)
	if err != nil {
		return FilterData{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare and return paginated response
	return FilterData{
		Data:  filters,
		Pages: totalPages,
		Page:  page,
	}, nil
}

// CreateOrUpdateFilte creates or updates a filter
func CreateOrUpdateFilte(db *gorm.DB, filter models.Filters) (bool, error) {
	return repositories.CreateOrUpdateFilter(db, &filter)
}

// RemoveFilter removes a filter by its ID
func RemoveFilter(db *gorm.DB, id string) (bool, error) {
	return repositories.RemoveFilter(db, id)
}
