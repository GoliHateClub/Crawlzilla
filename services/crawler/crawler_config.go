package crawler_config

import (
	"log"
	"os"
	"strconv"
)

func SetCrawlerConfig(crawlTime int, pageScrapTime int, adCount int, maxScroll int) {

	if crawlTime != 0 {
		err := os.Setenv("MAX_CRAWL_TIME", strconv.Itoa(crawlTime))
		if err != nil {
			log.Println("error setting crawlTime: ", err)
		}
	}

	if pageScrapTime != 0 {
		err := os.Setenv("MAX_SCRAP_TIME", strconv.Itoa(pageScrapTime))
		if err != nil {
			log.Println("error setting crawlTime: ", err)
		}
	}

	if adCount != 0 {
		err := os.Setenv("MAX_AD_COUNT", strconv.Itoa(adCount))
		if err != nil {
			log.Println("error setting crawlTime: ", err)
		}
	}

	if maxScroll != 0 {
		err := os.Setenv("MAX_PAGE", strconv.Itoa(maxScroll))
		if err != nil {
			log.Println("error setting crawlTime: ", err)
		}
	}
}
