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
	GetAllUsers() ([]models.User, error)
	GetByIDAndRole(id int64, role models.UserRole) (*models.User, error)
	UpdateRole(id int64, role models.UserRole) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) UpdateRole(id int64, role models.UserRole) error {
	// Обновляем только поле role для конкретного пользователя
	return r.db.Model(&models.User{}).Where("user_id = ?", id).Update("role", role).Error
}
func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Create(owner *models.User) error {
	return r.db.Create(owner).Error
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var owner models.User
	err := r.db.First(&owner, id).Error
	return &owner, err
}

func (r *userRepository) GetByIDAndRole(id int64, role models.UserRole) (*models.User, error) {
	var user models.User
	// Используем .Where для фильтрации по обоим полям
	err := r.db.Where("user_id = ? AND role = ?", id, role).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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
