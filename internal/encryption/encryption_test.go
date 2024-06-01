package encryption_test

import (
	"reflect"
	"testing"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptWithAES(t *testing.T) {
	masterPassword := "test-master-password"
	salt, err := encryption.GenerateSalt()
	require.NoError(t, err)
	correctKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, salt)

	testCases := []struct {
		name      string
		plaintext string
		key       []byte
		wantErr   bool
	}{
		{
			name:      "Valid plaintext with correct key",
			plaintext: "password123",
			key:       correctKey,
			wantErr:   false,
		},
		{
			name:      "Empty plaintext with correct key",
			plaintext: "",
			key:       correctKey,
			wantErr:   false,
		},
		{
			name:      "Valid plaintext with empty key",
			plaintext: "password123",
			key:       []byte{},
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := encryption.EncryptWithAES(tc.plaintext, tc.key)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				decrypted, err := encryption.DecryptWithAES(encrypted, tc.key)
				if reflect.DeepEqual(tc.key, correctKey) {
					assert.NoError(t, err)
					assert.Equal(t, tc.plaintext, decrypted)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestDecryptWithAES(t *testing.T) {
	masterPassword := "test-master-password"
	incorrectPassword := "incorrect-password"
	salt, err := encryption.GenerateSalt()
	require.NoError(t, err)
	correctKey := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, salt)
	incorrectKey := encryption.DeriveEncryptionKeyFromMasterPassword(incorrectPassword, salt)

	plaintext := "password123"
	validPassword, err := encryption.EncryptWithAES(plaintext, correctKey)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		password string
		key      []byte
		wantErr  bool
	}{
		{
			name:     "Valid password with correct key",
			password: validPassword,
			key:      correctKey,
			wantErr:  false,
		},
		{
			name:     "Valid password with incorrect key",
			password: validPassword,
			key:      incorrectKey,
			wantErr:  true,
		},
		{
			name:     "Empty password with correct key",
			password: "",
			key:      correctKey,
			wantErr:  true,
		},
		{
			name:     "Invalid password with correct key",
			password: "invalid-password",
			key:      correctKey,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := encryption.DecryptWithAES(tc.password, tc.key)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeriveEncryptionKey(t *testing.T) {
	masterPassword := "master-password"
	salt, err := encryption.GenerateSalt()
	require.NoError(t, err)

	key := encryption.DeriveEncryptionKeyFromMasterPassword(masterPassword, salt)
	assert.Len(t, key, 32)
}

func TestGenerateSalt(t *testing.T) {
	salt, err := encryption.GenerateSalt()
	assert.NoError(t, err)
	assert.Len(t, salt, 16)
}
