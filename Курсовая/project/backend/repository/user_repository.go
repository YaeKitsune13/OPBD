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
	// Добавляем метод в интерфейс
	Delete(id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Реализация удаления
func (r *userRepository) Delete(id int64) error {
	// Используем .Unscoped(), если нужно удалить запись физически из базы (даже если есть поле DeletedAt)
	// Либо просто .Delete(), если используете Soft Delete от GORM
	return r.db.Where("user_id = ?", id).Delete(&models.User{}).Error
}

func (r *userRepository) UpdateRole(id int64, role models.UserRole) error {
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
	// Если первичный ключ в модели назван user_id, GORM поймет это через теги.
	// Если нет, лучше использовать .Where("user_id = ?", id)
	err := r.db.Where("user_id = ?", id).First(&owner).Error
	return &owner, err
}

func (r *userRepository) GetByIDAndRole(id int64, role models.UserRole) (*models.User, error) {
	var user models.User
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
