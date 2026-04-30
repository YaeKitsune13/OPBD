package main

import (
	"example/project/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=1234 dbname=cursavoi_project port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Appointment{},       // 1
		&models.Doctor{},            // 2
		&models.Medication{},        // 3
		&models.Owner{},             // 4
		&models.Pet{},               // 5
		&models.Service{},           // 6
		&models.Visit{},             // 7
		&models.VisitPrescription{}, // 8
		&models.WeightHistory{},     // 9
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	return db, nil
}
