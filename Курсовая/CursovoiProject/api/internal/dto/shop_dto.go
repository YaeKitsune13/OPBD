package dto

type CartItemResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"productId"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Quantity  uint   `json:"quantity"`
}

type OrderRequest struct {
	UserID      uint `json:"userId"`
	TotalAmount uint `json:"totalAmount"`
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`
	Stock       int    `json:"stock"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	Date        string              `json:"date"`
	TotalAmount uint                `json:"totalAmount"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ProductName string `json:"name"`
	Price       int    `json:"price"`
	Quantity    uint   `json:"quantity"`
}
