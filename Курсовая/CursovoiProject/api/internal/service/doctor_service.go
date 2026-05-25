package service

import (
	"api/internal/dto"
	"api/internal/models"
	"api/internal/repository"
)

type DoctorService interface {
	GetSchedule(doctorID uint) []dto.DoctorScheduleDTO
	CompleteVisit(appID uint, req dto.CompleteVisitRequest) error
	GetPatients(search string) []dto.PatientDTO
	GetMedicalHistory(clientID uint) []dto.PetHistoryDTO
	UpdateStatus(appID uint, status string) error
}

type doctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo}
}

func (s *doctorService) GetSchedule(doctorID uint) []dto.DoctorScheduleDTO {
	apps, err := s.repo.GetSchedule(doctorID)
	if err != nil {
		return []dto.DoctorScheduleDTO{}
	}

	res := make([]dto.DoctorScheduleDTO, 0)
	for _, a := range apps {
		res = append(res, dto.DoctorScheduleDTO{
			ID:        a.ID,
			Time:      a.ScheduledAt.Format("15:04"),
			Date:      a.ScheduledAt.Format("2006-01-02"),
			PetName:   a.Pet.Name,
			OwnerName: a.Client.LastName + " " + a.Client.FirstName,
			Service:   a.Service.Name,
			Status:    string(a.Status),
		})
	}
	return res
}

func (s *doctorService) CompleteVisit(appID uint, req dto.CompleteVisitRequest) error {
	protocol := &models.MedicalProtocol{
		AppointmentID: appID,
		WeightAtVisit: req.Weight,
		Diagnosis:     req.Diagnosis,
		Treatment:     req.Treatment,
		Medications:   req.Medications,
	}
	return s.repo.CompleteAppointment(appID, protocol)
}

func (s *doctorService) GetPatients(search string) []dto.PatientDTO {
	clients, _ := s.repo.SearchClients(search)
	var res []dto.PatientDTO
	for _, c := range clients {
		res = append(res, dto.PatientDTO{
			ID:        c.ID,
			FullName:  c.LastName + " " + c.FirstName,
			Phone:     c.Phone,
			Email:     c.Email,
			PetsCount: len(c.Pets),
		})
	}
	return res
}

func (s *doctorService) GetMedicalHistory(clientID uint) []dto.PetHistoryDTO {
	pets, _ := s.repo.GetClientHistory(clientID)
	res := make([]dto.PetHistoryDTO, 0)
	for _, p := range pets {
		var visits []dto.VisitTimelineDTO
		for _, a := range p.Appointments {
			if a.Protocol != nil {
				visits = append(visits, dto.VisitTimelineDTO{
					ID:        a.ID,
					Date:      a.ScheduledAt.Format("02.01.2006"),
					Diagnosis: a.Protocol.Diagnosis,
					Treatment: a.Protocol.Treatment,
				})
			}
		}
		res = append(res, dto.PetHistoryDTO{PetName: p.Name, PetIcon: p.Avatar, Visits: visits})
	}
	return res
}

func (s *doctorService) UpdateStatus(appID uint, status string) error {
	return s.repo.UpdateAppointmentStatus(appID, status)
}
