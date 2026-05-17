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
	UpdateService(id int64, s *models.Service) error
	DeleteService(id int64) error

	CreateMedication(m *models.Medication) error
	UpdateMedication(id int64, m *models.Medication) error
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

func (r *inventoryRepository) CreateService(s *models.Service) error {
	return r.db.Create(s).Error
}

func (r *inventoryRepository) UpdateService(id int64, s *models.Service) error {
	return r.db.Model(&models.Service{}).
		Where("service_id = ?", id).
		Updates(map[string]interface{}{
			"name":        s.Name,
			"description": s.Description,
			"price":       s.Price,
		}).Error
}

func (r *inventoryRepository) DeleteService(id int64) error {
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *inventoryRepository) CreateMedication(m *models.Medication) error {
	return r.db.Create(m).Error
}

func (r *inventoryRepository) UpdateMedication(id int64, m *models.Medication) error {
	return r.db.Model(&models.Medication{}).
		Where("medication_id = ?", id).
		Updates(map[string]interface{}{
			"name":           m.Name,
			"description":    m.Description,
			"price_per_unit": m.PricePerUnit,
			"expiry_date":    m.ExpiryDate,
		}).Error
}

func (r *inventoryRepository) DeleteMedication(id int64) error {
	return r.db.Delete(&models.Medication{}, id).Error
}
