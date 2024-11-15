package repositories

import (
	"Crawlzilla/models"
	"fmt"

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
	if filter.ID == "" {
		// No ID provided, create a new filter
		if err := db.Create(filter).Error; err != nil {
			return false, err
		}
		return true, nil // Created new filter
	}

	// ID is provided, try to find the existing filter by ID
	var existingFilter models.Filters
	if err := db.Where("id = ?", filter.ID).First(&existingFilter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("id is not found: %v", err) // Created new filter
		}
		return false, err // An unexpected error occurred
	}

	// Record exists, update it
	if err := db.Model(&existingFilter).Updates(filter).Error; err != nil {
		return false, err
	}
	return true, nil // Updated existing filter
}

// RemoveFilter removes a filter by its ID
func RemoveFilter(db *gorm.DB, id string) (bool, error) {
	if err := db.Delete(&models.Filters{}, "id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}
