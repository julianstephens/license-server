package model

type License struct {
	Base
	ProductId string `json:"product_id"`
	Product   Product
	Value     string
}
