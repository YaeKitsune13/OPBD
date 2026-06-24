package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusAppointment string

const (
	StatusWaiting   StatusAppointment = "waiting"
	StatusConfirmed StatusAppointment = "confirmed"
	StatusRejected  StatusAppointment = "rejected"
	StatusDone      StatusAppointment = "done"
)

type Appointment struct {
	gorm.Model

	ClientID    uint              `json:"clientId"`
	PetID       uint              `json:"petId"`
	DoctorID    uint              `json:"doctorId"`
	ServiceID   *uint             `json:"serviceId"`
	ScheduledAt time.Time         `json:"scheduledAt"`
	Status      StatusAppointment `json:"status" gorm:"type:varchar(20);default:'waiting'"`
	Comment     string            `json:"comment"`
	Pet         Pet               `gorm:"foreignKey:PetID" json:"pet"`
	Client      User              `gorm:"foreignKey:ClientID" json:"client"`
	Doctor      User              `gorm:"foreignKey:DoctorID" json:"doctor"`
	Service     Service           `gorm:"foreignKey:ServiceID" json:"service"`
	Diagnoses   []Diagnosis       `gorm:"many2many:appointment_diagnoses;" json:"diagnoses"`
	Protocol    *MedicalProtocol  `gorm:"foreignKey:AppointmentID" json:"protocol"`
}
