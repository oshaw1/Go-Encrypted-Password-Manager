package encryption

import (
	"crypto/sha256"
	"fmt"
)

func EncryptWithSHA256(plaintext string) (string, error) {
	plaintextBytes := []byte(plaintext)

	hash := sha256.New()

	_, err := hash.Write(plaintextBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write plaintext to hash: %w", err)
	}

	hashedBytes := hash.Sum(nil)

	return fmt.Sprintf("%x", hashedBytes), nil
}
