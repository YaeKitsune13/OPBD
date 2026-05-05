package service

import (
	"example/project/backend/dto"
	"example/project/backend/repository"
	"time"
)

type InventoryService interface {
	GetServices() ([]dto.ServiceItemDTO, error)
	GetMedications() ([]dto.MedicationItemDTO, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) GetServices() ([]dto.ServiceItemDTO, error) {
	items, err := s.repo.GetAllServices()
	if err != nil {
		return nil, err
	}

	var result []dto.ServiceItemDTO
	for _, item := range items {
		result = append(result, dto.ServiceItemDTO{
			Id:    item.Code,
			Name:  item.Name,
			Desc:  item.Description,
			Price: int64(item.Price),
		})
	}
	return result, nil
}

func (s *inventoryService) GetMedications() ([]dto.MedicationItemDTO, error) {
	items, err := s.repo.GetAllMedications()
	if err != nil {
		return nil, err
	}

	var result []dto.MedicationItemDTO
	now := time.Now()

	for _, item := range items {
		status := dto.MedicationOk
		if item.ExpiryDate.Before(now) {
			status = dto.MedicationExpired
		}

		result = append(result, dto.MedicationItemDTO{
			Id:     item.Code,
			Name:   item.Name,
			Desc:   item.Description,
			Price:  int64(item.PricePerUnit),
			Expiry: item.ExpiryDate.Format("2006-01"),
			Status: status,
		})
	}
	return result, nil
}
