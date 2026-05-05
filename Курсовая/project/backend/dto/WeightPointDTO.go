package dto

type WeightPointDTO struct {
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Date       string  `json:"date"`
	DoctorName string  `json:"doctorName"`
}
