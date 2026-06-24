package service

import (
	"api/internal/models"
	"api/internal/repository"
)

type PetService interface {
	AddPet(pet *models.Pet) error
	GetOwnerPets(ownerID uint) ([]models.Pet, error)
	UpdatePet(pet *models.Pet) error
	RemovePet(petID uint) error
}

type petService struct {
	repo repository.PetRepository
}

func NewPetService(repo repository.PetRepository) PetService {
	return &petService{repo}
}

func (s *petService) AddPet(pet *models.Pet) error {
	return s.repo.Create(pet)
}

func (s *petService) GetOwnerPets(ownerID uint) ([]models.Pet, error) {
	return s.repo.GetByOwnerID(ownerID)
}

func (s *petService) UpdatePet(pet *models.Pet) error {
	return s.repo.Update(pet)
}

func (s *petService) RemovePet(petID uint) error {
	return s.repo.Delete(petID)
}
