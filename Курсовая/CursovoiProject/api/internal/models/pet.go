package models

import (
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	OwnerID      uint          `json:"ownerId"`
	Name         string        `json:"name"`
	Species      string        `json:"species"`
	Breed        string        `json:"breed"`
	BirthDate    time.Time     `json:"dob"`
	Weight       float64       `json:"weight"`
	Avatar       string        `json:"avatar"`
	Appointments []Appointment `gorm:"foreignKey:PetID" json:"appointments,omitempty"`
}
