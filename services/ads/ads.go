package ads

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"

	"gorm.io/gorm"
)

// AdData represents the paginated ad data returned by the service
type AdData struct {
	Data  []models.AdSummary `json:"data"`
	Pages int                `json:"pages"`
	Page  int                `json:"page"`
}

// GetAllAds retrieves paginated ads with only specific fields
func GetAllAds(db *gorm.DB, page int, pageSize int) (AdData, error) {
	// Fetch ads with pagination
	ads, totalRecords, err := repositories.GetAllAds(db, page, pageSize)
	if err != nil {
		return AdData{}, err
	}

	// Calculate total pages based on total records and page size
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	// Prepare and return paginated response
	return AdData{
		Data:  ads,
		Pages: totalPages,
		Page:  page,
	}, nil
}

func GetAdById(db *gorm.DB, id string) (models.Ads, error) {
	return repositories.GetAdByID(db, id)
}
