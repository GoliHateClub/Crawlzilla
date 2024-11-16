package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

// Constants for the roles
const (
	RoleUser       Role = "user"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "super-admin"
)

// Users struct definition as before
type Users struct {
	ID          string `gorm:"type:uuid;primary_key;"`
	Telegram_ID string `gorm:"type:varchar(10)"`
	Role        string `gorm:"type:varchar(15)"`

	Filers []Filters `gorm:"foreignKey:USER_ID"`
}

func (c *Users) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()

	return nil
}
