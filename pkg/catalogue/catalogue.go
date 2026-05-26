// Package catalogue provides access to the device catalogue.
// It indexes devices from two sources:
//  1. The Vault device catalogue (Vaults/devices/{user,shared,public}/) — documentation layer
//  2. The uCode code repos (uCode1-uCode4 in ~/Code/) — implementation layer
//
// The vault path is the primary source for curated device documentation, manuals, and configs.
// The uCode repos are scanned for emulators, firmware, and device implementations.
package catalogue

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Device represents a device entry in the catalogue
type Device struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"` // "vault" or "ucode"
	Repo        string `json:"repo"`
	Path        string `json:"path"`
	Type        string `json:"type"` // emulator, firmware, tool, manual, config
}

// Catalogue provides access to the device catalogue
type Catalogue struct {
	basePath    string
	vaultsPath  string
}

// New creates a new Catalogue instance
func New(basePath string) *Catalogue {
	if basePath == "" {
		basePath = guessBasePath()
	}
	vaultsPath := guessVaultsPath()
	return &Catalogue{basePath: basePath, vaultsPath: vaultsPath}
}

// ListDevices returns all devices in the catalogue
func (c *Catalogue) ListDevices() ([]Device, error) {
	var devices []Device

	// 1. Scan Vault device catalogue (documentation layer)
	vaultDevices, _ := c.scanVaultDevices()
	devices = append(devices, vaultDevices...)

	// 2. Scan uCode repos (implementation layer)
	codeDevices, _ := c.scanCodeRepos()
	devices = append(devices, codeDevices...)

	return devices, nil
}

// scanVaultDevices scans the Vault device catalogue directories
func (c *Catalogue) scanVaultDevices() ([]Device, error) {
	var devices []Device

	if c.vaultsPath == "" {
		return devices, nil
	}

	// Scan each level: user, shared, public
	for _, level := range []string{"user", "shared", "public"} {
		deviceDir := filepath.Join(c.vaultsPath, "devices", level)
		if _, err := os.Stat(deviceDir); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(deviceDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !strings.HasPrefix(entry.Name(), ".") {
				devices = append(devices, Device{
					Name:        entry.Name(),
					Source:      "vault",
					Repo:        fmt.Sprintf("devices/%s", level),
					Path:        filepath.Join(deviceDir, entry.Name()),
					Type:        classifyVaultDevice(entry.Name()),
				})
			}
		}
	}

	return devices, nil
}

// scanCodeRepos scans the uCode code repositories
func (c *Catalogue) scanCodeRepos() ([]Device, error) {
	var devices []Device

	for _, repo := range []string{"uCode1", "uCode2", "uCode3", "uCode4"} {
		repoPath := filepath.Join(c.basePath, repo)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(repoPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				devices = append(devices, Device{
					Name:   entry.Name(),
					Source: "ucode",
					Repo:   repo,
					Path:   filepath.Join(repoPath, entry.Name()),
					Type:   classifyCodeDevice(entry.Name()),
				})
			}
		}
	}

	return devices, nil
}

// FindDevice searches for a device by name across all sources
func (c *Catalogue) FindDevice(name string) (*Device, error) {
	devices, err := c.ListDevices()
	if err != nil {
		return nil, err
	}

	nameLower := strings.ToLower(name)
	for _, d := range devices {
		if strings.Contains(strings.ToLower(d.Name), nameLower) {
			return &d, nil
		}
	}

	return nil, fmt.Errorf("device '%s' not found in catalogue", name)
}

func classifyVaultDevice(name string) string {
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "manual") || strings.Contains(nameLower, "guide") || strings.Contains(nameLower, ".md") {
		return "manual"
	}
	if strings.Contains(nameLower, "config") || strings.Contains(nameLower, ".yaml") || strings.Contains(nameLower, ".json") {
		return "config"
	}
	if strings.Contains(nameLower, "spec") || strings.Contains(nameLower, "datasheet") {
		return "spec"
	}
	return "documentation"
}

func classifyCodeDevice(name string) string {
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "emu") || strings.Contains(nameLower, "emulator") {
		return "emulator"
	}
	if strings.Contains(nameLower, "firm") || strings.Contains(nameLower, "firmware") {
		return "firmware"
	}
	return "tool"
}

func guessBasePath() string {
	candidates := []string{
		"/Users/fredbook/Code",
		os.Getenv("HOME") + "/Code",
		os.Getenv("HOME") + "/uDos",
	}
	for _, p := range candidates {
		if _, err := os.Stat(filepath.Join(p, "uCode1")); err == nil {
			return p
		}
	}
	return os.Getenv("HOME") + "/Code"
}

func guessVaultsPath() string {
	// Try iCloud Drive vault location first (macOS)
	icloudVault := os.Getenv("HOME") + "/Library/Mobile Documents/com~apple~CloudDocs/Vaults"
	if _, err := os.Stat(icloudVault); err == nil {
		return icloudVault
	}

	// Fallback to local ~/Vaults
	localVault := os.Getenv("HOME") + "/Vaults"
	if _, err := os.Stat(localVault); err == nil {
		return localVault
	}

	return ""
}
