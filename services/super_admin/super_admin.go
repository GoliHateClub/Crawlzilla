package super_admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func IsSuperAdmin(userID int64) bool {
	superAdminId, err := strconv.ParseInt(os.Getenv("SUPER_ADMIN_ID"), 10, 64)
	if err != nil {
		return false
	}
	fmt.Println(superAdminId)
	return superAdminId == userID
}

// ValidateAdData validates the data fields in a Ads
func ValidateAdData(result *models.Ads) error {
	var validationErrors []string

	// Field-specific validation
	if err := validateTitle(result.Title); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateLocationURL(result.LocationURL); err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validateURL(result.LocationURL); err != nil { // Use validateURL here
		validationErrors = append(validationErrors, err.Error())
	}
	if err := validatePrice(result.Price); err != nil {
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

// validateLocationURL checks if the location URL is valid
func validateLocationURL(url string) error {
	if len(url) > 255 {
		return errors.New("location URL length exceeds 255 characters")
	}
	return nil
}

// validatePrice checks if the price is non-negative
func validatePrice(price int) error {
	if price < 0 {
		return errors.New("price cannot be negative")
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

// validateURL checks if the provided URL is a valid URL or not.
func validateURL(url string) error {
	if url == "" {
		return errors.New("URL cannot be empty")
	}
	if len(url) > 255 {
		return errors.New("URL length exceeds 255 characters")
	}
	return nil
}

// CreateAd attempts to save the ad, letting GORM handle model validation constraints
func CreateAd(result *models.Ads, database *gorm.DB) error {
	if result == nil {
		return fmt.Errorf("result cannot be nil")
	}

	if err := ValidateAdData(result); err != nil {
		return err
	}

	if _, err := repositories.CreateAd(database, result); err != nil {
		log.Fatalf("Failed to add data: %v", err)
	} else {
		fmt.Println("Data has been added to the DB successfully!")
	}
	return nil
}

// RemoveAdByID removes an advertisement by its ID
func RemoveAdByID(id string, database *gorm.DB) error {
	var ad models.Ads
	if err := database.First(&ad, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ad not found")
		}
		return err
	}

	if err := database.Delete(&ad).Error; err != nil {
		return err
	}

	return nil
}
