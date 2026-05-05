package dto

type MedicationStatus string

const (
	MedicationOk      MedicationStatus = "ok"
	MedicationExpired MedicationStatus = "expired"
)

type MedicationItemDTO struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Desc   string           `json:"desc"`
	Price  int64            `json:"price"`
	Expiry string           `json:"expiry"`
	Status MedicationStatus `json:"status"`
}
