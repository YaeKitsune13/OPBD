package models

import "time"

type Status string

const (
	StatusWaiting   Status = "waiting"
	StatusConfirmed Status = "confirmed"
	StatusRejected  Status = "rejected"
)

type Appointment struct {
	AppointmentId int64     `json:"appointment_id" gorm:"primaryKey;autoIncrement"`
	PetId         int64     `json:"pet_id" gorm:"not null"`
	Pet           Pet       `json:"-" gorm:"foreignKey:PetId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DoctorId      int64     `json:"doctor_id" gorm:"not null"`
	Doctor        Doctor    `json:"-" gorm:"foreignKey:DoctorId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduledAt   time.Time `json:"scheduled_at" gorm:"type:timestamp;not null"`
	Comment       string    `json:"comment" gorm:"type:text"`
	Status        Status    `json:"status" gorm:"type:varchar(20);default:waiting;check:status IN ('waiting','confirmed','rejected')"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}
