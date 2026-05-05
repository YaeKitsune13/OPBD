package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetAllServices() ([]models.Service, error)
	GetAllMedications() ([]models.Medication, error)
	GetMedicationByID(id int64) (models.Medication, error)
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
