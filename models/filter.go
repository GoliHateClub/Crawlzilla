package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Filters struct definition as before
type Filters struct {
	ID           string `gorm:"type:uuid;primary_key;"`
	USER_ID      string `gorm:"type:uuid;"`
	USER         Users  `gorm:"foreignKey:USER_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	City         string `gorm:"type:varchar(32)"`
	Neighborhood string `gorm:"type:varchar(32)"`
	Reference    string `gorm:"type:varchar(10)"`
	CategoryType string `gorm:"type:varchar(10)"`
	PropertyType string `gorm:"type:varchar(10)"`
	Sort         string `gorm:"type:varchar(10)"`
	Order        string `gorm:"type:varchar(10)"`
	Area         int    `gorm:"type:int"`
	Price        int    `gorm:"type:int"`
	Rent         int    `gorm:"type:int"`
	Room         int    `gorm:"type:int"`
	FloorNumber  int    `gorm:"type:int"`
	VisitCount   int    `gorm:"type:int"`
	HasElevator  bool   `gorm:"type:boolean"`
	HasStorage   bool   `gorm:"type:boolean"`
	HasParking   bool   `gorm:"type:boolean"`
	HasBalcony   bool   `gorm:"type:boolean"`
}

func (c *Filters) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()

	return err
}
