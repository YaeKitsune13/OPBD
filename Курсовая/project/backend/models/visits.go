package models

import "time"

type Visit struct {
	VisitId       int64       `json:"visit_id" gorm:"primaryKey;autoIncrement"`
	AppointmentID int64       `json:"appointment_id" gorm:"not null"`
	Appointment   Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`
	Anamnesis     string      `json:"anamnesis" gorm:"type:text;not null;check:length(anamnesis) >= 10 AND length(anamnesis) <= 2000"`
	Diagnosis     string      `json:"diagnosis" gorm:"type:varchar(500);not null"`
	TotalCost     float64     `json:"total_cost" gorm:"type:decimal(10,2);not null;default:0"`
	VisitDate     time.Time   `json:"visit_date" gorm:"autoCreateTime"`
}
