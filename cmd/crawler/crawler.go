package crawler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"Crawlzilla/database/repositories"
	"Crawlzilla/services/crawler/divar"
)

// WorkerPool size
const numWorkers = 1 //TODO: from .env

var adNumber int //TODO: from .env

// Worker function that processes jobs
func worker(ctx context.Context, jobs <-chan divar.Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Jobs channel closed, exit worker
				return
			}
			fmt.Println("Scraping Ad Number:", adNumber)
			data, err := divar.ScrapSellHousePage("https://divar.ir" + job.URL)
			if err != nil {
				log.Println("Can't add scrap divar page!", job.URL)
				continue
			}
			// Save the scrape data to the database
			if err := repositories.AddCrawlResult(&data); err != nil {
				log.Fatalf("Failed to add scrape result: %v", err)
			} else {
				fmt.Println("Add to DB successfully!\n")
			}
			adNumber++

		case <-ctx.Done():
			// Context canceled, exit worker
			fmt.Println("Worker received shutdown signal, stopping...")
			return
		}
	}
}

func StartCrawler() {
	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a signal channel to catch SIGINT (Ctrl+C) and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create the channel for jobs
	jobs := make(chan divar.Job)
	done := make(chan struct{})

	// WaitGroup to synchronize workers
	var wg sync.WaitGroup

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, &wg)
	}

	// Start a goroutine to fetch URLs and send them to the jobs channel
	go func() {
		divar.CrawlDivarAds(ctx, "https://divar.ir/s/iran/buy-apartment", jobs, done) //TODO: add gategories
	}()

	// Wait until scrolling is done and close jobs channel
	go func() {
		<-done
		close(jobs) // Close jobs channel when done sending
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("Received shutdown signal, closing down...")
	cancel() // Signal to cancel context

	// Wait for all workers to finish processing
	wg.Wait()
	fmt.Println("All workers stopped, exiting.")
}
