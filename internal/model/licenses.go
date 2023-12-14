package model

import (
	"gorm.io/datatypes"
)

type License struct {
	Base
	Name      string
	ProductId string `json:"product_id"`
}

type LicenseWithAttributes struct {
	Base
	ProductId  string         `json:"product_id"`
	Key        string         `json:"key"`
	Attributes datatypes.JSON `json:"metadata"`
}

type LicenseRequest struct {
	UserId    string           `binding:"omitempty"`
	ProductId string           `binding:"required"`
	Features  *map[string]bool `binding:"omitempty"`
	StartDate *int64           `binding:"omitempty"`
}
