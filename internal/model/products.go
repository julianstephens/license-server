package model

type Product struct {
	Base
	Name     string    `gorm:"not null;unique;uniqueIndex"`
	Version  string    `gorm:"default:1.0"`
	Licenses []License `gorm:"foreignKey:ProductId"`
	Rules    []*Rule   `gorm:"many2many:product_rules;"`
}

type ProductKeyPair struct {
	Id         string `json:"id"`
	ProductId  string `json:"product_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}
