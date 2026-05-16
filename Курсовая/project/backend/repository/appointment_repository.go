package repository

import (
	"example/project/backend/models"
	"time"

	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(appointment *models.Appointment) error
	GetUpcomingByOwnerId(ownerId int64) ([]models.Appointment, error)
	GetTodaySchedule(doctorId int64) ([]models.Appointment, error)
	UpdateStatus(appointmentId int64, status models.Status) error
	IsSlotTaken(doctorId int64, scheduledAt time.Time) (bool, error)
	GetByDoctorAndDateRange(doctorID int64, start, end time.Time) ([]models.Appointment, error)
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

// Реализация проверки занятости слота
func (r *appointmentRepository) IsSlotTaken(doctorId int64, scheduledAt time.Time) (bool, error) {
	var count int64
	// Проверяем наличие записей к конкретному врачу на конкретное время.
	// Исключаем отмененные записи (StatusCancelled), так как они освобождают время.
	err := r.db.Model(&models.Appointment{}).
		Where("doctor_id = ? AND scheduled_at = ? AND status != 'cancelled'",
			doctorId, scheduledAt).
		Count(&count).Error

	return count > 0, err
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) GetUpcomingByOwnerId(ownerId int64) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := r.db.Joins("JOIN pets ON pets.pet_id = appointments.pet_id").
		Where("pets.owner_id = ? AND appointments.scheduled_at >= ?", ownerId, time.Now()).
		Order("appointments.scheduled_at ASC").
		Find(&appointments).Error

	return appointments, err
}

func (r *appointmentRepository) GetTodaySchedule(doctorId int64) ([]models.Appointment, error) {
	var appointments []models.Appointment

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("doctor_id = ? AND scheduled_at BETWEEN ? AND ?", doctorId, startOfDay, endOfDay).
		Order("scheduled_at ASC").
		Find(&appointments).Error

	return appointments, err
}

func (r *appointmentRepository) UpdateStatus(appointmentId int64, status models.Status) error {
	return r.db.Model(&models.Appointment{}).
		Where("appointment_id = ?", appointmentId).
		Update("status", status).Error
}

func (r *appointmentRepository) GetByDoctorAndDateRange(doctorID int64, start, end time.Time) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("doctor_id = ? AND scheduled_at BETWEEN ? AND ?", doctorID, start, end).
		Find(&appointments).Error
	return appointments, err
}
