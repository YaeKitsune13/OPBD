package main

import (
	"fmt"
	"log"
)

func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
	fmt.Println("База подключена и миграция выполнена успешно!", db)
}
func printMessage() {
	fmt.Println("HElfods")
}
