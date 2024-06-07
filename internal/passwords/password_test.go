package passwords_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorePassword(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "password-manager-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	passwordFilePath := filepath.Join(tempDir, "passwords.json")

	masterPassword := "master-password"

	err = passwords.InitializePasswordManager(masterPassword, passwordFilePath)
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		title        string
		hyperlink    string
		password     string
		wantErr      bool
		errorMessage string
	}{
		{
			name:         "Valid password",
			title:        "Test Title",
			hyperlink:    "https://example.com",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
		{
			name:         "Empty password",
			title:        "Test Title",
			hyperlink:    "https://example.com",
			password:     "",
			wantErr:      true,
			errorMessage: "password cannot be empty",
		},
		{
			name:         "Empty title",
			title:        "",
			hyperlink:    "https://example.com",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
		{
			name:         "Empty hyperlink",
			title:        "Test Title",
			hyperlink:    "",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := passwords.StorePassword(tc.title, tc.hyperlink, tc.password, masterPassword, passwordFilePath, "username")
			if tc.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRetrievePassword(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "password-manager-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	passwordFilePath := filepath.Join(tempDir, "passwords.json")

	masterPassword := "master-password"

	err = passwords.InitializePasswordManager(masterPassword, passwordFilePath)
	assert.NoError(t, err)

	title := "Test Title"
	hyperlink := "https://example.com"
	password := "password123"
	err = passwords.StorePassword(title, hyperlink, password, masterPassword, passwordFilePath, "username")
	require.NoError(t, err)

	data, err := passwords.ReadPasswordManager(passwordFilePath)
	require.NoError(t, err)
	passwordID := data.Passwords[0].ID

	testCases := []struct {
		name           string
		id             string
		masterPassword string
		wantErr        bool
		errorMessage   string
	}{
		{
			name:           "Valid password retrieval",
			id:             passwordID,
			masterPassword: masterPassword,
			wantErr:        false,
			errorMessage:   "",
		},
		{
			name:           "Invalid password ID",
			id:             "invalid-id",
			masterPassword: masterPassword,
			wantErr:        true,
			errorMessage:   "password not found",
		},
		{
			name:           "Invalid master password",
			id:             passwordID,
			masterPassword: "wrong-password",
			wantErr:        true,
			errorMessage:   "invalid master password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			retrievedPassword, _, _, err := passwords.RetrievePassword(tc.id, tc.masterPassword, passwordFilePath)
			if tc.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, password, retrievedPassword)
			}
		})
	}
}

func TestCheckPasswordFileExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "password-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	passwordFilePath := filepath.Join(tempDir, "passwords.json")

	testCases := []struct {
		name          string
		setupFunc     func()
		expectedValue bool
	}{
		{
			name:          "File does not exist",
			setupFunc:     func() {},
			expectedValue: false,
		},
		{
			name: "File exists",
			setupFunc: func() {
				err := passwords.InitializePasswordManager("master-password", passwordFilePath)
				require.NoError(t, err)
			},
			expectedValue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupFunc()
			exists := passwords.CheckPasswordFileExistsInDataDirectory(passwordFilePath)
			assert.Equal(t, tc.expectedValue, exists)
		})
	}
}

func TestDeletePasswordByID(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "password-manager-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	passwordFilePath := filepath.Join(tempDir, "passwords.json")

	masterPassword := "master-password"

	err = passwords.InitializePasswordManager(masterPassword, passwordFilePath)
	assert.NoError(t, err)

	title := "Test Title"
	link := "https://example.com"
	password := "password123"
	err = passwords.StorePassword(title, link, password, masterPassword, passwordFilePath, "username")
	require.NoError(t, err)

	data, err := passwords.ReadPasswordManager(passwordFilePath)
	require.NoError(t, err)
	passwordID := data.Passwords[0].ID

	testCases := []struct {
		name           string
		id             string
		masterPassword string
		wantErr        bool
	}{
		{
			name:           "Valid password deletion",
			id:             passwordID,
			masterPassword: masterPassword,
			wantErr:        false,
		},
		{
			name:           "Invalid password ID",
			id:             "invalid-id",
			masterPassword: masterPassword,
			wantErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalData, err := passwords.ReadPasswordManager(passwordFilePath)
			require.NoError(t, err)

			err = passwords.DeletePasswordByID(tc.id, tc.masterPassword, passwordFilePath)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			updatedData, err := passwords.ReadPasswordManager(passwordFilePath)
			require.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, originalData, updatedData)
			} else {
				assert.Empty(t, updatedData.Passwords)
			}
		})
	}
}
