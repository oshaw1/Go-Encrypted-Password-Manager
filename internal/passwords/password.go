package passwords

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
)

type Password struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Link              string `json:"hyperlink"`
	EncryptedPassword string `json:"encrypted_password"`
}

type PasswordManager struct {
	MasterPasswordHash string     `json:"master_password_hash"`
	Salt               []byte     `json:"salt"`
	Passwords          []Password `json:"passwords"`
}

func StorePassword(title, link, password, masterPassword string, pathToPasswordFile string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	manager, err := ReadPasswordManager(pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to read password manager: %w", err)
	}

	encryptionKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, manager.Salt)

	encryptedPassword, err := encryption.EncryptWithAES(password, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	entry := Password{
		ID:                uuid.New().String(),
		Title:             title,
		Link:              link,
		EncryptedPassword: encryptedPassword,
	}

	manager.Passwords = append(manager.Passwords, entry)

	err = writePasswordManager(manager, pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func RetrievePassword(id, masterPassword string, pathToPasswordFile string) (string, error) {
	manager, err := ReadPasswordManager(pathToPasswordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password manager: %w", err)
	}

	if !VerifyMasterPasswordIsHashedPassword(masterPassword, manager.MasterPasswordHash) {
		return "", fmt.Errorf("invalid master password")
	}

	encryptionKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, manager.Salt)

	for _, password := range manager.Passwords {
		if password.ID == id {
			decryptedPassword, err := encryption.DecryptWithAES(password.EncryptedPassword, encryptionKey)
			if err != nil {
				return "", fmt.Errorf("failed to decrypt password: %w", err)
			}
			return decryptedPassword, nil
		}
	}

	return "", fmt.Errorf("password not found")
}

func DeletePasswordByID(id, masterPassword string, pathToPasswordFile string) error {
	manager, err := ReadPasswordManager(pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to read password manager: %w", err)
	}

	encryptionKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, manager.Salt)

	foundIndex := -1
	for i, password := range manager.Passwords {
		if password.ID == id {
			decryptedPassword, err := encryption.DecryptWithAES(password.EncryptedPassword, encryptionKey)
			if err != nil {
				return fmt.Errorf("failed to decrypt password: %w", err)
			}

			fmt.Printf("Deleting password for title: %s, link: %s, password: %s\n", password.Title, password.Link, decryptedPassword)
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return fmt.Errorf("password not found")
	}

	manager.Passwords = append(manager.Passwords[:foundIndex], manager.Passwords[foundIndex+1:]...)

	err = writePasswordManager(manager, pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func CheckPasswordFileExistsInDataDirectory(dataDir string) bool {
	_, err := os.Stat(dataDir)
	return !os.IsNotExist(err)
}
