package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
)

type AppointmentService interface {
	GetDoctorTodaySchedule(doctorId int64) ([]dto.TodayScheduleDTO, error)
	CreateAppointment(app models.Appointment) error
	GetUpcomingByOwner(ownerId int64) ([]dto.AppointmentRowDTO, error)
	UpdateStatus(appId int64, status models.Status) error
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	petRepo         repository.PetRepository
	ownerRepo       repository.OwnerRepository
	doctorRepo      repository.DoctorRepository
}

func NewAppointmentService(
	ar repository.AppointmentRepository,
	pr repository.PetRepository,
	or repository.OwnerRepository,
	dr repository.DoctorRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepo: ar,
		petRepo:         pr,
		ownerRepo:       or,
		doctorRepo:      dr,
	}
}

// 1. Расписание для врача (уже было)
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

		result = append(result, dto.TodayScheduleDTO{
			AppointmentId: app.AppointmentId,
			Time:          app.ScheduledAt.Format("15:04"),
			PetLabel:      fmt.Sprintf("%s %s", getEmojiAvatar(pet.Species), pet.Nickname),
			OwnerName:     fmt.Sprintf("%s %s", owner.FirstName, owner.LastName),
			Breed:         pet.Breed,
			Reason:        app.Comment,
			Status:        app.Status,
		})
	}
	return result, nil
}

// 2. Создание новой записи (Client)
func (s *appointmentService) CreateAppointment(app models.Appointment) error {
	app.Status = models.StatusWaiting // По умолчанию статус "Ожидание"
	return s.appointmentRepo.Create(&app)
}

// 3. Список записей для клиента (Экран "Записи")
func (s *appointmentService) GetUpcomingByOwner(ownerId int64) ([]dto.AppointmentRowDTO, error) {
	apps, err := s.appointmentRepo.GetUpcomingByOwnerId(ownerId)
	if err != nil {
		return nil, err
	}

	var result []dto.AppointmentRowDTO
	for _, a := range apps {
		pet, _ := s.petRepo.GetByID(a.PetID)
		doctor, _ := s.doctorRepo.GetByID(a.DoctorID)

		result = append(result, dto.AppointmentRowDTO{
			AppointmentId: a.AppointmentId,
			PetLabel:      fmt.Sprintf("%s %s", getEmojiAvatar(pet.Species), pet.Nickname),
			DoctorName:    fmt.Sprintf("%s %s.", doctor.LastName, string(doctor.FirstName[0])),
			Specialty:     doctor.Speciality,
			ScheduledDate: a.ScheduledAt.Format("02.01.2006"),
			ScheduledTime: a.ScheduledAt.Format("15:04"),
			Status:        a.Status,
		})
	}
	return result, nil
}

// 4. Обновление статуса (Confirm / Cancel)
func (s *appointmentService) UpdateStatus(appId int64, status models.Status) error {
	return s.appointmentRepo.UpdateStatus(appId, status)
}
