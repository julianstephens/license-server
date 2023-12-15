package model

import (
	"gorm.io/datatypes"
)

type License struct {
	Base
	ProductId  string         `json:"product_id"`
	Key        []byte         `gorm:"type:bytea" json:"key"`
	Attributes datatypes.JSON `json:"metadata"`
}

type LicenseRequest struct {
	Key     string `json:"key" binding:"required"`
	Machine string `json:"machine" binding:"omitempty"`
}
