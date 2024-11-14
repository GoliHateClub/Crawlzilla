package utils

import (
	"strconv"
	"strings"
)

// ConvertPersianNumber converts a string with Persian numbers to an integer.
func ConvertPersianNumber(persianNum string) (int, error) {
	// Define a map for Persian to English digit conversion
	persianToEnglish := map[rune]rune{
		'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
		'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
	}

	// Convert Persian digits to English digits
	var englishNum strings.Builder
	for _, r := range persianNum {
		// Skip commas or periods
		if r == ',' || r == '.' || r == '٬' {
			continue
		}
		if englishDigit, exists := persianToEnglish[r]; exists {
			englishNum.WriteRune(englishDigit)
		} else {
			englishNum.WriteRune(r) // Keep non-Persian characters as they are
		}
	}

	// Convert the result to an integer
	return strconv.Atoi(englishNum.String())
}
