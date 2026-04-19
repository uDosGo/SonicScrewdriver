package library

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Manager struct {
	indexPath string
	Index     Index
}

func NewManager(indexPath string) *Manager {
	return &Manager{indexPath: indexPath}
}

func (m *Manager) Load() error {
	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	var index Index
	if err := yaml.Unmarshal(data, &index); err != nil {
		return fmt.Errorf("failed to parse index: %w", err)
	}

	m.Index = index
	log.Printf("Loaded library index with %d games", len(index.Games))
	return nil
}

func (m *Manager) List() ([]Entry, error) {
	if len(m.Index.Games) == 0 {
		return nil, fmt.Errorf("no games in library")
	}
	return m.Index.Games, nil
}

func (m *Manager) GetGame(name string) (*Entry, error) {
	for _, game := range m.Index.Games {
		if game.Name == name {
			return &game, nil
		}
	}
	return nil, fmt.Errorf("game %s not found", name)
}

func (m *Manager) Update() error {
	return fmt.Errorf("library update not implemented")
}

func (m *Manager) ValidateManifest(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("manifest file does not exist: %s", path)
	}

	// Read and parse the manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	// Simple validation - check for required fields
	manifest := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Check required fields
	if _, ok := manifest["name"]; !ok {
		return fmt.Errorf("manifest missing required field: name")
	}
	if _, ok := manifest["version"]; !ok {
		return fmt.Errorf("manifest missing required field: version")
	}

	log.Printf("Validated manifest: %s", path)
	return nil
}

func GetDefaultIndexPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./library-index.yaml"
	}
	return filepath.Join(homeDir, ".sonic", "library-index.yaml")
}
