package dto

type ServiceItemDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"decs"`
	Price int64  `json:"price"`
}
