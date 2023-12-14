package model

type Product struct {
	Base
	Name     string           `gorm:"not null;unique;uniqueIndex"`
	Version  string           `gorm:"default:1.0"`
	Licenses []License        `gorm:"foreignKey:ProductId"`
	Features []ProductFeature `gorm:"foreignKey:ProductId"`
}

type ProductFeature struct {
	Base
	ProductId string `json:"product_id"`
	Name      string `gorm:"not null;index" json:"name"`
}

type ProductKeyPair struct {
	Id         string `json:"id" mapstructure:"id"`
	ProductId  string `json:"product_id" mapstructure:"product_id"`
	PrivateKey string `json:"private_key" mapstructure:"private_key"`
	PublicKey  string `json:"public_key" mapstructure:"public_key"`
}
