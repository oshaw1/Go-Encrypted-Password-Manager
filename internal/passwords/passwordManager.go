package passwords

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
)

func InitializePasswordManager(masterPassword string, pathToPasswordFile string) error {
	salt, err := encryption.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	passwordHash := encryption.HashMasterPassword(masterPassword)

	manager := PasswordManager{
		MasterPasswordHash: passwordHash,
		Salt:               salt,
		Passwords:          []Password{},
	}

	err = writePasswordManager(manager, pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func ReadPasswordManager(pathToPasswordFile string) (PasswordManager, error) {
	file, err := os.ReadFile(pathToPasswordFile)
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

func writePasswordManager(manager PasswordManager, pathToPasswordFile string) error {
	file, err := json.MarshalIndent(manager, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal password manager: %w", err)
	}

	dir := filepath.Dir(pathToPasswordFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	err = os.WriteFile(pathToPasswordFile, file, 0600)
	if err != nil {
		return fmt.Errorf("failed to write password manager file: %w", err)
	}

	return nil
}

func VerifyMasterPasswordIsHashedPassword(masterPassword, hashedPassword string) bool {
	return encryption.HashMasterPassword(masterPassword) == hashedPassword
}
