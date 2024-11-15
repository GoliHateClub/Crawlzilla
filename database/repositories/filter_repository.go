package repositories

import (
	"Crawlzilla/models"

	"gorm.io/gorm"
)

type FilterRepository interface {
	CreateOrUpdateFilter(db *gorm.DB, filter *models.Filters) error
	GetFiltersByUserID(db *gorm.DB, userID string, pageIndex, pageSize int) ([]models.Filters, int64, error)
	GetFiltersForAllUsers(db *gorm.DB, offset, limit int) ([]models.Filters, int64, error)
}

type filterRepository struct{}

func NewFilterRepository() FilterRepository {
	return &filterRepository{}
}

func (r *filterRepository) CreateOrUpdateFilter(db *gorm.DB, filter *models.Filters) error {
	return db.Save(filter).Error
}

// GetFiltersByUserID fetches filters for a user with pagination
func (r *filterRepository) GetFiltersByUserID(db *gorm.DB, userID string, pageIndex, pageSize int) ([]models.Filters, int64, error) {
	var filters []models.Filters
	var totalRecords int64

	// Count total filters for the user
	if err := db.Model(&models.Filters{}).Where("user_id = ?", userID).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Fetch filters with pagination
	offset := (pageIndex - 1) * pageSize
	if err := db.Where("user_id = ?", userID).Limit(pageSize).Offset(offset).Find(&filters).Error; err != nil {
		return nil, 0, err
	}

	return filters, totalRecords, nil
}

// GetFiltersForAllUsers fetches filters for all users with pagination
func (r *filterRepository) GetFiltersForAllUsers(db *gorm.DB, offset, limit int) ([]models.Filters, int64, error) {
	var filters []models.Filters
	var totalRecords int64

	// Count total filters
	if err := db.Model(&models.Filters{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Fetch filters with pagination
	if err := db.Limit(limit).Offset(offset).Find(&filters).Error; err != nil {
		return nil, 0, err
	}

	return filters, totalRecords, nil
}
