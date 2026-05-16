package service

import (
	"errors"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UsersService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByIdAndRole(id int64, role models.UserRole) (*models.User, error)
	UpdateUserRole(id int64, role models.UserRole) error
	DeleteUser(id int64) error
	CreateUser(user *models.User) error // Понадобится для создания врача
	UpdateUser(id int64, lastName, firstName, middleName, phone string) error
	ChangePassword(id int64, currentPassword, newPassword string) error
}
type usersService struct {
	userRepo repository.UserRepository
}

func NewUsersService(repo repository.UserRepository) UsersService {
	return &usersService{
		userRepo: repo,
	}
}
func (s *usersService) DeleteUser(id int64) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("service.DeleteUser: %w", err)
	}
	return nil
}

func (s *usersService) CreateUser(user *models.User) error {
	return s.userRepo.Create(user)
}
func (s *usersService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("service.GetAllUsers: %w", err)
	}
	return users, nil
}

func (s *usersService) GetUserByIdAndRole(id int64, role models.UserRole) (*models.User, error) {
	user, err := s.userRepo.GetByIDAndRole(id, role)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserByIdAndRole (id: %d): %w", id, err)
	}

	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	return user, nil
}

func (s *usersService) UpdateUserRole(id int64, role models.UserRole) error {
	return s.userRepo.UpdateRole(id, role)
}

func (s *usersService) UpdateUser(id int64, lastName, firstName, middleName, phone string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	user.LastName = lastName
	user.FirstName = firstName
	user.MiddleName = middleName
	user.Phone = phone

	return s.userRepo.Update(user)
}

func (s *usersService) ChangePassword(id int64, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	// Проверяем текущий пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("неверный текущий пароль")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.userRepo.Update(user)
}
