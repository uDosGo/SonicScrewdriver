package iso

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDistro(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantName  string
		wantErr   bool
	}{
		{"ubuntu", "ubuntu", "ubuntu", false},
		{"ubuntu-24.04", "ubuntu-24.04", "ubuntu", false},
		{"ubuntu2404", "ubuntu2404", "ubuntu", false},
		{"mint", "mint", "linuxmint", false},
		{"linuxmint", "linuxmint", "linuxmint", false},
		{"linux-mint", "linux-mint", "linuxmint", false},
		{"mint22", "mint22", "linuxmint", false},
		{"classicmodern", "classicmodern", "classicmodern", false},
		{"classic-modern", "classic-modern", "classicmodern", false},
		{"classic-modern-mint", "classic-modern-mint", "classicmodern", false},
		{"unknown", "unknown", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distro, err := GetDistro(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, distro.Name)
			assert.NotEmpty(t, distro.URL)
			assert.NotEmpty(t, distro.Arch)
		})
	}
}

func TestListDistros(t *testing.T) {
	distros := ListDistros()
	assert.Len(t, distros, 3)

	names := make(map[string]bool)
	for _, d := range distros {
		names[d.Name] = true
	}
	assert.True(t, names["ubuntu"])
	assert.True(t, names["linuxmint"])
	assert.True(t, names["classicmodern"])
}

func TestGetCacheDir(t *testing.T) {
	dir := GetCacheDir()
	assert.Contains(t, dir, ".sonic/iso-cache")
}

func TestVerifySHA256(t *testing.T) {
	// Create a temp file with known content
	content := []byte("test content for sha256 verification")
	tmpFile := filepath.Join(t.TempDir(), "test.iso")
	err := os.WriteFile(tmpFile, content, 0644)
	require.NoError(t, err)

	// Calculate expected hash
	h := sha256.New()
	h.Write(content)
	expectedHash := fmt.Sprintf("%x", h.Sum(nil))

	// Test valid hash
	assert.True(t, verifySHA256(tmpFile, expectedHash))

	// Test invalid hash
	assert.False(t, verifySHA256(tmpFile, "0000000000000000000000000000000000000000000000000000000000000000"))

	// Test empty hash (should pass)
	assert.True(t, verifySHA256(tmpFile, ""))

	// Test non-existent file
	assert.False(t, verifySHA256("/nonexistent/file.iso", expectedHash))
}

func TestDistroStruct(t *testing.T) {
	// Verify Ubuntu 24.04 distro
	assert.Equal(t, "ubuntu", Ubuntu2404.Name)
	assert.Equal(t, "24.04", Ubuntu2404.Version)
	assert.Contains(t, Ubuntu2404.URL, "releases.ubuntu.com")
	assert.Equal(t, "amd64", Ubuntu2404.Arch)

	// Verify Linux Mint 22 distro
	assert.Equal(t, "linuxmint", LinuxMint22.Name)
	assert.Equal(t, "22", LinuxMint22.Version)
	assert.Contains(t, LinuxMint22.URL, "linuxmint")
	assert.Equal(t, "amd64", LinuxMint22.Arch)

	// Verify Classic Modern uses Mint as base
	assert.Equal(t, "classicmodern", ClassicModern.Name)
	assert.Equal(t, LinuxMint22.URL, ClassicModern.URL)
	assert.Equal(t, LinuxMint22.Mirror, ClassicModern.Mirror)
}
