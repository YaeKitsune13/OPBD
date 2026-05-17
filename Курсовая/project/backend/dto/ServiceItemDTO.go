package dto

type ServiceItemDTO struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int64  `json:"price"`
}
