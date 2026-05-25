package database

import (
	"api/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	user := "root"
	pass := "password"
	host := "127.0.0.1"
	port := "3306"
	dbName := "vet_clinic"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.Pet{}, &models.Service{}, &models.Appointment{},
		&models.MedicalProtocol{}, &models.Product{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{})

	return db
}
