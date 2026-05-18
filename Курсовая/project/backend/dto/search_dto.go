package dto

type PatientSearchResultDTO struct {
	PetID      int64   `json:"pet_id"`
	PetName    string  `json:"pet_name"`
	Species    string  `json:"species"`
	Breed      string  `json:"breed"`
	Weight     float64 `json:"weight"`
	Age        int     `json:"age"`
	OwnerID    int64   `json:"owner_id"`
	OwnerName  string  `json:"owner_name"`
	OwnerPhone string  `json:"owner_phone"`
}
