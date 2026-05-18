package service

import (
	"example/project/backend/dto"
	"example/project/backend/repository"
	"time"
)

type SearchService interface {
	SearchPatients(query string) ([]dto.PatientSearchResultDTO, error)
}

type searchService struct {
	petRepo  repository.PetRepository
	userRepo repository.UserRepository
}

func NewSearchService(petRepo repository.PetRepository, userRepo repository.UserRepository) SearchService {
	return &searchService{petRepo: petRepo, userRepo: userRepo}
}

func (s *searchService) SearchPatients(query string) ([]dto.PatientSearchResultDTO, error) {
	var ownerIDSet = make(map[int64]bool)

	if query == "" {
		// Загружаем всех питомцев
		pets, err := s.petRepo.GetAllPets()
		if err != nil {
			return nil, err
		}
		for _, p := range pets {
			ownerIDSet[p.OwnerID] = true
		}
	} else {
		pets, err := s.petRepo.SearchByName(query)
		if err != nil {
			return nil, err
		}
		owners, err := s.userRepo.SearchByName(query)
		if err != nil {
			return nil, err
		}
		for _, p := range pets {
			ownerIDSet[p.OwnerID] = true
		}
		for _, u := range owners {
			ownerIDSet[u.UserID] = true
		}
	}

	var result []dto.PatientSearchResultDTO
	for ownerID := range ownerIDSet {
		ownerPets, err := s.petRepo.GetByOwnerID(ownerID)
		if err != nil || len(ownerPets) == 0 {
			continue
		}

		owner, err := s.userRepo.GetByID(ownerID)
		if err != nil {
			continue
		}

		for _, pet := range ownerPets {
			age := int(time.Since(pet.BirthDate).Hours() / 24 / 365)
			result = append(result, dto.PatientSearchResultDTO{
				PetID:      pet.PetID,
				PetName:    pet.Nickname, // было pet.Name
				Species:    pet.Species,
				Breed:      pet.Breed,
				Weight:     pet.CurrentWeight, // было pet.Weight
				Age:        age,               // было pet.Age
				OwnerID:    owner.UserID,
				OwnerName:  owner.LastName + " " + owner.FirstName + " " + owner.MiddleName,
				OwnerPhone: owner.Phone,
			})
		}
	}

	return result, nil
}
