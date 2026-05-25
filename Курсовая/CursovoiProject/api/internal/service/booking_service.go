package service

import (
	"api/internal/dto"
	"api/internal/models"
	"api/internal/repository"
	"time"
)

type BookingService interface {
	GetInitData(userID uint) (*dto.BookingInitResponse, error)
	GetBusySlots(doctorID uint, date string) ([]string, error)
	CreateAppointment(clientID uint, req dto.CreateAppointmentRequest) error
	GetClientHistory(userID uint) ([]dto.RecentAppointmentDTO, error)
}

type bookingService struct {
	repo    repository.AppointmentRepository
	petRepo repository.PetRepository
}

func NewBookingService(repo repository.AppointmentRepository, petRepo repository.PetRepository) BookingService {
	return &bookingService{repo, petRepo}
}

func (s *bookingService) GetInitData(userID uint) (*dto.BookingInitResponse, error) {
	pets, _ := s.petRepo.GetByOwnerID(userID)
	doctors, _ := s.repo.GetDoctors()
	services, _ := s.repo.GetServices()

	res := &dto.BookingInitResponse{}
	for _, p := range pets {
		res.Pets = append(res.Pets, dto.PetBriefDTO{ID: p.ID, Name: p.Name, Avatar: p.Avatar})
	}
	for _, d := range doctors {
		res.Doctors = append(res.Doctors, dto.DoctorDTO{ID: d.ID, FirstName: d.FirstName, LastName: d.LastName, Specialization: d.Specialization})
	}
	for _, sv := range services {
		res.Services = append(res.Services, dto.ServiceDTO{ID: sv.ID, Name: sv.Name, Price: sv.Price})
	}
	return res, nil
}

func (s *bookingService) GetBusySlots(doctorID uint, date string) ([]string, error) {
	return s.repo.GetBusySlots(doctorID, date)
}

func (s *bookingService) CreateAppointment(clientID uint, req dto.CreateAppointmentRequest) error {
	scheduledAt, err := time.ParseInLocation("2006-01-02T15:04:05", req.ScheduledAt, time.Local)
	if err != nil {
		return err
	}

	app := models.Appointment{
		ClientID:    clientID,
		PetID:       req.PetID,
		DoctorID:    req.DoctorID,
		ServiceID:   req.ServiceID,
		ScheduledAt: scheduledAt,
		Status:      models.StatusWaiting,
		Comment:     req.Comment,
	}
	return s.repo.Create(&app)
}

func (s *bookingService) GetClientHistory(userID uint) ([]dto.RecentAppointmentDTO, error) {
	apps, err := s.repo.GetByClientID(userID)
	if err != nil {
		return nil, err
	}

	var history []dto.RecentAppointmentDTO
	for _, a := range apps {
		item := dto.RecentAppointmentDTO{
			ID:         a.ID,
			PetName:    a.Pet.Name,
			DoctorName: a.Doctor.LastName + " " + a.Doctor.FirstName,
			Date:       a.ScheduledAt.Format("02.01.2006"),
			Time:       a.ScheduledAt.Format("15:04"),
			Status:     string(a.Status),
			Service:    a.Service.Name,
		}

		if a.Protocol != nil {
			item.Protocol = &dto.MedicalProtocolDTO{
				Weight:      a.Protocol.WeightAtVisit,
				Diagnosis:   a.Protocol.Diagnosis,
				Treatment:   a.Protocol.Treatment,
				Medications: a.Protocol.Medications,
			}
		}

		history = append(history, item)
	}
	return history, nil
}
