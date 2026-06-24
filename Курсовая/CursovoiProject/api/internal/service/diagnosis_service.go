package service

import (
	"api/internal/dto"
	"api/internal/repository"
)

type DiagnosisService interface {
	GetAll() ([]dto.DiagnosisDTO, error)
}

type diagnosisService struct {
	repo repository.DiagnosisRepository
}

func NewDiagnosisService(repo repository.DiagnosisRepository) DiagnosisService {
	return &diagnosisService{repo}
}

func (s *diagnosisService) GetAll() ([]dto.DiagnosisDTO, error) {
	diagnoses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]dto.DiagnosisDTO, 0, len(diagnoses))
	for _, d := range diagnoses {
		res = append(res, dto.DiagnosisDTO{ID: d.ID, Name: d.Name})
	}
	return res, nil
}
