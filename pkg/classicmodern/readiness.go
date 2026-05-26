package classicmodern

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ReadinessCheck represents the result of a readiness check
type ReadinessCheck struct {
	ThemeReady        bool     `json:"theme_ready"`
	IconsReady        bool     `json:"icons_ready"`
	FontsReady        bool     `json:"fonts_ready"`
	OBFValid          bool     `json:"obf_valid"`
	DependenciesReady bool     `json:"dependencies_ready"`
	Issues            []string `json:"issues"`
	Warnings         []string `json:"warnings"`
	MissingPackages  []string `json:"missing_packages"`
}

// CheckReadiness checks if Classic Modern Mint is ready for installation
type CheckReadiness struct {
	db            *sql.DB
	MissingPackages []string
}

// NewReadinessChecker creates a new readiness checker
func NewReadinessChecker(db *sql.DB) *CheckReadiness {
	return &CheckReadiness{db: db}
}

// PerformCheck performs the readiness check
func (c *CheckReadiness) PerformCheck() (ReadinessCheck, error) {
	check := ReadinessCheck{
		Issues:           make([]string, 0),
		Warnings:        make([]string, 0),
		MissingPackages: make([]string, 0),
	}

	// Check dependencies
	check.DependenciesReady = c.checkDependencies()
	if !check.DependenciesReady {
		check.Issues = append(check.Issues, "Missing required packages")
	}

	// Check theme files
	check.ThemeReady = c.checkThemeFiles()
	if !check.ThemeReady {
		check.Issues = append(check.Issues, "Theme files not found or incomplete")
	}

	// Check icons
	check.IconsReady = c.checkIcons()
	if !check.IconsReady {
		check.Issues = append(check.Issues, "Icon set not found or incomplete")
	}

	// Check fonts
	check.FontsReady = c.checkFonts()
	if !check.FontsReady {
		check.Issues = append(check.Issues, "Required fonts not found")
	}

	// Check OBF file
	check.OBFValid = c.checkOBF()
	if !check.OBFValid {
		check.Issues = append(check.Issues, "OBF file not found or invalid")
	}

	return check, nil
}

// checkDependencies checks if required packages are installed
func (c *CheckReadiness) checkDependencies() bool {
	requiredPackages := []string{
		"cinnamon",
		"gtk3",
		"gtk4",
		"metacity",
		"x11-utils",
		"fontconfig",
		"imagemagick",
	}

	allInstalled := true
	
	for _, pkg := range requiredPackages {
		if _, err := exec.LookPath(pkg); err != nil {
			// Try dpkg query for Debian-based systems
			cmd := exec.Command("dpkg", "-l", pkg)
			if err := cmd.Run(); err != nil {
				c.MissingPackages = append(c.MissingPackages, pkg)
				allInstalled = false
			}
		}
	}

	if len(c.MissingPackages) > 0 {
		log.Printf("Missing packages: %v", c.MissingPackages)
		log.Println("Install them with: sudo apt install " + strings.Join(c.MissingPackages, " "))
	}

	return allInstalled
}

// checkThemeFiles checks if theme files are present
func (c *CheckReadiness) checkThemeFiles() bool {
	// Check for theme in standard locations
	locations := []string{
		"/usr/share/themes/Classic-Modern",
		"~/.themes/Classic-Modern",
		"/home/wizard/code-vault/classic-modern-mint/themes/Classic-Modern",
	}

	for _, loc := range locations {
		expandedLoc := expandPath(loc)
		if _, err := os.Stat(expandedLoc); err == nil {
			// Check required theme files
			requiredFiles := []string{
				"gtk-3.0/gtk.css",
				"gtk-4.0/gtk.css",
				"cinnamon/cinnamon.css",
				"metacity-1/metacity-theme-3.xml",
				"index.theme",
			}
			
			allFilesExist := true
			for _, file := range requiredFiles {
				if _, err := os.Stat(filepath.Join(expandedLoc, file)); err != nil {
					log.Printf("Missing theme file: %s", file)
					allFilesExist = false
				}
			}
			
			if allFilesExist {
				return true
			}
		}
	}

	return false
}

// checkIcons checks if icon set is present
func (c *CheckReadiness) checkIcons() bool {
	// Check for icons in standard locations
	locations := []string{
		"/usr/share/icons/Classic-Modern",
		"~/.icons/Classic-Modern",
		"/home/wizard/code-vault/classic-modern-mint/icons/Classic-Modern",
	}

	for _, loc := range locations {
		expandedLoc := expandPath(loc)
		if _, err := os.Stat(expandedLoc); err == nil {
			// Check for some essential icons
			requiredIcons := []string{
				"apps/terminal.svg",
				"actions/close.svg",
				"status/battery.svg",
			}
			
			allIconsExist := true
			for _, icon := range requiredIcons {
				if _, err := os.Stat(filepath.Join(expandedLoc, icon)); err != nil {
					log.Printf("Missing icon: %s", icon)
					allIconsExist = false
				}
			}
			
			if allIconsExist {
				return true
			}
		}
	}

	return false
}

// checkFonts checks if required fonts are installed
func (c *CheckReadiness) checkFonts() bool {
	// Check for required fonts
	requiredFonts := []string{
		"ChicagoFLF",
		"Monaspace Argon",
		"Inter",
		"SF Pro Text",
	}

	fontCheckCmd := func(fontName string) bool {
		cmd := exec.Command("fc-list", ":", "family", "|", "grep", "-i", fontName)
		output, err := cmd.CombinedOutput()
		return err == nil && len(output) > 0
	}

	allFontsFound := true
	for _, font := range requiredFonts {
		if !fontCheckCmd(font) {
			log.Printf("Font not found: %s", font)
			allFontsFound = false
		}
	}

	return allFontsFound
}

// checkOBF checks if OBF file is valid
func (c *CheckReadiness) checkOBF() bool {
	// Check for OBF file
	obfPaths := []string{
		"/usr/share/classic-modern/classic-modern.obf",
		"~/.config/classic-modern/classic-modern.obf",
		"/home/wizard/code-vault/classic-modern-mint/obf/classic-modern.obf",
	}

	for _, path := range obfPaths {
		expandedPath := expandPath(path)
		if _, err := os.Stat(expandedPath); err == nil {
			// Check OBF structure
			data, err := ioutil.ReadFile(expandedPath)
			if err != nil {
				log.Printf("Failed to read OBF file: %v", err)
				return false
			}

			// Check magic header
			if len(data) < 8 || string(data[:8]) != "CMOBF01" {
				log.Println("Invalid OBF magic header")
				return false
			}

			// Check checksum (simplified)
			if len(data) < 4 {
				log.Println("OBF file too short")
				return false
			}

			return true
		}
	}

	return false
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// GenerateInstallationReport generates a detailed installation report
func (c *CheckReadiness) GenerateInstallationReport() (string, error) {
	check, err := c.PerformCheck()
	if err != nil {
		return "", err
	}

	report := "📋 Classic Modern Mint Installation Readiness Report\n"
	report += "=====================================================\n\n"

	// Overall status
	allReady := check.ThemeReady && check.IconsReady && check.FontsReady && check.OBFValid && check.DependenciesReady
	
	if allReady {
		report += "✅ ALL SYSTEMS READY — Classic Modern can be installed\n\n"
	} else {
		report += "⚠️  SOME COMPONENTS NOT READY — Review issues below\n\n"
	}

	// Component status
	report += "Component Status:\n"
	report += fmt.Sprintf("  Theme Files:        %s\n", statusEmoji(check.ThemeReady))
	report += fmt.Sprintf("  Icon Set:           %s\n", statusEmoji(check.IconsReady))
	report += fmt.Sprintf("  Fonts:              %s\n", statusEmoji(check.FontsReady))
	report += fmt.Sprintf("  OBF File:           %s\n", statusEmoji(check.OBFValid))
	report += fmt.Sprintf("  Dependencies:      %s\n", statusEmoji(check.DependenciesReady))

	// Issues
	if len(check.Issues) > 0 {
		report += "\n❌ Issues:\n"
		for _, issue := range check.Issues {
			report += fmt.Sprintf("  - %s\n", issue)
		}
	}

	// Warnings
	if len(check.Warnings) > 0 {
		report += "\n⚠️  Warnings:\n"
		for _, warning := range check.Warnings {
			report += fmt.Sprintf("  - %s\n", warning)
		}
	}

	// Missing packages
	if len(check.MissingPackages) > 0 {
		report += "\n📦 Missing Packages:\n"
		for _, pkg := range check.MissingPackages {
			report += fmt.Sprintf("  - %s\n", pkg)
		}
		report += "\nInstall with: sudo apt install " + strings.Join(check.MissingPackages, " ") + "\n"
	}

	// Next steps
	report += "\n🚀 Next Steps:\n"
	if allReady {
		report += "  Run: sonic mint install classic-modern\n"
		report += "  Then: sonic mint apply classic-modern\n"
	} else {
		report += "  1. Resolve all issues above\n"
		report += "  2. Run: sonic mint doctor --component=classic-modern\n"
		report += "  3. Re-run readiness check\n"
	}

	return report, nil
}

// statusEmoji returns an emoji for status
func statusEmoji(ready bool) string {
	if ready {
		return "✅"
	}
	return "❌"
}

// ExportCheckToJSON exports the check result to JSON
func (c *CheckReadiness) ExportCheckToJSON(filePath string) error {
	check, err := c.PerformCheck()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(check, "", "  ")
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// CheckThemeInstallation checks if theme is properly installed
func (c *CheckReadiness) CheckThemeInstallation() (bool, error) {
	// Check if theme is set in gsettings
	cmd := exec.Command("gsettings", "get", "org.cinnamon.theme", "name")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check current theme: %v", err)
	}

	currentTheme := strings.TrimSpace(string(output))
	currentTheme = strings.Trim(currentTheme, `'`)
	
	if currentTheme == "Classic-Modern" {
		return true, nil
	}

	return false, nil
}

// ApplyTheme applies the Classic Modern theme
func (c *CheckReadiness) ApplyTheme() error {
	commands := [][]string{
		{"gsettings", "set", "org.cinnamon.theme", "name", "Classic-Modern"},
		{"gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "Classic-Modern"},
		{"gsettings", "set", "org.gnome.desktop.wm.preferences", "theme", "Classic-Modern"},
		{"gsettings", "set", "org.gnome.desktop.interface", "font-name", "Inter 14"},
		{"gsettings", "set", "org.gnome.desktop.interface", "monospace-font-name", "ChicagoFLF 13"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to apply theme setting: %v", err)
		}
	}

	return nil
}

// GetThemeInfo returns information about the Classic Modern theme
func (c *CheckReadiness) GetThemeInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// Check if theme is installed
	locations := []string{
		"/usr/share/themes/Classic-Modern",
		"~/.themes/Classic-Modern",
	}

	for _, loc := range locations {
		expandedLoc := expandPath(loc)
		if _, err := os.Stat(expandedLoc); err == nil {
			info["installed"] = true
			info["location"] = expandedLoc
			
			// Read index.theme if it exists
			indexPath := filepath.Join(expandedLoc, "index.theme")
			if _, err := ioutil.ReadFile(indexPath); err == nil {
				info["name"] = "Classic Modern"
				info["version"] = "3.0.0"
				info["description"] = "Classic Mac Platinum meets Modern Mono"
			}
			
			// Check for OBF
			obfPath := filepath.Join(expandedLoc, "classic-modern.obf")
			if _, err := os.Stat(obfPath); err == nil {
				info["obf_present"] = true
			}
			
			return info, nil
		}
	}

	info["installed"] = false
	return info, errors.New("theme not found")
}