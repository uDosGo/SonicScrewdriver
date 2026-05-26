package ventoy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uDosGo/SonicScrewdriver/pkg/library"
)

func TestNewPackager(t *testing.T) {
	manager := library.NewManager("")
	p := NewPackager("/source", "/output", manager)
	assert.NotNil(t, p)
	assert.Equal(t, "/source", p.SourceDir)
	assert.Equal(t, "/output", p.OutputDir)
	assert.Equal(t, manager, p.Config)
}

func TestCreateBundle_NoSource(t *testing.T) {
	manager := library.NewManager("")
	p := NewPackager("", "/output", manager)
	path, err := p.CreateBundle("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source directory not specified")
	assert.Empty(t, path)
}

func TestCreateBundle_Valid(t *testing.T) {
	// Create a temp source directory with a test file
	sourceDir := t.TempDir()
	err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	outputDir := t.TempDir()
	manager := library.NewManager("")
	p := NewPackager(sourceDir, outputDir, manager)

	bundlePath, err := p.CreateBundle("test-bundle")
	require.NoError(t, err)
	assert.FileExists(t, bundlePath)
	assert.Equal(t, filepath.Join(outputDir, "test-bundle.she"), bundlePath)
}

func TestValidateBundle_Valid(t *testing.T) {
	// Create a valid bundle first
	sourceDir := t.TempDir()
	err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	outputDir := t.TempDir()
	manager := library.NewManager("")
	p := NewPackager(sourceDir, outputDir, manager)

	bundlePath, err := p.CreateBundle("test-validate")
	require.NoError(t, err)

	// Now validate it
	p2 := NewPackager("", "", manager)
	err = p2.ValidateBundle(bundlePath)
	assert.NoError(t, err)
}

func TestValidateBundle_Invalid(t *testing.T) {
	// Create a non-gzip file
	invalidPath := filepath.Join(t.TempDir(), "invalid.she")
	err := os.WriteFile(invalidPath, []byte("not a gzip file"), 0644)
	require.NoError(t, err)

	manager := library.NewManager("")
	p := NewPackager("", "", manager)
	err = p.ValidateBundle(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid gzip format")
}

func TestGetBundleInfo(t *testing.T) {
	// Create a valid bundle
	sourceDir := t.TempDir()
	err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0644)
	require.NoError(t, err)

	outputDir := t.TempDir()
	manager := library.NewManager("")
	p := NewPackager(sourceDir, outputDir, manager)

	bundlePath, err := p.CreateBundle("test-info")
	require.NoError(t, err)

	info, err := GetBundleInfo(bundlePath)
	require.NoError(t, err)
	assert.Equal(t, bundlePath, info["path"])
	assert.Equal(t, "she", info["extension"])
	assert.Equal(t, "Sonic Home Edition Bundle", info["type"])
}

func TestGetBundleInfo_NotFound(t *testing.T) {
	_, err := GetBundleInfo("/nonexistent/path.she")
	assert.Error(t, err)
}
