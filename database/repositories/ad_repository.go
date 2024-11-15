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

// GetAllAdsPaginated retrieves ads with pagination and selects specific fields
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
	err := database.Where("id = ?", id).First(&result).Error
	if err != nil {
		return result, err
	}
	err = database.Model(&models.Ads{}).Where("id = ?", id).UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error

	return result, err
}

// DeleteAd deletes a scrap result by ID
func DeleteAdById(database *gorm.DB, id string) error {
	return database.Where("id = ?", id).Delete(&models.Ads{}, id).Error
}
