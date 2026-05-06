package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetAllServices() ([]models.Service, error)
	GetAllMedications() ([]models.Medication, error)
	GetMedicationByID(id int64) (models.Medication, error)

	CreateService(s *models.Service) error
	UpdateService(s *models.Service) error
	DeleteService(id int64) error

	CreateMedication(m *models.Medication) error
	UpdateMedication(m *models.Medication) error
	DeleteMedication(id int64) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetAllServices() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Find(&services).Error
	return services, err
}

func (r *inventoryRepository) GetAllMedications() ([]models.Medication, error) {
	var medications []models.Medication
	err := r.db.Find(&medications).Error
	return medications, err
}

func (r *inventoryRepository) GetMedicationByID(id int64) (models.Medication, error) {
	var medication models.Medication
	err := r.db.First(&medication, id).Error
	return medication, err
}
func (r *inventoryRepository) CreateService(s *models.Service) error { return r.db.Create(s).Error }
func (r *inventoryRepository) UpdateService(s *models.Service) error { return r.db.Save(s).Error }
func (r *inventoryRepository) DeleteService(id int64) error {
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *inventoryRepository) CreateMedication(m *models.Medication) error {
	return r.db.Create(m).Error
}
func (r *inventoryRepository) UpdateMedication(m *models.Medication) error { return r.db.Save(m).Error }
func (r *inventoryRepository) DeleteMedication(id int64) error {
	return r.db.Delete(&models.Medication{}, id).Error
}
