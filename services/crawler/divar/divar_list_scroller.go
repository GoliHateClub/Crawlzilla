package divar

import (
	"Crawlzilla/logger"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

// Job represents a single URL to scrap
type Job struct {
	URL string
}

func CrawlDivarAds(ctx context.Context, url string, jobs chan<- Job) {

	configLogger := logger.ConfigLogger()
	crawlerLogger, _ := configLogger("crawler")

	// Create a ChromeDP context with a timeout
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Set timeout for the scraping task
	maxCrawlTime, err := strconv.Atoi(os.Getenv("MAX_CRAWL_TIME"))
	if err != nil {
		crawlerLogger.Error("Error reading MAX_CRAWL_TIME from .env:", zap.Error(err))
	}
	maxCrawlDuration := time.Duration(maxCrawlTime) * time.Minute

	ctx, cancel = context.WithTimeout(ctx, maxCrawlDuration)
	defer cancel()

	htmlChan := make(chan string)

	// Goroutine to scroll and load content
	go func() {
		defer close(htmlChan)

		page := 0
		maxPage, err := strconv.Atoi(os.Getenv("MAX_PAGE"))
		if err != nil {
			crawlerLogger.Error("Error reading MAX_PAGE from .env", zap.Error(err))
		}
		var lastHeight, newHeight int64

		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			crawlerLogger.Error("Navigation error:", zap.Error(err))
		}
		chromedp.Run(ctx, chromedp.Sleep(2*time.Second))
		for {
			select {
			case <-ctx.Done():
				// Gracefully stop scrolling if context is canceled
				crawlerLogger.Info("Crawler received shutdown signal, stopping...")
				return
			default:
				crawlerLogger.Info("LOADING PAGE:", zap.Int("page", page))
				fmt.Println()

				if err := chromedp.Run(ctx, chromedp.Evaluate(`document.body.scrollHeight`, &newHeight)); err != nil {
					crawlerLogger.Error("Error getting scroll height:", zap.Error(err))
					continue
				}

				if newHeight == lastHeight {
					crawlerLogger.Info("No more content to load.")
					continue
				}
				lastHeight = newHeight

				var buttonExists bool
				if err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('.post-list__load-more-btn-be092') !== null`, &buttonExists)); err != nil {
					crawlerLogger.Error("Error checking 'Load More' button:", zap.Error(err))
					continue
				}

				if buttonExists {
					if err := chromedp.Run(ctx, chromedp.Click(".post-list__load-more-btn-be092", chromedp.ByQuery), chromedp.Sleep(500*time.Millisecond)); err != nil {
						crawlerLogger.Error("Error clicking 'Load More':", zap.Error(err))
						continue
					}
					fmt.Println("\nClicked 'Load More'")
				} else {
					if err := chromedp.Run(ctx, chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil), chromedp.Sleep(500*time.Millisecond)); err != nil {
						crawlerLogger.Error("Error scrolling:", zap.Error(err))
						continue
					}
				}

				var html string
				if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &html)); err != nil {
					crawlerLogger.Error("Error getting HTML content:", zap.Error(err))
					continue
				}
				htmlChan <- html
				page++
				if page >= maxPage {
					crawlerLogger.Info("max page reached, stopping...")
					cancel() // Trigger context cancellation
					return
				}
			}
		}
	}()

	// Process HTML pages as they come in
	for html := range htmlChan {
		select {
		case <-ctx.Done():
			crawlerLogger.Info("HTML extraction received shutdown signal, stopping...")
			return
		default:
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
			if err != nil {
				crawlerLogger.Error("Error parsing HTML:", zap.Error(err))
				continue
			}

			// Extract links and send jobs to channel
			doc.Find("a.kt-post-card__action").Each(func(i int, s *goquery.Selection) {
				if href, exists := s.Attr("href"); exists {
					select {
					case jobs <- Job{URL: href}:
						crawlerLogger.Info("scrap started", zap.String("url", href))
					case <-ctx.Done():
						crawlerLogger.Info("Job sending received shutdown signal, stopping...")
						return
					}
				} else {
					crawlerLogger.Info("No href found in the link.")
				}
			})
		}
	}
	crawlerLogger.Info("Scrolling and extraction completed.")
}
