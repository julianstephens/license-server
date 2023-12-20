package model

type User struct {
	Base
	Name     string   `json:"name,omitempty"`
	Email    string   `gorm:"not null;unique;uniqueIndex" json:"email,omitempty" binding:"omitempty,email"`
	Password []byte   `gorm:"type:bytea" json:"-"`
	ApiKeys  []APIKey `gorm:"constraint:onDelete:CASCADE" json:"-"`
}

type UserWithScopes struct {
	User   User     `json:"user"`
	Scopes []string `json:"authentication_scopes"`
}

type AuthRequest struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}
