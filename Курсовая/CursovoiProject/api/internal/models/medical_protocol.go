package models

import "gorm.io/gorm"

type MedicalProtocol struct {
	gorm.Model
	AppointmentID uint    `gorm:"uniqueIndex" json:"appointmentId"`
	WeightAtVisit float64 `json:"weight"`
	Treatment     string  `json:"treatment"`
	Medications   string  `json:"medications"`
}
