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
	doctorRepo repository.DoctorRepository
	userRepo   repository.UserRepository
}

func NewAuthService(or repository.UserRepository, dr repository.DoctorRepository) AuthService {
	return &authService{userRepo: or, doctorRepo: dr}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Ищем в userRepository (бывший userRepo)
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// 2. Проверка пароля (bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	// 3. Генерируем токен (используем ID из таблицы User)
	tokenString, err := s.generateToken(user.UserID, string(user.Role))

	response := &dto.AuthResponse{
		Token:    tokenString,
		Role:     user.Role,
		UserName: user.FirstName,
		UserID:   user.UserID, // По умолчанию это ID клиента
	}

	// 4. Если это врач, подменяем UserID на DoctorID для удобства фронтенда
	if user.Role == models.RoleDoctor {
		doctor, err := s.doctorRepo.GetByUserID(user.UserID)
		if err == nil {
			response.UserID = doctor.DoctorID // Теперь фронт получит ID врача для расписания
		}
	}

	return response, nil
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Проверяем, не занят ли Email
	// Используем := так как это первое объявление err
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("пользователь с таким email уже зарегистрирован")
	}

	// 2. Проверяем, не занят ли Телефон
	// Используем просто = так как переменная err уже существует выше
	_, err = s.userRepo.GetByPhone(req.Phone)
	if err == nil {
		return nil, errors.New("этот номер телефона уже используется")
	}

	// 3. Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("ошибка при обработке пароля")
	}

	// 4. Создаем модель
	newOwner := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleClient, // Обязательно задаем роль в модель!
	}

	// 5. Сохраняем в базу
	if err := s.userRepo.Create(newOwner); err != nil {
		return nil, err
	}

	// 6. Генерируем токен для нового пользователя
	tokenString, err := s.generateToken(newOwner.UserID, string(newOwner.Role))
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	return &dto.AuthResponse{
		Token:    tokenString,
		Role:     newOwner.Role,
		UserName: newOwner.FirstName,
		UserID:   newOwner.UserID,
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
