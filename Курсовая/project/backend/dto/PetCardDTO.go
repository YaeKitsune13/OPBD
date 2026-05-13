package dto

type PetCardDTO struct {
	PetId   int64   `json:"petId"`
	Name    string  `json:"name" binding:"required"`
	Species string  `json:"species"`
	Breed   string  `json:"breed"`
	Dob     string  `json:"dob"`
	Weight  float64 `json:"weight"`
	Avatar  string  `json:"avatar"`
}
