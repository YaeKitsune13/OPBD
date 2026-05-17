package service

import (
	"errors"
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type InventoryService interface {
	GetServices() ([]dto.ServiceItemDTO, error)
	GetMedications() ([]dto.MedicationItemDTO, error)
	DeleteService(id int64) error
	CreateService(serv dto.ServiceItemDTO) error
	UpdateServices(id int64, serv dto.ServiceItemDTO) error
	CreateMedication(med dto.MedicationItemDTO) error
	UpdateMedication(id string, med dto.MedicationItemDTO) error
	DeleteMedication(id string) error
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
			Id:    item.ServiceID, // ← был потерян ID, из-за этого delete/update не знали какую запись трогать
			Name:  item.Name,
			Desc:  item.Description,
			Price: int64(item.Price),
		})
	}
	return result, nil
}

func (s *inventoryService) DeleteService(id int64) error {
	return s.repo.DeleteService(id)
}

func (s *inventoryService) CreateService(serv dto.ServiceItemDTO) error {
	if serv.Price < 0 {
		return errors.New("цена не может быть меньше 0")
	}

	service := models.Service{
		Name:        strings.TrimSpace(serv.Name),
		Description: strings.TrimSpace(serv.Desc),
		Price:       float64(serv.Price),
	}

	if err := s.repo.CreateService(&service); err != nil {
		log.Printf("ошибка при создании сервиса: %v", err)
		return fmt.Errorf("ошибка при создании сервиса: %w", err)
	}

	return nil
}

func (s *inventoryService) UpdateServices(id int64, serv dto.ServiceItemDTO) error {
	if serv.Price < 0 {
		return errors.New("цена не может быть меньше 0")
	}

	service := models.Service{
		ServiceID:   id,
		Name:        strings.TrimSpace(serv.Name), // ← добавлен TrimSpace
		Description: strings.TrimSpace(serv.Desc), // ← добавлен TrimSpace
		Price:       float64(serv.Price),
	}

	if err := s.repo.UpdateService(id, &service); err != nil {
		log.Printf("ошибка при обновлении сервиса: %v", err)
		return fmt.Errorf("ошибка при обновлении сервиса: %w", err) // ← было %W
	}

	return nil
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
			Id:     item.MedicationID, // ← было item.Code (строка)
			Name:   item.Name,
			Desc:   item.Description,
			Price:  int64(item.PricePerUnit),
			Expiry: item.ExpiryDate.Format("2006-01"),
			Status: status,
		})
	}
	return result, nil
}

func (s *inventoryService) CreateMedication(med dto.MedicationItemDTO) error {
	expiry, err := time.Parse("2006-01", med.Expiry)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}

	m := models.Medication{
		Name:         strings.TrimSpace(med.Name),
		Description:  strings.TrimSpace(med.Desc),
		PricePerUnit: float64(med.Price),
		ExpiryDate:   expiry,
		// Code убран, Status не передаём — БД проставит через default/check
	}

	if err := s.repo.CreateMedication(&m); err != nil {
		log.Printf("ошибка при создании медикамента: %v", err)
		return fmt.Errorf("ошибка при создании медикамента: %w", err)
	}
	return nil
}

func (s *inventoryService) UpdateMedication(id string, med dto.MedicationItemDTO) error {
	expiry, err := time.Parse("2006-01", med.Expiry)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный id: %w", err)
	}

	m := models.Medication{
		Name:         strings.TrimSpace(med.Name),
		Description:  strings.TrimSpace(med.Desc),
		PricePerUnit: float64(med.Price),
		ExpiryDate:   expiry,
	}

	if err := s.repo.UpdateMedication(idInt, &m); err != nil {
		log.Printf("ошибка при обновлении медикамента: %v", err)
		return fmt.Errorf("ошибка при обновлении медикамента: %w", err)
	}
	return nil
}
func (s *inventoryService) DeleteMedication(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный id: %w", err)
	}
	return s.repo.DeleteMedication(idInt)
}
