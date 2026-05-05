package dto

import (
	"example/project/backend/models"
	"time"
)

type NextAppointmentData struct {
	DateTime   time.Time `json:"date_time"`
	PetName    string    `json:"pet_name"`
	DoctorName string    `json:"doctor_name"`
}

type Pet struct {
	Name    string  `json:"name"`
	Species string  `json:"species"`
	Breed   string  `json:"breed"`
	Weight  float64 `json:"weight"`
}

type DashboardDTO struct {
	ClientName         string               `json:"client_name" binding:"required"`
	PetsCount          int8                 `json:"pets_count"`
	PendingApps        int32                `json:"pending_apps"`
	TotalVisits        int32                `json:"total_visits"`
	NextAppointment    NextAppointmentData  `json:"next_appointment"`
	RecentAppointments []models.Appointment `json:"recent_appointments"`
	PetsShort          []Pet                `json:"pets_short"`
}
