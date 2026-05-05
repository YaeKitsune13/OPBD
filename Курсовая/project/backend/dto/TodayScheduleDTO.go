package dto

import "example/project/backend/models"

type TodayScheduleDTO struct {
	AppointmentId int64         `json:"appointmentId"`
	Time          string        `json:"time"`
	PetLabel      string        `json:"petLabel"`
	OwnerName     string        `json:"ownerName"`
	Breed         string        `json:"breed"`
	Reason        string        `json:"reason"`
	Status        models.Status `json:"status"`
}
