package repositories

import (
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// GetAllFiltersPaginated retrieves filters with pagination and optionally includes user data based on role
func GetAllFiltersPaginated(db *gorm.DB, page int, pageSize int, includeUserData bool) ([]models.Filters, int64, error) {
	var filters []models.Filters
	var totalRecords int64

	// Count total records for pagination
	if err := db.Model(&models.Filters{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Query to get paginated results
	query := db.Offset(offset).Limit(pageSize)
	if includeUserData {
		query = query.Preload("USER") // Preload user data for super admin
	}
	err := query.Find(&filters).Error

	return filters, totalRecords, err
}

// CreateOrUpdateFilter creates a new filter or updates an existing one based on the ID
func CreateOrUpdateFilter(db *gorm.DB, filter *models.Filters) (bool, error) {
	// Fetch the existing filter using the ID; if none found, GORM will automatically insert a new record.
	tx := db.Where("id = ?", filter.ID).FirstOrCreate(&filter)

	if tx.Error != nil {
		return false, tx.Error
	}

	if tx.RowsAffected == 0 {
		// The filter was found, and no new row was created; update the existing record.
		if err := db.Save(filter).Error; err != nil {
			return false, err
		}
		return true, nil // Updated existing filter
	}
	return true, nil // Created new filter, or updated existing one without changes
}

// RemoveFilter removes a filter by its ID
func RemoveFilter(db *gorm.DB, id string) (bool, error) {
	if err := db.Delete(&models.Filters{}, "id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}
