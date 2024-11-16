package database

import (
	"log"
	"os"
	"strconv"

	"Crawlzilla/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// SetupDB initializes and returns a database connection
func SetupDB() (*gorm.DB, error) {
	// Retrieve the database URL from environment variables
	databaseURL := os.Getenv("DB_URL")
	if databaseURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Run migrations for the Ads, Filters, and Users models
	err = db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Check if Super Admin exists, and create one if not
	var count int64
	err = db.Model(&models.Users{}).Where("role = ?", "super-admin").Count(&count).Error
	if err != nil {
		log.Fatalf("Failed to check for initial user: %v", err)
	}
	super_admin_id_str := os.Getenv("SUPER_ADMIN_ID")
	super_admin_id, err := strconv.Atoi(super_admin_id_str)
	// If no user found with the 'super-admin' role, create an initial super-admin
	if count == 0 && super_admin_id_str != "" {
		initialUser := models.Users{
			ID:          uuid.NewString(),
			Telegram_ID: int64(super_admin_id),
			Role:        "super-admin",
		}
		err = db.Create(&initialUser).Error
		if err != nil {
			log.Fatalf("Failed to create initial admin user: %v", err)
		}
		log.Println("Initial admin user created")
	}

	DB = db
	return db, nil
}
