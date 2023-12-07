package model

type APIKey struct {
	Base
	UserId string `gorm:"index"`
	User   User
	Mask   string `gorm:"index"`
	Key    string
	Scopes string
}
