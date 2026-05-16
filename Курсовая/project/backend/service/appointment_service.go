package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
	"time"
)

type AppointmentService interface {
	GetDoctorTodaySchedule(doctorId int64) ([]dto.TodayScheduleDTO, error)
	CreateAppointment(app models.Appointment) error
	GetUpcomingByOwner(ownerId int64) ([]dto.AppointmentRowDTO, error)
	UpdateStatus(appId int64, status models.Status) error
	GetOccupiedTimeSlots(doctorID int64, dateStr string) ([]string, error)
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	petRepo         repository.PetRepository
	userRepo        repository.UserRepository
	doctorRepo      repository.DoctorRepository
}

func NewAppointmentService(
	ar repository.AppointmentRepository,
	pr repository.PetRepository,
	or repository.UserRepository,
	dr repository.DoctorRepository,
) AppointmentService {
	return &appointmentService{
		appointmentRepo: ar,
		petRepo:         pr,
		userRepo:        or,
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
		owner, err := s.userRepo.GetByID(pet.OwnerID)
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
	app.Status = models.StatusWaiting
	isBusy, err := s.appointmentRepo.IsSlotTaken(app.DoctorID, app.ScheduledAt)
	if err != nil {
		return fmt.Errorf("ошибка при проверке доступности времени: %w", err)
	}

	if isBusy {
		return fmt.Errorf("выбранное время (%s) уже занято этим врачом", app.ScheduledAt.Format("15:04"))
	}

	// 3. Если время свободно — создаем запись
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
		pet, errPet := s.petRepo.GetByID(a.PetID)
		doctor, errDoc := s.doctorRepo.GetByID(a.DoctorID)

		// Формируем PetLabel с проверкой на ошибку
		petLabel := "Питомец удален"
		if errPet == nil && pet != nil {
			petLabel = fmt.Sprintf("%s %s", getEmojiAvatar(pet.Species), pet.Nickname)
		}

		// Формируем Имя врача и специальность
		doctorName := "Врач не назначен"
		specialty := "Общая"
		if errDoc == nil && doctor != nil {
			// ВАЖНО: берем данные из вложенного User (вашего Owner)
			doctorName = fmt.Sprintf("%s %s.",
				doctor.User.LastName,
				string(doctor.User.FirstName[0]),
			)
			specialty = doctor.Speciality
		}

		result = append(result, dto.AppointmentRowDTO{
			AppointmentId: a.AppointmentId,
			PetLabel:      petLabel,
			DoctorName:    doctorName,
			Specialty:     specialty,
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

func (s *appointmentService) GetOccupiedTimeSlots(doctorID int64, dateStr string) ([]string, error) {
	// Парсим дату из строки "2026-05-13"
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты: %w", err)
	}

	// Определяем границы суток
	start := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	// Вызываем репозиторий (создадим этот метод ниже)
	apps, err := s.appointmentRepo.GetByDoctorAndDateRange(doctorID, start, end)
	if err != nil {
		return nil, err
	}

	// Собираем только время занятых слотов
	var occupied []string
	for _, a := range apps {
		// Добавляем в список только те, что не отменены/не отклонены
		if a.Status != models.StatusRejected {
			occupied = append(occupied, a.ScheduledAt.Format("15:04"))
		}
	}

	return occupied, nil
}
