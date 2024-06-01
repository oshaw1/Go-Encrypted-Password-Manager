package main

import (
	"fmt"
	"os"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

func main() {
	masterPassword := "master-password"
	pathToPasswordFile := "data/passwords.json"

	if !passwords.CheckPasswordFileExistsInDataDirectory(pathToPasswordFile) {
		err := passwords.InitializePasswordManager(masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Println("Failed to initialize password data:", err)
			return
		}
		fmt.Println("Password data initialized successfully")
	}

	err := passwords.StorePassword("Example Title", "https://example.com", "password123", masterPassword, pathToPasswordFile)
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	passwordID := "e27acad0-f57f-4d5f-a431-7d2263f026e5"
	retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, pathToPasswordFile)
	if err != nil {
		fmt.Printf("Failed to retrieve password: %v\n", err)
		os.Exit(1)
	}

	passwordID = "03acf7ae-6076-42c1-94e9-d7114bcf1be0"
	err = passwords.DeletePassword(passwordID, masterPassword, pathToPasswordFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Retrieved password: %s\n", retrievedPassword)
}
