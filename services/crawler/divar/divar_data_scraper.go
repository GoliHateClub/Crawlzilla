package divar

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"Crawlzilla/models/ads"
	"Crawlzilla/utils"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ScrapSellHousePage scraps the given URL, fills the CrawlResult struct, and returns it
func ScrapSellHousePage(pageURL string) (ads.CrawlResult, error) {
	result := ads.CrawlResult{}

	fmt.Println("SCRAPING: ", pageURL)

	DIVAR_TOKEN := os.Getenv("DIVAR_TOKEN")
	if DIVAR_TOKEN == "" {
		log.Println("DIVAR_TOKEN environment variable is not set")
	}

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout for the scraping task
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Define the cookie parameters.
	cookie := &network.CookieParam{
		Name:   "token", // Replace "token" with the cookie name
		Value:  DIVAR_TOKEN,
		Domain: ".divar.ir", // Replace with the domain for your cookie
		Path:   "/",
	}

	// Use an ActionFunc to set the cookie.
	setCookie := chromedp.ActionFunc(func(ctx context.Context) error {
		return network.SetCookie(cookie.Name, cookie.Value).
			WithDomain(cookie.Domain).
			WithPath(cookie.Path).
			Do(ctx)
	})

	// Run the Chromedp tasks
	err := chromedp.Run(ctx,
		network.Enable(),           // Enable the network domain to apply cookies
		setCookie,                  // Set the cookie
		chromedp.Navigate(pageURL)) // Navigate to the page

	if err != nil {
		log.Println("Cant navigate URL:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Title
	err = chromedp.Run(ctx,
		chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section:nth-child(1) div.kt-page-title div h1`, &result.Title),
	)
	if err != nil {
		log.Println("Cant convert or get Total Floor value:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Area
	var stringArea string

	err = chromedp.Run(ctx,
		chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section:nth-child(1) div.post-page__section--padded table:nth-child(1) tbody tr td:nth-child(1)`, &stringArea),
	)
	if err != nil {
		log.Println("Cant get Area:", err)
	}
	// Convert extracted Persian Area text to integer
	result.Area, _ = utils.ConvertPersianNumber(stringArea)

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Price
	var stringPrice string

	err = chromedp.Run(ctx,
		chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section:nth-child(1) div.post-page__section--padded div:nth-child(3) div.kt-base-row__end.kt-unexpandable-row__value-box p`, &stringPrice),
	)
	if err != nil {
		log.Println("Cant get Price:", err)
	}
	// remove price text
	stringPrice = strings.Split(stringPrice, " ")[0]
	if stringPrice == "مجانی" {
		stringPrice = "0"
	}
	priceInt, err := utils.ConvertPersianNumber(stringPrice) // Fill the Area field
	if err != nil {
		log.Println("Cant convert or get Price value:", err)
	}
	result.Price = priceInt // Fill the Price field

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Room
	var stringRoom string

	err = chromedp.Run(ctx,
		chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > table:nth-child(1) > tbody > tr > td:nth-child(3)`, &stringRoom),
	)
	if err != nil {
		log.Println("Cant get Room:", err)
	}
	// Convert extracted Persian Room text to integer
	result.Room, err = utils.ConvertPersianNumber(stringRoom)
	if err != nil {
		log.Println("Cant convert room string to int:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Floor and Total Floors
	// Extract Floor and Total Floors
	var floorExists bool
	var stringFloors string

	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > div:nth-child(7) > div.kt-base-row__end.kt-unexpandable-row__value-box > p") !== null`, &floorExists),
	)
	if err != nil {
		log.Println("Cant get Contact Number element:", err)
	}
	if floorExists {
		err = chromedp.Run(ctx,
			chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section:nth-child(1) div.post-page__section--padded div:nth-child(7) div.kt-base-row__end.kt-unexpandable-row__value-box p`, &stringFloors),
		)
		if err != nil {
			log.Println("Cant Extract Floor and Total Floors:", err)
		}

		// separate floor number and total floors
		floorsSplit := strings.Split(stringFloors, " ")
		if len(floorsSplit) == 1 {
			// Check if the floor value is "همکف"
			if floorsSplit[0] == "همکف" {
				result.FloorNumber = 0
				result.TotalFloors = 0
			} else {
				floorNumberInt, err := utils.ConvertPersianNumber(floorsSplit[0]) // Convert the floor number
				if err != nil {
					log.Println("Cant convert or get Floor Number value:", err)
				}
				result.FloorNumber = floorNumberInt
				result.TotalFloors = 0
			}
		} else {
			// Check if the floor value is "همکف"
			if floorsSplit[0] == "همکف" {
				result.FloorNumber = 0
			} else {
				floorNumberInt, err := utils.ConvertPersianNumber(floorsSplit[0]) // Convert the floor number
				if err != nil {
					log.Println("Cant convert or get Floor Number value:", err)
				}
				result.FloorNumber = floorNumberInt
			}

			// Check if the total floors value is "همکف"
			if floorsSplit[2] == "همکف" {
				result.TotalFloors = 0
			} else {
				totalFloorInt, err := utils.ConvertPersianNumber(floorsSplit[2]) // Convert the total floors
				if err != nil {
					log.Println("Cant convert or get Total Floor value:", err)
				}
				result.TotalFloors = totalFloorInt
			}
		}
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Description
	err = chromedp.Run(ctx,
		chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section.post-page__section--padded div div.kt-base-row.kt-base-row--large.kt-description-row div p`, &result.Description),
	)
	if err != nil {
		log.Println("Cant get Description:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Check Elevator
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			document.querySelector('#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > table:nth-child(10) > tbody > tr > td:nth-child(1)') !== null
		`, &result.HasElevator),
	)
	if err != nil {
		log.Println("Cant get Elevator:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Check Parking
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			document.querySelector('#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > table:nth-child(10) > tbody > tr > td:nth-child(2)') !== null
		`, &result.HasParking),
	)
	if err != nil {
		log.Println("Cant get Parking:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Check Storage
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			document.querySelector('#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > table:nth-child(10) > tbody > tr > td:nth-child(3)') !== null
		`, &result.HasStorage),
	)
	if err != nil {
		log.Println("Cant get Storage:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Contact number
	var contactExists bool
	var stringContactNumber string

	err = chromedp.Run(ctx,
		chromedp.Click(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-actions > button.kt-button.kt-button--primary.post-actions__get-contact`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.EvaluateAsDevTools(`document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.expandable-box > div.copy-row > div > div.kt-base-row__end.kt-unexpandable-row__value-box > p") !== null`, &contactExists),
	)
	if err != nil {
		log.Println("Cant get Contact Number element:", err)
	}

	if contactExists {
		err = chromedp.Run(ctx,
			chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.expandable-box > div.copy-row > div > div.kt-base-row__end.kt-unexpandable-row__value-box > p`, &stringContactNumber),
		)
		if err != nil {
			log.Println("Cant get Contact Number:", err)
		}
		// Convert extracted Persian Contact number text to integer
		result.ContactNumber, err = utils.ConvertPersianNumber(stringContactNumber)
		if err != nil {
			log.Println("Cant convert Contact string to int:", err)
		}
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Location
	// Check if Location URL element exists before extracting
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var exists bool
			// Check if the selector exists
			err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(
				`document.querySelector("a.map-cm__attribution.map-cm__button") !== null`, &exists))
			if err != nil || !exists {
				// Element not found or error occurred, skip extraction
				return nil
			}
			// Extract href attribute if element exists
			return chromedp.AttributeValue(`a.map-cm__attribution.map-cm__button`, "href", &result.LocationURL, nil).Do(ctx)
		}),
	)
	if err != nil {
		log.Println("Cant get Location:", err)
	}

	// Parse latitude and longitude from the location URL
	if result.LocationURL != "" {
		parsedURL, err := url.Parse(result.LocationURL)
		if err == nil {
			queryParams := parsedURL.Query()
			lat, err1 := strconv.ParseFloat(queryParams.Get("latitude"), 64)
			long, err2 := strconv.ParseFloat(queryParams.Get("longitude"), 64)
			if err1 != nil || err2 != nil {
				log.Println("Cant convert Latitude or Longitude:", err)
			}
			result.Latitude = lat
			result.Longitude = long

		}
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Image
	// Check if Image element exists before extracting
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var exists bool
			// Check if the selector exists
			err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(
				`document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-6.kt-offset-1 > section:nth-child(1) > div > div > div.keen-slider.kt-base-carousel__slides.slides-d6304 > div:nth-child(1) > figure > div > picture > img") !== null`, &exists))
			if err != nil || !exists {
				// Element not found or error occurred, skip extraction
				return nil
			}
			// Extract src attribute if element exists
			return chromedp.AttributeValue(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-6.kt-offset-1 > section:nth-child(1) > div > div > div.keen-slider.kt-base-carousel__slides.slides-d6304 > div:nth-child(1) > figure > div > picture > img`, "src", &result.ImageURL, nil).Do(ctx)
		}),
	)
	if err != nil {
		log.Println("Cant get Image:", err)
	}

	// Add URL and Reference divar
	result.Reference = "divar"
	result.URL = pageURL
	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	return result, nil
}