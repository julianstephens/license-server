package model

type Product struct {
	Base
	Name     string    `gorm:"not null;unique;uniqueIndex"`
	Version  string    `gorm:"default:1.0"`
	Licenses []License `gorm:"foreignKey:ProductId"`
	Rules    []*Rule   `gorm:"many2many:product_rules;"`
}
