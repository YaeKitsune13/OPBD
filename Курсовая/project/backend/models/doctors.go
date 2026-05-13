package models

type Doctor struct {
	DoctorID   int64  `json:"doctor_id" gorm:"primaryKey;autoIncrement"`
	UserID     int64  `json:"user_id" gorm:"unique;not null"` // Ссылка на UserID
	Speciality string `json:"speciality" gorm:"type:varchar(50);not null"`
	User       User   `json:"user" gorm:"foreignKey:UserID"`
}
