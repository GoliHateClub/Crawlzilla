package divar

import (
	"Crawlzilla/services/crawler/divar"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestScrapPropertyPage(t *testing.T) {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("No .env file found: %v", err)
	}

	// Define a mock page URL (replace with a valid test URL if possible)
	mockURL := "https://divar.ir/v/%DB%B8%DB%B0%D9%85%D8%AA%D8%B1-%D8%A2%D9%81%D8%AA%D8%A7%D8%A8%DA%AF%DB%8C%D8%B1-%D9%81%D9%88%D9%84-%D8%A7%D9%85%DA%A9%D8%A7%D9%86%D8%A7%D8%AA-%D8%B4%D9%87%D8%B1%D8%B2%DB%8C%D8%A8%D8%A7/wZ10kKqk"

	// Call the ScrapPropertyPage function
	result, err := divar.ScrapPropertyPage(mockURL)

	// Validate the result and error handling
	assert.NoError(t, err, "Expected no error during scraping")
	assert.NotEmpty(t, result.Title, "Expected non-empty Title")
	assert.NotEmpty(t, result.City, "Expected non-empty City")
	assert.NotEmpty(t, result.CategoryType, "Expected non-empty CategoryType")
	assert.NotEmpty(t, result.PropertyType, "Expected non-empty PropertyType")

	// Further assertions on `result` fields if needed
	assert.GreaterOrEqual(t, result.Price, 0, "Price should be greater than or equal to 0")
	assert.GreaterOrEqual(t, result.Area, 0, "Area should be greater than or equal to 0")
}
