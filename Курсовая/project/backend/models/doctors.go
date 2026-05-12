package models

type Doctor struct {
	DoctorID   int64  `json:"doctor_id" gorm:"primaryKey;autoIncrement"`
	UserID     int64  `json:"user_id" gorm:"unique;not null"` // Связь с таблицей Owner/User
	LastName   string `json:"last_name" gorm:"type:varchar(50);not null"`
	FirstName  string `json:"first_name" gorm:"type:varchar(50);not null"`
	MiddleName string `json:"middle_name" gorm:"type:varchar(50)"`
	Speciality string `json:"speciality" gorm:"type:varchar(50);not null"`
}
