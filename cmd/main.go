package main

import (
	"fmt"
	"github.com/KaranJagtiani/go-logstash"
	"log"

	"Crawlzilla/config"
	"Crawlzilla/database"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load logger
	logger := logstash_logger.Init("localhost", 5228, "udp", 5)

	// Initialize the database
	db, err := database.SetupDB()
	if err != nil {
		logger.Error(map[string]interface{}{
			"message": fmt.Sprintf("Database setup error: %v", err),
			"error":   true,
		})

		return
	}

	_ = db //TODO: delete this line if use db
}
