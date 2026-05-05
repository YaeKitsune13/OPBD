package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type VisitRepository interface {
	Create(visit *models.Visit) error
	GetByID(id int64) (*models.Visit, error)
	GetByPetID(petId int64) ([]models.Visit, error)
}

type visitRepository struct {
	db *gorm.DB
}

func NewVisitRepository(db *gorm.DB) VisitRepository {
	return &visitRepository{db: db}
}

func (r *visitRepository) Create(visit *models.Visit) error {
	return r.db.Create(visit).Error
}

func (r *visitRepository) GetByID(id int64) (*models.Visit, error) {
	var visit models.Visit
	err := r.db.Preload("Appointment").First(&visit, id).Error
	return &visit, err
}

func (r *visitRepository) GetByPetID(petId int64) ([]models.Visit, error) {
	var visits []models.Visit

	err := r.db.
		Joins("JOIN appointments ON appointments.appointment_id = visits.appointment_id").
		Where("appointments.pet_id = ?", petId).
		Order("visits.visit_date DESC").
		Find(&visits).Error

	return visits, err
}
