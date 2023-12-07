package model

type APIKey struct {
	Base
	UserId string `gorm:"index"`
	User   User   `json:"-"`
	Mask   string `gorm:"index" json:"-"`
	Key    string
	Scopes string
}
