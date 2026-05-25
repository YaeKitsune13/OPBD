package models

import "gorm.io/gorm"

type UserRole string

const (
	UserClient UserRole = "client"
	UserDoctor UserRole = "doctor"
)

type User struct {
	gorm.Model
	Email          string   `gorm:"type:varchar(191);uniqueIndex" json:"email"`
	Password       string   `json:"-"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Phone          string   `gorm:"type:varchar(191);unique" json:"phone"`
	Role           UserRole `gorm:"type:enum('client','doctor')" json:"role"`
	Specialization string   `json:"specialization,omitempty"`
	Pets           []Pet    `gorm:"foreignKey:OwnerID" json:"pets,omitempty"`
}
