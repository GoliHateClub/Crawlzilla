package search

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"fmt"

	"gorm.io/gorm"
)

type PaginatedAds struct {
	Data  []models.Ads `json:"data"`  // Array of filtered ads
	Pages int          `json:"pages"` // Total number of pages
	Page  int          `json:"page"`  // Current page number
}

// GetFilteredAdsPaginatedService retrieves filtered ads with pagination, sorting, and filtering.
func GetFilteredAds(db *gorm.DB, filterID string, page, pageSize int) (PaginatedAds, error) {
	// Retrieve the filter by ID
	filter, err := repositories.GetFilterByID(db, filterID)
	if err != nil {
		return PaginatedAds{}, err
	}

	// Start building the query for ads
	query := db.Model(&models.Ads{})
	if filter.City != "" {
		query = query.Where("city = ?", filter.City)
	}
	if filter.Neighborhood != "" {
		query = query.Where("neighborhood = ?", filter.Neighborhood)
	}
	if filter.Reference != "" {
		query = query.Where("reference = ?", filter.Reference)
	}
	if filter.CategoryType != "" {
		query = query.Where("category_type = ?", filter.CategoryType)
	}
	if filter.PropertyType != "" {
		query = query.Where("property_type = ?", filter.PropertyType)
	}
	if filter.MinArea > 0 {
		query = query.Where("area >= ?", filter.MinArea)
	}
	if filter.MaxArea > 0 {
		query = query.Where("area <= ?", filter.MaxArea)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.MinRent > 0 {
		query = query.Where("rent >= ?", filter.MinRent)
	}
	if filter.MaxRent > 0 {
		query = query.Where("rent <= ?", filter.MaxRent)
	}
	if filter.MinRoom > 0 {
		query = query.Where("room >= ?", filter.MinRoom)
	}
	if filter.MaxRoom > 0 {
		query = query.Where("room <= ?", filter.MaxRoom)
	}
	if filter.MinFloorNumber > 0 {
		query = query.Where("floor_number >= ?", filter.MinFloorNumber)
	}
	if filter.MaxFloorNumber > 0 {
		query = query.Where("floor_number <= ?", filter.MaxFloorNumber)
	}
	if filter.HasElevator {
		query = query.Where("has_elevator = ?", true)
	}
	if filter.HasStorage {
		query = query.Where("has_storage = ?", true)
	}
	if filter.HasParking {
		query = query.Where("has_parking = ?", true)
	}
	if filter.HasBalcony {
		query = query.Where("has_balcony = ?", true)
	}

	// Add sorting if specified in the filter
	if filter.Sort != "" && filter.Order != "" {
		validSortColumns := map[string]bool{
			"price":        true,
			"rent":         true,
			"area":         true,
			"room":         true,
			"floor_number": true,
			"visit_count":  true,
			"created_at":   true,
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		// Validate sort column and order
		if validSortColumns[filter.Sort] && validOrders[filter.Order] {
			sortExpression := fmt.Sprintf("%s %s", filter.Sort, filter.Order)
			query = query.Order(sortExpression)
		} else {
			return PaginatedAds{}, fmt.Errorf("invalid sort column or order")
		}
	}

	// Count total records matching the filter
	totalRecords, err := repositories.CountFilteredAds(db, query)
	if err != nil {
		return PaginatedAds{}, err
	}

	// Fetch the paginated ads
	ads, err := repositories.GetFilteredAds(db, query, page, pageSize)
	if err != nil {
		return PaginatedAds{}, err
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare the paginated response
	result := PaginatedAds{
		Data:  ads,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}

// GetMostFilteredAds retrieves ads based on the most-used filter
func GetMostFilteredAds(db *gorm.DB, page, pageSize int) (PaginatedAds, error) {
	var mostUsedFilter models.Filters

	// Step 1: Find the most-used filter
	err := db.Model(&models.Filters{}).
		Select("*").
		Order("usage_count DESC"). // Assuming a usage_count column exists
		First(&mostUsedFilter).Error
	if err != nil {
		return PaginatedAds{}, fmt.Errorf("failed to find the most-used filter: %w", err)
	}

	// Step 2: Build the ads query using the most-used filter
	query := db.Model(&models.Ads{})
	if mostUsedFilter.City != "" {
		query = query.Where("city = ?", mostUsedFilter.City)
	}
	if mostUsedFilter.Neighborhood != "" {
		query = query.Where("neighborhood = ?", mostUsedFilter.Neighborhood)
	}
	if mostUsedFilter.Reference != "" {
		query = query.Where("reference = ?", mostUsedFilter.Reference)
	}
	if mostUsedFilter.CategoryType != "" {
		query = query.Where("category_type = ?", mostUsedFilter.CategoryType)
	}
	if mostUsedFilter.PropertyType != "" {
		query = query.Where("property_type = ?", mostUsedFilter.PropertyType)
	}
	if mostUsedFilter.MinArea > 0 {
		query = query.Where("area >= ?", mostUsedFilter.MinArea)
	}
	if mostUsedFilter.MaxArea > 0 {
		query = query.Where("area <= ?", mostUsedFilter.MaxArea)
	}
	if mostUsedFilter.MinPrice > 0 {
		query = query.Where("price >= ?", mostUsedFilter.MinPrice)
	}
	if mostUsedFilter.MaxPrice > 0 {
		query = query.Where("price <= ?", mostUsedFilter.MaxPrice)
	}
	if mostUsedFilter.MinRent > 0 {
		query = query.Where("rent >= ?", mostUsedFilter.MinRent)
	}
	if mostUsedFilter.MaxRent > 0 {
		query = query.Where("rent <= ?", mostUsedFilter.MaxRent)
	}
	if mostUsedFilter.MinRoom > 0 {
		query = query.Where("room >= ?", mostUsedFilter.MinRoom)
	}
	if mostUsedFilter.MaxRoom > 0 {
		query = query.Where("room <= ?", mostUsedFilter.MaxRoom)
	}
	if mostUsedFilter.MinFloorNumber > 0 {
		query = query.Where("floor_number >= ?", mostUsedFilter.MinFloorNumber)
	}
	if mostUsedFilter.MaxFloorNumber > 0 {
		query = query.Where("floor_number <= ?", mostUsedFilter.MaxFloorNumber)
	}
	if mostUsedFilter.HasElevator {
		query = query.Where("has_elevator = ?", true)
	}
	if mostUsedFilter.HasStorage {
		query = query.Where("has_storage = ?", true)
	}
	if mostUsedFilter.HasParking {
		query = query.Where("has_parking = ?", true)
	}
	if mostUsedFilter.HasBalcony {
		query = query.Where("has_balcony = ?", true)
	}

	// Step 3: Count total records matching the most-used filter
	totalRecords, err := repositories.CountFilteredAds(db, query)
	if err != nil {
		return PaginatedAds{}, err
	}

	// Step 4: Fetch paginated ads
	ads, err := repositories.GetFilteredAds(db, query, page, pageSize)
	if err != nil {
		return PaginatedAds{}, err
	}

	// Step 5: Calculate total pages
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Step 6: Prepare the paginated response
	result := PaginatedAds{
		Data:  ads,
		Pages: totalPages,
		Page:  page,
	}

	return result, nil
}
