package repository

import (
	"api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetProfileData(id uint) (*models.User, int64, int64, error)
	Update(user *models.User) error
	GetByPhone(phone string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepo) GetProfileData(id uint) (*models.User, int64, int64, error) {
	var user models.User
	var pCount, vCount int64

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, 0, 0, err
	}

	r.db.Model(&models.Pet{}).Where("owner_id = ?", id).Count(&pCount)
	r.db.Model(&models.Appointment{}).Where("client_id = ? AND status = ?", id, "done").Count(&vCount)

	return &user, pCount, vCount, nil
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}
