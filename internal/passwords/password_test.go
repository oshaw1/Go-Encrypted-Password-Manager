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

	testCases := []struct {
		name     string
		title    string
		link     string
		password string
		wantErr  bool
	}{
		{"Valid password", "Test Title", "https://example.com", "password123", false},
		{"Empty password", "Test Title", "https://example.com", "", true},
		{"Empty title", "", "https://example.com", "password123", false},
		{"Empty title", "Test title", "", "password123", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := passwords.StorePassword(tc.title, tc.link, tc.password)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestWritePasswords(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		passwords []passwords.Password
		wantErr   bool
	}{
		{"Valid passwords", []passwords.Password{{ID: "1", Title: "Test", Link: "https://example.com", HashedPassword: "hashed"}}, false},
		{"Empty passwords", []passwords.Password{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := passwords.WritePasswords(tc.passwords)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	err = os.RemoveAll("data")
	require.NoError(t, err)
}

func TestReadPasswords(t *testing.T) {
	err := os.MkdirAll("data", os.ModePerm)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		setup   func()
		want    []passwords.Password
		wantErr bool
	}{
		{"Valid passwords", func() {
			passwords.WritePasswords([]passwords.Password{{ID: "1", Title: "Test", Link: "https://example.com", HashedPassword: "hashed"}})
		}, []passwords.Password{{ID: "1", Title: "Test", Link: "https://example.com", HashedPassword: "hashed"}}, false},
		{"Empty passwords", func() {
			passwords.WritePasswords([]passwords.Password{})
		}, []passwords.Password{}, false},
		{"Non-existent file", func() {
			os.Remove("data/passwords.json")
		}, []passwords.Password{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			got, err := passwords.ReadPasswords()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	err = os.RemoveAll("data")
	require.NoError(t, err)
}
