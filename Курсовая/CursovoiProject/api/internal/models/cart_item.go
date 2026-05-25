package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint    `json:"userId"`
	ProductID uint    `json:"productId"`
	Quantity  uint    `json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}
