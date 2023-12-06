package model

type Product struct {
	Base
	Name     string
	Licenses []License `gorm:"foreignKey:ProductId"`
	Rules    []*Rule   `gorm:"many2many:product_rules;"`
}
