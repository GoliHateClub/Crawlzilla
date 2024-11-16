package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func MapToStruct(data interface{}, output interface{}) error {
	// Convert map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	// Convert JSON to struct
	err = json.Unmarshal(jsonData, output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON to struct: %w", err)
	}

	return nil
}

// GenerateRandomNumber generates a random number with the specified number of digits
func GenerateRandomNumber(length int) string {

	// Calculate the range for the random number
	min := int64(math.Pow10(length - 1))
	max := int64(math.Pow10(length)) - 1
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))
	// Generate the random number
	randomNumber := rand.Int63n(max-min+1) + min
	return strconv.Itoa(int(randomNumber))
}
