package dto

type BookingInitResponse struct {
	Pets     []PetBriefDTO `json:"pets"`
	Doctors  []DoctorDTO   `json:"doctors"`
	Services []ServiceDTO  `json:"services"`
}

type DoctorDTO struct {
	ID             uint   `json:"doctorId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
}

type ServiceDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type CreateAppointmentRequest struct {
	PetID       uint   `json:"pet_id" binding:"required"`
	DoctorID    uint   `json:"doctor_id" binding:"required"`
	ServiceID   *uint  `json:"service_id"`
	ScheduledAt string `json:"scheduled_at" binding:"required"` // "2024-05-20T09:00:00"
	Comment     string `json:"comment"`
}
