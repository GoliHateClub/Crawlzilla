package database

import (
	"log"
	"os"

	"Crawlzilla/models/ads"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// SetupDB initializes and returns a database connection
func SetupDB() error {
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

	// Run migrations for the CrawlResult model
	err = db.AutoMigrate(&ads.CrawlResult{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	return nil
}
