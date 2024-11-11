package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// Create a context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start Crawler
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Crawler...")
		crawler.StartCrawler(ctx)
		fmt.Println("Crawler stopped.")
	}()

	// Start Bot
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Bot...")
		bot.StartBot(ctx)
		fmt.Println("Bot stopped.")
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Server received shutdown signal, waiting for components to stop...")
	return
}
