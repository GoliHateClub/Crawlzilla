package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgBot "Crawlzilla/services/bot"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"Crawlzilla/cmd/bot"
	"Crawlzilla/cmd/crawler"
	"Crawlzilla/config"
	"Crawlzilla/database"
	"Crawlzilla/logger"
	"Crawlzilla/utils"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Config logger
	configLogger := logger.ConfigLogger()
	dbLogger, _ := configLogger("database")

	/// Initialize the database
	_, err := database.SetupDB()
	if err != nil {
		dbLogger.Error("Database setup error", zap.Error(err))
		return
	}

	// Create a context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	ctx = context.WithValue(ctx, "configLogger", configLogger)

	defer stop()

	// Initialize bot and attach it to context
	botInstance := tgBot.Init()
	ctx = context.WithValue(ctx, "bot", botInstance)

	var wg sync.WaitGroup
	defer wg.Wait()

	c := cron.New()
	c.AddFunc("@daily", func() {
		log.Println("Starting Crawler...")
		utils.MeasureExecutionStats(func() { crawler.StartDivarCrawler(ctx) })
		log.Println("Crawler stopped.")
	})

	c.Start()
	defer c.Stop()

	// Start Bot
	if config.GetBoolean("IS_BOT_ACTIVE") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Starting Bot...")
			bot.StartBot(ctx)
			fmt.Println("Bot stopped.")
		}()
	}

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Server received shutdown signal, waiting for components to stop...")
	return
}
