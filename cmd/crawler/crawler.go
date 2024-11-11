package crawler

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/services/crawler/divar"
	"context"
	"fmt"
	"log"
	"sync"
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

func StartCrawler(ctx context.Context) {
	jobs := make(chan divar.Job)
	var wg sync.WaitGroup
	defer wg.Wait()

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, &wg)
	}

	// Start a goroutine to fetch URLs and send them to the jobs channel
	go func() {
		divar.CrawlDivarAds(ctx, "https://divar.ir/s/iran/buy-apartment", jobs)
		close(jobs)
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Received shutdown signal, closing down...")
	return
}
