package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type License struct {
	Base
	ExternalId uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();unique;uniqueIndex" json:"external_id"`
	ProductId  string         `json:"product_id"`
	Key        []byte         `gorm:"type:bytea" json:"key"`
	Attributes datatypes.JSON `json:"metadata"`
	Revoked    bool           `gorm:"default:false" json:"revoked" binding:"omitempty"`
}

type LicenseRequest struct {
	Key     string `json:"key" binding:"required"`
	Machine string `json:"machine" binding:"omitempty"`
}

type ActivationData struct {
	Id          uuid.UUID `json:"license_id"`
	Product     string    `json:"product"`
	Key         string    `json:"product_key"`
	IssueDate   int64     `json:"issue_date"`
	RefreshDate int64     `json:"refresh_date"`
	EndDate     int64     `json:"expiration_date"`
}

type LicenseRevokeRequest struct {
	Id uuid.UUID `json:"id"`
}
