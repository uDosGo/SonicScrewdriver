// Package driver provides driver tracking and version management.
// It maintains a catalogue of device drivers, their versions, and
// compatibility information across different operating systems.
package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DriverType classifies the type of driver
type DriverType string

const (
	TypeUSB     DriverType = "usb"     // USB device driver (CDC, DFU, etc.)
	TypeUART    DriverType = "uart"    // UART/serial driver (CP210x, CH340, FTDI)
	TypeJTAG    DriverType = "jtag"    // Debug probe driver (ST-Link, J-Link)
	TypeNetwork DriverType = "network" // Network interface driver
	TypeStorage DriverType = "storage" // Storage controller driver
	TypeGPIO    DriverType = "gpio"    // GPIO/expansion driver
	TypeWiFi    DriverType = "wifi"    // Wireless driver
	TypeBT      DriverType = "bt"      // Bluetooth driver
	TypeAudio   DriverType = "audio"   // Audio driver
	TypeVideo   DriverType = "video"   // Video/display driver
)

// Platform represents a target operating system
type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformMacOS   Platform = "macos"
	PlatformWindows Platform = "windows"
	PlatformAll     Platform = "all"
)

// Driver represents a device driver with version tracking
type Driver struct {
	Name         string            `json:"name"`
	Type         DriverType        `json:"type"`
	Version      string            `json:"version"`
	Installed    bool              `json:"installed"`
	InstallPath  string            `json:"install_path,omitempty"`
	Platforms    []Platform        `json:"platforms"`
	Devices      []string          `json:"devices"`       // Compatible device names
	Vendor       string            `json:"vendor"`        // Manufacturer/vendor
	Description  string            `json:"description"`
	URL          string            `json:"url"`           // Download URL
	Checksum     string            `json:"checksum"`      // SHA256
	InstalledAt  string            `json:"installed_at,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// DriverManager manages driver installation and tracking
type DriverManager struct {
	dbPath    string
	drivers   map[string]*Driver
}

// NewManager creates a new DriverManager
func NewManager() *DriverManager {
	home, _ := os.UserHomeDir()
	return &DriverManager{
		dbPath:  filepath.Join(home, ".sonic", "drivers.json"),
		drivers: make(map[string]*Driver),
	}
}

// Load reads the driver database from disk
func (m *DriverManager) Load() error {
	if _, err := os.Stat(m.dbPath); os.IsNotExist(err) {
		return nil // Fresh start
	}

	data, err := os.ReadFile(m.dbPath)
	if err != nil {
		return fmt.Errorf("failed to read driver db: %w", err)
	}

	var drivers []*Driver
	if err := json.Unmarshal(data, &drivers); err != nil {
		return fmt.Errorf("failed to parse driver db: %w", err)
	}

	for _, d := range drivers {
		m.drivers[d.Name] = d
	}
	return nil
}

// Save writes the driver database to disk
func (m *DriverManager) Save() error {
	drivers := make([]*Driver, 0, len(m.drivers))
	for _, d := range m.drivers {
		drivers = append(drivers, d)
	}

	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].Name < drivers[j].Name
	})

	data, err := json.MarshalIndent(drivers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal drivers: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(m.dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create driver db dir: %w", err)
	}

	return os.WriteFile(m.dbPath, data, 0644)
}

// Register adds a new driver to the database
func (m *DriverManager) Register(driver *Driver) error {
	if driver.Name == "" {
		return fmt.Errorf("driver name is required")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if existing, ok := m.drivers[driver.Name]; ok {
		driver.InstalledAt = existing.InstalledAt
		driver.UpdatedAt = now
	} else {
		driver.InstalledAt = now
		driver.UpdatedAt = now
	}

	m.drivers[driver.Name] = driver
	return m.Save()
}

// Unregister removes a driver from the database
func (m *DriverManager) Unregister(name string) error {
	if _, ok := m.drivers[name]; !ok {
		return fmt.Errorf("driver '%s' not found", name)
	}
	delete(m.drivers, name)
	return m.Save()
}

// Get returns a driver by name
func (m *DriverManager) Get(name string) (*Driver, error) {
	driver, ok := m.drivers[name]
	if !ok {
		return nil, fmt.Errorf("driver '%s' not found", name)
	}
	return driver, nil
}

// List returns all registered drivers, optionally filtered by type
func (m *DriverManager) List(driverType DriverType) []*Driver {
	var result []*Driver
	for _, d := range m.drivers {
		if driverType == "" || d.Type == driverType {
			result = append(result, d)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// ListByDevice returns all drivers compatible with a given device
func (m *DriverManager) ListByDevice(deviceName string) []*Driver {
	var result []*Driver
	deviceLower := strings.ToLower(deviceName)

	for _, d := range m.drivers {
		for _, dev := range d.Devices {
			if strings.Contains(strings.ToLower(dev), deviceLower) {
				result = append(result, d)
				break
			}
		}
	}

	return result
}

// CheckInstalled checks if a driver is installed on the system
func (m *DriverManager) CheckInstalled(name string) (bool, error) {
	driver, ok := m.drivers[name]
	if !ok {
		return false, fmt.Errorf("driver '%s' not found", name)
	}

	// Check if the driver binary/tool exists
	if driver.InstallPath != "" {
		if _, err := os.Stat(driver.InstallPath); os.IsNotExist(err) {
			driver.Installed = false
			m.Save()
			return false, nil
		}
	}

	driver.Installed = true
	return true, nil
}

// GetVersion returns the installed version of a driver
func (m *DriverManager) GetVersion(name string) (string, error) {
	driver, ok := m.drivers[name]
	if !ok {
		return "", fmt.Errorf("driver '%s' not found", name)
	}

	installed, err := m.CheckInstalled(name)
	if err != nil {
		return "", err
	}

	if !installed {
		return "", fmt.Errorf("driver '%s' is not installed", name)
	}

	return driver.Version, nil
}

// DetectInstalledDrivers scans the system for known drivers and updates status
func (m *DriverManager) DetectInstalledDrivers() ([]string, error) {
	var found []string

	// Common driver paths and tools to check
	checks := map[string]string{
		"dfu-util":    "/usr/local/bin/dfu-util",
		"openocd":     "/usr/local/bin/openocd",
		"st-flash":    "/usr/local/bin/st-flash",
		"esptool":     "/usr/local/bin/esptool.py",
		"avrdude":     "/usr/local/bin/avrdude",
		"JLinkExe":    "/usr/local/bin/JLinkExe",
		"bossac":      "/usr/local/bin/bossac",
		"uf2conv":     "/usr/local/bin/uf2conv",
		"stm32flash":  "/usr/local/bin/stm32flash",
		"balena-etcher": "/usr/local/bin/balena-etcher",
	}

	for name, path := range checks {
		if _, err := os.Stat(path); err == nil {
			if driver, ok := m.drivers[name]; ok {
				driver.Installed = true
				driver.InstallPath = path
			}
			found = append(found, name)
		}
	}

	if len(found) > 0 {
		m.Save()
	}

	return found, nil
}

// GetDefaultDrivers returns the built-in driver catalogue
func GetDefaultDrivers() []*Driver {
	return []*Driver{
		{
			Name: "dfu-util", Type: TypeUSB,
			Version: "0.11", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"STM32", "RP2040", "nRF52", "iCE40"},
			Vendor:      "DFU-Utils",
			Description: "Device Firmware Update utility for USB DFU standard",
			URL:         "http://dfu-util.sourceforge.net/",
		},
		{
			Name: "openocd", Type: TypeJTAG,
			Version: "0.12", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"STM32", "nRF52", "LPC", "SAMD", "RISC-V"},
			Vendor:      "OpenOCD",
			Description: "Open On-Chip Debugger - JTAG/SWD debugging and flashing",
			URL:         "https://openocd.org/",
		},
		{
			Name: "st-flash", Type: TypeJTAG,
			Version: "1.7", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"STM32F0", "STM32F1", "STM32F3", "STM32F4", "STM32F7", "STM32H7", "STM32G0", "STM32G4", "STM32L0", "STM32L4", "STM32WB"},
			Vendor:      "STMicroelectronics",
			Description: "ST-Link flash tool for STM32 microcontrollers",
			URL:         "https://github.com/stlink-org/stlink",
		},
		{
			Name: "esptool", Type: TypeUART,
			Version: "4.7", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"ESP32", "ESP32-S2", "ESP32-S3", "ESP32-C3", "ESP32-C6", "ESP8266"},
			Vendor:      "Espressif",
			Description: "ESP32/ESP8266 serial bootloader tool",
			URL:         "https://github.com/espressif/esptool",
		},
		{
			Name: "avrdude", Type: TypeUART,
			Version: "7.3", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"ATmega328P", "ATmega2560", "ATtiny85", "ATtiny13", "ATmega32U4"},
			Vendor:      "AVRDUDE",
			Description: "AVR Downloader/UploaDEr - AVR microcontroller programmer",
			URL:         "https://github.com/avrdudes/avrdude",
		},
		{
			Name: "JLinkExe", Type: TypeJTAG,
			Version: "V7.96", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"ARM7", "ARM9", "ARM11", "Cortex-M", "Cortex-A", "Cortex-R", "RISC-V"},
			Vendor:      "SEGGER",
			Description: "J-Link Commander - JTAG/SWD debug probe for ARM/RISC-V",
			URL:         "https://www.segger.com/products/debug-probes/j-link/",
		},
		{
			Name: "bossac", Type: TypeUSB,
			Version: "1.9", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"SAMD21", "SAMD51", "SAME51"},
			Vendor:      "Arduino",
			Description: "BOSSAC - SAM BA bootloader flash tool for Arduino Zero/MKR",
			URL:         "https://github.com/arduino/arduino-core-avr",
		},
		{
			Name: "uf2conv", Type: TypeUSB,
			Version: "1.0", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"RP2040", "RP2350", "nRF52840", "nRF5340", "SAMD51"},
			Vendor:      "Microsoft",
			Description: "UF2 file converter and flasher for UF2 bootloader devices",
			URL:         "https://github.com/microsoft/uf2",
		},
		{
			Name: "stm32flash", Type: TypeUART,
			Version: "0.7", Platforms: []Platform{PlatformLinux, PlatformMacOS},
			Devices:     []string{"STM32F0", "STM32F1", "STM32F3", "STM32F4"},
			Vendor:      "STMicroelectronics",
			Description: "STM32 serial bootloader flash tool",
			URL:         "https://sourceforge.net/projects/stm32flash/",
		},
		{
			Name: "CH340", Type: TypeUART,
			Version: "1.5", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"CH340G", "CH340C", "CH341A"},
			Vendor:      "WCH",
			Description: "CH340/CH341 USB-to-UART bridge driver",
			URL:         "https://www.wch.cn/download/CH341SER_MAC_ZIP.html",
		},
		{
			Name: "CP210x", Type: TypeUART,
			Version: "6.0", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"CP2102", "CP2104", "CP2105", "CP2108", "CP2109"},
			Vendor:      "Silicon Labs",
			Description: "CP210x USB-to-UART bridge driver",
			URL:         "https://www.silabs.com/developers/usb-to-uart-bridge-vcp-drivers",
		},
		{
			Name: "FTDI VCP", Type: TypeUART,
			Version: "2.5", Platforms: []Platform{PlatformLinux, PlatformMacOS, PlatformWindows},
			Devices:     []string{"FT232R", "FT230X", "FT231X", "FT2232H", "FT4232H"},
			Vendor:      "FTDI",
			Description: "FTDI Virtual COM Port driver for USB-to-UART adapters",
			URL:         "https://ftdichip.com/drivers/vcp-drivers/",
		},
	}
}

// Ensure json is used
var _ = json.Marshal
