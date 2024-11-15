package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleAdmin      = "admin"
	RoleUser       = "user"
	RoleSuperAdmin = "super-admin"
)

// Users struct definition as before
type Users struct {
	ID          string `gorm:"type:uuid;primary_key;"`
	Telegram_ID string `gorm:"type:varchar(10)"`
	Role        string `gorm:"type:enum('admin', 'user', 'super-admin');not null;default:'user'"`

	Filers []Filters `gorm:"foreignKey:USER_ID"`
}

func (c *Users) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the ID to a new UUID
	c.ID = uuid.NewString()

	return nil
}
