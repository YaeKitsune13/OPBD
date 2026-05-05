package dto

import (
	"example/project/backend/models"
	"time"
)

type NextAppointmentData struct {
	DateTime   time.Time `json:"dateTime"`
	PetName    string    `json:"petName"`
	DoctorName string    `json:"doctorName"`
}

type Pet struct {
	Name    string  `json:"name"`
	Species string  `json:"species"`
	Breed   string  `json:"breed"`
	Weight  float64 `json:"weight"`
}

type DashboardDTO struct {
	ClientName         string               `json:"client_name" binding:"required"`
	PetsCount          int8                 `json:"petsCount"`
	PendingApps        int32                `json:"pendingApps"`
	TotalVisits        int32                `json:"totalVisits"`
	NextAppointment    NextAppointmentData  `json:"nextAppointment"`
	RecentAppointments []models.Appointment `json:"recentAppointments"`
	PetsShort          []Pet                `json:"petsShort"`
}
