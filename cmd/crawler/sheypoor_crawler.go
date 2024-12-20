package crawler

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"Crawlzilla/services/crawler/sheypoor"
)

func StartSheypoorWorker(wg *sync.WaitGroup) {
	urlChannel := make(chan sheypoor.AdURL, 5)
	categories := map[string]string{
		"house-apartment-for-rent":   "house-apartment-for-rent",
		"houses-apartments-for-sale": "houses-apartments-for-sale",
		"villa-for-sale":             "villa-for-sale",
	}
	batchSize := 1

	// Channel to manage graceful shutdown
	stopChannel := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Listen for shutdown signal in a separate goroutine
	go func() {
		<-sigs
		log.Println("Shutdown signal received. Cleaning up...")
		close(stopChannel) // Notify all goroutines to stop
	}()

	// Start category scrapers
	for name, ctg := range categories {

		go func(categoryName, ctg string) {
			maxCrawlTime, err := strconv.Atoi(os.Getenv("MAX_CRAWL_TIME"))
			if err != nil {
				log.Fatalf("Error reading MAX_CRAWL_TIME from .env: %v", err)
			}
			maxCrawlDuration := time.Duration(maxCrawlTime) * time.Minute
			ctx := sheypoor.CreateChromeContext(maxCrawlDuration)
			defer ctx.Cancel()
			sheypoor.ScrapeCategory(categoryName, ctx.Ctx, ctg, urlChannel, batchSize)
		}(name, ctg)
	}

	// Start consumers for each category
	for category := range categories {
		categoryChannel := make(chan sheypoor.AdURL, 5)
		go func(category string) {
			sheypoor.StartConsumer(category, categoryChannel)
		}(category)

		// Forward URLs to specific category channels
		go func(category string) {
			for ad := range urlChannel {
				if ad.Category == category {
					categoryChannel <- ad
				}
			}
			close(categoryChannel) // Close when forwarding is done
		}(category)
	}

	// Wait for shutdown signal
	<-stopChannel
	wg.Done()
	// Close the main URL channel to stop forwarding
	close(urlChannel)
	log.Println("Main URL channel closed.")

	// Wait until all goroutines are done
	log.Println("Program terminated gracefully.")
}
