package model

type User struct {
	Base
	Name     string       `json:"name"`
	Email    string       `gorm:"not null;unique;uniqueIndex" json:"email" binding:"required,email"`
	Password string       `gorm:"uniqueIndex" json:"-"`
	Groups   []*UserGroup `gorm:"many2many:groups;" json:"groups"`
	Rules    []*Rule      `gorm:"many2many:user_rules;" json:"-"`
}
