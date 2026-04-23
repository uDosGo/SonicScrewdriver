package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SecretStore manages encrypted secrets
type SecretStore struct {
	masterKey []byte
	filePath string
	secrets   map[string]string
	policies  map[string]SecretPolicy
	mu        sync.RWMutex
	cache     map[string]string
	cachePath string
}

// SecretPolicy defines access control for a secret
type SecretPolicy struct {
	AllowedNodes []string `json:"allowed_nodes"`
	AllowedRoles  []string `json:"allowed_roles"`
	RateLimit     string   `json:"rate_limit"`
}

// StoreData represents the encrypted data structure
type StoreData struct {
	Version  int                     `json:"version"`
	Secrets  map[string]string      `json:"secrets"`
	Policies map[string]SecretPolicy `json:"policies"`
}

// NewSecretStore creates a new secret store instance
func NewSecretStore(masterKey []byte, filePath string) (*SecretStore, error) {
	if len(masterKey) != 32 {
		return nil, errors.New("master key must be 32 bytes")
	}

	cachePath := getCachePath()
	
	store := &SecretStore{
		masterKey: masterKey,
		filePath:  filePath,
		cachePath: cachePath,
		secrets:   make(map[string]string),
		policies:  make(map[string]SecretPolicy),
		cache:     make(map[string]string),
	}

	// Load existing secrets if file exists
	if _, err := os.Stat(filePath); err == nil {
		if err := store.load(); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// Load cache if it exists
	if _, err := os.Stat(cachePath); err == nil {
		if err := store.loadCache(); err != nil {
			log.Printf("Warning: Failed to load cache: %v", err)
		}
	}

	return store, nil
}

// load decrypts and loads secrets from file
func (s *SecretStore) load() error {
	encryptedData, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	decryptedData, err := s.decrypt(encryptedData)
	if err != nil {
		return err
	}

	var storeData StoreData
	if err := json.Unmarshal(decryptedData, &storeData); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.secrets = storeData.Secrets
	s.policies = storeData.Policies

	return nil
}

// save encrypts and saves secrets to file
func (s *SecretStore) save() error {
	s.mu.RLock()
	storeData := StoreData{
		Version:  2,
		Secrets:  s.secrets,
		Policies: s.policies,
	}
	s.mu.RUnlock()

	data, err := json.Marshal(storeData)
	if err != nil {
		return err
	}

	encryptedData, err := s.encrypt(data)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, encryptedData, 0600)
}

// encrypt encrypts data using AES-256-GCM
func (s *SecretStore) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// decrypt decrypts data using AES-256-GCM
func (s *SecretStore) decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, errors.New("invalid encrypted data")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AddSecret adds a new secret to the store
func (s *SecretStore) AddSecret(name, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.secrets[name] = value
	return s.save()
}

// GetSecret retrieves a secret by name
func (s *SecretStore) GetSecret(name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.secrets[name]
	if !exists {
		return "", errors.New("secret not found")
	}
	return value, nil
}

// ListSecrets returns a list of secret names
func (s *SecretStore) ListSecrets() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	secrets := make([]string, 0, len(s.secrets))
	for name := range s.secrets {
		secrets = append(secrets, name)
	}
	return secrets, nil
}

// SetPolicy sets access policy for a secret
func (s *SecretStore) SetPolicy(name string, policy SecretPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[name] = policy
	return s.save()
}

// GetPolicy retrieves access policy for a secret
func (s *SecretStore) GetPolicy(name string) (SecretPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	policy, exists := s.policies[name]
	if !exists {
		return SecretPolicy{}, errors.New("policy not found")
	}
	return policy, nil
}

// GenerateMasterKey generates a new 32-byte master key
func GenerateMasterKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// SaveMasterKey saves the master key to a file
func SaveMasterKey(key []byte, filePath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, key, 0600)
}

// getCachePath returns the path to the cache file
func getCachePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./secrets.cache"
	}
	return filepath.Join(homeDir, ".sonic", "secrets.cache")
}

// loadCache loads the cache from file
func (s *SecretStore) loadCache() error {
	data, err := os.ReadFile(s.cachePath)
	if err != nil {
		return err
	}

	var cache map[string]string
	if err := json.Unmarshal(data, &cache); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = cache

	return nil
}

// saveCache saves the cache to file
func (s *SecretStore) saveCache() error {
	s.mu.RLock()
	data, err := json.Marshal(s.cache)
	s.mu.RUnlock()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.cachePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.cachePath, data, 0600)
}

// RotateSecret rotates a secret and maintains version history
func (s *SecretStore) RotateSecret(name, newValue string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get current value if it exists
	currentValue, exists := s.secrets[name]
	
	// Store previous value in history (simple implementation)
	historyKey := name + "_history"
	if !exists {
		s.secrets[historyKey] = ""
	} else {
		// Store previous value with timestamp
		historyValue := fmt.Sprintf("%s|%s|%s", currentValue, time.Now().Format("2006-01-02"), "rotated")
		s.secrets[historyKey] = historyValue
	}

	// Update to new value
	s.secrets[name] = newValue
	
	// Save to file
	return s.save()
}

// GetSecretHistory retrieves rotation history for a secret
func (s *SecretStore) GetSecretHistory(name string) ([]map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	historyKey := name + "_history"
	historyValue, exists := s.secrets[historyKey]
	
	if !exists || historyValue == "" {
		return nil, errors.New("no history available")
	}

	// Parse history (simple format: value|date|action)
	historyEntries := make([]map[string]string, 0)
	
	// For this simple implementation, we store one history entry
	// In a production system, this would be more sophisticated
	parts := strings.Split(historyValue, "|")
	if len(parts) >= 3 {
		historyEntries = append(historyEntries, map[string]string{
			"value":   parts[0],
			"date":    parts[1],
			"action":  parts[2],
		})
	}

	return historyEntries, nil
}

// GetSecretWithCache retrieves a secret, using cache if available
func (s *SecretStore) GetSecretWithCache(name string, allowCached bool) (string, bool, error) {
	// First try to get from main store
	value, err := s.GetSecret(name)
	if err == nil {
		// Update cache
		s.mu.Lock()
		s.cache[name] = value
		s.mu.Unlock()
		go s.saveCache() // Save cache asynchronously
		return value, false, nil
	}

	// If not found and caching is allowed, try cache
	if allowCached {
		s.mu.RLock()
		cachedValue, exists := s.cache[name]
		s.mu.RUnlock()
		if exists {
			return cachedValue, true, nil
		}
	}

	return "", false, err
}

// LoadMasterKey loads the master key from a file
func LoadMasterKey(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// Backup creates a backup of the secret store
func (s *SecretStore) Backup(backupPath string) error {
	s.mu.RLock()
	storeData := StoreData{
		Version:  2,
		Secrets:  s.secrets,
		Policies: s.policies,
	}
	s.mu.RUnlock()

	data, err := json.Marshal(storeData)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0600)
}

// Restore restores a backup to the secret store
func (s *SecretStore) Restore(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	var storeData StoreData
	if err := json.Unmarshal(data, &storeData); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.secrets = storeData.Secrets
	s.policies = storeData.Policies

	// Save to the main file
	return s.save()
}

// ExportBackup exports an encrypted backup
func (s *SecretStore) ExportBackup(backupPath string) error {
	s.mu.RLock()
	storeData := StoreData{
		Version:  2,
		Secrets:  s.secrets,
		Policies: s.policies,
	}
	s.mu.RUnlock()

	data, err := json.Marshal(storeData)
	if err != nil {
		return err
	}

	encryptedData, err := s.encrypt(data)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(backupPath, encryptedData, 0600)
}

// ImportBackup imports an encrypted backup
func (s *SecretStore) ImportBackup(backupPath string) error {
	encryptedData, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	decryptedData, err := s.decrypt(encryptedData)
	if err != nil {
		return err
	}

	var storeData StoreData
	if err := json.Unmarshal(decryptedData, &storeData); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.secrets = storeData.Secrets
	s.policies = storeData.Policies

	// Save to the main file
	return s.save()
}