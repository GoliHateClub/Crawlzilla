package repositories

import (
	"Crawlzilla/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// AddAd adds a new scrap result to the database if it doesn't already exist
func CreateAd(result *models.Ads, database *gorm.DB) (string, error) {
	// Manually call BeforeCreate to generate the hash before querying the database
	if err := result.BeforeCreate(database); err != nil {
		return "", fmt.Errorf("failed to generate hash: %v", err)
	}

	// Check if a record with the same hash already exists
	var existing models.Ads
	err := database.Where("hash = ?", result.Hash).First(&existing).Error

	if err == nil {
		// Hash exists, log that it already exists and skip insertion
		fmt.Println("Record with hash already exists, skipping insert.")
		return "", errors.New("hash of data existed!")
	}

	// If no record with the hash was found, proceed with the creation
	if err == gorm.ErrRecordNotFound {
		if err2 := database.Create(result).Error; err2 != nil {
			// Handle error if insert fails
			return "", errors.New("cant't add to database")
		}
		fmt.Println("Record added to DB successfully!")
		return result.ID, nil
	}
	return "", err
}

// GetAllAds retrieves all scrap results
func GetAllAds(database *gorm.DB) ([]models.Ads, error) {
	var results []models.Ads
	err := database.Find(&results).Error
	return results, err
}

// GetAdByID retrieves a scrap result by ID
func GetAdByID(id string, database *gorm.DB) (models.Ads, error) {
	var result models.Ads
	err := database.First(&result, id).Error
	return result, err
}

// DeleteAd deletes a scrap result by ID
func DeleteAdById(id string, database *gorm.DB) error {
	return database.Delete(&models.Ads{}, id).Error
}
