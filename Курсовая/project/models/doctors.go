package models

type Doctor struct {
	DoctorID     int64  `json:"doctor_id" gorm:"primaryKey;autoIncrement"`
	LastName     string `json:"last_name" gorm:"type:varchar(50);not null"`
	FirstName    string `json:"first_name" gorm:"type:varchar(50);not null"`
	MiddleName   string `json:"middle_name" gorm:"type:varchar(50)"`
	Speciality   string `json:"speciality" gorm:"type:varchar(50);not null"`
	Phone        string `json:"phone" gorm:"type:varchar(15)"`
	Email        string `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string `json:"password_hash" gorm:"type:varchar(255);not null"`
}
