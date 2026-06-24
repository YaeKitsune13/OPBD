package service

import (
	"api/internal/dto"
	"api/internal/repository"
)

type DashboardService interface {
	GetDashboardData(userID uint) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	repo    repository.DashboardRepository
	petRepo repository.PetRepository
}

func NewDashboardService(repo repository.DashboardRepository, petRepo repository.PetRepository) DashboardService {
	return &dashboardService{repo, petRepo}
}

func (s *dashboardService) GetDashboardData(userID uint) (*dto.DashboardResponse, error) {
	pCount, aCount, vCount, err := s.repo.GetSummary(userID)
	if err != nil {
		return nil, err
	}

	pets, _ := s.petRepo.GetByOwnerID(userID)
	petDTOs := make([]dto.PetBriefDTO, 0)
	for _, p := range pets {
		petDTOs = append(petDTOs, dto.PetBriefDTO{
			ID:     p.ID,
			Name:   p.Name,
			Breed:  p.Breed,
			Avatar: p.Avatar,
			Weight: p.Weight,
			Dob:    p.BirthDate.Format("02.01.2006"),
		})
	}

	nextApp, _ := s.repo.GetNextAppointment(userID)
	var nextDTO *dto.NextAppointmentDTO
	if nextApp != nil && nextApp.ID != 0 {
		nextDTO = &dto.NextAppointmentDTO{
			Date:    nextApp.ScheduledAt.Format("02.01.2006"),
			Time:    nextApp.ScheduledAt.Format("15:04"),
			PetName: nextApp.Pet.Name,
		}
	}

	recentApps, _ := s.repo.GetRecentAppointments(userID, 5)
	recentDTOs := make([]dto.RecentAppointmentDTO, 0)
	for _, a := range recentApps {
		docName := "Врач"
		if a.Doctor.LastName != "" {
			docName = a.Doctor.LastName
			if a.Doctor.FirstName != "" {
				docName += " " + string([]rune(a.Doctor.FirstName)[0]) + "."
			}
		}

		serviceName := "Приём"
		if a.Service.Name != "" {
			serviceName = a.Service.Name
		}

		recentDTOs = append(recentDTOs, dto.RecentAppointmentDTO{
			ID:         a.ID,
			PetName:    a.Pet.Name,
			DoctorName: docName,
			Service:    serviceName,
			Date:       a.ScheduledAt.Format("02.01.2006"),
			Time:       a.ScheduledAt.Format("15:04"),
			Status:     string(a.Status),
		})
	}

	return &dto.DashboardResponse{
		PetsCount:          pCount,
		AppointmentsCount:  aCount,
		VisitsCount:        vCount,
		NextAppointment:    nextDTO,
		RecentAppointments: recentDTOs,
		Pets:               petDTOs,
	}, nil
}
