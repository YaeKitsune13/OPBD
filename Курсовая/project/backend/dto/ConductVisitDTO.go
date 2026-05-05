package dto

type SelectPet struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Breed  string `json:"breed"`
	Owner  string `json:"owner"`
}

type VisitAssignments struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type ConductVisitDTO struct {
	SelectedPet SelectPet          `json:"selectedPet"`
	Anamnesis   string             `json:"anamnesis"`
	Diagnosis   string             `json:"diagnosis"`
	Assignments []VisitAssignments `json:"assignments"`
	TotalCost   int64              `json:"totalCost"`
}
