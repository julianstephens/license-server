package model

type UserGroup struct {
	Base
	Users []*User `gorm:"many2many:group_members;"`
	Rules []*Rule `gorm:"many2many:group_rules;"`
}
