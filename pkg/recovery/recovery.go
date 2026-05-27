// Package recovery provides device backup and restore capabilities.
// It supports backing up firmware, configurations, and device state
// to the SonicScrewdriver vault for safe keeping and recovery.
package recovery

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BackupType classifies the type of backup
type BackupType string

const (
	BackupFirmware    BackupType = "firmware"     // Device firmware dump
	BackupConfig      BackupType = "config"       // Device configuration
	BackupState       BackupType = "state"        // Device state/settings
	BackupFull        BackupType = "full"         // Full device backup
	BackupPartition   BackupType = "partition"    // Single partition backup
	BackupBootloader  BackupType = "bootloader"   // Bootloader backup
	BackupFilesystem  BackupType = "filesystem"   // Filesystem backup
)

// BackupRecord represents a single backup entry
type BackupRecord struct {
	ID          string            `json:"id"`
	Device      string            `json:"device"`
	Type        BackupType        `json:"type"`
	Description string            `json:"description"`
	FilePath    string            `json:"file_path"`
	Size        int64             `json:"size"`
	SHA256      string            `json:"sha256"`
	CreatedAt   string            `json:"created_at"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RecoveryManager handles device backup and restore operations
type RecoveryManager struct {
	backupDir string
	indexPath string
	index     map[string]*BackupRecord
}

// NewManager creates a new RecoveryManager
func NewManager() *RecoveryManager {
	home, _ := os.UserHomeDir()
	backupDir := filepath.Join(home, ".sonic", "backups")
	return &RecoveryManager{
		backupDir: backupDir,
		indexPath: filepath.Join(backupDir, "index.json"),
		index:     make(map[string]*BackupRecord),
	}
}

// Load reads the backup index from disk
func (m *RecoveryManager) Load() error {
	if _, err := os.Stat(m.indexPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return fmt.Errorf("failed to read backup index: %w", err)
	}

	var records []*BackupRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse backup index: %w", err)
	}

	for _, r := range records {
		m.index[r.ID] = r
	}
	return nil
}

// Save writes the backup index to disk
func (m *RecoveryManager) Save() error {
	records := make([]*BackupRecord, 0, len(m.index))
	for _, r := range m.index {
		records = append(records, r)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup index: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(m.indexPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup dir: %w", err)
	}

	return os.WriteFile(m.indexPath, data, 0644)
}

// BackupDevice creates a backup of a device
func (m *RecoveryManager) BackupDevice(devicePath string, backupType BackupType, description string) (*BackupRecord, error) {
	// Validate source
	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("device path not found: %s", devicePath)
	}

	// Generate backup ID
	timestamp := time.Now().UTC().Format("20060102-150405")
	deviceName := filepath.Base(devicePath)
	id := fmt.Sprintf("%s-%s-%s", deviceName, string(backupType), timestamp)

	// Create backup directory
	backupPath := filepath.Join(m.backupDir, id)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup dir: %w", err)
	}

	// Perform backup based on type
	var backupFile string
	var err error

	switch backupType {
	case BackupFirmware:
		backupFile, err = m.backupFirmware(devicePath, backupPath)
	case BackupConfig:
		backupFile, err = m.backupConfig(devicePath, backupPath)
	case BackupState:
		backupFile, err = m.backupState(devicePath, backupPath)
	case BackupFull:
		backupFile, err = m.backupFull(devicePath, backupPath)
	case BackupPartition:
		backupFile, err = m.backupPartition(devicePath, backupPath)
	case BackupBootloader:
		backupFile, err = m.backupBootloader(devicePath, backupPath)
	case BackupFilesystem:
		backupFile, err = m.backupFilesystem(devicePath, backupPath)
	default:
		return nil, fmt.Errorf("unsupported backup type: %s", backupType)
	}

	if err != nil {
		return nil, fmt.Errorf("backup failed: %w", err)
	}

	// Calculate checksum
	checksum, err := calculateSHA256(backupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Get file size
	info, err := os.Stat(backupFile)
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	record := &BackupRecord{
		ID:          id,
		Device:      devicePath,
		Type:        backupType,
		Description: description,
		FilePath:    backupFile,
		Size:        info.Size(),
		SHA256:      checksum,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		Tags:        []string{string(backupType), deviceName},
		Metadata: map[string]string{
			"device": devicePath,
			"type":   string(backupType),
		},
	}

	m.index[id] = record
	if err := m.Save(); err != nil {
		return nil, fmt.Errorf("failed to save backup index: %w", err)
	}

	return record, nil
}

// backupFirmware creates a firmware backup by reading from the device
func (m *RecoveryManager) backupFirmware(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "firmware.bin")

	// Read firmware from device (using dd)
	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read device: %w", err)
	}

	if err := os.WriteFile(outputFile, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return outputFile, nil
}

// backupConfig creates a configuration backup
func (m *RecoveryManager) backupConfig(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "config.json")

	// Read config from device (first 1MB typically contains config)
	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read device config: %w", err)
	}

	// Limit to first 1MB for config
	configSize := len(input)
	if configSize > 1024*1024 {
		configSize = 1024 * 1024
	}

	config := map[string]interface{}{
		"device":      devicePath,
		"backup_time": time.Now().UTC().Format(time.RFC3339),
		"size":        configSize,
		"data":        fmt.Sprintf("%x", input[:configSize]),
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write config backup: %w", err)
	}

	return outputFile, nil
}

// backupState creates a device state backup
func (m *RecoveryManager) backupState(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "state.json")

	// Read device state
	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read device state: %w", err)
	}

	state := map[string]interface{}{
		"device":      devicePath,
		"backup_time": time.Now().UTC().Format(time.RFC3339),
		"size":        len(input),
		"checksum":    fmt.Sprintf("%x", sha256.Sum256(input)),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write state backup: %w", err)
	}

	return outputFile, nil
}

// backupFull creates a full device backup (disk image)
func (m *RecoveryManager) backupFull(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "full-backup.img")

	// Full device backup - copy entire device
	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read device for full backup: %w", err)
	}

	if err := os.WriteFile(outputFile, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write full backup: %w", err)
	}

	return outputFile, nil
}

// backupPartition creates a partition backup
func (m *RecoveryManager) backupPartition(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "partition-backup.bin")

	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read partition: %w", err)
	}

	if err := os.WriteFile(outputFile, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write partition backup: %w", err)
	}

	return outputFile, nil
}

// backupBootloader creates a bootloader backup
func (m *RecoveryManager) backupBootloader(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "bootloader.bin")

	// Bootloader is typically in the first few sectors
	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read device: %w", err)
	}

	// First 512KB typically contains bootloader
	bootloaderSize := len(input)
	if bootloaderSize > 512*1024 {
		bootloaderSize = 512 * 1024
	}

	if err := os.WriteFile(outputFile, input[:bootloaderSize], 0644); err != nil {
		return "", fmt.Errorf("failed to write bootloader backup: %w", err)
	}

	return outputFile, nil
}

// backupFilesystem creates a filesystem backup
func (m *RecoveryManager) backupFilesystem(devicePath, backupPath string) (string, error) {
	outputFile := filepath.Join(backupPath, "filesystem.img")

	input, err := os.ReadFile(devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to read filesystem: %w", err)
	}

	if err := os.WriteFile(outputFile, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write filesystem backup: %w", err)
	}

	return outputFile, nil
}

// RestoreDevice restores a device from a backup
func (m *RecoveryManager) RestoreDevice(backupID, targetDevice string) error {
	record, ok := m.index[backupID]
	if !ok {
		return fmt.Errorf("backup '%s' not found", backupID)
	}

	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", record.FilePath)
	}

	// Verify checksum
	checksum, err := calculateSHA256(record.FilePath)
	if err != nil {
		return fmt.Errorf("failed to verify backup: %w", err)
	}

	if checksum != record.SHA256 {
		return fmt.Errorf("backup checksum mismatch - file may be corrupted")
	}

	// Write backup to target device
	input, err := os.ReadFile(record.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	if err := os.WriteFile(targetDevice, input, 0644); err != nil {
		return fmt.Errorf("failed to write to device: %w", err)
	}

	return nil
}

// ListBackups returns all backup records, optionally filtered by device
func (m *RecoveryManager) ListBackups(deviceFilter string) []*BackupRecord {
	var records []*BackupRecord
	for _, r := range m.index {
		if deviceFilter == "" || strings.Contains(r.Device, deviceFilter) {
			records = append(records, r)
		}
	}
	return records
}

// GetBackup returns a specific backup record
func (m *RecoveryManager) GetBackup(id string) (*BackupRecord, error) {
	record, ok := m.index[id]
	if !ok {
		return nil, fmt.Errorf("backup '%s' not found", id)
	}
	return record, nil
}

// DeleteBackup removes a backup record and its files
func (m *RecoveryManager) DeleteBackup(id string) error {
	record, ok := m.index[id]
	if !ok {
		return fmt.Errorf("backup '%s' not found", id)
	}

	// Remove backup directory
	backupDir := filepath.Dir(record.FilePath)
	if err := os.RemoveAll(backupDir); err != nil {
		return fmt.Errorf("failed to remove backup files: %w", err)
	}

	delete(m.index, id)
	return m.Save()
}

// ExportBackup exports a backup to a specified location
func (m *RecoveryManager) ExportBackup(id, exportPath string) error {
	record, ok := m.index[id]
	if !ok {
		return fmt.Errorf("backup '%s' not found", id)
	}

	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", record.FilePath)
	}

	// Copy backup file to export path
	input, err := os.ReadFile(record.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	if err := os.WriteFile(exportPath, input, 0644); err != nil {
		return fmt.Errorf("failed to export backup: %w", err)
	}

	return nil
}

// ImportBackup imports a backup file into the manager
func (m *RecoveryManager) ImportBackup(filePath, deviceName string, backupType BackupType) (*BackupRecord, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	timestamp := time.Now().UTC().Format("20060102-150405")
	id := fmt.Sprintf("%s-%s-%s-imported", deviceName, string(backupType), timestamp)

	backupPath := filepath.Join(m.backupDir, id)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup dir: %w", err)
	}

	destFile := filepath.Join(backupPath, filepath.Base(filePath))
	input, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	if err := os.WriteFile(destFile, input, 0644); err != nil {
		return nil, fmt.Errorf("failed to copy import file: %w", err)
	}

	checksum, err := calculateSHA256(destFile)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	info, _ := os.Stat(destFile)

	record := &BackupRecord{
		ID:          id,
		Device:      deviceName,
		Type:        backupType,
		Description: fmt.Sprintf("Imported from %s", filePath),
		FilePath:    destFile,
		Size:        info.Size(),
		SHA256:      checksum,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		Tags:        []string{"imported", string(backupType)},
	}

	m.index[id] = record
	if err := m.Save(); err != nil {
		return nil, fmt.Errorf("failed to save backup index: %w", err)
	}

	return record, nil
}

// calculateSHA256 calculates the SHA256 checksum of a file
func calculateSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Ensure json is used
var _ = json.Marshal
