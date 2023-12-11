package model

type APIKey struct {
	Base
	UserId string `gorm:"index"`
	User   User   `json:"-"`
	Mask   string `gorm:"index" json:"-"`
	Key    []byte `gorm:"type:bytea"`
	Scopes string
}

type DisplayAPIKey struct {
	Base
	UserId string
	User   User   `json:"-"`
	Mask   string `json:"-"`
	Key    string
	Scopes string
}
