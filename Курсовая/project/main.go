package main

import "fmt"

type data struct {
	ID int8 `json:"id"`
}

func main() {
	fmt.Println("Hello, World!")

	user := data{ID: 3}

	fmt.Printf("User ID: %d\n", user.ID)
	printMessage()
}

func printMessage() {
	fmt.Println("HElfods")
}
