package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type OwnerRepository interface {
	Create(owner *models.Owner) error
	GetByID(id int64) (*models.Owner, error)
	GetByEmail(email string) (*models.Owner, error)
	Update(owner *models.Owner) error
}

type ownerRepository struct {
	db *gorm.DB
}

func NewOwnerRepository(db *gorm.DB) OwnerRepository {
	return &ownerRepository{db: db}
}

func (r *ownerRepository) Create(owner *models.Owner) error {
	return r.db.Create(owner).Error
}

func (r *ownerRepository) GetByID(id int64) (*models.Owner, error) {
	var owner models.Owner
	err := r.db.First(&owner, id).Error
	return &owner, err
}

func (r *ownerRepository) GetByEmail(email string) (*models.Owner, error) {
	var owner models.Owner
	err := r.db.Where("email = ?", email).First(&owner).Error
	return &owner, err
}

func (r *ownerRepository) Update(owner *models.Owner) error {
	return r.db.Save(owner).Error
}
