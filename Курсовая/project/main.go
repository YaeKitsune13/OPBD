package main

import "fmt"

type data struct {
	ID int8 `json:"id"` // Тег пишется в обратных кавычках
}

func main() {
	fmt.Println("Hello, World!")

	// Используем := для объявления и соблюдаем регистр поля ID
	user := data{ID: 3}

	fmt.Printf("User ID: %d\n", user.ID)
	printMessage() // Переименовано, чтобы не путать со встроенными именами
}

func printMessage() {
	fmt.Println("HElfods")
}
