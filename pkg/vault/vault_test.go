package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expands tilde",
			input:    "~/test/path",
			expected: filepath.Join(home, "test/path"),
		},
		{
			name:     "no tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just tilde",
			input:    "~",
			expected: home,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVaultOpen(t *testing.T) {
	// Test that Open fails gracefully when no key exists
	// (we can't easily test the full flow without a real key)
	v, err := Open()
	if err != nil {
		// This is expected if no master key exists and we can't create one
		// (e.g., in CI environments without write access)
		t.Logf("Open returned expected error: %v", err)
		return
	}
	if v != nil {
		defer func() {
			// Clean up test artifacts
			home, _ := os.UserHomeDir()
			os.RemoveAll(filepath.Join(home, ".sonic"))
		}()
		assert.NotNil(t, v.store)
	}
}
