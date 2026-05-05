package service

import (
	"example/project/backend/dto"
	"example/project/backend/repository"
	"fmt"
)

type AppointmentService interface {
	GetDoctorTodaySchedule(doctorId int64) ([]dto.TodayScheduleDTO, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	petRepo         repository.PetRepository
	ownerRepo       repository.OwnerRepository
}

func NewAppointmentService(
	ar repository.AppointmentRepository,
	pr repository.PetRepository,
	or repository.OwnerRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepo: ar,
		petRepo:         pr,
		ownerRepo:       or,
	}
}

func (s *appointmentService) GetDoctorTodaySchedule(doctorId int64) ([]dto.TodayScheduleDTO, error) {
	schedules, err := s.appointmentRepo.GetTodaySchedule(doctorId)
	if err != nil {
		return nil, err
	}

	var result []dto.TodayScheduleDTO

	for _, app := range schedules {
		pet, err := s.petRepo.GetByID(app.PetID)
		if err != nil {
			continue
		}

		owner, err := s.ownerRepo.GetByID(pet.OwnerID)
		if err != nil {
			continue
		}

		row := dto.TodayScheduleDTO{
			AppointmentId: app.AppointmentId,
			Time:          app.ScheduledAt.Format("15:04"),
			PetLabel:      fmt.Sprintf("%s %s", getEmojiAvatar(pet.Species), pet.Nickname),
			OwnerName:     fmt.Sprintf("%s %s", owner.FirstName, owner.LastName),
			Breed:         pet.Breed,
			Reason:        app.Comment,
			Status:        app.Status,
		}

		result = append(result, row)
	}

	return result, nil
}
