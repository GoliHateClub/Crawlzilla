package filters

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"errors"
	"math"

	"gorm.io/gorm"
)

// PaginatedFilters represents the response structure
type PaginatedFilters struct {
	Data       []models.Filters `json:"data"`
	TotalPages int              `json:"total_pages"`
	PageIndex  int              `json:"page_index"`
}

// GetFiltersByUserID retrieves filters for a user with pagination
func GetFiltersByUserID(db *gorm.DB, userID string, pageIndex, pageSize int) (PaginatedFilters, error) {
	// Validate page index
	if pageIndex < 1 {
		return PaginatedFilters{}, errors.New("pageIndex must be greater than 0")
	}

	// Fetch filters and total records from the repository
	filters, totalRecords, err := repositories.GetFiltersByUserID(db, userID, pageIndex, pageSize)
	if err != nil {
		return PaginatedFilters{}, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return PaginatedFilters{
		Data:       filters,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
	}, nil
}

// GetAllFilters retrieves filters based on the user's role (SUPER ADMIN - ADMIN)
func GetAllFilters(db *gorm.DB, userID string, pageIndex, pageSize int) (PaginatedFilters, error) {
	// Validate page index
	if pageIndex < 1 {
		return PaginatedFilters{}, errors.New("pageIndex must be greater than 0")
	}

	// Fetch the user to determine their role
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return PaginatedFilters{}, err
	}

	// Role-based logic
	var filters []models.Filters
	var totalRecords int64
	offset := (pageIndex - 1) * pageSize

	if user.Role == models.RoleSuperAdmin {
		// Fetch all filters for all users
		filters, totalRecords, err = repositories.GetFiltersForAllUsers(db, offset, pageSize)
		if err != nil {
			return PaginatedFilters{}, err
		}

	} else if user.Role == models.RoleAdmin {
		// Fetch all filters but hide the USER_ID field
		filters, totalRecords, err = repositories.GetFiltersForAllUsers(db, offset, pageSize)
		if err != nil {
			return PaginatedFilters{}, err
		}

		// Hide the USER_ID field
		for i := range filters {
			filters[i].USER_ID = ""
		}

	} else if user.Role == models.RoleUser {
		// Regular users are not authorized
		return PaginatedFilters{}, errors.New("unauthorized access for regular users")
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return PaginatedFilters{
		Data:       filters,
		TotalPages: totalPages,
		PageIndex:  pageIndex,
	}, nil
}

func CreateOrUpdateFilter(db *gorm.DB, filter models.Filters) (string, error) {
	// Step 1: Validate all fields
	if err := validateFilterFields(filter); err != nil {
		return "", err
	}
	if filter.CategoryType == "فروش" {
		filter.CategoryType = "sell"
	} else if filter.CategoryType == "اجاره" {
		filter.CategoryType = "rent"
	} else {
		filter.CategoryType = ""
	}

	if filter.PropertyType == "آپارتمانی" {
		filter.PropertyType = "apartment"
	} else if filter.PropertyType == "ویلایی" {
		filter.PropertyType = "vila"
	} else {
		filter.PropertyType = ""
	}

	if filter.Reference == "دیوار" {
		filter.Reference = "divar"
	} else if filter.Reference == "شیپور" {
		filter.Reference = "sheypoor"
	} else if filter.Reference == "ادمین" {
		filter.Reference = "admin"
	} else {
		filter.Reference = ""
	}

	err := repositories.CreateOrUpdateFilter(db, &filter)
	if err != nil {
		return "", err
	}

	return filter.ID, nil // Return the created filter ID
}

// RemoveFilter removes a filter based on the user's role and filter ownership
func RemoveFilter(db *gorm.DB, userID, filterID string) error {
	// Fetch the user to determine their role
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return err
	}

	// Fetch the filter to check ownership
	filter, err := repositories.GetFilterByID(db, filterID)
	if err != nil {
		return err
	}
	if filter == nil {
		return errors.New("filter not found")
	}

	// Role-based logic for deletion
	if user.Role == models.RoleSuperAdmin {
		// Super-admin can delete any filter
		return repositories.RemoveFilter(db, filterID)
	} else if user.Role == models.RoleAdmin || user.Role == models.RoleUser {
		// Admin or user can delete only their own filters
		if filter.USER_ID != userID {
			return errors.New("unauthorized to delete this filter")
		}
		return repositories.RemoveFilter(db, filterID)
	}

	return errors.New("role not authorized to delete filters")
}

// RemoveAllFilters removes all user's filters (Clear History)
func RemoveAllFilters(db *gorm.DB, userID string) error {
	return repositories.RemoveAllFilters(db, userID)
}

func GetFilterByID(db *gorm.DB, filterID string) (models.Filters, error) {
	var filter models.Filters
	err := db.Where("id = ?", filterID).First(&filter).Error
	if err != nil {
		return models.Filters{}, err
	}
	return filter, nil
}

// validateFilterFields validates all fields in the filter
func validateFilterFields(filter models.Filters) error {
	if err := validateTitle(filter.Title); err != nil {
		return err
	}
	if err := validateArea(filter.MinArea, filter.MaxArea); err != nil {
		return err
	}
	if err := validatePrice(filter.MinPrice, filter.MaxPrice); err != nil {
		return err
	}
	if err := validateRent(filter.MinRent, filter.MaxRent); err != nil {
		return err
	}
	if err := validateRoom(filter.MinRoom, filter.MaxRoom); err != nil {
		return err
	}
	if err := validateFloorNumber(filter.MinFloorNumber, filter.MaxFloorNumber); err != nil {
		return err
	}
	// Validate optional string fields only if they are provided
	if filter.City != "" {
		if err := validateCity(filter.City); err != nil {
			return err
		}
	}
	if filter.Neighborhood != "" {
		if err := validateNeighborhood(filter.Neighborhood); err != nil {
			return err
		}
	}
	if filter.Reference != "" {
		if err := validateReference(filter.Reference); err != nil {
			return err
		}
	}
	if filter.CategoryType != "" {
		if err := validateCategoryType(filter.CategoryType); err != nil {
			return err
		}
	}
	if filter.PropertyType != "" {
		if err := validatePropertyType(filter.PropertyType); err != nil {
			return err
		}
	}
	if filter.Sort != "" || filter.Order != "" {
		if err := validateSortOrder(filter.Sort, filter.Order); err != nil {
			return err
		}
	}

	// No validation needed for booleans (HasElevator, HasStorage, HasParking, HasBalcony)

	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	return nil
}

// Individual validation functions
func validateCity(city string) error {
	if city == "" {
		return errors.New("city cannot be empty")
	}
	return nil
}

func validateNeighborhood(neighborhood string) error {
	if neighborhood == "" {
		return errors.New("neighborhood cannot be empty")
	}
	return nil
}

func validateReference(reference string) error {
	if reference == "" {
		return errors.New("reference cannot be empty")
	}
	return nil
}

func validateCategoryType(categoryType string) error {
	if categoryType == "" {
		return errors.New("category type cannot be empty")
	}
	return nil
}

func validatePropertyType(propertyType string) error {
	if propertyType == "" {
		return errors.New("property type cannot be empty")
	}
	return nil
}

func validateSortOrder(sort, order string) error {
	if sort == "" || order == "" {
		return errors.New("both sort and order must be provided if one is specified")
	}
	return nil
}

func validateArea(minArea, maxArea int) error {
	if minArea != 0 && maxArea != 0 {
		if minArea < 0 || maxArea < 0 {
			return errors.New("area values cannot be negative")
		}
		if minArea > maxArea {
			return errors.New("minArea cannot be greater than maxArea")
		}
	}
	return nil
}

func validatePrice(minPrice, maxPrice int) error {
	if minPrice != 0 && maxPrice != 0 {
		if minPrice < 0 || maxPrice < 0 {
			return errors.New("price values cannot be negative")
		}
		if minPrice > maxPrice {
			return errors.New("minPrice cannot be greater than maxPrice")
		}
	}
	return nil
}

func validateRent(minRent, maxRent int) error {
	if minRent != 0 && maxRent != 0 {
		if minRent < 0 || maxRent < 0 {
			return errors.New("rent values cannot be negative")
		}
		if minRent > maxRent {
			return errors.New("minRent cannot be greater than maxRent")
		}
	}
	return nil
}

func validateRoom(minRoom, maxRoom int) error {
	if minRoom != 0 && maxRoom != 0 {
		if minRoom < 0 || maxRoom < 0 {
			return errors.New("room values cannot be negative")
		}
		if minRoom > maxRoom {
			return errors.New("minRoom cannot be greater than maxRoom")
		}
	}
	return nil
}

func validateFloorNumber(minFloor, maxFloor int) error {
	if minFloor != 0 && maxFloor != 0 {
		if minFloor < 0 || maxFloor < 0 {
			return errors.New("floor number values cannot be negative")
		}
		if minFloor > maxFloor {
			return errors.New("minFloorNumber cannot be greater than maxFloorNumber")
		}
	}
	return nil
}
