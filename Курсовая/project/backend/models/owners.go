package models

type Owner struct {
	OwnerID      int64  `json:"owner_id" gorm:"primaryKey;autoIncrement"`
	LastName     string `json:"last_name" gorm:"type:varchar(50);not null"`
	FirstName    string `json:"first_name" gorm:"type:varchar(50);not null"`
	MiddleName   string `json:"middle_name" gorm:"type:varchar(50)"`
	Address      string `json:"address" gorm:"type:varchar(255)"`
	Phone        string `json:"phone" gorm:"type:varchar(15);unique;not null"`
	Email        string `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string `json:"password" gorm:"type:varchar(255);not null"`
}
