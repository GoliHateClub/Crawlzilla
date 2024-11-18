package models

import (
	"time"

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
	Telegram_ID int64
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Role        Role      `gorm:"type:varchar(15)"`
	ChatID      int64
	Filers      []Filters `gorm:"foreignKey:USER_ID"`
}

func (c *Users) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()

	return nil
}
