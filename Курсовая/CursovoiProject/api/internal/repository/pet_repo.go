package repository

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type PetRepository interface {
	Create(pet *models.Pet) error
	GetByOwnerID(ownerID uint) ([]models.Pet, error)
	GetByID(id uint) (*models.Pet, error)
	Update(pet *models.Pet) error
	Delete(id uint) error
}

type petRepo struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) PetRepository {
	return &petRepo{db}
}

func (r *petRepo) Create(pet *models.Pet) error {
	return r.db.Create(pet).Error
}

func (r *petRepo) GetByOwnerID(ownerID uint) ([]models.Pet, error) {
	pets := make([]models.Pet, 0)
	err := r.db.Where("owner_id = ?", ownerID).Find(&pets).Error
	return pets, err
}

func (r *petRepo) GetByID(id uint) (*models.Pet, error) {
	var pet models.Pet
	err := r.db.First(&pet, id).Error
	return &pet, err
}

func (r *petRepo) Update(pet *models.Pet) error {
	return r.db.Model(&models.Pet{}).
		Where("id = ?", pet.ID).
		Updates(map[string]interface{}{
			"name":       pet.Name,
			"species":    pet.Species,
			"breed":      pet.Breed,
			"birth_date": pet.BirthDate,
			"weight":     pet.Weight,
			"avatar":     pet.Avatar,
		}).Error
}

func (r *petRepo) Delete(id uint) error {
	return r.db.Delete(&models.Pet{}, id).Error
}
