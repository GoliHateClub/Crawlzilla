package ads

import "gorm.io/gorm"

type CrawlResult struct {
	gorm.Model
	Title         string  `gorm:"type:varchar(50);not null"`
	Description   string  `gorm:"type:text"`
	LocationURL   string  `gorm:"type:varchar(255)"`
	ImageURL      string  `gorm:"type:varchar(255)"`
	URL           string  `gorm:"type:varchar(255)"`
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
	ContactNumber int     `gorm:"type:int"`
	VisitCount    int     `gorm:"type:int"`
	HasElevator   bool    `gorm:"type:boolean"`
	HasStorage    bool    `gorm:"type:boolean"`
	HasParking    bool    `gorm:"type:boolean"`
	HasBalcony    bool    `gorm:"type:boolean"`
}
