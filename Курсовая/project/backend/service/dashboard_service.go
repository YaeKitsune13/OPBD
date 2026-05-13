package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
)

type DashboardService interface {
	GetDashboardData(ownerId int64) (*dto.DashboardDTO, error)
}

type dashboardService struct {
	userRepo        repository.UserRepository
	petRepo         repository.PetRepository
	appointmentRepo repository.AppointmentRepository
	visitRepo       repository.VisitRepository // Добавили для подсчета посещений
}

func NewDashboardService(
	or repository.UserRepository,
	pr repository.PetRepository,
	ar repository.AppointmentRepository,
	vr repository.VisitRepository,
) DashboardService {
	return &dashboardService{
		userRepo:        or,
		petRepo:         pr,
		appointmentRepo: ar,
		visitRepo:       vr,
	}
}

func (s *dashboardService) GetDashboardData(ownerID int64) (*dto.DashboardDTO, error) {
	// 1. Получаем владельца
	owner, err := s.userRepo.GetByID(ownerID)
	if err != nil {
		return nil, err
	}

	// 2. Получаем всех питомцев
	pets, err := s.petRepo.GetByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	// 3. Получаем будущие записи
	upcomingApps, err := s.appointmentRepo.GetUpcomingByOwnerId(ownerID)
	if err != nil {
		return nil, err
	}

	// Начальные данные для DTO
	dashboard := &dto.DashboardDTO{
		ClientName:  owner.FirstName,
		PetsCount:   int8(len(pets)),
		TotalVisits: 0,
		PendingApps: 0,
	}

	// 4. Считаем общее кол-во визитов по всем питомцам
	for _, pet := range pets {
		history, _ := s.visitRepo.GetByPetID(pet.PetID)
		dashboard.TotalVisits += int32(len(history))

		dashboard.PetsShort = append(dashboard.PetsShort, dto.Pet{
			Name:    pet.Nickname,
			Species: pet.Species,
			Breed:   pet.Breed,
			Weight:  pet.CurrentWeight,
		})
	}

	// 5. Считаем записи в статусе "Ожидание" и выбираем ближайшую
	for _, app := range upcomingApps {
		if app.Status == models.StatusWaiting {
			dashboard.PendingApps++
		}
	}

	// 6. Заполняем данные о ближайшем приеме (если есть)
	if len(upcomingApps) > 0 {
		next := upcomingApps[0]
		dashboard.NextAppointment = dto.NextAppointmentData{
			DateTime:   next.ScheduledAt,
			PetName:    "Питомец",
			DoctorName: "Врач",
		}
	}

	// 7. Последние записи для виджета (например, последние 3)
	dashboard.RecentAppointments = upcomingApps
	if len(upcomingApps) > 3 {
		dashboard.RecentAppointments = upcomingApps[:3]
	}

	return dashboard, nil
}
