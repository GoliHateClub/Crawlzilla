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

type FilterService struct {
	filterRepo repositories.FilterRepository
}

func NewFilterService(repo repositories.FilterRepository) *FilterService {
	return &FilterService{filterRepo: repo}
}

// GetFiltersByUserID retrieves filters for a user with pagination
func (s *FilterService) GetFiltersByUserID(db *gorm.DB, userID string, pageIndex, pageSize int) (PaginatedFilters, error) {
	// Validate page index
	if pageIndex < 1 {
		return PaginatedFilters{}, errors.New("pageIndex must be greater than 0")
	}

	// Fetch filters and total records from the repository
	filters, totalRecords, err := s.filterRepo.GetFiltersByUserID(db, userID, pageIndex, pageSize)
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

// GetAllFilters retrieves filters based on the user's role
func (s *FilterService) GetAllFilters(db *gorm.DB, userID string, pageIndex, pageSize int) (PaginatedFilters, error) {
	// Validate page index
	if pageIndex < 1 {
		return PaginatedFilters{}, errors.New("pageIndex must be greater than 0")
	}

	// Fetch the user to determine their role
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return PaginatedFilters{}, err
	}
	if user == nil {
		return PaginatedFilters{}, errors.New("user not found")
	}

	// Role-based logic
	var filters []models.Filters
	var totalRecords int64
	offset := (pageIndex - 1) * pageSize

	if user.Role == "super-admin" {
		// Fetch all filters for all users
		filters, totalRecords, err = s.filterRepo.GetFiltersForAllUsers(db, offset, pageSize)
		if err != nil {
			return PaginatedFilters{}, err
		}

	} else if user.Role == "admin" {
		// Fetch all filters but hide the USER_ID field
		filters, totalRecords, err = s.filterRepo.GetFiltersForAllUsers(db, offset, pageSize)
		if err != nil {
			return PaginatedFilters{}, err
		}

		// Hide the USER_ID field
		for i := range filters {
			filters[i].USER_ID = ""
		}

	} else if user.Role == "user" {
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
func (s *FilterService) CreateOrUpdateFilter(db *gorm.DB, filter models.Filters) (bool, error) {
	// Step 1: Validate all fields
	if err := s.validateFilterFields(filter); err != nil {
		return false, err
	}

	// Step 2: Check if the user exists
	user, err := repositories.GetUserByID(db, filter.USER_ID)
	if err != nil {
		return false, err
	}

	// Step 3: If the user does not exist, create a new user with "user" role
	if user == nil {
		newUser := &models.Users{
			ID:          filter.USER_ID,
			Telegram_ID: "", // Set Telegram_ID as needed
			Role:        "user",
		}
		if err := repositories.CreateUser(db, newUser); err != nil {
			return false, err
		}
	}

	// Step 4: Add the filter to the database
	if err := s.filterRepo.CreateOrUpdateFilter(db, &filter); err != nil {
		return false, err
	}

	return true, nil
}

// RemoveFilter removes a filter based on the user's role and filter ownership
func (s *FilterService) RemoveFilter(db *gorm.DB, userID, filterID string) error {
	// Fetch the user to determine their role
	user, err := repositories.GetUserByID(db, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Fetch the filter to check ownership
	filter, err := s.filterRepo.GetFilterByID(db, filterID)
	if err != nil {
		return err
	}
	if filter == nil {
		return errors.New("filter not found")
	}

	// Role-based logic for deletion
	if user.Role == "super-admin" {
		// Super-admin can delete any filter
		return s.filterRepo.RemoveFilter(db, filterID)
	} else if user.Role == "admin" || user.Role == "user" {
		// Admin or user can delete only their own filters
		if filter.USER_ID != userID {
			return errors.New("unauthorized to delete this filter")
		}
		return s.filterRepo.RemoveFilter(db, filterID)
	}

	return errors.New("role not authorized to delete filters")
}

// validateFilterFields validates all fields in the filter
func (s *FilterService) validateFilterFields(filter models.Filters) error {
	if err := s.validateArea(filter.MinArea, filter.MaxArea); err != nil {
		return err
	}
	if err := s.validatePrice(filter.MinPrice, filter.MaxPrice); err != nil {
		return err
	}
	if err := s.validateRent(filter.MinRent, filter.MaxRent); err != nil {
		return err
	}
	if err := s.validateRoom(filter.MinRoom, filter.MaxRoom); err != nil {
		return err
	}
	if err := s.validateFloorNumber(filter.MinFloorNumber, filter.MaxFloorNumber); err != nil {
		return err
	}
	// Validate optional string fields only if they are provided
	if filter.City != "" {
		if err := s.validateCity(filter.City); err != nil {
			return err
		}
	}
	if filter.Neighborhood != "" {
		if err := s.validateNeighborhood(filter.Neighborhood); err != nil {
			return err
		}
	}
	if filter.Reference != "" {
		if err := s.validateReference(filter.Reference); err != nil {
			return err
		}
	}
	if filter.CategoryType != "" {
		if err := s.validateCategoryType(filter.CategoryType); err != nil {
			return err
		}
	}
	if filter.PropertyType != "" {
		if err := s.validatePropertyType(filter.PropertyType); err != nil {
			return err
		}
	}
	if filter.Sort != "" || filter.Order != "" {
		if err := s.validateSortOrder(filter.Sort, filter.Order); err != nil {
			return err
		}
	}

	// No validation needed for booleans (HasElevator, HasStorage, HasParking, HasBalcony)

	return nil
}

// Individual validation functions
func (s *FilterService) validateCity(city string) error {
	if city == "" {
		return errors.New("city cannot be empty")
	}
	return nil
}

func (s *FilterService) validateNeighborhood(neighborhood string) error {
	if neighborhood == "" {
		return errors.New("neighborhood cannot be empty")
	}
	return nil
}

func (s *FilterService) validateReference(reference string) error {
	if reference == "" {
		return errors.New("reference cannot be empty")
	}
	return nil
}

func (s *FilterService) validateCategoryType(categoryType string) error {
	if categoryType == "" {
		return errors.New("category type cannot be empty")
	}
	return nil
}

func (s *FilterService) validatePropertyType(propertyType string) error {
	if propertyType == "" {
		return errors.New("property type cannot be empty")
	}
	return nil
}

func (s *FilterService) validateSortOrder(sort, order string) error {
	if sort == "" || order == "" {
		return errors.New("both sort and order must be provided if one is specified")
	}
	return nil
}

func (s *FilterService) validateArea(minArea, maxArea int) error {
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

func (s *FilterService) validatePrice(minPrice, maxPrice int) error {
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

func (s *FilterService) validateRent(minRent, maxRent int) error {
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

func (s *FilterService) validateRoom(minRoom, maxRoom int) error {
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

func (s *FilterService) validateFloorNumber(minFloor, maxFloor int) error {
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
