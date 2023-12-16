package model

type APIKey struct {
	Base
	UserId    string `gorm:"index" json:"user_id"`
	User      User   `json:"-"`
	Mask      string `gorm:"index" json:"-"`
	Key       []byte `gorm:"type:bytea" json:"key"`
	ExpiresAt int64  `gorm:"default:extract(epoch from (now() + interval '1 week'))" json:"expires_at"`
	Scopes    string `json:"authentication_scopes"`
}

type DisplayAPIKey struct {
	Base
	UserId    string `json:"user_id"`
	User      User   `json:"-"`
	Mask      string `json:"-"`
	Key       string `json:"key"`
	ExpiresAt string `json:"expires_at"`
	Scopes    string `json:"authentication_scopes"`
}
