package repositories

import (
	"Crawlzilla/database"
	"Crawlzilla/models/ads"
	"errors"
	"fmt"
)

// AddCrawlResult adds a new scrap result to the database
func AddCrawlResult(result *ads.CrawlResult) error {
	return database.DB.Create(result).Error
}

// GetAllCrawlResults retrieves all scrap results
func GetAllCrawlResults() ([]ads.CrawlResult, error) {
	var results []ads.CrawlResult
	err := database.DB.Find(&results).Error
	return results, err
}

// GetCrawlResultByID retrieves a scrap result by ID
func GetCrawlResultByID(id uint) (ads.CrawlResult, error) {
	var result ads.CrawlResult
	err := database.DB.First(&result, id).Error
	return result, err
}

// UpdateCrawlResultById updates specific fields of an existing crawl result in the database
func UpdateCrawlResultById(id uint, updatedData *ads.CrawlResult) error {
	if updatedData == nil {
		return errors.New("updated data cannot be nil")
	}

	if err := database.DB.Model(&ads.CrawlResult{}).Where("id = ?", id).Updates(updatedData).Error; err != nil {
		return fmt.Errorf("failed to update ad: %v", err)
	}

	return nil
}

// DeleteCrawlResult deletes a scrap result by ID
func DeleteCrawlResult(id uint) error {
	return database.DB.Delete(&ads.CrawlResult{}, id).Error
}
