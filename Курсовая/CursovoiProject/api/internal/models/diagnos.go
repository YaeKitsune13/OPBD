package models

import "gorm.io/gorm"

type Diagnosis struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}
