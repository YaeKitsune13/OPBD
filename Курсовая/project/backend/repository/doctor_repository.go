package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type DoctorRepository interface {
	Create(doctor *models.Doctor) error
	GetByID(id int64) (*models.Doctor, error)
	GetByEmail(email string) (*models.Doctor, error)
	GetAll() ([]models.Doctor, error)
	GetBySpecialty(specialty string) ([]models.Doctor, error)
}

type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

func (r *doctorRepository) Create(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *doctorRepository) GetByID(id int64) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.First(&doctor, id).Error
	return &doctor, err
}

func (r *doctorRepository) GetByEmail(email string) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Where("email = ?", email).First(&doctor).Error
	return &doctor, err
}

func (r *doctorRepository) GetAll() ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := r.db.Find(&doctors).Error
	return doctors, err
}
func (r *doctorRepository) GetBySpecialty(specialty string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := r.db.Where("speciality = ?", specialty).Find(&doctors).Error
	return doctors, err
}
