package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type PetRepository interface {
	Create(pet *models.Pet) error
	GetByID(id int64) (*models.Pet, error)
	GetByOwnerID(ownerID int64) ([]models.Pet, error)
	Update(pet *models.Pet) error
	AddWeightRecord(record *models.WeightHistory) error
}

type petRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) PetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(pet *models.Pet) error {
	return r.db.Create(pet).Error
}

func (r *petRepository) GetByID(id int64) (*models.Pet, error) {
	var pet models.Pet
	err := r.db.First(&pet, id).Error
	return &pet, err
}

func (r *petRepository) GetByOwnerID(ownerID int64) ([]models.Pet, error) {
	var pets []models.Pet
	err := r.db.Where("owner_id = ?", ownerID).Find(&pets).Error
	return pets, err
}

func (r *petRepository) Update(pet *models.Pet) error {
	return r.db.Save(pet).Error
}

func (r *petRepository) AddWeightRecord(record *models.WeightHistory) error {
	return r.db.Create(record).Error
}
