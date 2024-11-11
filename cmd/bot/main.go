package main

import (
	"Crawlzilla/config"
	"Crawlzilla/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	logConfig := logger.CreateLogger("crawls", "bot", "database")
	crawlsLogger, _ := logConfig("crawls")

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		crawlsLogger.Error("messages", zap.Error(err))

		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load logger

	// Initialize the database
	/*db, err := database.SetupDB()
	if err != nil {

		return
	}

	_ = db //TODO: delete this line if use db*/
}
