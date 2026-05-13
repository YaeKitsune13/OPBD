package repository

import (
	"example/project/backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(owner *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(owner *models.User) error
	GetByPhone(phone string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewOwnerRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(owner *models.User) error {
	return r.db.Create(owner).Error
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var owner models.User
	err := r.db.First(&owner, id).Error
	return &owner, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var owner models.User
	err := r.db.Where("email = ?", email).First(&owner).Error
	return &owner, err
}

func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	var owner models.User
	err := r.db.Where("phone = ?", phone).First(&owner).Error
	return &owner, err
}

func (r *userRepository) Update(owner *models.User) error {
	return r.db.Save(owner).Error
}
