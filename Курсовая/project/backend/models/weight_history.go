package models

import "time"

type WeightHistory struct {
	RecordID   int64     `json:"record_id" gorm:"primaryKey;autoIncrement"`
	PetID      int64     `json:"pet_id" gorm:"not null"`
	VisitID    *int64    `json:"visit_id" gorm:"default:null"`
	Weight     float64   `json:"weight_kg" gorm:"column:weight_kg;type:decimal(5,2);not null;check:weight_kg >= 0.1 AND weight_kg <= 300"`
	MeasuredAt time.Time `json:"measured_at" gorm:"type:timestamp;not null;autoCreateTime"`
	DoctorID   *int64    `json:"doctor_id" gorm:"default:null"`
}
