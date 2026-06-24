package dto

type DoctorScheduleDTO struct {
	ID        uint   `json:"id"`
	Time      string `json:"time"`
	Date      string `json:"date"`
	PetName   string `json:"petName"`
	OwnerName string `json:"ownerName"`
	Service   string `json:"service"`
	Status    string `json:"status"`
}

type CompleteVisitRequest struct {
	Weight       float64 `json:"weight" binding:"required"`
	DiagnosisIDs []uint  `json:"diagnosis_ids" binding:"required"`
	Treatment    string  `json:"treatment"`
	Medications  string  `json:"medications"`
}

type PatientDTO struct {
	ID        uint   `json:"id"`
	FullName  string `json:"fullName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	PetsCount int    `json:"petsCount"`
	LastVisit string `json:"lastVisit"`
}

type PetHistoryDTO struct {
	PetName string             `json:"petName"`
	PetIcon string             `json:"petIcon"`
	Visits  []VisitTimelineDTO `json:"visits"`
}

type VisitTimelineDTO struct {
	ID        uint           `json:"id"`
	Date      string         `json:"date"`
	Diagnoses []DiagnosisDTO `json:"diagnoses"`
	Treatment string         `json:"treatment"`
}
