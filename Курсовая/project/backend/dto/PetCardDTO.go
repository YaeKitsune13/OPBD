package dto

import "time"

type PetCardDTO struct {
	PetId   int64     `json:"petId" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Species string    `json:"species"`
	Breed   string    `json:"breed"`
	Dob     time.Time `json:"dob"`
	Weight  float64   `json:"weight"`
	Avatar  string    `json:"avatar"`
}
