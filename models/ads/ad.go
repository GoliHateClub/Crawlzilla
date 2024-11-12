package ads

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type CrawlResult struct {
	gorm.Model
	Title            string  `gorm:"type:varchar(255);not null"`
	PropertyType     string  `gorm:"type:varchar(255)"`
	Description      string  `gorm:"type:text"`
	LocationURL      string  `gorm:"type:varchar(255)"`
	ImageURL         string  `gorm:"type:varchar(255)"`
	URL              string  `gorm:"type:varchar(255)"`
	Reference        string  `gorm:"type:varchar(16)"`
	BuildingAgeType  string  `gorm:"type:varchar(255)"`
	City             string  `gorm:"type:varchar(255)"`
	District         string  `gorm:"type:varchar(255)"`
	Latitude         float64 `gorm:"type:decimal(9,6)"`
	Longitude        float64 `gorm:"type:decimal(9,6)"`
	BuildingAgeValue int     `gorm:"type:int"`
	Area             int     `gorm:"type:int"`
	Price            int     `gorm:"type:int"`
	Rent             int     `gorm:"type:int"`
	Room             int     `gorm:"type:int"`
	FloorNumber      int     `gorm:"type:int"`
	TotalFloors      int     `gorm:"type:int"`
	ContactNumber    int     `gorm:"type:int"`
	VisitCount       int     `gorm:"type:int"`
	HasElevator      bool    `gorm:"type:boolean"`
	HasStorage       bool    `gorm:"type:boolean"`
	HasParking       bool    `gorm:"type:boolean"`
}

func (cr CrawlResult) String() string {
	var result strings.Builder
	val := reflect.ValueOf(cr)
	typ := reflect.TypeOf(cr)

	result.WriteString("CrawlResult {\n")
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		// Only include non-null fields in the output
		if !isNullValue(fieldValue) {
			result.WriteString(fmt.Sprintf("  %s: %v\n", fieldType.Name, fieldValue.Interface()))
		}
	}
	result.WriteString("}")
	return result.String()
}

// Helper function to check if a value is null
func isNullValue(v reflect.Value) bool {
	// Check for null values based on sql.Null* types
	switch v.Interface().(type) {
	case sql.NullString:
		return !v.Interface().(sql.NullString).Valid
	case sql.NullInt64:
		return !v.Interface().(sql.NullInt64).Valid
	case sql.NullFloat64:
		return !v.Interface().(sql.NullFloat64).Valid
	case sql.NullBool:
		return !v.Interface().(sql.NullBool).Valid
	default:
		return false // For non-nullable types, consider it as non-null
	}
}
