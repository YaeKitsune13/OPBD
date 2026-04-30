package models

import "time"

type WeightHistory struct {
	RecordID   int64     `json:"record_id" gorm:"primaryKey;autoIncrement"`
	PetID      int64     `json:"pet_id" gorm:"not null"`
	Pet        Pet       `json:"-" gorm:"foreignKey:PetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	VisitID    *int64    `json:"visit_id" gorm:"default:null"`
	Visit      *Visit    `json:"-" gorm:"foreignKey:VisitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Weight     float64   `json:"weight_kg" gorm:"type:decimal(5,2);not null;check:weight_kg >= 0.1 AND weight_kg <= 300"`
	MeasuredAt time.Time `json:"measured_at" gorm:"type:timestamp;not null;autoCreateTime"`
	DoctorID   *int64    `json:"doctor_id" gorm:"default:null"`
	Doctor     *Doctor   `json:"-" gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
