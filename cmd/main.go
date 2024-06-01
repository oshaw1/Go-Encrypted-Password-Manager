package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
)

type PasswordData struct {
	MasterPasswordHash string `json:"master_password_hash"`
	Salt               string `json:"salt"`
	Passwords          []struct {
		ID                string `json:"id"`
		Title             string `json:"title"`
		Hyperlink         string `json:"hyperlink"`
		EncryptedPassword string `json:"encrypted_password"`
	} `json:"passwords"`
}

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

	err := passwords.StorePassword("Example Title", "https://example.com", "pnewhdsfusdf", masterPassword, pathToPasswordFile)
	if err != nil {
		fmt.Println("Failed to store password:", err)
	} else {
		fmt.Println("Password stored successfully")
	}

	data, err := os.ReadFile(pathToPasswordFile)
	if err != nil {
		fmt.Printf("Failed to read password file: %v\n", err)
		os.Exit(1)
	}

	var passwordData PasswordData
	err = json.Unmarshal(data, &passwordData)
	if err != nil {
		fmt.Printf("Failed to parse JSON data: %v\n", err)
		os.Exit(1)
	}

	for _, password := range passwordData.Passwords {
		passwordTitle := password.Title
		passwordID := password.ID
		retrievedPassword, err := passwords.RetrievePassword(passwordID, masterPassword, pathToPasswordFile)
		if err != nil {
			fmt.Printf("Failed to retrieve password for ID %s: %v\n", passwordID, err)
			continue
		}
		fmt.Printf("Password for %s: %s\n", passwordTitle, retrievedPassword)
	}

	// passwordID = "03acf7ae-6076-42c1-94e9-d7114bcf1be0"
	// err = passwords.DeletePasswordByID(passwordID, masterPassword, pathToPasswordFile)
	// if err != nil {
	// 	fmt.Println("Error deleting password:", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Retrieved password: %s\n", retrievedPassword)
}
