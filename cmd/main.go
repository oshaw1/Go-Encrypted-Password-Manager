package main

import (
	"fmt"
	"os"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func main() {
	masterPassword := "master-password"
	dataDirectory := "data/passwords.json"

	if !passwords.CheckPasswordFileExistsInDataDirectory(dataDirectory) {
		err := passwords.InitializePasswordManager(masterPassword, dataDirectory)
		if err != nil {
			fmt.Println("Failed to initialize password data:", err)
			return
		}
		fmt.Println("Password data initialized successfully")
	}

	err := passwords.StorePassword("Example Title", "https://example.com", "password123", masterPassword, dataDirectory)
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	passwordID := "e27acad0-f57f-4d5f-a431-7d2263f026e5"
	retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, dataDirectory)
	if err != nil {
		fmt.Printf("Failed to retrieve password: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved password: %s\n", retrievedPassword)
}
