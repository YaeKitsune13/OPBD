package service

import (
	"example/project/backend/dto"
	"example/project/backend/repository"
	"fmt"
	"time"
)

type AnalyticsService interface {
	GetSummary() (*dto.AnalyticsSummaryDTO, error)
	GetRevenueReport() (*dto.RevenueReportDTO, error)
}

type analyticsService struct {
	visitRepo     repository.VisitRepository
	doctorRepo    repository.DoctorRepository
	inventoryRepo repository.InventoryRepository
}

func NewAnalyticsService(
	vr repository.VisitRepository,
	dr repository.DoctorRepository,
	ir repository.InventoryRepository,
) AnalyticsService {
	return &analyticsService{
		visitRepo:     vr,
		doctorRepo:    dr,
		inventoryRepo: ir,
	}
}

func (s *analyticsService) GetSummary() (*dto.AnalyticsSummaryDTO, error) {
	// 1. Определяем временные рамки
	now := time.Now()
	// Начало и конец текущего месяца
	currentStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	currentEnd := currentStart.AddDate(0, 1, 0)

	// Начало и конец прошлого месяца (для сравнения)
	prevStart := currentStart.AddDate(0, -1, 0)
	prevEnd := currentStart

	// 2. Получаем визиты за оба периода
	currentVisits, err := s.visitRepo.GetByPeriod(currentStart, currentEnd)
	if err != nil {
		return nil, err
	}

	prevVisits, err := s.visitRepo.GetByPeriod(prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	// 3. Считаем основные метрики текущего месяца
	var currentRevenue float64
	uniqueClients := make(map[int64]bool)
	doctorVisitCount := make(map[int64]int64)

	for _, v := range currentVisits {
		currentRevenue += v.TotalCost
		uniqueClients[v.Appointment.PetID] = true
		doctorVisitCount[v.Appointment.DoctorID]++
	}

	// 4. Считаем выручку прошлого месяца для вычисления изменения (%)
	var prevRevenue float64
	for _, v := range prevVisits {
		prevRevenue += v.TotalCost
	}

	revenueChange := 0.0
	if prevRevenue > 0 {
		revenueChange = ((currentRevenue - prevRevenue) / prevRevenue) * 100
	}

	// 5. Вычисляем средний чек
	avgCheck := 0.0
	if len(currentVisits) > 0 {
		avgCheck = currentRevenue / float64(len(currentVisits))
	}

	// 6. Формируем список нагрузки на врачей
	doctors, _ := s.doctorRepo.GetAll()
	var docLoads []dto.DoctorLoad
	for _, d := range doctors {
		count := doctorVisitCount[d.DoctorID]
		status := "Low"
		if count > 20 {
			status = "High"
		} else if count > 10 {
			status = "Medium"
		}

		// Инициал имени (первая буква)
		firstInitial := ""
		if len(d.User.FirstName) > 0 {
			firstInitial = string([]rune(d.User.FirstName)[0]) // Используем rune для корректной работы с кириллицей
		}

		docLoads = append(docLoads, dto.DoctorLoad{
			// ИСПРАВЛЕНИЕ: обращаемся через d.User
			Name:       fmt.Sprintf("%s %s.", d.User.LastName, firstInitial),
			VisitCount: count,
			LoadStatus: status,
		})
	}

	// 7. Собираем итоговый DTO
	summary := &dto.AnalyticsSummaryDTO{
		MonthlyVisits: int64(len(currentVisits)),
		TotalRevenue:  fmt.Sprintf("%.0f ₽", currentRevenue),
		RevenueChange: fmt.Sprintf("%+.1f%%", revenueChange),
		ActiveClients: int64(len(uniqueClients)),
		AvgCheck:      int64(avgCheck),
		DoctorLoad:    docLoads,
	}

	// 8. Получаем популярные услуги динамически
	popular, err := s.visitRepo.GetPopularServices(currentStart, currentEnd)
	if err == nil {
		summary.PopularServices = popular
	} else {
		summary.PopularServices = []dto.PopularServices{}
	}

	return summary, nil
}
func (s *analyticsService) GetRevenueReport() (*dto.RevenueReportDTO, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)

	visits, err := s.visitRepo.GetByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	report := &dto.RevenueReportDTO{}
	dailyMap := make(map[string]*dto.DailyRows)

	for _, v := range visits {
		dateStr := v.VisitDate.Format("02.01")
		report.PeriodTotal += int64(v.TotalCost)

		if _, ok := dailyMap[dateStr]; !ok {
			dailyMap[dateStr] = &dto.DailyRows{Date: dateStr}
		}
		// Здесь можно добавить детальный расчет по услугам/медикаментам
		dailyMap[dateStr].Total = fmt.Sprintf("%.0f", v.TotalCost)
	}

	for _, row := range dailyMap {
		report.DailyRows = append(report.DailyRows, *row)
	}

	return report, nil
}
