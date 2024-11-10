package main

import (
	"fmt"
	"log"
	"sync"

	"Crawlzilla/cmd/bot"
	"Crawlzilla/cmd/crawler"
	"Crawlzilla/config"
	"Crawlzilla/database"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database
	database.SetupDB()
	fmt.Println("Database Setup Successfully!")

	var wg sync.WaitGroup

	wg.Add(1)
	go crawler.StartCrawler()

	wg.Add(1)
	go bot.StartBot()

	wg.Wait()
}
