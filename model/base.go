package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string     `gorm:"primary_key;" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	guid := xid.New()
	// if guid == nil {
	// 	return fmt.Errorf("Unable to generate DB UID")
	// }

	base.ID = guid.String()
	return nil
}
