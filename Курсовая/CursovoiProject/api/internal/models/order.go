package models

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderPaid      OrderStatus = "paid"
	OrderConfirmed OrderStatus = "confirmed"
	OrderDelivered OrderStatus = "delivered"
)

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"orderId"`
	ProductID uint    `json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  uint    `json:"quantity"`
	Price     int     `json:"price"`
}

type Order struct {
	gorm.Model
	UserID      uint        `json:"userId"`
	TotalAmount uint        `json:"totalAmount"`
	Status      OrderStatus `gorm:"type:enum('paid','confirmed','delivered')" json:"status"`
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}
