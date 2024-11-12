package sheypoor

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

type AdURL struct {
	URL      string
	Category string
}

func ScrapeCategory(categoryName string, ctx context.Context, ctg string, urlChannel chan<- AdURL, batchSize int) {
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.sheypoor.com/s/iran/"+ctg),
	)
	if err != nil {
		log.Fatalf("Failed to navigate to category %s: %v", ctg, err)
	}
	for {
		var urls []string
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			continue
		}
		err = chromedp.Run(ctx,
			chromedp.WaitVisible("[data-index]", chromedp.ByQuery),
		)
		if err != nil {
			log.Printf("Error waiting for new ads in category %s: %v", ctg, err)
			continue
		}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`
                (() => {
                    let urls = [];
                    let adCards = document.querySelectorAll("div[data-index][data-item-index]");
					if(adCards.length == 0) adCards = document.querySelectorAll("section[data-index][data-item-index]");
                    adCards.forEach(card => {
                        const link = card.querySelectorAll("a");
                        link.forEach(ad => {
                            if (ad && ad.href) {
                                urls.push(ad.href);
                            }
                        });
                    });
                    return urls;
                })()
            `, &urls),
		)
		if err != nil {
			log.Printf("Error extracting URLs in category %s: %v", ctg, err)
			continue
		}
		for {
			availableCapacity := cap(urlChannel) - len(urlChannel)
			if availableCapacity >= batchSize {
				break
			}
			log.Printf("Waiting for channel capacity in category %s", ctg)
			time.Sleep(500 * time.Millisecond)
		}
		for _, link := range urls {
			decodedURL, _ := url.QueryUnescape(link)
			urlChannel <- AdURL{URL: decodedURL, Category: ctg}
		}

		time.Sleep(3 * time.Second)
	}
}

type ChromeContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func CreateChromeContext(timeout time.Duration) *ChromeContext {
	// Create base context for Chrome
	ctx, cancelChrome := chromedp.NewContext(context.Background())

	// Wrap context with timeout

	ctx, cancelTimeout := context.WithTimeout(ctx, timeout)

	return &ChromeContext{
		Ctx: ctx,
		Cancel: func() {
			cancelTimeout() // Cancel the timeout
			cancelChrome()  // Cancel the Chrome context
		},
	}
}
