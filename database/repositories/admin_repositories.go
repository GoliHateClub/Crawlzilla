package repositories

import (
	"Crawlzilla/models"
	"errors"
	"gorm.io/gorm"
)

// CreateAdmin creates an admin if it doesn't exist or make an existing user to be an admin.
func CreateAdmin(db *gorm.DB, telegramID int64) (models.Role, error) {
	var user models.Users

	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAdmin := models.Users{
				Telegram_ID: telegramID,
				Role:        models.RoleAdmin,
			}
			if err := db.Create(&newAdmin).Error; err != nil {
				return "", err
			}
			return newAdmin.Role, nil
		}
		return "", err
	}

	if user.Role != models.RoleAdmin {
		user.Role = models.RoleAdmin
		if err := db.Save(&user).Error; err != nil {
			return "", err
		}
	}

	return user.Role, nil
}

// GetAdminByID retrieves a user with the "admin" role by their Telegram ID
func GetAdminByID(db *gorm.DB, telegramID int64) (*models.Users, error) {
	var user models.Users
	err := db.Where("telegram_id = ? AND role = ?", telegramID, models.RoleAdmin).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// RemoveAdminByID removes a user with the admin role by their Telegram ID.
func RemoveAdminByID(db *gorm.DB, telegramID int64) error {
	var user models.Users
	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	if user.Role != "admin" {
		return errors.New("user is not an admin")
	}
	user.Role = "user"
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetAdminsPaginated retrieves paginated admin data from the database.
func GetAdminsPaginated(db *gorm.DB, page int, pageSize int) ([]models.Users, int64, error) {
	var admins []models.Users
	var totalRecords int64

	if err := db.Model(&models.Users{}).
		Where("role = ?", "admin").
		Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := db.Where("role = ?", "admin").
		Offset(offset).
		Limit(pageSize).
		Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, totalRecords, nil
}