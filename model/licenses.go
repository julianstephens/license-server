package model

import (
	"gorm.io/datatypes"
)

type License struct {
	Base
	UserId    string `json:"user_id"`
	User      User
	ProductId string  `json:"product_id"`
	Product   Product `gorm:"not null"`
	Value     string  `gorm:"not null;unique"`
}

type LicenseWithAttributes struct {
	Base
	Key        string `json:"key"`
	Attributes datatypes.JSON
}
