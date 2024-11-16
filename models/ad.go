package models

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdSummary represents a lightweight ad with only the necessary fields
type AdSummary struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
}

// Ads struct definition as before
type Ads struct {
	ID            string  `gorm:"type:uuid;primary_key;"`
	Hash          string  `gorm:"type:char(64);uniqueIndex"` // Unique hash to prevent duplicates
	Title         string  `gorm:"type:varchar(50);not null"`
	Description   string  `gorm:"type:text"`
	LocationURL   string  `gorm:"type:varchar(255)"`
	ImageURL      string  `gorm:"type:varchar(255)"`
	URL           string  `gorm:"type:varchar(255)"`
	City          string  `gorm:"type:varchar(32)"`
	Neighborhood  string  `gorm:"type:varchar(32)"`
	ContactNumber string  `gorm:"type:varchar(32)"`
	Reference     string  `gorm:"type:varchar(10)"`
	CategoryType  string  `gorm:"type:varchar(10)"`
	PropertyType  string  `gorm:"type:varchar(10)"`
	Latitude      float64 `gorm:"type:decimal(9,6)"`
	Longitude     float64 `gorm:"type:decimal(9,6)"`
	Area          int     `gorm:"type:int"`
	Price         int     `gorm:"type:int"`
	Rent          int     `gorm:"type:int"`
	Room          int     `gorm:"type:int"`
	FloorNumber   int     `gorm:"type:int"`
	TotalFloors   int     `gorm:"type:int"`
	VisitCount    int     `gorm:"type:int"`
	HasElevator   bool    `gorm:"type:boolean"`
	HasStorage    bool    `gorm:"type:boolean"`
	HasParking    bool    `gorm:"type:boolean"`
	HasBalcony    bool    `gorm:"type:boolean"`
}

func (c *Ads) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()

	// Create a variable to store the concatenated string
	var hashInput string

	// Use reflection to iterate over the fields of the struct
	val := reflect.ValueOf(c).Elem()

	// Iterate over all fields of the struct (excluding the ID field)
	for i := 0; i < val.NumField(); i++ {
		// Get the field name (ID is excluded)
		fieldName := val.Type().Field(i).Name

		// Skip the "ID" field
		if fieldName == "ID" {
			continue
		}

		// Get the field value and append it to the hash input string
		fieldValue := val.Field(i).String()
		hashInput += fieldValue
	}

	// Create a new SHA-256 hash
	hash := sha256.New()
	hash.Write([]byte(hashInput))

	// Set the hash field with the generated hash
	c.Hash = hex.EncodeToString(hash.Sum(nil))

	// fmt.Println("Generated Hash:", c.Hash) // Log the hash to confirm it's being set correctly

	return nil
}
