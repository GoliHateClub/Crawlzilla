package divar

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// Job represents a single URL to scrap
type Job struct {
	URL string
}

func CrawlDivarAds(ctx context.Context, url string, jobs chan<- Job) {
	// Create a ChromeDP context with a timeout
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Set timeout for the scraping task
	maxCrawlTime, err := strconv.Atoi(os.Getenv("MAX_CRAWL_TIME"))
	if err != nil {
		log.Printf("Error reading MAX_CRAWL_TIME from .env: %v", err)
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
			log.Println("Error reading MAX_PAGE from .env: ", err)
		}
		var lastHeight, newHeight int64

		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			log.Println("Navigation error:", err)
		}

		for {
			select {
			case <-ctx.Done():
				// Gracefully stop scrolling if context is canceled
				fmt.Println("Crawler received shutdown signal, stopping...")
				return
			default:
				fmt.Println("\nLOADING PAGE:", page)
				fmt.Println()

				if err := chromedp.Run(ctx, chromedp.Evaluate(`document.body.scrollHeight`, &newHeight)); err != nil {
					log.Println("Error getting scroll height:", err)
					continue
				}

				if newHeight == lastHeight {
					fmt.Println("No more content to load.")
					continue
				}
				lastHeight = newHeight

				var buttonExists bool
				if err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('.post-list__load-more-btn-be092') !== null`, &buttonExists)); err != nil {
					log.Println("Error checking 'Load More' button:", err)
					continue
				}

				if buttonExists {
					if err := chromedp.Run(ctx, chromedp.Click(".post-list__load-more-btn-be092", chromedp.ByQuery), chromedp.Sleep(500*time.Millisecond)); err != nil {
						log.Println("Error clicking 'Load More':", err)
						continue
					}
					fmt.Println("\nClicked 'Load More'")
				} else {
					if err := chromedp.Run(ctx, chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil), chromedp.Sleep(500*time.Millisecond)); err != nil {
						log.Println("Error scrolling:", err)
						continue
					}
				}

				var html string
				if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &html)); err != nil {
					log.Println("Error getting HTML content:", err)
					continue
				}
				htmlChan <- html
				page++
				if page >= maxPage {
					fmt.Println("max page reached, stopping...")
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
			fmt.Println("HTML extraction received shutdown signal, stopping...")
			return
		default:
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
			if err != nil {
				log.Println("Error parsing HTML:", err)
				continue
			}

			// Extract links and send jobs to channel
			doc.Find("a.kt-post-card__action").Each(func(i int, s *goquery.Selection) {
				if href, exists := s.Attr("href"); exists {
					select {
					case jobs <- Job{URL: href}:
						fmt.Println("scrap started")
					case <-ctx.Done():
						fmt.Println("Job sending received shutdown signal, stopping...")
						return
					}
				} else {
					fmt.Println("No href found in the link.")
				}
			})
		}
	}
	fmt.Println("Scrolling and extraction completed.")
}
