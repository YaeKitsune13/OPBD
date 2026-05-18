package repository

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"time"

	"gorm.io/gorm"
)

type VisitRepository interface {
	Create(visit *models.Visit) error
	GetByID(id int64) (*models.Visit, error)
	GetByPetID(petId int64) ([]models.Visit, error)
	GetByPeriod(start, end time.Time) ([]models.Visit, error)
	GetPopularServices(start, end time.Time) ([]dto.PopularServices, error)
	AddPrescription(p *models.VisitPrescription) error
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

func (r *visitRepository) GetByPeriod(start, end time.Time) ([]models.Visit, error) {
	var visits []models.Visit
	err := r.db.Preload("Appointment").
		Where("visit_date BETWEEN ? AND ?", start, end).
		Find(&visits).Error
	return visits, err
}

func (r *visitRepository) GetPopularServices(start, end time.Time) ([]dto.PopularServices, error) {
	var result []dto.PopularServices
	err := r.db.Table("visit_prescriptions").
		Select("services.name as name, count(*) as count, sum(visit_prescriptions.unit_price * visit_prescriptions.quantity) as revenue").
		Joins("join services on services.service_id = visit_prescriptions.service_id").
		Joins("join visits on visits.visit_id = visit_prescriptions.visit_id").
		Where("visits.visit_date BETWEEN ? AND ?", start, end).
		Group("services.name").
		Order("count DESC").
		Limit(5).
		Scan(&result).Error
	return result, err
}

func (r *visitRepository) AddPrescription(p *models.VisitPrescription) error {
	return r.db.Create(p).Error
}
