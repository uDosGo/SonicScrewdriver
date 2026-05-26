package library

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	m := NewManager("/test/path.yaml")
	assert.NotNil(t, m)
	assert.Equal(t, "/test/path.yaml", m.indexPath)
}

func TestManagerLoad_FileNotFound(t *testing.T) {
	m := NewManager("/nonexistent/path.yaml")
	err := m.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read index")
}

func TestManagerLoad_Valid(t *testing.T) {
	// Create a valid YAML index file
	content := `version: 1
last_updated: "2026-01-01"
games:
  - name: game1
    path: /games/game1
    status: available
  - name: game2
    path: /games/game2
    status: installed
`
	indexPath := filepath.Join(t.TempDir(), "library-index.yaml")
	err := os.WriteFile(indexPath, []byte(content), 0644)
	require.NoError(t, err)

	m := NewManager(indexPath)
	err = m.Load()
	require.NoError(t, err)
	assert.Len(t, m.Index.Games, 2)
	assert.Equal(t, 1, m.Index.Version)
}

func TestManagerList(t *testing.T) {
	m := NewManager("")
	m.Index = Index{
		Version: 1,
		Games: []Entry{
			{Name: "game1", Path: "/games/game1", Status: "available"},
			{Name: "game2", Path: "/games/game2", Status: "installed"},
		},
	}

	entries, err := m.List()
	require.NoError(t, err)
	assert.Len(t, entries, 2)
}

func TestManagerList_Empty(t *testing.T) {
	m := NewManager("")
	m.Index = Index{Version: 1}

	_, err := m.List()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no games in library")
}

func TestManagerGetGame(t *testing.T) {
	m := NewManager("")
	m.Index = Index{
		Version: 1,
		Games: []Entry{
			{Name: "game1", Path: "/games/game1", Status: "available"},
			{Name: "game2", Path: "/games/game2", Status: "installed"},
		},
	}

	entry, err := m.GetGame("game1")
	require.NoError(t, err)
	assert.Equal(t, "game1", entry.Name)
	assert.Equal(t, "available", entry.Status)
}

func TestManagerGetGame_NotFound(t *testing.T) {
	m := NewManager("")
	m.Index = Index{
		Version: 1,
		Games: []Entry{
			{Name: "game1", Path: "/games/game1"},
		},
	}

	_, err := m.GetGame("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestManagerValidateManifest_FileNotFound(t *testing.T) {
	m := NewManager("")
	err := m.ValidateManifest("/nonexistent/manifest.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "manifest file does not exist")
}

func TestManagerValidateManifest_InvalidYAML(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "manifest.yaml")
	err := os.WriteFile(manifestPath, []byte("invalid: yaml: [bad"), 0644)
	require.NoError(t, err)

	m := NewManager("")
	err = m.ValidateManifest(manifestPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse manifest YAML")
}

func TestManagerValidateManifest_Valid(t *testing.T) {
	manifestContent := `name: test-game
version: 1.0.0
container:
  image: test/game
`
	manifestPath := filepath.Join(t.TempDir(), "manifest.yaml")
	err := os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	require.NoError(t, err)

	m := NewManager("")
	err = m.ValidateManifest(manifestPath)
	assert.NoError(t, err)
}

func TestManagerValidateManifest_MissingRequired(t *testing.T) {
	// Missing required "container" field
	manifestContent := `name: test-game
version: 1.0.0
`
	manifestPath := filepath.Join(t.TempDir(), "manifest.yaml")
	err := os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	require.NoError(t, err)

	m := NewManager("")
	err = m.ValidateManifest(manifestPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "schema validation errors")
}

func TestManagerUpdate(t *testing.T) {
	m := NewManager("")
	err := m.Update()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestGetDefaultIndexPath(t *testing.T) {
	path := GetDefaultIndexPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "library-index.yaml")
}

func TestEntryStruct(t *testing.T) {
	e := Entry{
		Name:   "test-game",
		Path:   "/games/test-game",
		Status: "available",
	}
	assert.Equal(t, "test-game", e.Name)
	assert.Equal(t, "/games/test-game", e.Path)
	assert.Equal(t, "available", e.Status)
}

func TestIndexStruct(t *testing.T) {
	idx := Index{
		Version:     1,
		LastUpdated: "2026-01-01",
		Games:       []Entry{},
	}
	assert.Equal(t, 1, idx.Version)
	assert.Equal(t, "2026-01-01", idx.LastUpdated)
	assert.Empty(t, idx.Games)
}
