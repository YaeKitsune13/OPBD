package dto

type DashboardResponse struct {
	PetsCount          int64                  `json:"petsCount"`
	AppointmentsCount  int64                  `json:"appointmentsCount"`
	VisitsCount        int64                  `json:"visitsCount"`
	NextAppointment    *NextAppointmentDTO    `json:"nextAppointment"`
	RecentAppointments []RecentAppointmentDTO `json:"recentAppointments"`
	Pets               []PetBriefDTO          `json:"pets"`
}

type NextAppointmentDTO struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	PetName string `json:"petName"`
}

type MedicalProtocolDTO struct {
	Weight      float64        `json:"weight"`
	Diagnoses   []DiagnosisDTO `json:"diagnoses"`
	Treatment   string         `json:"treatment"`
	Medications string         `json:"medications"`
}
type RecentAppointmentDTO struct {
	ID         uint                `json:"id"`
	PetName    string              `json:"petName"`
	DoctorName string              `json:"doctorName"`
	Service    string              `json:"service"`
	Date       string              `json:"date"`
	Time       string              `json:"time"`
	Status     string              `json:"status"`
	Protocol   *MedicalProtocolDTO `json:"protocol"`
}

type PetBriefDTO struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Breed  string  `json:"breed"`
	Avatar string  `json:"avatar"`
	Weight float64 `json:"weight"`
	Dob    string  `json:"dob"`
}
