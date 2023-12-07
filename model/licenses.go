package model

type License struct {
	Base
	UserId    string `json:"user_id"`
	User      User
	ProductId string  `json:"product_id"`
	Product   Product `gorm:"not null"`
	Value     string  `gorm:"not null;unique" binding:"alphanum"`
}
