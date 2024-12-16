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
	Username          string `json:"username/account"`
	EncryptedPassword string `json:"encrypted_password"`
}

type PasswordManager struct {
	MasterPasswordHash string     `json:"master_password_hash"`
	Salt               []byte     `json:"salt"`
	Passwords          []Password `json:"passwords"`
}

func StorePassword(title, link, username, password, masterPassword string, pathToPasswordFile string) error {
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
	encryptedLink, err := encryption.EncryptWithAES(link, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	encryptedUsername, err := encryption.EncryptWithAES(username, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	entry := Password{
		ID:                uuid.New().String(),
		Title:             title,
		Link:              encryptedLink,
		Username:          encryptedUsername,
		EncryptedPassword: encryptedPassword,
	}

	manager.Passwords = append(manager.Passwords, entry)

	err = writePasswordManager(manager, pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to write password manager: %w", err)
	}

	return nil
}

func RetrievePassword(id, masterPassword string, pathToPasswordFile string) (decryptedLink string, decryptedUsername string, decryptedPassword string, err error) {
	manager, err := ReadPasswordManager(pathToPasswordFile)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read password manager: %w", err)
	}

	encryptionKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, manager.Salt)

	for _, password := range manager.Passwords {
		if password.ID == id {
			decryptedLink, err := encryption.DecryptWithAES(password.Link, encryptionKey)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to decrypt link: %w", err)
			}
			decryptedUsername, err := encryption.DecryptWithAES(password.Username, encryptionKey)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to decrypt username: %w", err)
			}
			decryptedPassword, err := encryption.DecryptWithAES(password.EncryptedPassword, encryptionKey)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to decrypt password: %w", err)
			}
			return decryptedLink, decryptedUsername, decryptedPassword, nil
		}
	}

	return "", "", "", fmt.Errorf("password not found")
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

func EditPassword(id, newTitle, newLink, newUsername, newPassword, masterPassword string, pathToPasswordFile string) error {
	manager, err := ReadPasswordManager(pathToPasswordFile)
	if err != nil {
		return fmt.Errorf("failed to read password manager: %w", err)
	}

	if !VerifyMasterPasswordIsHashedPassword(masterPassword, manager.MasterPasswordHash) {
		return fmt.Errorf("invalid master password")
	}

	encryptionKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, manager.Salt)

	for i, password := range manager.Passwords {
		if password.ID == id {
			if newPassword != "" {
				encryptedPassword, err := encryption.EncryptWithAES(newPassword, encryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt new password: %w", err)
				}
				manager.Passwords[i].EncryptedPassword = encryptedPassword
			}

			if newLink != "" {
				encryptedLink, err := encryption.EncryptWithAES(newLink, encryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt new link: %w", err)
				}
				manager.Passwords[i].Link = encryptedLink
			}

			if newUsername != "" {
				encryptedUsername, err := encryption.EncryptWithAES(newUsername, encryptionKey)
				if err != nil {
					return fmt.Errorf("failed to encrypt new username: %w", err)
				}
				manager.Passwords[i].Username = encryptedUsername
			}

			if newTitle != "" {
				manager.Passwords[i].Title = newTitle
			}

			err = writePasswordManager(manager, pathToPasswordFile)
			if err != nil {
				return fmt.Errorf("failed to write password manager: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("password not found")
}

func CheckPasswordFileExistsInDataDirectory(dataDir string) bool {
	_, err := os.Stat(dataDir)
	return !os.IsNotExist(err)
}
