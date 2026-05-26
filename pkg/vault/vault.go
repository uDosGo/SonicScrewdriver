// Package vault provides access to the uServer encrypted secret store.
// It wraps the uServer/pkg/secrets package to provide a SonicScrewdriver-specific
// interface for managing secrets (API keys, tokens, credentials).
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uDosGo/uServer/pkg/secrets"
)

// Default paths for the SonicScrewdriver vault
const (
	DefaultKeyPath  = "~/.sonic/master.key"
	DefaultDataPath = "~/.sonic/secrets.enc"
)

// Vault wraps the uServer secret store with SonicScrewdriver defaults
type Vault struct {
	store *secrets.SecretStore
}

// Open opens the SonicScrewdriver vault, creating it if needed
func Open() (*Vault, error) {
	keyPath := expandPath(DefaultKeyPath)
	dataPath := expandPath(DefaultDataPath)

	// Load or generate master key
	var masterKey []byte
	if _, err := os.Stat(keyPath); err == nil {
		masterKey, err = secrets.LoadMasterKey(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load master key: %w", err)
		}
	} else {
		masterKey, err = secrets.GenerateMasterKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate master key: %w", err)
		}
		if err := secrets.SaveMasterKey(masterKey, keyPath); err != nil {
			return nil, fmt.Errorf("failed to save master key: %w", err)
		}
	}

	store, err := secrets.NewSecretStore(masterKey, dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open secret store: %w", err)
	}

	return &Vault{store: store}, nil
}

// Get retrieves a secret value by key
func (v *Vault) Get(key string) (string, error) {
	return v.store.GetSecret(key)
}

// Set stores a secret value by key
func (v *Vault) Set(key, value string) error {
	return v.store.AddSecret(key, value)
}

// List returns all secret keys
func (v *Vault) List() ([]string, error) {
	return v.store.ListSecrets()
}

// Rotate replaces a secret value, preserving history
func (v *Vault) Rotate(key, newValue string) error {
	return v.store.RotateSecret(key, newValue)
}

// History returns the rotation history for a secret
func (v *Vault) History(key string) ([]map[string]string, error) {
	return v.store.GetSecretHistory(key)
}

func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
