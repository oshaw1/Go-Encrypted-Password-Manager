package passwords

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func StorePassword(title, link, password, masterPassword string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	manager, err := ReadPasswordManager()
	if err != nil {
		return fmt.Errorf("failed to read password manager: %w", err)
	}

	encryptionKey := encryption.DeriveEncryptionKey(masterPassword, manager.Salt)

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

	err = writePasswordManager(manager)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func RetrievePassword(id, masterPassword string) (string, error) {
	manager, err := ReadPasswordManager()
	if err != nil {
		return "", fmt.Errorf("failed to read password manager: %w", err)
	}

	if !verifyMasterPassword(masterPassword, manager.MasterPasswordHash) {
		return "", fmt.Errorf("invalid master password")
	}

	encryptionKey := encryption.DeriveEncryptionKey(masterPassword, manager.Salt)

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

func InitializePasswordManager(masterPassword string) error {
	salt, err := encryption.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	passwordHash := hashMasterPassword(masterPassword)

	manager := PasswordManager{
		MasterPasswordHash: passwordHash,
		Salt:               salt,
		Passwords:          []Password{},
	}

	err = writePasswordManager(manager)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func ReadPasswordManager() (PasswordManager, error) {
	file, err := os.ReadFile("data/passwords.json")
	if err != nil {
		if os.IsNotExist(err) {
			return PasswordManager{}, nil
		}
		return PasswordManager{}, fmt.Errorf("failed to read password manager file: %w", err)
	}

	var manager PasswordManager
	err = json.Unmarshal(file, &manager)
	if err != nil {
		return PasswordManager{}, fmt.Errorf("failed to unmarshal password manager: %w", err)
	}

	return manager, nil
}

func writePasswordManager(manager PasswordManager) error {
	file, err := json.MarshalIndent(manager, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal password manager: %w", err)
	}

	err = os.WriteFile("data/passwords.json", file, 0600)
	if err != nil {
		return fmt.Errorf("failed to write password manager file: %w", err)
	}

	return nil
}

func hashMasterPassword(masterPassword string) string {
	hash := sha256.Sum256([]byte(masterPassword))
	return fmt.Sprintf("%x", hash)
}

func verifyMasterPassword(masterPassword, hashedPassword string) bool {
	return hashMasterPassword(masterPassword) == hashedPassword
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
