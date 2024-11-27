package sheypoor

import (
	"Crawlzilla/database"
	"Crawlzilla/database/repositories"
	"Crawlzilla/models"
	"Crawlzilla/utils"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
)

// CategoryHandler is a function type that defines the signature of handlers for each category
type CategoryHandler func(context.Context, AdURL) error

// StartConsumer starts a consumer for a specific category with a shared Chrome context
func StartConsumer(category string, adChannel <-chan AdURL) {

	maxCrawlTime, err := strconv.Atoi(os.Getenv("MAX_CRAWL_TIME"))
	if err != nil {
		log.Fatalf("Error reading MAX_CRAWL_TIME from .env: %v", err)
	}
	maxCrawlDuration := time.Duration(maxCrawlTime) * time.Minute
	ctx := CreateChromeContext(maxCrawlDuration)
	defer ctx.Cancel()

	// Define handlers for each category
	handlers := map[string]CategoryHandler{
		"villa-for-sale":             handleVillaForSale,
		"house-apartment-for-rent":   handleHouseApartmentForRent,
		"houses-apartments-for-sale": handleHouseApartmentForSale,
	}

	// Get the handler for the current category
	handler, exists := handlers[category]
	if !exists {
		log.Fatalf("No handler found for category %s", category)
	}

	// Process URLs from the adChannel
	for ad := range adChannel {
		err := chromedp.Run(ctx.Ctx, chromedp.Navigate(ad.URL))
		if err != nil {
			log.Printf("Failed to navigate to URL %s: %v", ad.URL, err)
			continue
		}

		// Call the specific handler for the category
		if err := handler(ctx.Ctx, ad); err != nil {
			log.Printf("Error handling ad for category %s at URL %s: %v", category, ad.URL, err)
		}
	}
}

func handleVillaForSale(ctx context.Context, ad AdURL) error {
	// fmt.Println("-------------------------------------فروش ویلا-------------------------------------")
	// Extract title
	title, err := utils.ExtractTitle(ctx)
	if err != nil {
		return err
	}

	// Extract attributes
	attributes, err := utils.ExtractVillaForSale(ctx)
	if err != nil {
		return err
	}

	// Extract image URLs
	imageURL, err := utils.ExtractImageURL(ctx)
	if err != nil {
		return err
	}

	// Extract city and district
	city, district, err := utils.ExtractCityAndDistrict(ctx)
	if err != nil {
		return err
	}

	// Extract description
	description, err := utils.ExtractDescription(ctx)
	if err != nil {
		log.Printf("error extracting description: %v", err)
	}
	price, err := utils.ExtractPrice(ctx)
	if err != nil {
		log.Printf("error extracting price: %v", err)
	}

	// Construct CrawlResult
	crawlResult := models.Ads{
		Reference:    "sheypoor",
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		URL:          ad.URL,
		PropertyType: attributes.PropertyType,
		Area:         attributes.Area,
		Room:         attributes.Room,
		Price:        price,
		City:         city,
		Neighborhood: district,
		// BuildingAgeType:  attributes.BuildingAgeType,
		// BuildingAgeValue: attributes.BuildingAgeValue,
		HasElevator: attributes.HasElevator,
		HasParking:  attributes.HasParking,
		HasStorage:  attributes.HasStorage,
	}

	// log.Println(crawlResult.String())
	_, createErr := repositories.CreateAd(database.DB, &crawlResult)
	if createErr != nil {
		return err
	}
	return nil
}
func handleHouseApartmentForSale(ctx context.Context, ad AdURL) error {
	// fmt.Println("-------------------------------------فروش خانه و آپارتمان-------------------------------------")
	// Extract title
	title, err := utils.ExtractTitle(ctx)
	if err != nil {
		return err
	}

	// Extract attributes
	attributes, err := utils.ExtractVillaForSale(ctx)
	if err != nil {
		return err
	}

	// Extract image URLs
	imageURL, err := utils.ExtractImageURL(ctx)
	if err != nil {
		return err
	}

	// Extract city and district
	city, district, err := utils.ExtractCityAndDistrict(ctx)
	if err != nil {
		return err
	}

	// Extract description
	description, err := utils.ExtractDescription(ctx)
	if err != nil {
		log.Printf("error extracting description: %v", err)
	}
	price, err := utils.ExtractPrice(ctx)
	if err != nil {
		log.Printf("error extracting price: %v", err)
	}

	crawlResult := models.Ads{
		Reference:    "sheypoor",
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		URL:          ad.URL,
		PropertyType: attributes.PropertyType,
		Area:         attributes.Area,
		Room:         attributes.Room,
		Price:        price,
		City:         city,
		Neighborhood: district,
		// BuildingAgeType:  attributes.BuildingAgeType,
		// BuildingAgeValue: attributes.BuildingAgeValue,
		HasElevator: attributes.HasElevator,
		HasParking:  attributes.HasParking,
		HasStorage:  attributes.HasStorage,
	}

	// log.Println(crawlResult.String())
	_, createErr := repositories.CreateAd(database.DB, &crawlResult)
	if createErr != nil {
		return err
	}
	return nil
}

func handleHouseApartmentForRent(ctx context.Context, ad AdURL) error {
	// fmt.Println("-------------------------------------رهن و اجاره-------------------------------------")
	// Extract title
	title, err := utils.ExtractTitle(ctx)
	if err != nil {
		return err
	}

	// Extract attributes
	attributes, err := utils.ExtractVillaForSale(ctx)
	if err != nil {
		return err
	}

	// Extract image URLs
	imageURL, err := utils.ExtractImageURL(ctx)
	if err != nil {
		return err
	}

	// Extract city and district
	city, district, err := utils.ExtractCityAndDistrict(ctx)
	if err != nil {
		return err
	}

	// Extract description
	description, err := utils.ExtractDescription(ctx)
	if err != nil {
		log.Printf("error extracting description: %v", err)
	}
	crawlResult := models.Ads{
		Reference:    "sheypoor",
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
		URL:          ad.URL,
		PropertyType: attributes.PropertyType,
		Area:         attributes.Area,
		Room:         attributes.Room,
		Price:        attributes.Price,
		Rent:         attributes.Rent,
		City:         city,
		Neighborhood: district,
		// BuildingAgeType:  attributes.BuildingAgeType,
		// BuildingAgeValue: attributes.BuildingAgeValue,
		HasElevator: attributes.HasElevator,
		HasParking:  attributes.HasParking,
		HasStorage:  attributes.HasStorage,
	}
	// log.Println(crawlResult.String())
	_, createErr := repositories.CreateAd(database.DB, &crawlResult)
	if createErr != nil {
		return err
	}
	return nil
}
