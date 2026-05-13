package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
)

type HealthJournalService interface {
	GetPetHistory(petId int64) ([]dto.HealthJournalDTO, error)
	SaveVisit(data dto.ConductVisitDTO) error
}

type healthJournalService struct {
	visitRepo  repository.VisitRepository
	doctorRepo repository.DoctorRepository
}

func NewHealthJournalService(vr repository.VisitRepository, dr repository.DoctorRepository) HealthJournalService {
	return &healthJournalService{
		visitRepo:  vr,
		doctorRepo: dr,
	}
}

func (s *healthJournalService) GetPetHistory(petId int64) ([]dto.HealthJournalDTO, error) {
	// 1. Получаем список визитов для питомца
	visits, err := s.visitRepo.GetByPetID(petId)
	if err != nil {
		return nil, err
	}

	var history []dto.HealthJournalDTO

	for _, v := range visits {
		// 1. Получаем данные врача
		doctor, err := s.doctorRepo.GetByID(v.Appointment.DoctorID)
		doctorName := "Врач не указан"

		if err == nil && doctor != nil {
			// Безопасно берем первую букву имени (работает с кириллицей)
			firstNameRunes := []rune(doctor.User.FirstName)
			initial := ""
			if len(firstNameRunes) > 0 {
				initial = string(firstNameRunes[0])
			}

			// ИСПРАВЛЕНИЕ: Обращаемся к doctor.User.LastName и initial
			doctorName = fmt.Sprintf("%s %s.", doctor.User.LastName, initial)
		}

		// 2. Собираем DTO
		item := dto.HealthJournalDTO{
			VisitId:         v.VisitId,
			Date:            v.Appointment.ScheduledAt.Format("02.01.2006"),
			Time:            v.Appointment.ScheduledAt.Format("15:04"),
			Doctor:          doctorName,
			Diagnosis:       v.Diagnosis,
			Details:         v.Anamnesis,
			Analysis:        v.Analysis,
			Recommendations: v.Recommendations,
			Price:           fmt.Sprintf("%.2f ₽", v.TotalCost),
		}

		history = append(history, item)
	}

	return history, nil
}

func (s *healthJournalService) SaveVisit(data dto.ConductVisitDTO) error {
	// Создаем модель визита на основе DTO
	newVisit := &models.Visit{
		AppointmentID: data.SelectedPet.Id, // Предположим, тут передаем ID записи
		Anamnesis:     data.Anamnesis,
		Diagnosis:     data.Diagnosis,
		TotalCost:     float64(data.TotalCost),
	}

	return s.visitRepo.Create(newVisit)
}
