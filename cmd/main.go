package main

import (
	"fmt"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func main() {
	if !passwords.CheckPasswordFileExists() {
		err := passwords.CreateEmptyPasswordFile()
		if err != nil {
			fmt.Println("Failed to create passwords file:", err)
			return
		}
		fmt.Println("Passwords file created successfully")
	}

	err := passwords.StorePassword("Example Title", "https://example.com", "password123")
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	err = passwords.StorePassword("", "", "password123")
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}
}
