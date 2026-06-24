package repository

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	GetSchedule(doctorID uint) ([]models.Appointment, error)
	CompleteAppointment(appID uint, protocol *models.MedicalProtocol) error
	SearchClients(query string) ([]models.User, error)
	GetClientHistory(clientID uint) ([]models.Pet, error)
	UpdateAppointmentStatus(appID uint, status string) error
	SetDiagnoses(appID uint, diagnosisIDs []uint) error
}

type doctorRepo struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db}
}

func (r *doctorRepo) GetSchedule(doctorID uint) ([]models.Appointment, error) { // Добавлен error
	apps := make([]models.Appointment, 0)
	err := r.db.Preload("Pet").Preload("Client").Preload("Service").
		Where("doctor_id = ? AND status != ?", doctorID, "rejected").
		Order("scheduled_at asc").Find(&apps).Error

	return apps, err
}

func (r *doctorRepo) CompleteAppointment(appID uint, protocol *models.MedicalProtocol) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(protocol).Error; err != nil {
			return err
		}
		return tx.Model(&models.Appointment{}).Where("id = ?", appID).Update("status", "done").Error
	})
}

func (r *doctorRepo) SearchClients(query string) ([]models.User, error) {
	users := make([]models.User, 0)
	dbQuery := r.db.Preload("Pets").Where("role = ?", "client")
	if query != "" {
		q := "%" + query + "%"
		dbQuery = dbQuery.Where("first_name LIKE ? OR last_name LIKE ? OR phone LIKE ?", q, q, q)
	}
	err := dbQuery.Find(&users).Error
	return users, err
}

func (r *doctorRepo) GetClientHistory(clientID uint) ([]models.Pet, error) {
	pets := make([]models.Pet, 0)
	err := r.db.Preload("Appointments", "status = ?", "done").
		Preload("Appointments.Protocol").
		Preload("Appointments.Diagnoses").
		Where("owner_id = ?", clientID).Find(&pets).Error
	return pets, err
}

func (r *doctorRepo) UpdateAppointmentStatus(appID uint, status string) error {
	return r.db.Model(&models.Appointment{}).Where("id = ?", appID).Update("status", status).Error
}

func (r *doctorRepo) SetDiagnoses(appID uint, diagnosisIDs []uint) error {
	var appointment models.Appointment

	if err := r.db.First(&appointment, appID).Error; err != nil {
		return err
	}

	var diagnoses []models.Diagnosis

	if err := r.db.Where("id IN ?", diagnosisIDs).
		Find(&diagnoses).Error; err != nil {
		return err
	}

	return r.db.Model(&appointment).
		Association("Diagnoses").
		Replace(diagnoses)
}
