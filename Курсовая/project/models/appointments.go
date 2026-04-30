package models

import "time"

type Appointments struct {
	AppointmentId int64     `json:"appointment_id" gorm:"primaryKey;autoIncrement"`
	PetId         int64     `json:"pet_id" gorm:"not null"`
	Pet           Pet       `json:"-" gorm:"foreignKey:PetId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DoctorId      int64     `json:"doctor_id" gorm:"not null"`
	Doctor        Doctor    `json:"-" gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduledAt   time.Time `json:"scheduled_at" gorm:"type:datetime;not null"`
	Comment       string    `json:"comment"`
}
