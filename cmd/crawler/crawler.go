package crawler

import (
	"Crawlzilla/database"
	"Crawlzilla/database/repositories"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/notification"
	"Crawlzilla/services/crawler/divar"
	"Crawlzilla/utils"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"
)

// WorkerPool size
const numWorkers = 1 //TODO: from .env

type CrawlerState struct {
	SuccessAdCount int
	FailAdCount    int
	mu             sync.Mutex // To avoid race conditions
}

func worker(ctx context.Context, jobs <-chan divar.Job, maxAdCount int, state *CrawlerState, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	crawlerLogger, _ := configLogger("crawler")

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Jobs channel closed, exit worker
				return
			}

			crawlerLogger.Info("Scraping Ad Number", zap.Int("successAdCounter", state.SuccessAdCount+1))
			data, err := divar.ScrapPropertyPage("https://divar.ir" + job.URL)
			if err != nil {
				crawlerLogger.Warn("page passed! property type is not house or vila", zap.String("passedURL", "https://divar.ir"+job.URL))
				state.mu.Lock()
				state.FailAdCount++
				state.mu.Unlock()
				continue
			}

			// Save the scrape data to the database
			if id, err := repositories.CreateAd(database.DB, &data); err != nil {
				crawlerLogger.Warn("Data Exists:", zap.String("Existed URL", "https://divar.ir"+data.URL))
			} else {
				crawlerLogger.Info("added to db successfully", zap.String("Ad ID", id))
			}

			// Increment the counter and check if we reached maxAdCount
			state.mu.Lock()
			state.SuccessAdCount++
			if state.SuccessAdCount >= maxAdCount {
				state.mu.Unlock()
				cancel() // Trigger context cancellation
				return
			}
			state.mu.Unlock()

		case <-ctx.Done():
			// Context canceled, exit worker
			fmt.Println("Worker received shutdown signal, stopping...")
			crawlerLogger.Info("worker received shutdown signal, stopping...")
			return
		}
	}
}

func StartDivarCrawler(ctx context.Context, state *CrawlerState) {
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	crawlerLogger, _ := configLogger("crawler")

	crawlerLogger.Info("crawler started successfully")

	jobs := make(chan divar.Job)
	var wg sync.WaitGroup
	defer wg.Wait()

	// Get maxAdCount from environment
	maxAdCount, err := strconv.Atoi(os.Getenv("MAX_AD_COUNT"))
	if err != nil {
		log.Printf("Error reading MAX_AD_COUNT from .env: %v", err)
		crawlerLogger.Error("Error reading MAX_AD_COUNT from .env:", zap.Error(err))
	}

	// Create a cancellable context for controlled shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, maxAdCount, state, &wg, cancel)
	}

	// Start a goroutine to fetch URLs and send them to the jobs channel
	go func() {
		defer close(jobs)
		divar.CrawlDivarAds(ctx, "https://divar.ir/s/iran/real-estate", jobs)
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Received shutdown signal, closing down...")
	crawlerLogger.Info("Received shutdown signal, closing down...")
}

func RunCrawler(ctx context.Context) {
	// Create shared state for success and fail counts
	state := &CrawlerState{}

	metrics := utils.MeasureExecutionStats(func() { StartDivarCrawler(ctx, state) })

	// Access shared state after crawler finishes
	state.mu.Lock()
	successAdCount := state.SuccessAdCount
	failAdCount := state.FailAdCount
	state.mu.Unlock()

	metrics = metrics + fmt.Sprintf("Success Crawled Ad Count: %v\nFailed Crawled Ad Count: %v\n", successAdCount, failAdCount)
	notification.NotifySuperAdmin(ctx, metrics)
}
