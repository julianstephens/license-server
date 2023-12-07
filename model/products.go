package model

type Product struct {
	Base
	Name     string    `gorm:"not null;uniqueIndex" binding:"alpha"`
	Licenses []License `gorm:"foreignKey:ProductId"`
	Rules    []*Rule   `gorm:"many2many:product_rules;"`
}
