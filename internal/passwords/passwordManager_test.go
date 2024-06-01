package passwords_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializePasswordManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "password-manager-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	passwordFilePath := filepath.Join(tempDir, "passwords.json")

	masterPassword := "test-master-password"

	err = passwords.InitializePasswordManager(masterPassword, passwordFilePath)
	assert.NoError(t, err)

	_, err = os.Stat(passwordFilePath)
	assert.NoError(t, err)

	data, err := passwords.ReadPasswordManager(passwordFilePath)
	require.NoError(t, err)
	assert.NotEmpty(t, data.MasterPasswordHash)
	assert.NotEmpty(t, data.Salt)
	assert.Empty(t, data.Passwords)

	err = os.RemoveAll("data")
	require.NoError(t, err)
}
