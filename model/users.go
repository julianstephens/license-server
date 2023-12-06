package model

type User struct {
	Base
	Name   string
	Groups []*UserGroup `gorm:"many2many:groups;"`
	Rules  []*Rule      `gorm:"many2many:user_rules;"`
}
