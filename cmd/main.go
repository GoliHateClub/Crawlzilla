package main

import (
	"log"

	"github.com/GoliHateClub/Crawlzilla/config"
	"github.com/GoliHateClub/Crawlzilla/database"
)

func main() {
	// Load configuration
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database
	db, err := database.SetupDB()
	if err != nil {
		log.Fatalf("Database setup error: %v", err)
	}

	_ = db //TODO: delete this line if use db
}
