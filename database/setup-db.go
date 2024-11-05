package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDB initializes and returns a database connection
func SetupDB() (*gorm.DB, error) {
	// Build the database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Verify the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection is not alive: %v", err)
	}

	log.Println("Connected to the database successfully!")
	return db, nil
}
