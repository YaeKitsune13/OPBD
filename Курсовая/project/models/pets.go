package models

import "time"

type Pet struct {
	PetID         int64     `json:"pet_id" gorm:"primaryKey;autoIncrement"`
	OwnerID       int64     `json:"owner_id" gorm:"not null"`
	Nickname      string    `json:"nickname" gorm:"type:varchar(30);not null"`
	Species       string    `json:"species" gorm:"type:varchar(50);not null"`
	Breed         string    `json:"breed" gorm:"type:varchar(100)"`
	BirthDate     time.Time `json:"birth_date" gorm:"type:date;not null"`
	Photo         string    `json:"photo" gorm:"type:varchar(255)"`
	CurrentWeight float64   `json:"current_weight" gorm:"type:decimal(5,2);check:current_weight > 0.1 AND current_weight < 300"`
}
