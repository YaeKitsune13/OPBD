package service

import (
	"errors"
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	// 1. Ищем пользователя в базе по Email
	user, err := s.ownerRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// 2. Проверяем пароль
	// Внимание: для реальных проектов пароли сравнивают через bcrypt.CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Если пароли не совпали, вернется ошибка
		return nil, errors.New("неверный email или пароль")
	}

	tokenString, err := s.generateToken(user.OwnerID, string(user.Role))
	if err != nil {
		return nil, errors.New("ошибка создания токена")
	}

	// 5. ФОРМИРОВАНИЕ ОТВЕТА
	response := &dto.AuthResponse{
		Token:    tokenString, // Теперь здесь реальная строка "eyJhbG..."
		Role:     user.Role,
		UserName: user.FirstName,
		UserID:   user.OwnerID, // По умолчанию ID владельца
	}

	// 6. СПЕЦИФИЧНАЯ ЛОГИКА ДЛЯ ВРАЧЕЙ
	// Если вошел врач, нам нужно вернуть в фронтенд его DoctorID,
	// чтобы фронт мог запрашивать расписание именно этого врача.
	if user.Role == models.RoleDoctor {
		doctor, err := s.doctorRepo.GetByUserID(user.OwnerID)
		if err == nil {
			response.UserID = doctor.DoctorID
		}
	}

	return response, nil
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// По умолчанию регистрация через общую форму создает КЛИЕНТА (Owner)
	// Докторов обычно создает Админ в своей панели
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("ошибка при обработке пароля")
	}
	newOwner := &models.Owner{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword), // Сохраняем хэш-строку
	}

	if err := s.ownerRepo.Create(newOwner); err != nil {
		return nil, err
	}
	tokenString, _ := s.generateToken(newOwner.OwnerID, "client")
	return &dto.AuthResponse{
		Token:    tokenString,
		Role:     models.RoleClient,
		UserName: newOwner.FirstName,
		UserID:   newOwner.OwnerID,
	}, nil
}

// Вспомогательная функция (пишется внизу файла auth_service.go)
func (s *authService) generateToken(userID int64, role string) (string, error) {
	claims := models.MyCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenObj.SignedString(models.GetJWTSecret())
}
