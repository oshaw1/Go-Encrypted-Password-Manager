package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/xdg-go/pbkdf2"
)

func EncryptWithAES(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

func DecryptWithAES(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	plaintext := string(ciphertextBytes)
	if !utf8.ValidString(plaintext) {
		return "", fmt.Errorf("invalid plaintext")
	}

	return string(ciphertextBytes), nil
}

func DeriveEncryptionKey(masterPassword string, salt []byte) []byte {
	iterations := 100000
	keyLength := 32

	return pbkdf2.Key([]byte(masterPassword), salt, iterations, keyLength, sha256.New)
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
