package models

type Doctor struct {
	DoctorID   int64  `json:"doctor_id" gorm:"primaryKey;autoIncrement;column:doctor_id"`
	UserID     int64  `json:"user_id" gorm:"column:user_id;unique;not null"`
	User       User   `json:"user" gorm:"foreignKey:UserID;references:UserID"`
	Speciality string `json:"speciality" gorm:"type:varchar(50);not null"`
}
