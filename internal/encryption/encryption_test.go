package encryption_test

import (
	"testing"

	"github.com/oshaw1/Encrypted-Password-Manager/internal/encryption"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWithSHA256(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		expected  string
		wantErr   bool
	}{
		{
			name:      "Valid plaintext",
			plaintext: "password123",
			expected:  "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			wantErr:   false,
		},
		{
			name:      "Empty plaintext",
			plaintext: "",
			expected:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encryption.EncryptWithSHA256(tc.plaintext)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
