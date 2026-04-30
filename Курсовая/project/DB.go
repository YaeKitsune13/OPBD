package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель вашей таблицы
type User struct {
	gorm.Model
	Name string
}

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=1234 dbname=cursavoi_project port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	return db, nil
}
