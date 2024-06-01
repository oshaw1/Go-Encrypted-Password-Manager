package passwords

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
)

type Password struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Link           string `json:"hyperlink"`
	HashedPassword string `json:"hashed_password"`
}

func StorePassword(title, link, password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	hashedPassword, err := encryption.EncryptWithSHA256(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	entry := Password{
		ID:             uuid.New().String(),
		Title:          title,
		Link:           link,
		HashedPassword: hashedPassword,
	}

	passwords, err := ReadPasswords()
	if err != nil {
		return fmt.Errorf("failed to read passwords: %w", err)
	}

	passwords = append(passwords, entry)

	err = WritePasswords(passwords)
	if err != nil {
		return fmt.Errorf("failed to write passwords: %w", err)
	}

	return nil
}

func WritePasswords(passwords []Password) error {
	filePath := "data/passwords.json"

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open passwords.json: %w", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(passwords, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal passwords: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write passwords to file: %w", err)
	}

	return nil
}

func ReadPasswords() ([]Password, error) {
	filePath := "data/passwords.json"

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Password{}, nil
		}
		return nil, fmt.Errorf("failed to read passwords file: %w", err)
	}

	var passwords []Password
	err = json.Unmarshal(fileData, &passwords)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal passwords: %w", err)
	}

	return passwords, nil
}

func CheckPasswordFileExists() bool {
	dataDir := "data"
	filePath := filepath.Join(dataDir, "passwords.json")
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func CreateEmptyPasswordFile() error {
	dataDir := "data"
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dataDir, "passwords.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode([]Password{})
	if err != nil {
		return err
	}

	return nil
}
