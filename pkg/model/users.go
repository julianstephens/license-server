package model

type User struct {
	Base
	Name     string   `json:"name,omitempty"`
	Email    string   `gorm:"not null;unique;uniqueIndex" json:"email,omitempty" binding:"omitempty,email"`
	Password []byte   `gorm:"type:bytea" json:"-"`
	ApiKeys  []APIKey `gorm:"constraint:onDelete:CASCADE"`
}

type UserWithScopes struct {
	User   User     `json:"user"`
	Scopes []string `json:"authentication_scopes"`
}

type AuthRequest struct {
	Name     string `binding:"omitempty"`
	Email    string `binding:"required,email"`
	Password string
}
