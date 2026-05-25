package repository

import (
	"api/internal/dto"
	"api/internal/models"

	"gorm.io/gorm"
)

type StatsRepository interface {
	GetWeightHistory(petID uint) ([]dto.WeightDataDTO, error)
}

type statsRepo struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepo{db}
}

func (r *statsRepo) GetWeightHistory(petID uint) ([]dto.WeightDataDTO, error) {
	results := make([]dto.WeightDataDTO, 0)

	err := r.db.Model(&models.MedicalProtocol{}).
		Select("appointments.scheduled_at as date, medical_protocols.weight_at_visit as weight").
		Joins("join appointments on appointments.id = medical_protocols.appointment_id").
		Where("appointments.pet_id = ? AND appointments.status = ? AND appointments.deleted_at IS NULL", petID, "done").
		Order("appointments.scheduled_at asc").
		Scan(&results).Error

	return results, err
}
