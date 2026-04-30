package models

type ItemType string

const (
	TypeService    ItemType = "service"
	TypeMedication ItemType = "medication"
)

type VisitPrescription struct {
	PrescriptionID int64    `json:"prescription_id" gorm:"primaryKey;autoIncrement"`
	VisitID        int64    `json:"visit_id" gorm:"not null"`
	ItemType       ItemType `json:"item_type" gorm:"type:varchar(20);not null;check:item_type IN ('service', 'medication')"`
	ServiceID      *int64   `json:"service_id" gorm:"default:null"`
	MedicationID   *int64   `json:"medication_id" gorm:"default:null"`
	Quantity       int64    `json:"quantity" gorm:"not null;default:1"`
	UnitPrice      float64  `json:"unit_price" gorm:"type:decimal(10,2);not null"`
}
