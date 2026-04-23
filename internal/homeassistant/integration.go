package homeassistant

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// HAConfig represents Home Assistant configuration
type HAConfig struct {
	BaseURL      string `json:"base_url"`
	APIToken     string `json:"api_token"`
	EmbedURL     string `json:"embed_url"`
	KioskMode    bool   `json:"kiosk_mode"`
	RefreshRate  int    `json:"refresh_rate"`
	LastChecked  string `json:"last_checked"`
	Status       string `json:"status"`
}

// HAIntegration manages Home Assistant integration
type HAIntegration struct {
	config       HAConfig
	httpClient   *http.Client
	iframeHTML  string
	kioskConfig string
}

// NewHAIntegration creates a new HA integration instance
func NewHAIntegration(baseURL, apiToken string) *HAIntegration {
	return &HAIntegration{
		config: HAConfig{
			BaseURL:     baseURL,
			APIToken:    apiToken,
			EmbedURL:    baseURL + "/lovelace/default_view",
			KioskMode:   true,
			RefreshRate: 60,
			Status:      "not_configured",
		},
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Configure sets up HA integration
func (ha *HAIntegration) Configure() error {
	// Validate base URL
	if _, err := url.ParseRequestURI(ha.config.BaseURL); err != nil {
		return fmt.Errorf("invalid Home Assistant URL: %v", err)
	}

	// Test API connection
	if err := ha.testAPIConnection(); err != nil {
		return fmt.Errorf("failed to connect to Home Assistant: %v", err)
	}

	// Generate iframe HTML
	if err := ha.generateIframeHTML(); err != nil {
		return fmt.Errorf("failed to generate iframe HTML: %v", err)
	}

	// Generate kiosk configuration
	if err := ha.generateKioskConfig(); err != nil {
		return fmt.Errorf("failed to generate kiosk config: %v", err)
	}

	ha.config.Status = "configured"
	ha.config.LastChecked = time.Now().Format("2006-01-02 15:04:05")

	return nil
}

// testAPIConnection tests if HA API is accessible
func (ha *HAIntegration) testAPIConnection() error {
	apiURL := ha.config.BaseURL + "/api/"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+ha.config.APIToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ha.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("HA API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// generateIframeHTML generates the iframe embed code
func (ha *HAIntegration) generateIframeHTML() error {
	iframeHTML := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Home Assistant Embed</title>
    <style>
        body, html {
            margin: 0;
            padding: 0;
            height: 100%;
            overflow: hidden;
        }
        .ha-container {
            width: 100%;
            height: 100%;
            border: none;
        }
    </style>
</head>
<body>
    <iframe
        class="ha-container"
        src="%s"
        frameborder="0"
        allowfullscreen>
    </iframe>
    <script>
        // Auto-refresh if kiosk mode enabled
        %s
    </script>
</body>
</html>`

	refreshScript := ""
	if ha.config.KioskMode && ha.config.RefreshRate > 0 {
		refreshScript = fmt.Sprintf("setInterval(function() { location.reload(); }, %d * 60 * 1000);", ha.config.RefreshRate)
	}

	ha.iframeHTML = fmt.Sprintf(iframeHTML, ha.config.EmbedURL, refreshScript)

	return nil
}

// generateKioskConfig generates HA kiosk mode configuration
func (ha *HAIntegration) generateKioskConfig() error {
	kioskConfig := `{
  "hide_header": true,
  "hide_sidebar": true,
  "hide_overflow": true,
  "fullscreen": true,
  "background": "#E8E8E8",
  "text_color": "#111111",
  "border_color": "#222222",
  "button_color": "#3A7BD5",
  "refresh_interval": %d
}`

	ha.kioskConfig = fmt.Sprintf(kioskConfig, ha.config.RefreshRate)

	return nil
}

// GetIframeHTML returns the generated iframe HTML
func (ha *HAIntegration) GetIframeHTML() string {
	return ha.iframeHTML
}

// GetKioskConfig returns the kiosk configuration
func (ha *HAIntegration) GetKioskConfig() string {
	return ha.kioskConfig
}

// GetStatus returns the current status
func (ha *HAIntegration) GetStatus() HAConfig {
	return ha.config
}

// SaveConfig saves the HA configuration
func (ha *HAIntegration) SaveConfig(filePath string) error {
	configDir := filepath.Dir(filePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ha.config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// LoadConfig loads HA configuration from file
func (ha *HAIntegration) LoadConfig(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var config HAConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	ha.config = config
	return nil
}

// ApplyKioskMode applies kiosk mode configuration to HA
func (ha *HAIntegration) ApplyKioskMode() error {
	// This would typically use the HA API to apply kiosk mode settings
	// For now, we'll document the manual steps
	log.Println("⚠️  Kiosk mode application requires manual setup in Home Assistant")
	log.Println("   Add the following to your configuration.yaml:")
	log.Println(ha.kioskConfig)
	
	return nil
}

// GetEmbedURL returns the URL to embed
func (ha *HAIntegration) GetEmbedURL() string {
	return ha.config.EmbedURL
}

// SetEmbedURL sets the embed URL
func (ha *HAIntegration) SetEmbedURL(embedURL string) {
	ha.config.EmbedURL = embedURL
}

// SetKioskMode enables or disables kiosk mode
func (ha *HAIntegration) SetKioskMode(enabled bool) {
	ha.config.KioskMode = enabled
}

// SetRefreshRate sets the refresh rate in minutes
func (ha *HAIntegration) SetRefreshRate(rate int) {
	ha.config.RefreshRate = rate
}

// CheckHAStatus checks if Home Assistant is running
func (ha *HAIntegration) CheckHAStatus() (bool, error) {
	apiURL := ha.config.BaseURL + "/api/"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+ha.config.APIToken)
	
	resp, err := ha.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetHAInfo retrieves information about the HA instance
func (ha *HAIntegration) GetHAInfo() (map[string]interface{}, error) {
	apiURL := ha.config.BaseURL + "/api/config"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ha.config.APIToken)
	
	resp, err := ha.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HA API returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// InstallHAAddon installs a Home Assistant addon
func (ha *HAIntegration) InstallHAAddon(addonName string) error {
	log.Printf("⚠️  Addon installation is not yet implemented in this version")
	log.Printf("   You can manually install %s through the HA UI", addonName)
	return nil
}

// GetHAAddons lists installed Home Assistant addons
func (ha *HAIntegration) GetHAAddons() ([]string, error) {
	log.Println("⚠️  Addon listing is not yet implemented in this version")
	return []string{}, nil
}

// GetHAVersion returns the Home Assistant version
func (ha *HAIntegration) GetHAVersion() (string, error) {
	info, err := ha.GetHAInfo()
	if err != nil {
		return "", err
	}

	if version, ok := info["version"].(string); ok {
		return version, nil
	}
	
	return "unknown", errors.New("version not found in HA config")
}

// GenerateEmbedFile generates an HTML file for embedding
func (ha *HAIntegration) GenerateEmbedFile(outputPath string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// Generate the iframe HTML
	if err := ha.generateIframeHTML(); err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(outputPath, []byte(ha.iframeHTML), 0644)
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig(baseURL, apiToken string) HAConfig {
	return HAConfig{
		BaseURL:     baseURL,
		APIToken:    apiToken,
		EmbedURL:    baseURL + "/lovelace/default_view",
		KioskMode:   true,
		RefreshRate: 60,
		Status:      "not_configured",
	}
}

// CheckDependencies checks if required dependencies are installed
func CheckDependencies() ([]string, []string) {
	required := []string{"curl", "jq", "sqlite3"}

	
	var installed []string
	var missing []string
	
	for _, cmd := range required {
		if _, err := exec.LookPath(cmd); err == nil {
			installed = append(installed, cmd)
		} else {
			missing = append(missing, cmd)
		}
	}
	
	// Check if Home Assistant is accessible
	// This is a placeholder - in real implementation, we'd check the actual HA instance
	
	return installed, missing
}

// InstallDependencies attempts to install missing dependencies
func InstallDependencies(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	
	log.Printf("Installing missing dependencies: %v\n", missing)
	
	// This is platform-specific and would need to be implemented
	// for different operating systems
	log.Println("⚠️  Automatic dependency installation is not implemented")
	log.Println("   Please install the following packages manually:")
	for _, pkg := range missing {
		log.Printf("   - %s\n", pkg)
	}
	
	return nil
}