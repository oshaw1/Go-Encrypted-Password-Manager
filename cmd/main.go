package main

import (
	"fmt"
	"os"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func main() {
	masterPassword := "master-password"

	if !passwords.CheckPasswordFileExists() {
		err := passwords.InitializePasswordManager(masterPassword)
		if err != nil {
			fmt.Println("Failed to initialize password data:", err)
			return
		}
		fmt.Println("Password data initialized successfully")
	}

	err := passwords.StorePassword("Example Title", "https://example.com", "password123", masterPassword)
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	err = passwords.StorePassword("", "", "password456", masterPassword)
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	passwordID := "16659de2-c2c5-4715-8fb7-c664eb77021a"
	retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword)
	if err != nil {
		fmt.Printf("Failed to retrieve password: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved password: %s\n", retrievedPassword)
}
