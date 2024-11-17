package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Filters struct definition as before
type Filters struct {
	ID             string    `gorm:"type:uuid;primary_key;"`
	USER_ID        string    `gorm:"type:uuid;"`
	USER           Users     `gorm:"foreignKey:USER_ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	Title          string    `gorm:"type:varchar(32)"`
	City           string    `gorm:"type:varchar(32)"`
	Neighborhood   string    `gorm:"type:varchar(32)"`
	Reference      string    `gorm:"type:varchar(10)"`
	CategoryType   string    `gorm:"type:varchar(10)"`
	PropertyType   string    `gorm:"type:varchar(10)"`
	Sort           string    `gorm:"type:varchar(10)"`
	Order          string    `gorm:"type:varchar(10)"`
	MinArea        int       `gorm:"type:int"`
	MaxArea        int       `gorm:"type:int"`
	MinPrice       int       `gorm:"type:int"`
	MaxPrice       int       `gorm:"type:int"`
	MinRent        int       `gorm:"type:int"`
	MaxRent        int       `gorm:"type:int"`
	MinRoom        int       `gorm:"type:int"`
	MaxRoom        int       `gorm:"type:int"`
	MinFloorNumber int       `gorm:"type:int"`
	MaxFloorNumber int       `gorm:"type:int"`
	HasElevator    bool      `gorm:"type:boolean"`
	HasStorage     bool      `gorm:"type:boolean"`
	HasParking     bool      `gorm:"type:boolean"`
	HasBalcony     bool      `gorm:"type:boolean"`
}

func (c *Filters) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()
	return nil
}
