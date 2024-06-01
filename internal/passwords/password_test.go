package passwords_test

import (
	"os"
	"testing"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorePassword(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	masterPassword := "master-password"
	err = passwords.InitializePasswordManager(masterPassword)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		title        string
		link         string
		password     string
		wantErr      bool
		errorMessage string
	}{
		{
			name:         "Valid password",
			title:        "Test Title",
			link:         "https://example.com",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
		{
			name:         "Empty password",
			title:        "Test Title",
			link:         "https://example.com",
			password:     "",
			wantErr:      true,
			errorMessage: "password cannot be empty",
		},
		{
			name:         "Empty title",
			title:        "",
			link:         "https://example.com",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
		{
			name:         "Empty link",
			title:        "Test Title",
			link:         "",
			password:     "password123",
			wantErr:      false,
			errorMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := passwords.StorePassword(tc.title, tc.link, tc.password, masterPassword)
			if tc.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestRetrievePassword(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	masterPassword := "master-password"
	err = passwords.InitializePasswordManager(masterPassword)
	require.NoError(t, err)

	title := "Test Title"
	link := "https://example.com"
	password := "password123"
	err = passwords.StorePassword(title, link, password, masterPassword)
	require.NoError(t, err)

	data, err := passwords.ReadPasswordManager()
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
			retrievedPassword, err := passwords.RetrievePassword(tc.id, tc.masterPassword)
			if tc.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, password, retrievedPassword)
			}
		})
	}

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestInitializePasswordManager(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	masterPassword := "test-master-password"

	err = passwords.InitializePasswordManager(masterPassword)
	assert.NoError(t, err)

	_, err = os.Stat("data/passwords.json")
	assert.NoError(t, err)

	data, err := passwords.ReadPasswordManager()
	require.NoError(t, err)
	assert.NotEmpty(t, data.MasterPasswordHash)
	assert.NotEmpty(t, data.Salt)
	assert.Empty(t, data.Passwords)

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestCheckPasswordFileExists(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	exists := passwords.CheckPasswordFileExists()
	assert.False(t, exists)

	err = passwords.CreateEmptyPasswordFile()
	require.NoError(t, err)

	exists = passwords.CheckPasswordFileExists()
	assert.True(t, exists)

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestCreateEmptyPasswordFile(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	err = passwords.CreateEmptyPasswordFile()
	assert.NoError(t, err)

	_, err = os.Stat("data/passwords.json")
	assert.NoError(t, err)

	data, err := os.ReadFile("data/passwords.json")
	require.NoError(t, err)
	assert.Equal(t, "[]\n", string(data))

	err = os.RemoveAll("data")
	require.NoError(t, err)
}
