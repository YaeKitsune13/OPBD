package main

import (
	"fmt"
	"log"

	"example/project/backend"
)

func main() {
	db, err := backend.ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	fmt.Println("База подключена и миграция выполнена успешно!", db)
}
func printMessage() {
	fmt.Println("HElfods")
}
