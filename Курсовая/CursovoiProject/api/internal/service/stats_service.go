package service

import (
	"api/internal/dto"
	"api/internal/repository"
	"time"
)

type StatsService interface {
	GetPetWeightStats(petID uint) ([]dto.WeightDataDTO, error)
}

type statsService struct {
	repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository) StatsService {
	return &statsService{repo}
}

func (s *statsService) GetPetWeightStats(petID uint) ([]dto.WeightDataDTO, error) {
	rawStats, err := s.repo.GetWeightHistory(petID)
	if err != nil {
		return make([]dto.WeightDataDTO, 0), nil
	}

	if rawStats == nil {
		return make([]dto.WeightDataDTO, 0), nil
	}

	for i, item := range rawStats {
		t, err := time.Parse(time.RFC3339, item.Date)
		if err != nil {
			t, _ = time.Parse("2006-01-02 15:04:05", item.Date)
		}

		if !t.IsZero() {
			rawStats[i].Date = t.Format("02.01")
		}
	}

	return rawStats, nil
}
