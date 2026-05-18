package models

import "time"

type Visit struct {
	VisitId int64 `json:"visit_id" gorm:"primaryKey;autoIncrement"`

	// Указатель позволяет полю быть NULL в базе данных
	AppointmentID *int64       `json:"appointment_id" gorm:"index"`
	Appointment   *Appointment `json:"appointment" gorm:"foreignKey:AppointmentID"`

	// Прямая связь с питомцем (теперь она главная для истории болезни)
	PetID int64 `json:"pet_id" gorm:"not null"`
	Pet   Pet   `json:"pet" gorm:"foreignKey:PetID"`

	Anamnesis string    `json:"anamnesis" gorm:"type:text;not null"`
	Diagnosis string    `json:"diagnosis" gorm:"type:varchar(500);not null"`
	TotalCost float64   `json:"total_cost" gorm:"type:decimal(10,2);not null;default:0"`
	VisitDate time.Time `json:"visit_date" gorm:"autoCreateTime"`

	Analysis        *string `json:"analysis" gorm:"type:text"`
	Recommendations *string `json:"recommendations" gorm:"type:text"`
}
