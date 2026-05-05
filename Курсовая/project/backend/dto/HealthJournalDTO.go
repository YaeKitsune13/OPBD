package dto

type HealthJournalDTO struct {
	VisitId         int64   `json:"visitId"`
	Date            string  `json:"date"`
	Time            string  `json:"time"`
	Doctor          string  `json:"doctor"`
	Diagnosis       string  `json:"diagnosis"`
	Details         string  `json:"details"`
	Analysis        *string `json:"analysis"`
	Recommendations *string `json:"recommendations"`
	Price           string  `json:"price"`
}
