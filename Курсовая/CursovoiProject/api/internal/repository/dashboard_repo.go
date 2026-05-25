package repository

import (
	"api/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetSummary(userID uint) (int64, int64, int64, error)
	GetNextAppointment(userID uint) (*models.Appointment, error)
	GetRecentAppointments(userID uint, limit int) ([]models.Appointment, error)
}

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepo{db}
}

func (r *dashboardRepo) GetSummary(userID uint) (int64, int64, int64, error) {
	var pCount, aCount, vCount int64

	r.db.Model(&models.Pet{}).Where("owner_id = ?", userID).Count(&pCount)

	r.db.Model(&models.Appointment{}).
		Where("client_id = ? AND status IN ?", userID, []string{"waiting", "confirmed"}).
		Count(&aCount)

	r.db.Model(&models.Appointment{}).
		Where("client_id = ? AND status = ?", userID, "done").
		Count(&vCount)

	return pCount, aCount, vCount, nil
}

func (r *dashboardRepo) GetNextAppointment(userID uint) (*models.Appointment, error) {
	var app models.Appointment
	err := r.db.Preload("Pet").
		Where("client_id = ? AND scheduled_at > ? AND status = ?", userID, time.Now(), "confirmed").
		Order("scheduled_at asc").
		Limit(1).
		Find(&app).Error

	return &app, err
}

func (r *dashboardRepo) GetRecentAppointments(userID uint, limit int) ([]models.Appointment, error) {
	apps := make([]models.Appointment, 0)

	err := r.db.Preload("Pet").
		Preload("Doctor").
		Preload("Service").
		Where("client_id = ?", userID).
		Order("scheduled_at desc").
		Limit(limit).
		Find(&apps).Error

	return apps, err
}
