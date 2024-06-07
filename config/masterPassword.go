package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
)

type MasterPassword struct {
	MasterPasswordHash string `json:"master_password_hash"`
}

func StoreMasterPassword(masterPassword string, pathToMasterPassword string) error {
	if masterPassword == "" {
		return fmt.Errorf("password cannot be empty")
	}

	hashedMasterPassword := encryption.HashMasterPassword(masterPassword)

	entry := MasterPassword{
		MasterPasswordHash: hashedMasterPassword,
	}

	err := writeMasterPassword(entry, pathToMasterPassword)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func ReadMasterPassword(pathToMasterPassword string) (MasterPassword, error) {
	file, err := os.ReadFile(pathToMasterPassword)
	if err != nil {
		if os.IsNotExist(err) {
			return MasterPassword{}, nil
		}
		return MasterPassword{}, fmt.Errorf("failed to read password manager file: %w", err)
	}

	var manager MasterPassword
	err = json.Unmarshal(file, &manager)
	if err != nil {
		return MasterPassword{}, fmt.Errorf("failed to unmarshal password manager: %w", err)
	}

	return manager, nil
}

func writeMasterPassword(manager MasterPassword, pathToMasterPassword string) error {
	file, err := json.MarshalIndent(manager, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal password manager: %w", err)
	}

	err = os.WriteFile(pathToMasterPassword, file, 0600)
	if err != nil {
		return fmt.Errorf("failed to write password manager file: %w", err)
	}

	return nil
}
