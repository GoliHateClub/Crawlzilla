package divar

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"Crawlzilla/models"
	"Crawlzilla/utils"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ScrapPropertyPage scraps the given URL, fills the Ads struct, and returns it
func ScrapPropertyPage(pageURL string) (models.Ads, error) {
	result := models.Ads{}

	fmt.Println("SCRAPING: ", pageURL)

	DIVAR_TOKEN := os.Getenv("DIVAR_TOKEN")
	if DIVAR_TOKEN == "" {
		log.Println("DIVAR_TOKEN environment variable is not set")
	}

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout for the scraping task
	maxScrapTime, err := strconv.Atoi(os.Getenv("MAX_SCRAP_TIME"))
	if err != nil {
		log.Printf("Error reading MAX_CRAWL_TIME from .env: %v", err)
	}
	maxScrapDuration := time.Duration(maxScrapTime) * time.Second
	// Set timeout for the scraping task
	ctx, cancel = context.WithTimeout(ctx, maxScrapDuration)
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
	err = chromedp.Run(ctx,
		network.Enable(),           // Enable the network domain to apply cookies
		setCookie,                  // Set the cookie
		chromedp.Navigate(pageURL)) // Navigate to the page

	if err != nil {
		log.Println("Cant navigate URL:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Category
	var categoryText string
	err = chromedp.Run(ctx,
		chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > div > nav > div > a > button > span`, &categoryText),
	)
	if err != nil {
		log.Println("cant get category string:", err)
		return result, errors.New("")
	}

	category_property := strings.Split(categoryText, " ")
	category := category_property[0]
	property := category_property[1]

	if category == "فروش" {
		result.CategoryType = "sell"
		if property == "خانه" {
			result.PropertyType = "vila"
		} else if property == "آپارتمان" {
			result.PropertyType = "apartment"
		} else {
			return result, errors.New("property type not found")
		}
	} else if category == "اجارهٔ" {
		result.CategoryType = "rent"
		if property == "خانه" {
			result.PropertyType = "vila"
		} else if property == "آپارتمان" {
			result.PropertyType = "apartment"
		} else {
			return result, errors.New("property type not found")
		}
	} else {
		return result, errors.New("category not found")
	}
	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Title
	err = chromedp.Run(ctx,
		chromedp.Text(`#app div.container--has-footer-d86a9.kt-container div main article div div.kt-col-5 section:nth-child(1) div.kt-page-title div h1`, &result.Title),
	)
	if err != nil {
		log.Println("Cant get Title:", err)
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract City
	var stringCity string

	err = chromedp.Run(ctx,
		chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.kt-page-title > div > div`, &stringCity),
	)
	if err != nil {
		log.Println("Cant Get City:", err)
	}

	// Find the index of "در"
	index := strings.Index(stringCity, "در ")
	if index == -1 {
		fmt.Println("The text does not contain 'در'")
		return result, errors.New("")
	}

	// Get the part of the text after "در" and trim any leading or trailing spaces
	locationPart := strings.TrimSpace(stringCity[index+len("در "):])

	// Split the location part by spaces
	locationParts := strings.Split(locationPart, "، ")

	// Check if there is at least one part for the city
	if len(locationParts) > 0 {
		// The first part is the city
		result.City = locationParts[0]
	}

	// If there are more parts, join them as the neighborhood
	if len(locationParts) > 1 {
		result.Neighborhood = strings.Join(locationParts[1:], " ")
	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Extract Room
	var stringRoom string

	err = chromedp.Run(ctx,
		chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > table:nth-child(1) > tbody > tr > td:nth-child(3)`, &stringRoom),
	)
	if err != nil {
		log.Println("Cant get Room:", err)
	}
	if stringRoom == "بدون اتاق" {
		result.Room = 0
	} else {
		// Convert extracted Persian Room text to integer
		result.Room, err = utils.ConvertPersianNumber(stringRoom)
		if err != nil {
			log.Println("Cant convert room string to int:", err)
		}
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
	// Extract Rent, Price, Floors
	var stringPrice string
	var stringRent string
	var stringFloors string

	err = chromedp.Run(ctx,

		// Retrieve all rows within the container
		chromedp.ActionFunc(func(ctx context.Context) error {
			var rows []*cdp.Node
			// Find all rows that match the class selector
			err := chromedp.Nodes(`div.post-page__section--padded div.kt-base-row.kt-base-row--large.kt-unexpandable-row`, &rows, chromedp.ByQueryAll).Do(ctx)
			if err != nil {
				return err
			}

			// Loop through each row and extract title and value
			for _, row := range rows {
				var title, value string

				// Extract the title text from within each row
				if err := chromedp.Text(`.kt-base-row__title.kt-unexpandable-row__title`, &title, chromedp.ByQuery, chromedp.FromNode(row)).Do(ctx); err != nil {
					log.Printf("Failed to extract title for row: %v", err)
					continue
				}
				title = strings.TrimSpace(title)

				// Extract the value text from within each row
				if err := chromedp.Text(`.kt-unexpandable-row__value`, &value, chromedp.ByQuery, chromedp.FromNode(row)).Do(ctx); err != nil {
					log.Printf("Failed to extract value for row: %v", err)
					continue
				}
				value = strings.TrimSpace(value)

				// Assign value based on title
				switch title {
				case "قیمت کل":
					stringPrice = value
					// remove price text
					stringPrice = strings.Split(stringPrice, " ")[0]
					if stringPrice == "مجانی" {
						stringPrice = "0"
					}
					priceInt, err := utils.ConvertPersianNumber(stringPrice)
					if err != nil {
						log.Println("Cant convert or get Price value:", err)
					}
					result.Price = priceInt // Fill the Price field

				case "ودیعه":
					stringPrice = value
					// remove price text
					stringPrice = strings.Split(stringPrice, " ")[0]
					if stringPrice == "مجانی" {
						stringPrice = "0"
					}
					priceInt, err := utils.ConvertPersianNumber(stringPrice)
					if err != nil {
						log.Println("Cant convert or get Price value:", err)
					}
					result.Price = priceInt // Fill the Price field

				case "اجارهٔ ماهانه":
					stringRent = value
					// remove rent text
					stringRent = strings.Split(stringRent, " ")[0]
					if stringRent == "مجانی" {
						stringRent = "0"
					}
					rentInt, err := utils.ConvertPersianNumber(stringRent) // Fill the Area field
					if err != nil {
						log.Println("Cant convert or get Price value:", err)
					}
					result.Rent = rentInt // Fill the Price field

				case "طبقه":
					stringFloors = value
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
			}
			return nil
		}),
	)

	if err != nil {
		log.Println("cant extract floor or prirce or rent", err)

	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Check Price Slider Existed
	var sliderExist bool

	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > div.convert-slider > table > tbody > tr") !== null
		`, &sliderExist),
	)
	if err != nil {
		log.Println("Cant get Parking:", err)
	}

	if sliderExist {
		err = chromedp.Run(ctx,
			chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > div.convert-slider > table > tbody > tr > td:nth-child(1)`, &stringPrice),
		)
		if err != nil {
			log.Println("Cant Extract Slider price and rent:", err)
		}
		err = chromedp.Run(ctx,
			chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-page__section--padded > div.convert-slider > table > tbody > tr > td:nth-child(2)`, &stringRent),
		)
		if err != nil {
			log.Println("Cant Extract Slider price and rent:", err)
		}

		// remove price text
		stringPriceParts := strings.Split(stringPrice, " ")

		var stringPriceValue, currencyPrice string

		if len(stringPriceParts) == 2 {
			// If there are two parts, assign each part separately
			stringPriceValue = stringPriceParts[0]
			currencyPrice = stringPriceParts[1]
		} else if len(stringPriceParts) == 1 {
			// If there is only one part, assign it to stringPriceValue and leave currencyPrice empty
			stringPriceValue = stringPriceParts[0]
			currencyPrice = "" // or some default value if needed
		} else {
			// Handle any unexpected cases, such as an empty string
			stringPriceValue = "0"
			currencyPrice = "0"
		}
		if stringPriceValue == "مجانی" {
			stringPriceValue = "0"
		}
		if stringPriceValue == "رایگان" {
			stringPriceValue = "0"
		}
		if stringPriceValue == "توافقی" {
			stringPriceValue = "0"
		}
		priceInt, err := utils.ConvertPersianNumber(stringPriceValue)
		if err != nil {
			log.Println("Cant convert or get Price value:", err)
		}
		if currencyPrice == "میلیارد" {
			result.Price = priceInt * 1000000000
		} else if currencyPrice == "میلیون" {
			result.Price = priceInt * 1000000
		} else {
			result.Price = priceInt // Fill the Price field
		}

		// remove rent text
		stringRentParts := strings.Split(stringRent, " ")

		var stringRentValue, currencyRent string

		if len(stringRentParts) == 2 {
			// If there are two parts, assign each part separately
			stringRentValue = stringRentParts[0]
			currencyRent = stringRentParts[1]
		} else if len(stringRentParts) == 1 {
			// If there is only one part, assign it to stringRentValue and leave currencyRent empty
			stringRentValue = stringRentParts[0]
			currencyRent = "" // or some default value if needed
		} else {
			// Handle any unexpected cases, such as an empty string
			stringRentValue = "0"
			currencyRent = "0"
		}
		if stringRentValue == "مجانی" {
			stringRentValue = "0"
		}
		if stringRentValue == "رایگان" {
			stringRentValue = "0"
		}
		if stringRentValue == "توافقی" {
			stringRentValue = "0"
		}
		rentInt, err := utils.ConvertPersianNumber(stringRentValue) // Fill the Area field
		if err != nil {
			log.Println("Cant convert or get Rent value:", err)
		}
		if currencyRent == "میلیارد" {
			result.Rent = rentInt * 1000000000
		} else if currencyRent == "میلیون" {
			result.Rent = rentInt * 1000000
		} else {
			result.Rent = rentInt // Fill the Price field
		}

	}

	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// Check Elevator, Balcony, Storage, Parking

	// Run Chromedp tasks
	err = chromedp.Run(ctx,

		// Retrieve the row and check its child td elements
		chromedp.ActionFunc(func(ctx context.Context) error {
			var rowNodes []*cdp.Node
			// Find the main row element
			err := chromedp.Nodes(`tr.kt-group-row__data-row td.kt-group-row-item__value`, &rowNodes, chromedp.ByQueryAll).Do(ctx)
			if err != nil {
				return err
			}

			// Loop through each <td> element and check for keywords
			for _, node := range rowNodes {
				var featureText string
				// Extract the text of the current <td> using its NodeID
				if err := chromedp.Text(node.FullXPath(), &featureText, chromedp.BySearch).Do(ctx); err != nil {
					log.Printf("Failed to extract text for td element: %v", err)
					continue
				}
				featureText = strings.TrimSpace(featureText)

				// Switch case to check if the feature exists and set the corresponding variable
				switch featureText {
				case "پارکینگ":
					result.HasParking = true
				case "انباری":
					result.HasStorage = true
				case "بالکن":
					result.HasBalcony = true
				case "آسانسور":
					result.HasElevator = true
				}
			}
			return nil
		}),
	)

	if err != nil {
		log.Printf("cant get elevator balcony storage parking: %v", err)
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
	// Contact number
	var contactExists bool
	var stringContactNumber string

	err = chromedp.Run(ctx,
		chromedp.Click(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.post-actions > button.kt-button.kt-button--primary.post-actions__get-contact`, chromedp.NodeVisible),
		chromedp.Sleep(2*time.Second),
		chromedp.EvaluateAsDevTools(`document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.expandable-box > div.copy-row > div > div.kt-base-row__end.kt-unexpandable-row__value-box > a") !== null`, &contactExists),
	)
	if err != nil {
		log.Println("Cant get Contact Number element:", err)
	}

	if contactExists {
		err = chromedp.Run(ctx,
			chromedp.Text(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-5 > section:nth-child(1) > div.expandable-box > div.copy-row > div > div.kt-base-row__end.kt-unexpandable-row__value-box > a`, &stringContactNumber),
		)
		if err != nil {
			log.Println("Cant get Contact Number:", err)
		}
		// Convert extracted Persian Contact number text to integer
		result.ContactNumber = stringContactNumber

	} else {
		log.Println("phone number is not exist")
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
				return errors.New("location not found")
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
				`document.querySelector("#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-6.kt-offset-1 > section:nth-child(1) > div > div > div.keen-slider.kt-base-carousel__slides.slides-d6304 > div:nth-child(2) > figure > div > picture > img") !== null`, &exists))
			if err != nil || !exists {
				// Element not found or error occurred, skip extraction
				return errors.New("image not found")
			}
			// Extract src attribute if element exists
			return chromedp.AttributeValue(`#app > div.container--has-footer-d86a9.kt-container > div > main > article > div > div.kt-col-6.kt-offset-1 > section:nth-child(1) > div > div > div.keen-slider.kt-base-carousel__slides.slides-d6304 > div:nth-child(2) > figure > div > picture > img`, "src", &result.ImageURL, nil).Do(ctx)
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
