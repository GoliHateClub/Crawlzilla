package utils

import (
	"Crawlzilla/models"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// AttributeData struct to store ad attribute information
type AttributeData struct {
	Title string
	Value string
}

// Helper function to convert Go slice of titles to JavaScript array syntax
func TitlesToJSArray(titles []string) string {
	jsArray := "["
	for i, title := range titles {
		jsArray += `"` + title + `"`
		if i < len(titles)-1 {
			jsArray += ", "
		}
	}
	jsArray += "]"
	return jsArray
}

func MapAttributesToCrawlResult(attributes []AttributeData) models.Ads {
	var result models.Ads
	for _, attr := range attributes {
		switch attr.Title {
		case "متراژ":
			result.Area = parseInt(attr.Value)
		case "تعداد اتاق":
			result.Room = parseInt(attr.Value)
		case "پارکینگ":
			result.HasParking = parseBool(attr.Value)
		case "انباری":
			result.HasStorage = parseBool(attr.Value)
		case "بالکن":
			// Assuming a similar bool mapping for balcony (not included in CrawlResult, so consider adding if needed)
		// case "سن بنا":
		// 	if age, err := extractAge(attr.Value); err == nil {
		// 		result.BuildingAgeValue = age
		// 	} else {
		// 		result.BuildingAgeType = attr.Value
		// 	}
		case "رهن":
			result.Price = parseInt(attr.Value)
		case "اجاره":
			result.Rent = parseInt(attr.Value)
		case "نوع ملک":
			if attr.Value == "آپارتمان" {
				result.PropertyType = "house"
			} else {

				result.PropertyType = "vila"
			}
		case "آسانسور":
			result.HasElevator = parseBool(attr.Value)
		}
	}

	return result
}

// Updated parseInt using ConvertPersianNumber
func parseInt(value string) int {
	// Remove "تومان" if present
	cleanedValue := strings.ReplaceAll(value, "تومان", "")
	cleanedValue = strings.TrimSpace(cleanedValue)

	// Convert Persian number to an integer
	parsedValue, err := ConvertPersianNumber(cleanedValue)
	if err != nil {
		return 0
	}
	return parsedValue
}

// Helper function to parse boolean values based on Persian terms
func parseBool(value string) bool {
	// "دارد" means true, "ندارد" means false
	return strings.TrimSpace(value) == "دارد"
}

// func extractAge(persianNum string) (int, error) {
// 	// Define a map for Persian to English digit conversion
// 	persianToEnglish := map[rune]rune{
// 		'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
// 		'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
// 	}

// 	// Convert Persian digits to English digits, ignoring any non-numeric characters
// 	var englishNum strings.Builder
// 	for _, r := range persianNum {
// 		// Skip non-Persian and non-numeric characters
// 		if englishDigit, exists := persianToEnglish[r]; exists {
// 			englishNum.WriteRune(englishDigit)
// 		}
// 	}

// 	// Convert the result to an integer
// 	return strconv.Atoi(englishNum.String())
// }

// ExtractVillaForSale extracts specific attributes for the "villa-for-sale" category
func ExtractVillaForSale(ctx context.Context) (models.Ads, error) {
	attrs := []string{"متراژ", "تعداد اتاق", "پارکینگ", "انباری", "بالکن", "سن بنا", "رهن", "اجاره", "نوع ملک", "آسانسور"}
	var attributes []AttributeData

	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			((titles) => {
				let results = [];
				let elements = document.querySelectorAll('div');
				for (let elem of elements) {
					let siblingElems = Array.from(elem.querySelectorAll('p'));
					if (siblingElems.length === 2) {
						let [titleElem, valueElem] = siblingElems;
						if (titles.includes(titleElem.innerText.trim())) {
							results.push({
								title: titleElem.innerText.trim(),
								value: valueElem.innerText.trim()
							});
						}
					}
				}
				return results;
			})(`+TitlesToJSArray(attrs)+`)`, &attributes),
	)
	return MapAttributesToCrawlResult(attributes), err
}
func ExtractPrice(ctx context.Context) (int, error) {
	var priceText string
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			(() => {
				// Find the svg element with the specific path
				const svg = document.querySelector('svg path[d="M5.422 18.114c-.528 0-.977-.11-1.347-.33a2.127 2.127 0 0 1-.814-.92A3.185 3.185 0 0 1 3 15.545v-4.893h1.727v4.893c0 .198.01.337.033.418.029.073.087.12.173.142.095.022.257.033.49.033h.542l.065 1.02-.065.955h-.543Z"]');
				
				if (!svg) return ""; // If SVG not found, return empty string
				
				// Traverse up the DOM tree until a strong tag is found
				let parent = svg.closest("span"); // start from closest span parent
				while (parent) {
					const strongTag = parent.querySelector("strong");
					if (strongTag && strongTag.innerText) {
						return strongTag.innerText.trim(); // Return the price text
					}
					parent = parent.parentElement; // Go to the next parent element
				}
				return ""; // Return empty if no strong tag with price is found
			})()
		`, &priceText),
	)
	if err != nil {
		return 0, err
	}
	// Convert the extracted Persian price text to an integer
	price, err := ConvertPersianNumber(priceText)
	if err != nil {
		return 0, err
	}
	return price, nil
}
func ExtractCityAndDistrict(ctx context.Context) (city, district string, err error) {
	var result []string
	err = chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			(() => {
				const container = document.querySelector("#imooo");
				if (!container) return [];
				const links = container.querySelectorAll("li > a");
				const names = Array.from(links).map(link => link.innerText.trim());
				return [names[1], names.length > 4 ? names[names.length - 1] : ""];
			})()
		`, &result),
	)

	if err != nil {
		return "", "", err
	}

	if len(result) > 0 {
		city = result[0]
	}
	if len(result) > 1 {
		district = result[1]
	}
	return city, district, nil
}
func ExtractListingID(adURL string) string {
	// Regular expression to match only the numeric ID before .html
	re := regexp.MustCompile(`(\d+)\.html$`)
	match := re.FindStringSubmatch(adURL)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ExtractPhoneNumber fetches the phone number for a given ad URL by querying the API
func ExtractPhoneNumber(ctx context.Context, adURL string) (string, error) {
	fmt.Println("we are in consumer")
	// Extract the listing ID from the ad URL to construct the API URL

	listingID := ExtractListingID(adURL)
	apiURL := fmt.Sprintf("https://www.sheypoor.com/api/v10.0.0/listings/%s/number", listingID)

	setCookie := func(name, value string) chromedp.ActionFunc {
		return chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetCookie(name, value).
				WithDomain("www.sheypoor.com").
				WithPath("/").
				WithHTTPOnly(true).
				Do(ctx)
		})
	}
	refreshToken := os.Getenv("SHEYPOOR_TOKEN")

	var responseText string
	err := chromedp.Run(ctx,
		network.Enable(),
		setCookie("refresh_token", refreshToken),
		chromedp.Navigate(apiURL),
		chromedp.Text("body", &responseText, chromedp.ByQuery),
	)
	if err != nil {
		return "", fmt.Errorf("failed to set cookies and navigate to API: %v", err)
	}

	return responseText, nil
}

// ExtractTitle extracts the listing title from the page
func ExtractTitle(ctx context.Context) (string, error) {
	var title string
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#listing-title`, chromedp.ByID),
		chromedp.Evaluate(`document.querySelector('#listing-title')?.innerText`, &title),
	)
	return title, err
}

// ExtractImageURL extracts a single image URL from the ad page
func ExtractImageURL(ctx context.Context) (string, error) {
	var url string
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
            (() => {
                const img = document.querySelector('.swiper .swiper-slide img');
                return img && img.src ? img.src : "";
            })()
        `, &url),
	)
	if err != nil {
		return "", err
	}
	return url, nil
}
func ExtractDescription(ctx context.Context) (string, error) {
	var descriptionHTML string

	// Run JavaScript in the browser context to find and retrieve the description HTML
	err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(`
			(() => {
				const descriptionDivs = document.querySelectorAll("div");
				for (let i = 0; i < descriptionDivs.length; i++) {
					if (descriptionDivs[i].innerText.trim() === "توضیحات:") {
						const nextDiv = descriptionDivs[i].nextElementSibling;
						if (nextDiv) {
							return nextDiv.innerHTML; // Return the HTML of the following div
						}
					}
				}
				return ""; // Return an empty string if no match is found
			})()
		`, &descriptionHTML),
	)
	if err != nil {
		return "", err
	}

	// Replace <br> tags with newline characters
	description := strings.ReplaceAll(descriptionHTML, "<br>", "\n")

	// Use a regular expression to remove <span> tags and their content
	re := regexp.MustCompile(`<span[^>]*>.*?</span>`)
	description = re.ReplaceAllString(description, "")

	// Remove any remaining HTML tags (if needed)
	description = removeAllHTMLTags(description)

	return strings.TrimSpace(description), nil
}

// Helper function to remove all HTML tags, except for <br> tags which have already been replaced
func removeAllHTMLTags(input string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(input, "")
}
