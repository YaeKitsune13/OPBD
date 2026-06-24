package repository

import (
	"api/internal/models"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(app *models.Appointment) error
	GetBusySlots(doctorID uint, date string) ([]string, error)
	GetByClientID(clientID uint) ([]models.Appointment, error)
	GetDoctors() ([]models.User, error)
	GetServices() ([]models.Service, error)
	SetDiagnoses(appointmentID uint, diagnosisIDs []uint) error
}

type appointmentRepo struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepo{db}
}

func (r *appointmentRepo) Create(app *models.Appointment) error {
	return r.db.Create(app).Error
}

func (r *appointmentRepo) GetBusySlots(doctorID uint, date string) ([]string, error) {
	times := make([]time.Time, 0)
	err := r.db.Model(&models.Appointment{}).
		Where("doctor_id = ? AND DATE(scheduled_at) = ? AND status != ?", doctorID, date, "rejected").
		Pluck("scheduled_at", &times).Error

	slots := make([]string, 0)
	for _, t := range times {
		slots = append(slots, t.Format("15:04"))
	}
	return slots, err
}

func (r *appointmentRepo) GetByClientID(clientID uint) ([]models.Appointment, error) {
	apps := []models.Appointment{}

	err := r.db.Preload("Pet").
		Preload("Protocol").
		Preload("Doctor").
		Preload("Service").
		Preload("Diagnoses").
		Where("client_id = ?", clientID).
		Order("scheduled_at desc").
		Find(&apps).Error

	return apps, err
}

func (r *appointmentRepo) GetDoctors() ([]models.User, error) {
	doctors := make([]models.User, 0)
	err := r.db.Where("role = ?", "doctor").Find(&doctors).Error
	return doctors, err
}

func (r *appointmentRepo) GetServices() ([]models.Service, error) {
	services := make([]models.Service, 0)
	err := r.db.Find(&services).Error
	return services, err
}

func (r *appointmentRepo) SetDiagnoses(appointmentID uint, diagnosisIDs []uint) error {
	var appointment models.Appointment

	if err := r.db.First(&appointment, appointmentID).Error; err != nil {
		return err
	}

	var diagnoses []models.Diagnosis

	if err := r.db.Where("id IN ?", diagnosisIDs).Find(&diagnoses).Error; err != nil {
		return err
	}

	return r.db.Model(&appointment).
		Association("Diagnoses").
		Replace(diagnoses)
}
