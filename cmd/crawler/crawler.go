package crawler

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/services/crawler/divar"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

// WorkerPool size
const numWorkers = 1 //TODO: from .env

var adCounter int

// Worker function that processes jobs
func worker(ctx context.Context, jobs <-chan divar.Job, maxAdCount int, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Jobs channel closed, exit worker
				return
			}

			fmt.Println("Scraping Ad Number:", adCounter)
			data, err := divar.ScrapPropertyPage("https://divar.ir" + job.URL)
			if err != nil {
				log.Println("page passed! property type is not house or vila\n")
				continue
			}
			// Save the scrape data to the database
			if err := repositories.AddCrawlResult(&data); err != nil {
				log.Printf("Failed to add scrape result: %v", err)
			} else {
				fmt.Println("Added to DB successfully!\n")
			}

			// Increment the counter and check if we reached maxAdCount
			adCounter++
			if adCounter >= maxAdCount {
				cancel() // Trigger context cancellation
				return
			}

		case <-ctx.Done():
			// Context canceled, exit worker
			fmt.Println("Worker received shutdown signal, stopping...")
			return
		}
	}
}

func StartDivarCrawler(ctx context.Context) {
	jobs := make(chan divar.Job)
	var wg sync.WaitGroup
	defer wg.Wait()

	// Get maxAdCount from environment
	maxAdCount, err := strconv.Atoi(os.Getenv("MAX_AD_COUNT"))
	if err != nil {
		log.Printf("Error reading MAX_AD_COUNT from .env: %v", err)
	}

	// Create a cancellable context for controlled shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, maxAdCount, &wg, cancel)
	}

	// Start a goroutine to fetch URLs and send them to the jobs channel
	go func() {
		defer close(jobs)
		divar.CrawlDivarAds(ctx, "https://divar.ir/s/iran/real-estate", jobs)
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Received shutdown signal, closing down...")
	return
}
