package repositories

import (
	"Crawlzilla/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// CreateAd adds a new scrap result to the database if it doesn't already exist
func CreateAd(database *gorm.DB, result *models.Ads) (string, error) {
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
		return "", errors.New("hash of data existed")
	}

	// If no record with the hash was found, proceed with the creation
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err2 := database.Create(&result).Error; err2 != nil {
			// Handle error if insert fails
			return "", errors.New("cant't add to database")
		}
		fmt.Println("Record added to DB successfully!")
		return result.ID, nil
	}
	return "", err
}

// GetAllAds retrieves ads with pagination and selects specific fields
func GetAllAds(db *gorm.DB, page int, pageSize int) ([]models.AdSummary, int64, error) {
	var ads []models.AdSummary
	var totalRecords int64

	// Count total records for pagination info
	if err := db.Model(&models.Ads{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset based on page and pageSize
	offset := (page - 1) * pageSize

	// Retrieve paginated records with specific fields
	err := db.Model(&models.Ads{}).Select("ID", "Title", "ImageURL").Offset(offset).Limit(pageSize).Find(&ads).Error
	return ads, totalRecords, err
}

// GetAdByID retrieves a scrap result by ID
func GetAdByID(database *gorm.DB, id string) (models.Ads, error) {
	var result models.Ads

	// Fetch the ad by ID
	err := database.Where("id = ?", id).First(&result).Error
	if err != nil {
		return result, err
	}

	// Increment the visit_count column
	if err := database.Model(&models.Ads{}).Where("id = ?", id).UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error; err != nil {
		return result, err
	}

	// Fetch the updated record to include the incremented visit_count
	err = database.Where("id = ?", id).First(&result).Error
	return result, err
}

// DeleteAdById deletes a scrap result by ID
func DeleteAdById(database *gorm.DB, id string) error {
	return database.Where("id = ?", id).Delete(&models.Ads{}).Error
}
