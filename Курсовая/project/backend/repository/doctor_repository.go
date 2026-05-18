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
	GetByUserID(user_id int64) (*models.Doctor, error)
}

type doctorRepository struct {
	db *gorm.DB
}

func (r *doctorRepository) GetByID(id int64) (*models.Doctor, error) {
	var doctor models.Doctor
	// Обязательно подгружаем данные пользователя (имя, фамилию)
	err := r.db.Preload("User").First(&doctor, id).Error
	return &doctor, err
}

func (r *doctorRepository) GetByUserID(userID int64) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&doctor).Error
	return &doctor, err
}

func (r *doctorRepository) GetAll() ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := r.db.Preload("User").Find(&doctors).Error
	return doctors, err
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{db: db}
}

func (r *doctorRepository) Create(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *doctorRepository) GetByEmail(email string) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Where("email = ?", email).First(&doctor).Error
	return &doctor, err
}

func (r *doctorRepository) GetBySpecialty(specialty string) ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := r.db.Where("speciality = ?", specialty).Preload("User").Find(&doctors).Error
	return doctors, err
}
