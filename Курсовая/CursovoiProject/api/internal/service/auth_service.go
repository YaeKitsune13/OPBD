package service

import (
	"api/internal/dto"
	"api/internal/models"
	"api/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("super-secret-key")

type AuthService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(req dto.RegisterRequest) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := models.User{
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserClient,
	}
	_, err := s.repo.GetByEmail(user.Email)
	if err == nil {
		return errors.New("Данная почта уже зарегистрирована")
	}
	_, err = s.repo.GetByPhone(user.Phone)
	if err == nil {
		return errors.New("Данный номер телефона уже зарегистрирован")
	}
	return s.repo.Create(&user)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("неверный пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	return &dto.LoginResponse{
		Token:    tokenString,
		UserId:   user.ID,
		Role:     string(user.Role),
		UserName: user.FirstName,
		LastName: user.LastName,
		Phone:    user.Phone,
	}, nil
}
