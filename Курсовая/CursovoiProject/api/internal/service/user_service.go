package service

import (
	"api/internal/dto"
	"api/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(id uint) (*dto.UserProfileResponse, error)
	UpdateProfile(id uint, req dto.UpdateProfileRequest) error
	ChangePassword(id uint, req dto.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetProfile(id uint) (*dto.UserProfileResponse, error) {
	user, pCount, vCount, err := s.repo.GetProfileData(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID: user.ID, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, Phone: user.Phone, Role: string(user.Role),
		PetsCount: int(pCount), VisitsCount: int(vCount),
	}, nil
}

func (s *userService) UpdateProfile(id uint, req dto.UpdateProfileRequest) error {
	user, _ := s.repo.GetByEmail("")
	user, _, _, _ = s.repo.GetProfileData(id)

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	return s.repo.Update(user)
}

func (s *userService) ChangePassword(id uint, req dto.ChangePasswordRequest) error {
	user, _, _, _ := s.repo.GetProfileData(id)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Current)); err != nil {
		return errors.New("неверный текущий пароль")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Next), 12)
	user.Password = string(hashed)
	return s.repo.Update(user)
}
