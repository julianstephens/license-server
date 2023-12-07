package model

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"primary_key;" json:"id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	DeletedAt *int   `gorm:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	guid := xid.New()
	base.ID = guid.String()
	return nil
}
