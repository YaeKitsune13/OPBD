package dto

type PetRequest struct {
	Name      string  `json:"name" binding:"required"`
	Species   string  `json:"species" binding:"required"`
	Breed     string  `json:"breed"`
	BirthDate string  `json:"dob" binding:"required"`
	Weight    float64 `json:"weight"`
	Avatar    string  `json:"avatar"`
}

type PetResponse struct {
	ID        uint    `json:"petId"`
	Name      string  `json:"name"`
	Species   string  `json:"species"`
	Breed     string  `json:"breed"`
	BirthDate string  `json:"dob"`
	Weight    float64 `json:"weight"`
	Avatar    string  `json:"avatar"`
}
