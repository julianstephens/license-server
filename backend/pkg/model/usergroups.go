package model

type UserGroup struct {
	Base
	Name  string  `gorm:"uniqueIndex" binding:"alpha"`
	Users []*User `gorm:"many2many:group_members;"`
	Rules []*Rule `gorm:"many2many:group_rules;"`
}
