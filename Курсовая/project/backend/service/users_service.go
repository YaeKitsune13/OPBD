package service

import (
	"errors"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
)

type UsersService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByIdAndRole(id int64, role models.UserRole) (*models.User, error)
	UpdateUserRole(id int64, role models.UserRole) error
}

type usersService struct {
	userRepo repository.UserRepository
}

func NewUsersService(repo repository.UserRepository) UsersService {
	return &usersService{
		userRepo: repo,
	}
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
