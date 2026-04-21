package library

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/qri-io/jsonschema"
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

	// Parse YAML to JSON for schema validation
	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	// Convert to JSON for schema validation
	jsonData, err := json.Marshal(yamlData)
	if err != nil {
		return fmt.Errorf("failed to convert manifest to JSON: %w", err)
	}

	// Validate against schema
	if err := validateAgainstSchema(jsonData); err != nil {
		return fmt.Errorf("manifest validation failed: %w", err)
	}

	log.Printf("Validated manifest: %s", path)
	return nil
}

func validateAgainstSchema(data []byte) error {
	// Parse the schema
	schema := jsonschema.Schema{}
	if err := json.Unmarshal([]byte(ManifestSchema), &schema); err != nil {
		return fmt.Errorf("failed to parse validation schema: %w", err)
	}

	// Prepare data as interface{}
	var dataInterface interface{}
	if err := json.Unmarshal(data, &dataInterface); err != nil {
		return fmt.Errorf("failed to parse manifest data: %w", err)
	}

	// Validate the data
	ctx := context.Background()
	state := schema.Validate(ctx, dataInterface)

	if !state.IsValid() {
		var errorMessages []string
		if state.Errs != nil {
			for _, err := range *state.Errs {
				errorMessages = append(errorMessages, err.Message)
			}
		}
		return fmt.Errorf("schema validation errors: %v", errorMessages)
	}

	return nil
}

func GetDefaultIndexPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./library-index.yaml"
	}
	// Use uDos vendor directory for games
	udosVendorPath := filepath.Join(homeDir, "uDos", "vendor", "games")
	if _, err := os.Stat(udosVendorPath); err == nil {
		return filepath.Join(udosVendorPath, "library-index.yaml")
	}
	return filepath.Join(homeDir, ".sonic", "library-index.yaml")
}
