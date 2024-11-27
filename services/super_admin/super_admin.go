package super_admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func IsSuperAdmin(telegram_id int64) bool {
	superAdminId, err := strconv.ParseInt(os.Getenv("SUPER_ADMIN_ID"), 10, 64)
	if err != nil {
		return false
	}
	return superAdminId == telegram_id
}

// ValidateAdData validates the data fields in a Ads
func ValidateAdData(result *models.Ads) error {
	var validationErrors []string

	// Field-specific validation
	if err := validateTitle(result.Title); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.Price); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.Rent); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.Area); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.Room); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.FloorNumber); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateNumber(result.TotalFloors); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateCategory(result.CategoryType); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateProperty(result.PropertyType); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validatePhoneNumber(result.ContactNumber); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateCoordinates(result.Latitude, result.Longitude); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}

	if len(validationErrors) > 0 {
		return errors.New("validation failed: " + strings.Join(validationErrors, "; "))
	}

	return nil
}

// validateTitle checks if the title is valid
func validateTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if len(title) > 50 {
		return errors.New("title length exceeds 50 characters")
	}
	return nil
}

// validateNumber checks if the number is non-negative
func validateNumber(number int) error {
	if number < 0 {
		return errors.New("number cannot be negative")
	}
	return nil
}

// validateCategory checks if the category is valid
func validateCategory(categoryType string) error {
	if categoryType != "فروش" && categoryType != "رهن اجاره" {
		return errors.New("category type is wrong")
	}
	return nil
}

// validateCategory checks if the category is valid
func validateProperty(propertyType string) error {
	if propertyType != "آپارتمانی" && propertyType != "ویلایی" {
		return errors.New("property type is wrong")
	}
	return nil
}

// validatePhoneNumber checks if phoneNumber is all numbers, starts with 0, and is exactly 11 characters long
func validatePhoneNumber(phoneNumber string) error {
	// Define the regex pattern
	pattern := `^0(\d{10})?$` // Starts with 0, followed by exactly 10 digits (total 11 characters)

	// Compile the regex
	re := regexp.MustCompile(pattern)

	// Validate the phone number
	if !re.MatchString(phoneNumber) {
		return errors.New("phone number is invalid")
	}
	return nil
}

// validateCoordinates checks if latitude and longitude are within valid ranges
func validateCoordinates(lat, long float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude %v is out of range (-90 to 90)", lat)
	}
	if long < -180 || long > 180 {
		return fmt.Errorf("longitude %v is out of range (-180 to 180)", long)
	}
	return nil
}

// CreateAd attempts to save the ad, letting GORM handle model validation constraints
func CreateAd(database *gorm.DB, result *models.Ads) error {
	if result == nil {
		return fmt.Errorf("result cannot be nil")
	}

	if err := ValidateAdData(result); err != nil {
		return err
	}

	// generate URL
	result.URL = "super-admin-" + utils.GenerateRandomNumber(5)
	// generate locationURL
	result.LocationURL = fmt.Sprintf("https://balad.ir/location?latitude=%v&longitude=%v", result.Latitude, result.Longitude)
	// generate reference
	result.Reference = "admin"

	// generate Category and Property Type
	if result.CategoryType == "فروش" {
		result.CategoryType = "sell"
	} else {
		result.CategoryType = "rent"
	}
	if result.PropertyType == "آپارتمانی" {
		result.PropertyType = "house"
	} else {
		result.PropertyType = "vila"
	}

	if _, err := repositories.CreateAd(database, result); err != nil {
		log.Printf("Failed to add data: %v", err)
	} else {
		fmt.Println("Data has been added to the DB successfully!")
	}
	return nil
}

// RemoveAdByID removes an advertisement by its ID
func RemoveAdByID(database *gorm.DB, id string) error {
	ad, err := repositories.GetAdByID(database, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ad not found")
		}
		return err
	}

	if err := repositories.DeleteAdById(database, ad.ID); err != nil {
		return fmt.Errorf("failed to delete ad: %v", err)
	}

	log.Printf("Ad with ID %s successfully deleted", id)
	return nil
}
