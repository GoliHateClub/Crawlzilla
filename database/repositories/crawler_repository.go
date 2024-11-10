package repositories

import (
	"Crawlzilla/database"
	"Crawlzilla/models/ads"
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

// DeleteCrawlResult deletes a scrap result by ID
func DeleteCrawlResult(id uint) error {
	return database.DB.Delete(&ads.CrawlResult{}, id).Error
}
