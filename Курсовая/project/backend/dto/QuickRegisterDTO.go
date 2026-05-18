package dto

type QuickRegisterRequest struct {
	OwnerName   string `json:"owner_name"`
	OwnerPhone  string `json:"owner_phone"`
	PetName     string `json:"pet_name"`
	Species     string `json:"species"`
	Breed       string `json:"breed"`
	IsAnonymous bool   `json:"is_anonymous"` // Флаг для анонимного входа
}

type QuickRegisterResponse struct {
	PetID     int64  `json:"pet_id"`
	PetName   string `json:"pet_name"`
	OwnerName string `json:"owner_name"`
	Breed     string `json:"breed"`
}
