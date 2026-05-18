package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
	"time"
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
	visits, err := s.visitRepo.GetByPetID(petId)
	if err != nil {
		return nil, err
	}

	var history []dto.HealthJournalDTO
	for _, v := range visits {
		item := dto.HealthJournalDTO{
			VisitId:   v.VisitId,
			Date:      v.VisitDate.Format("02.01.2006"),
			Time:      v.VisitDate.Format("15:04"),
			Diagnosis: v.Diagnosis,
			Details:   v.Anamnesis,
			Price:     fmt.Sprintf("%.2f ₽", v.TotalCost),
		}
		history = append(history, item)
	}
	return history, nil
}

func (s *healthJournalService) SaveVisit(data dto.ConductVisitDTO) error {
	newVisit := &models.Visit{
		PetID:     data.SelectedPet.PetId,
		Anamnesis: data.Anamnesis,
		Diagnosis: data.Diagnosis,
		TotalCost: float64(data.TotalCost),
		VisitDate: time.Now(),
	}

	// Если Id (AppointmentID) равен 0, в базу запишется NULL (благодаря указателю в модели)
	if data.SelectedPet.Id > 0 {
		id := data.SelectedPet.Id
		newVisit.AppointmentID = &id
	} else {
		newVisit.AppointmentID = nil
	}

	if err := s.visitRepo.Create(newVisit); err != nil {
		return fmt.Errorf("ошибка создания визита: %w", err)
	}

	for _, a := range data.Assignments {
		prescription := &models.VisitPrescription{
			VisitID:   newVisit.VisitId,
			ItemType:  models.TypeService,
			ServiceID: &a.Id,
			Quantity:  int64(a.Qty),
			UnitPrice: a.Price,
		}
		if err := s.visitRepo.AddPrescription(prescription); err != nil {
			return fmt.Errorf("ошибка сохранения услуги %s: %w", a.Name, err)
		}
	}

	return nil
}
