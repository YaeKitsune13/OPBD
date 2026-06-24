package repository

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type DiagnosisRepository interface {
	GetAll() ([]models.Diagnosis, error)
	SetDiagnoses(appointment *models.Appointment, diagnoses []models.Diagnosis) error
}

type diagnosisRepo struct {
	db *gorm.DB
}

func NewDiagnosisRepository(db *gorm.DB) DiagnosisRepository {
	return &diagnosisRepo{
		db: db,
	}
}

func (r *diagnosisRepo) GetAll() ([]models.Diagnosis, error) {
	diagnoses := make([]models.Diagnosis, 0)

	err := r.db.Find(&diagnoses).Error

	return diagnoses, err
}

func (r *diagnosisRepo) SetDiagnoses(
	appointment *models.Appointment,
	diagnoses []models.Diagnosis,
) error {
	return r.db.Model(appointment).
		Association("Diagnoses").
		Replace(diagnoses)
}
