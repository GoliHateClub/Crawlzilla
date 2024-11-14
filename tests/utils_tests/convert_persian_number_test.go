package tests

import (
	"Crawlzilla/utils"
	"testing"
)

func TestConvertPersianNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"۱۲۳۴۵۶", 123456, false},
		{"۱٬۲۳۴٬۵۶۷", 1234567, false}, // with commas
		{"۰", 0, false},
		{"۹۸۷۶۵۴۳۲۱", 987654321, false},
		{"۵۰۰۰", 5000, false},
		{"۱۰۰.۵۰", 10050, false},    // with a period
		{"۴۵۶abc۷۸۹", 456789, true}, // mixed characters
		{"xyz", 0, true},            // non-numeric input
	}

	for _, test := range tests {
		result, err := utils.ConvertPersianNumber(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("For input %s, expected %d, but got %d", test.input, test.expected, result)
			}
		}
	}
}
