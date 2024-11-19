package repositories

import (
	"Crawlzilla/models"
	"errors"

	"gorm.io/gorm"
)

func CreateOrUpdateFilter(db *gorm.DB, filter *models.Filters) error {

	// Step 2: Check if the filter already exists only if filter.ID is not empty
	if filter.ID != "" {
		var existingFilter models.Filters
		if err := db.Where("id = ?", filter.ID).First(&existingFilter).Error; err == nil {
			// If the filter exists, update its fields
			existingFilter.City = filter.City
			existingFilter.Neighborhood = filter.Neighborhood
			existingFilter.Reference = filter.Reference
			existingFilter.CategoryType = filter.CategoryType
			existingFilter.PropertyType = filter.PropertyType
			existingFilter.Sort = filter.Sort
			existingFilter.Order = filter.Order
			existingFilter.MinArea = filter.MinArea
			existingFilter.MaxArea = filter.MaxArea
			existingFilter.MinPrice = filter.MinPrice
			existingFilter.MaxPrice = filter.MaxPrice
			existingFilter.MinRent = filter.MinRent
			existingFilter.MaxRent = filter.MaxRent
			existingFilter.MinRoom = filter.MinRoom
			existingFilter.MaxRoom = filter.MaxRoom
			existingFilter.MinFloorNumber = filter.MinFloorNumber
			existingFilter.MaxFloorNumber = filter.MaxFloorNumber
			existingFilter.HasElevator = filter.HasElevator
			existingFilter.HasStorage = filter.HasStorage
			existingFilter.HasParking = filter.HasParking
			existingFilter.HasBalcony = filter.HasBalcony

			// Save the updated filter
			if err := db.Save(&existingFilter).Error; err != nil {
				return err
			}
			return err
		}

	}

	return db.Create(&filter).Error
}

// GetFiltersByUserID fetches filters for a user with pagination
func GetFiltersByUserID(db *gorm.DB, userID string, pageIndex, pageSize int) ([]models.Filters, int64, error) {
	var filters []models.Filters
	var totalRecords int64

	// Count total filters for the user
	if err := db.Model(&models.Filters{}).Where("user_id = ?", userID).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Fetch filters with pagination
	offset := (pageIndex - 1) * pageSize
	if err := db.Where("user_id = ?", userID).Limit(pageSize).Offset(offset).Find(&filters).Error; err != nil {
		return nil, 0, err
	}

	return filters, totalRecords, nil
}

// GetFiltersForAllUsers fetches filters for all users with pagination
func GetFiltersForAllUsers(db *gorm.DB, offset, limit int) ([]models.Filters, int64, error) {
	var filters []models.Filters
	var totalRecords int64

	// Count total filters
	if err := db.Model(&models.Filters{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Fetch filters with pagination
	if err := db.Limit(limit).Offset(offset).Find(&filters).Error; err != nil {
		return nil, 0, err
	}

	return filters, totalRecords, nil
}

// RemoveFilter deletes a filter by ID
func RemoveFilter(db *gorm.DB, filterID string) error {
	if err := db.Delete(&models.Filters{}, "id = ?", filterID).Error; err != nil {
		return err
	}
	return nil
}

// RemoveAllFilters deletes all filters by userID
func RemoveAllFilters(db *gorm.DB, userID string) error {
	// Use a WHERE query to target rows with the given userID
	if err := db.Where("user_id = ?", userID).Delete(&models.Filters{}).Error; err != nil {
		return err
	}
	return nil
}

// GetFilterByID retrieves a filter by its ID
func GetFilterByID(db *gorm.DB, filterID string) (*models.Filters, error) {
	var filter models.Filters
	if err := db.Where("id = ?", filterID).First(&filter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	// Increment the usage_count column
	if err := db.Model(&models.Filters{}).Where("id = ?", filterID).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error; err != nil {
		return nil, err
	}
	return &filter, nil
}
