package models

import "time"

type MedicationStatus string

const (
	MedicationOK      MedicationStatus = "ok"
	MedicationExpired MedicationStatus = "expired"
)

type Medication struct {
	MedicationID int64            `json:"medication_id" gorm:"primaryKey;autoIncrement"`
	Code         string           `json:"code" gorm:"type:varchar(10);unique;not null"`
	Name         string           `json:"name" gorm:"type:varchar(150);not null"`
	Description  string           `json:"description" gorm:"type:text"`
	PricePerUnit float64          `json:"price_per_unit" gorm:"type:decimal(10,2);not null;check:price_per_unit > 0"`
	ExpiryDate   time.Time        `json:"expiry_date" gorm:"type:date;not null"`
	Status       MedicationStatus `json:"status" gorm:"type:varchar(10);default:expired;check:status IN ('ok', 'expired')"`
}
