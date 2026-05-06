package service

import (
	"errors"
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
}

type authService struct {
	ownerRepo  repository.OwnerRepository
	doctorRepo repository.DoctorRepository
}

func NewAuthService(or repository.OwnerRepository, dr repository.DoctorRepository) AuthService {
	return &authService{ownerRepo: or, doctorRepo: dr}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Проверка на статического Админа (логин: admin, пароль: admin)
	if req.Email == "admin" && req.Password == "admin" {
		return &dto.AuthResponse{
			Token:    "admin-token-secret",
			Role:     "admin",
			UserName: "Администратор",
			UserID:   0,
		}, nil
	}

	// 2. Ищем в таблице Докторов
	doctor, err := s.doctorRepo.GetByEmail(req.Email)
	if err == nil { // Если нашли доктора
		if doctor.PasswordHash == req.Password {
			return &dto.AuthResponse{
				Token:    "doctor-token-secret",
				Role:     "doctor",
				UserName: doctor.FirstName,
				UserID:   doctor.DoctorID,
			}, nil
		}
	}

	// 3. Ищем в таблице Владельцев (Клиентов)
	owner, err := s.ownerRepo.GetByEmail(req.Email)
	if err == nil { // Если нашли владельца
		if owner.PasswordHash == req.Password {
			return &dto.AuthResponse{
				Token:    "client-token-secret",
				Role:     "client",
				UserName: owner.FirstName,
				UserID:   owner.OwnerID,
			}, nil
		}
	}

	return nil, errors.New("неверный email или пароль")
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// По умолчанию регистрация через общую форму создает КЛИЕНТА (Owner)
	// Докторов обычно создает Админ в своей панели

	newOwner := &models.Owner{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: req.Password, // В идеале bcrypt.GenerateFromPassword
	}

	if err := s.ownerRepo.Create(newOwner); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:    "new-client-token",
		Role:     "client",
		UserName: newOwner.FirstName,
		UserID:   newOwner.OwnerID,
	}, nil
}
