package service

import (
	"example/project/backend/models"
	"example/project/backend/repository"
)

type DoctorService interface {
	GetBySpecialty(specialty string) ([]models.Doctor, error)
	GetAll() ([]models.Doctor, error)
	CreateDoctor(doctor *models.Doctor) error
	GetByUserID(userID int64) (*models.Doctor, error)
}

type doctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo: repo}
}

func (s *doctorService) GetBySpecialty(specialty string) ([]models.Doctor, error) {
	return s.repo.GetBySpecialty(specialty)
}

func (s *doctorService) GetAll() ([]models.Doctor, error) {
	return s.repo.GetAll()
}

func (s *doctorService) CreateDoctor(doctor *models.Doctor) error {
	return s.repo.Create(doctor)
}
func (s *doctorService) GetByUserID(userID int64) (*models.Doctor, error) {
	return s.repo.GetByUserID(userID)
}
