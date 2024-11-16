package utils

import (
	"encoding/json"
	"fmt"
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
