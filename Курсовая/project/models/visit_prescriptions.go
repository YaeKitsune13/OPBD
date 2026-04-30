package models

type ItemType string

const (
	TypeService    ItemType = "service"
	TypeMedication ItemType = "medication"
)

type VisitPrescription struct {
	PrescriptionID int64       `json:"prescription_id" gorm:"primaryKey;autoIncrement"`
	VisitID        int64       `json:"visit_id" gorm:"not null"`
	Visit          Visit       `json:"-" gorm:"foreignKey:VisitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemType       ItemType    `json:"item_type" gorm:"type:varchar(20);not null;check:item_type IN ('service', 'medication')"`
	ServiceID      *int64      `json:"service_id" gorm:"default:null"`
	Service        *Service    `json:"-" gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MedicationID   *int64      `json:"medication_id" gorm:"default:null"`
	Medication     *Medication `json:"-" gorm:"foreignKey:MedicationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity       int64       `json:"quantity" gorm:"not null;default:1"`
	UnitPrice      float64     `json:"unit_price" gorm:"type:decimal(10,2);not null"`
}
