package dto

import (
	"example/project/backend/models"
)

type AppointmentRowDTO struct {
	AppointmentId int64         `json:"appointmentId" binding:"required"`
	PetLabel      string        `json:"petLabel" binding:"required"`
	DoctorName    string        `json:"doctorName" binding:"required"`
	Specialty     string        `json:"specialty"`
	ScheduledDate string        `json:"scheduledDate"`
	ScheduledTime string        `json:"scheduledTime"`
	Status        models.Status `json:"status"`
}
