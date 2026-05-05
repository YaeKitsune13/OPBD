package service

import (
	"example/project/backend/dto"
	"example/project/backend/models"
	"example/project/backend/repository"
	"strings"
	"time"
)

type PetService interface {
	GetOwnerPets(ownerId int64) ([]dto.PetCardDTO, error)
	GetPetDetails(petId int64) (*dto.PetCardDTO, error)
	AddPet(ownerId int64, data dto.PetCardDTO) error
	UpdateWeight(petId int64, newWeight float64, doctorId *int64) error
}

type petService struct {
	petRepo   repository.PetRepository
	ownerRepo repository.OwnerRepository
}

func NewPetService(pr repository.PetRepository, or repository.OwnerRepository) PetService {
	return &petService{
		petRepo:   pr,
		ownerRepo: or,
	}
}

// 1. Получаем список всех карточек владельца
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
			Dob:     p.BirthDate,
			Weight:  p.CurrentWeight,
			Avatar:  getEmojiAvatar(p.Species),
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
		Dob:     p.BirthDate,
		Weight:  p.CurrentWeight,
		Avatar:  getEmojiAvatar(p.Species),
	}, nil
}

// 3. Добавление нового питомца
func (s *petService) AddPet(ownerId int64, data dto.PetCardDTO) error {
	if _, err := s.ownerRepo.GetByID(ownerId); err != nil {
		return err
	}

	newPet := &models.Pet{
		OwnerID:       ownerId,
		Nickname:      data.Name,
		Species:       data.Species,
		Breed:         data.Breed,
		BirthDate:     data.Dob,
		CurrentWeight: data.Weight,
	}

	return s.petRepo.Create(newPet)
}

// 4. Обновление веса (Двойное действие: в профиль и в историю)
func (s *petService) UpdateWeight(petId int64, newWeight float64, doctorId *int64) error {
	// 1. Находим питомца
	pet, err := s.petRepo.GetByID(petId)
	if err != nil {
		return err
	}

	// 2. Обновляем текущий вес в основной модели
	pet.CurrentWeight = newWeight
	if err := s.petRepo.Update(pet); err != nil {
		return err
	}

	// 3. Создаем запись в историю веса (для графика)
	historyRecord := &models.WeightHistory{
		PetID:      petId,
		Weight:     newWeight,
		MeasuredAt: time.Now(),
		DoctorID:   doctorId, // Может быть null, если взвесил владелец дома
	}

	return s.petRepo.AddWeightRecord(historyRecord)
}

// Вспомогательная логика
func getEmojiAvatar(species string) string {
	s := strings.ToLower(species)
	switch {
	case strings.Contains(s, "кот") || strings.Contains(s, "кош"):
		return "🐱"
	case strings.Contains(s, "соб") || strings.Contains(s, "пес"):
		return "🐶"
	case strings.Contains(s, "крол"):
		return "🐇"
	default:
		return "🐾"
	}
}
