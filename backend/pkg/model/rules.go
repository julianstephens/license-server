package model

type Rule struct {
	Base
	Name       string       `gorm:"uniqueIndex"`
	Products   []*Product   `gorm:"many2many:product_rules;"`
	Users      []*User      `gorm:"many2many:user_rules;"`
	UserGroups []*UserGroup `gorm:"many2many:group_rules;" json:"user_groups"`
}
