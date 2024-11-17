package repositories

import (
	"Crawlzilla/models"
	"fmt"

	"gorm.io/gorm"
)

// CountFilteredAds counts the total number of ads matching the filter criteria.
func CountFilteredAds(db *gorm.DB, query *gorm.DB) (int64, error) {
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return 0, fmt.Errorf("failed to count ads: %w", err)
	}
	return totalRecords, nil
}

// GetFilteredAds retrieves filtered ads with pagination applied.
func GetFilteredAds(db *gorm.DB, query *gorm.DB, page, pageSize int) ([]models.Ads, error) {
	offset := (page - 1) * pageSize
	query = query.Limit(pageSize).Offset(offset)

	var ads []models.Ads
	if err := query.Find(&ads).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve ads: %w", err)
	}
	return ads, nil
}
