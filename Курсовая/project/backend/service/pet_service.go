package service

import (
	"errors"
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PetService interface {
	GetOwnerPets(ownerId int64) ([]dto.PetCardDTO, error)
	GetPetDetails(petId int64) (*dto.PetCardDTO, error)
	AddPet(ownerId int64, data dto.PetCardDTO) error
	UpdateWeight(petId int64, newWeight float64, doctorId *int64) error
	UpdatePet(data dto.PetCardDTO) error
	DeletePet(petId int64) error
	GetWeightChartData(petId int64) ([]dto.WeightPointDTO, error)
	QuickRegister(req dto.QuickRegisterRequest) (*dto.QuickRegisterResponse, error)
}

type petService struct {
	petRepo  repository.PetRepository
	userRepo repository.UserRepository
}

func NewPetService(pr repository.PetRepository, or repository.UserRepository) PetService {
	return &petService{
		petRepo:  pr,
		userRepo: or,
	}
}

func (s *petService) GetOwnerPets(ownerId int64) ([]dto.PetCardDTO, error) {
	pets, err := s.petRepo.GetByOwnerID(ownerId)
	if err != nil {
		return nil, err
	}

	var result []dto.PetCardDTO
	for _, p := range pets {
		result = append(result, dto.PetCardDTO{
			PetId:   p.PetID,
			Name:    p.Nickname,
			Species: p.Species,
			Breed:   p.Breed,
			// ИСПРАВЛЕНИЕ ТУТ: конвертируем time.Time в string
			Dob:    p.BirthDate.Format("2006-01-02"),
			Weight: p.CurrentWeight,
			Avatar: getEmojiAvatar(p.Species),
		})
	}
	return result, nil
}

// 2. Детальная информация об одном питомце
func (s *petService) GetPetDetails(petId int64) (*dto.PetCardDTO, error) {
	p, err := s.petRepo.GetByID(petId)
	if err != nil {
		return nil, err
	}

	return &dto.PetCardDTO{
		PetId:   p.PetID,
		Name:    p.Nickname,
		Species: p.Species,
		Breed:   p.Breed,
		// ИСПРАВЛЕНИЕ ТУТ: конвертируем time.Time в string
		Dob:    p.BirthDate.Format("2006-01-02"),
		Weight: p.CurrentWeight,
		Avatar: getEmojiAvatar(p.Species),
	}, nil
}

func (s *petService) AddPet(ownerId int64, data dto.PetCardDTO) error {
	// 1. Парсим дату из строки (теперь это снова нужно)
	dobTime, err := time.Parse("2006-01-02", data.Dob)
	if err != nil {
		return errors.New("неверный формат даты (ожидается ГГГГ-ММ-ДД)")
	}

	// 2. Валидация
	if dobTime.After(time.Now()) {
		return errors.New("дата рождения не может быть в будущем")
	}
	if data.Weight < 0.1 || data.Weight > 300.0 {
		return errors.New("вес должен быть в диапазоне от 0.1 до 300 кг")
	}

	// 3. Создание модели (мапим DTO на Model)
	newPet := &models.Pet{
		OwnerID:       ownerId,
		Nickname:      data.Name,
		Species:       data.Species,
		Breed:         data.Breed,
		BirthDate:     dobTime, // Передаем спарсенное время
		CurrentWeight: data.Weight,
		Photo:         getEmojiAvatar(data.Species),
	}

	if err := s.petRepo.Create(newPet); err != nil {
		return err
	}

	return s.UpdateWeight(newPet.PetID, data.Weight, nil)
}

// 4. Обновление веса (Двойное действие: в профиль и в историю)
func (s *petService) UpdateWeight(petId int64, newWeight float64, doctorId *int64) error {
	// --- ВАЛИДАЦИЯ ВЕСА ---
	if newWeight < 0.1 || newWeight > 300.0 {
		return errors.New("некорректный вес (0.1 - 300 кг)")
	}

	pet, err := s.petRepo.GetByID(petId)
	if err != nil {
		return err
	}

	pet.CurrentWeight = newWeight
	if err := s.petRepo.Update(pet); err != nil {
		return err
	}

	historyRecord := &models.WeightHistory{
		PetID:      petId,
		Weight:     newWeight,
		MeasuredAt: time.Now(),
		DoctorID:   doctorId,
	}

	return s.petRepo.AddWeightRecord(historyRecord)
}

// 5. Обновление основных данных
func (s *petService) UpdatePet(data dto.PetCardDTO) error {
	// --- ВАЛИДАЦИЯ ИМЕНИ ---
	if len(data.Name) < 1 || len(data.Name) > 30 {
		return errors.New("кличка должна быть от 1 до 30 символов")
	}

	pet, err := s.petRepo.GetByID(data.PetId)
	if err != nil {
		return err
	}

	dobTime, err := time.Parse("2006-01-02", data.Dob)
	if err != nil {
		return errors.New("неверный формат даты рождения (ГГГГ-ММ-ДД)")
	}

	pet.Nickname = data.Name
	pet.Species = data.Species
	pet.Breed = data.Breed
	pet.BirthDate = dobTime // Теперь присваиваем объект времени, а не строку

	return s.petRepo.Update(pet)
}

// 6. Удаление питомца
func (s *petService) DeletePet(petId int64) error {
	return s.petRepo.Delete(petId)
}

// 7. Данные для графика веса
func (s *petService) GetWeightChartData(petId int64) ([]dto.WeightPointDTO, error) {
	history, err := s.petRepo.GetWeightHistory(petId)
	if err != nil {
		return nil, err
	}

	var chartData []dto.WeightPointDTO
	for _, h := range history {
		point := dto.WeightPointDTO{
			Label: h.MeasuredAt.Format("Jan"), // Формат месяца для графика
			Value: h.Weight,
			Date:  h.MeasuredAt.Format("02.01.2006"),
		}
		point.DoctorName = "Ветклиника"
		chartData = append(chartData, point)
	}
	return chartData, nil
}

// Вспомогательная логика аватарок
func getEmojiAvatar(species string) string {
	s := strings.ToLower(species)
	switch {
	case strings.Contains(s, "кот") || strings.Contains(s, "кош") || strings.Contains(s, "cat"):
		return "🐱"
	case strings.Contains(s, "соб") || strings.Contains(s, "пес") || strings.Contains(s, "dog"):
		return "🐶"
	case strings.Contains(s, "крол") || strings.Contains(s, "rabbit"):
		return "🐇"
	case strings.Contains(s, "птиц") || strings.Contains(s, "попуг") || strings.Contains(s, "bird"):
		return "🦜"
	default:
		return "🐾"
	}
}

func (s *petService) QuickRegister(req dto.QuickRegisterRequest) (*dto.QuickRegisterResponse, error) {
	var owner *models.User
	var err error

	if req.IsAnonymous {
		// Логика Анонима: ищем системного пользователя "Гость"
		owner, err = s.userRepo.GetByPhone("0000") // Системный номер
		if err != nil {
			// Если гостя нет в базе — создаем один раз
			pass, _ := bcrypt.GenerateFromPassword([]byte("guest_pass"), 10)
			owner = &models.User{
				FirstName: "Анонимный", LastName: "Владелец",
				Phone: "0000", Email: "anonymous@clinic.local",
				PasswordHash: string(pass), Role: models.RoleClient,
			}
			s.userRepo.Create(owner)
		}
		req.PetName = "Пациент " + time.Now().Format("15:04") // Генерим имя, если пусто
	} else {
		// Логика обычного быстрого клиента
		owner, err = s.userRepo.GetByPhone(req.OwnerPhone)
		if err != nil {
			// Создаем нового владельца
			names := strings.Split(req.OwnerName, " ")
			fName := names[0]
			lName := "Клиент"
			if len(names) > 1 {
				lName = names[1]
			}

			pass, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)
			owner = &models.User{
				FirstName: fName, LastName: lName,
				Phone:        req.OwnerPhone,
				Email:        fmt.Sprintf("user_%s@clinic.temp", req.OwnerPhone),
				PasswordHash: string(pass), Role: models.RoleClient,
			}
			if err := s.userRepo.Create(owner); err != nil {
				return nil, err
			}
		}
	}

	// Создаем питомца
	newPet := &models.Pet{
		OwnerID:       owner.UserID,
		Nickname:      req.PetName,
		Species:       req.Species,
		Breed:         req.Breed,
		BirthDate:     time.Now(),
		CurrentWeight: 0.1, // ФИКС: обходим ограничение базы (вес > 0)
		Photo:         getEmojiAvatar(req.Species),
	}

	if err := s.petRepo.Create(newPet); err != nil {
		return nil, err
	}

	return &dto.QuickRegisterResponse{
		PetID:     newPet.PetID,
		PetName:   newPet.Nickname,
		OwnerName: owner.FirstName + " " + owner.LastName,
		Breed:     newPet.Breed,
	}, nil
}
