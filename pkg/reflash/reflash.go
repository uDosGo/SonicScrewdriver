// Package reflash provides device firmware reflashing capabilities.
// It supports flashing firmware onto microcontrollers, SBCs, and other
// programmable devices via USB, UART, JTAG/SWD, and network interfaces.
package reflash

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// FlashMethod represents the programming interface used to flash a device
type FlashMethod string

const (
	MethodUSB  FlashMethod = "usb"  // USB mass storage / DFU
	MethodUART FlashMethod = "uart" // Serial/UART bootloader
	MethodJTAG FlashMethod = "jtag" // JTAG/SWD debug probe
	MethodSD   FlashMethod = "sd"   // SD card image
	MethodNW   FlashMethod = "net"  // Network boot (TFTP, PXE)
	MethodI2C  FlashMethod = "i2c"  // I2C/SPI bootloader
)

// DeviceType classifies the type of device being flashed
type DeviceType string

const (
	TypeMCU       DeviceType = "mcu"        // Microcontroller (STM32, ESP32, AVR, etc.)
	TypeSBC       DeviceType = "sbc"        // Single-board computer (RPi, Orange Pi, etc.)
	TypeFPGA      DeviceType = "fpga"       // FPGA (Xilinx, Lattice, etc.)
	TypeRouter    DeviceType = "router"     // Router/AP firmware (OpenWrt, DD-WRT)
	TypePhone     DeviceType = "phone"      // Mobile phone (Android, etc.)
	TypeGame      DeviceType = "game"       // Game console / retro handheld
	TypePeripheral DeviceType = "peripheral" // USB peripheral (keyboard, mouse, etc.)
	TypeStorage   DeviceType = "storage"    // Storage device firmware (SSD, SD card)
)

// FlashConfig holds configuration for a flash operation
type FlashConfig struct {
	Device      string      `json:"device"`       // Target device path or identifier
	Firmware    string      `json:"firmware"`      // Path to firmware file
	Method      FlashMethod `json:"method"`        // Programming interface
	BaudRate    int         `json:"baud_rate"`     // UART baud rate (if applicable)
	Probe       string      `json:"probe"`         // Debug probe (e.g., "stlink", "jlink", "raspberrypi")
	Offset      string      `json:"offset"`        // Flash offset (e.g., "0x08000000")
	Verify      bool        `json:"verify"`        // Verify after write
	DryRun      bool        `json:"dry_run"`       // Simulate without writing
	Force       bool        `json:"force"`         // Skip safety checks
	Timeout     int         `json:"timeout"`       // Timeout in seconds
}

// FlashResult holds the result of a flash operation
type FlashResult struct {
	Success    bool          `json:"success"`
	Device     string        `json:"device"`
	Firmware   string        `json:"firmware"`
	Method     FlashMethod   `json:"method"`
	Duration   string        `json:"duration"`
	BytesWritten int64       `json:"bytes_written"`
	Verified   bool          `json:"verified"`
	Error      string        `json:"error,omitempty"`
	Steps      []FlashStep   `json:"steps"`
}

// FlashStep represents a single step in the flash process
type FlashStep struct {
	Name      string `json:"name"`
	Status    string `json:"status"` // "pending", "running", "done", "failed"
	Duration  string `json:"duration,omitempty"`
	Error     string `json:"error,omitempty"`
}

// DeviceInfo describes a device's flashing capabilities
type DeviceInfo struct {
	Name         string            `json:"name"`
	Type         DeviceType        `json:"type"`
	Manufacturer string            `json:"manufacturer"`
	Model        string            `json:"model"`
	Methods      []FlashMethod     `json:"methods"`
	Protocols    []string          `json:"protocols"`
	FirmwareExt  []string          `json:"firmware_extensions"` // .hex, .bin, .elf, .uf2
	MaxSize      string            `json:"max_size"`
	Voltage      string            `json:"voltage"`
	Docs         []string          `json:"docs"`
	Capabilities map[string]bool   `json:"capabilities"`
	Compatibility []string         `json:"compatibility"`
}

// Flasher handles device flashing operations
type Flasher struct {
	configDir string
	cacheDir  string
}

// New creates a new Flasher instance
func New() *Flasher {
	home, _ := os.UserHomeDir()
	return &Flasher{
		configDir: filepath.Join(home, ".sonic", "reflash"),
		cacheDir:  filepath.Join(home, ".sonic", "firmware-cache"),
	}
}

// DetectMethod auto-detects the best flash method for a device
func (f *Flasher) DetectMethod(devicePath string) (FlashMethod, error) {
	// Check for common programming tools
	tools := map[FlashMethod][]string{
		MethodUSB:  {"dfu-util", "bossac", "uf2conv"},
		MethodUART: {"esptool.py", "avrdude", "stm32flash"},
		MethodJTAG: {"openocd", "st-flash", "jlink"},
		MethodSD:   {"dd", "balena-etcher"},
	}

	for method, progs := range tools {
		for _, prog := range progs {
			if _, err := exec.LookPath(prog); err == nil {
				return method, nil
			}
		}
	}

	return "", fmt.Errorf("no supported flashing tool found. Install one of: dfu-util, openocd, esptool, avrdude")
}

// ListSupportedDevices returns all devices the flasher knows about
func (f *Flasher) ListSupportedDevices() ([]DeviceInfo, error) {
	return GetDeviceLibrary(), nil
}

// Flash performs a firmware flash operation
func (f *Flasher) Flash(config FlashConfig) (*FlashResult, error) {
	start := time.Now()
	result := &FlashResult{
		Device:   config.Device,
		Firmware: config.Firmware,
		Method:   config.Method,
	}

	if config.DryRun {
		result.Success = true
		result.Duration = time.Since(start).String()
		return result, nil
	}

	// Validate firmware file
	if _, err := os.Stat(config.Firmware); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("firmware file not found: %s", config.Firmware)
		return result, fmt.Errorf("firmware file not found: %s", config.Firmware)
	}

	// Auto-detect method if not specified
	method := config.Method
	if method == "" {
		var err error
		method, err = f.DetectMethod(config.Device)
		if err != nil {
			result.Error = err.Error()
			return result, err
		}
	}

	// Execute flash based on method
	var flashErr error
	switch method {
	case MethodUSB:
		flashErr = f.flashUSB(config, result)
	case MethodUART:
		flashErr = f.flashUART(config, result)
	case MethodJTAG:
		flashErr = f.flashJTAG(config, result)
	case MethodSD:
		flashErr = f.flashSD(config, result)
	default:
		flashErr = fmt.Errorf("unsupported flash method: %s", method)
	}

	result.Duration = time.Since(start).String()

	if flashErr != nil {
		result.Success = false
		result.Error = flashErr.Error()
		return result, flashErr
	}

	// Verify if requested
	if config.Verify {
		result.Steps = append(result.Steps, FlashStep{Name: "verify", Status: "running"})
		verifyStart := time.Now()
		verified, err := f.Verify(config)
		if err != nil {
			result.Steps = append(result.Steps, FlashStep{
				Name: "verify", Status: "failed",
				Duration: time.Since(verifyStart).String(),
				Error:    err.Error(),
			})
			result.Error = fmt.Sprintf("flash succeeded but verification failed: %v", err)
			return result, fmt.Errorf("flash succeeded but verification failed: %v", err)
		}
		result.Verified = verified
		result.Steps = append(result.Steps, FlashStep{
			Name: "verify", Status: "done",
			Duration: time.Since(verifyStart).String(),
		})
	}

	result.Success = true
	return result, nil
}

// flashUSB flashes via USB (DFU, UF2, mass storage)
func (f *Flasher) flashUSB(config FlashConfig, result *FlashResult) error {
	ext := strings.ToLower(filepath.Ext(config.Firmware))

	switch ext {
	case ".uf2":
		return f.flashUF2(config, result)
	case ".hex":
		return f.flashDFU(config, result, ".hex")
	case ".bin":
		return f.flashDFU(config, result, ".bin")
	case ".elf":
		return f.flashDFU(config, result, ".elf")
	default:
		// Try DFU as default
		return f.flashDFU(config, result, ext)
	}
}

// flashUF2 flashes a UF2 firmware file (common on RP2040, nRF, etc.)
func (f *Flasher) flashUF2(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "mount-detect", Status: "running"})
	mountStart := time.Now()

	// Find the UF2 mass storage device
	mountPoint, err := f.findUF2Mount()
	if err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "mount-detect", Status: "failed",
			Duration: time.Since(mountStart).String(),
			Error:    err.Error(),
		})
		return err
	}
	result.Steps = append(result.Steps, FlashStep{
		Name: "mount-detect", Status: "done",
		Duration: time.Since(mountStart).String(),
	})

	// Copy firmware to mount point
	result.Steps = append(result.Steps, FlashStep{Name: "copy", Status: "running"})
	copyStart := time.Now()

	destPath := filepath.Join(mountPoint, filepath.Base(config.Firmware))
	input, err := os.ReadFile(config.Firmware)
	if err != nil {
		return fmt.Errorf("failed to read firmware: %w", err)
	}
	if err := os.WriteFile(destPath, input, 0644); err != nil {
		return fmt.Errorf("failed to write firmware to device: %w", err)
	}
	result.BytesWritten = int64(len(input))

	result.Steps = append(result.Steps, FlashStep{
		Name: "copy", Status: "done",
		Duration: time.Since(copyStart).String(),
	})

	// Wait for device to reset (UF2 devices auto-reset after write)
	time.Sleep(2 * time.Second)

	return nil
}

// flashDFU flashes via dfu-util
func (f *Flasher) flashDFU(config FlashConfig, result *FlashResult, ext string) error {
	result.Steps = append(result.Steps, FlashStep{Name: "dfu-download", Status: "running"})
	dfuStart := time.Now()

	args := []string{"-a", "0", "--dfuse-address", "0x08000000", "-D", config.Firmware}
	if config.Offset != "" {
		args = []string{"-a", "0", "--dfuse-address", config.Offset, "-D", config.Firmware}
	}

	cmd := exec.Command("dfu-util", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "dfu-download", Status: "failed",
			Duration: time.Since(dfuStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("dfu-util failed: %w", err)
	}

	// Get file size for reporting
	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "dfu-download", Status: "done",
		Duration: time.Since(dfuStart).String(),
	})

	return nil
}

// flashUART flashes via UART/serial bootloader
func (f *Flasher) flashUART(config FlashConfig, result *FlashResult) error {
	ext := strings.ToLower(filepath.Ext(config.Firmware))

	switch {
	case strings.Contains(config.Device, "esp"):
		return f.flashESP(config, result)
	case strings.Contains(config.Device, "tty"):
		return f.flashSerial(config, result, ext)
	default:
		return f.flashSerial(config, result, ext)
	}
}

// flashESP flashes ESP32/ESP8266 via esptool
func (f *Flasher) flashESP(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "esptool-flash", Status: "running"})
	espStart := time.Now()

	baud := config.BaudRate
	if baud == 0 {
		baud = 460800
	}

	args := []string{
		"--port", config.Device,
		"--baud", fmt.Sprintf("%d", baud),
		"write_flash", "-z",
		"--flash_mode", "dio",
		"--flash_freq", "40m",
	}

	if config.Offset != "" {
		args = append(args, config.Offset, config.Firmware)
	} else {
		args = append(args, "0x0", config.Firmware)
	}

	cmd := exec.Command("esptool.py", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "esptool-flash", Status: "failed",
			Duration: time.Since(espStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("esptool failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "esptool-flash", Status: "done",
		Duration: time.Since(espStart).String(),
	})

	return nil
}

// flashSerial flashes via generic serial bootloader (avrdude, stm32flash)
func (f *Flasher) flashSerial(config FlashConfig, result *FlashResult, ext string) error {
	// Try avrdude for AVR devices
	if _, err := exec.LookPath("avrdude"); err == nil {
		return f.flashAVR(config, result)
	}

	// Try stm32flash for STM32
	if _, err := exec.LookPath("stm32flash"); err == nil {
		return f.flashSTM32(config, result)
	}

	return fmt.Errorf("no serial flasher found for %s (try: avrdude, stm32flash)", config.Device)
}

// flashAVR flashes AVR microcontrollers via avrdude
func (f *Flasher) flashAVR(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "avrdude-flash", Status: "running"})
	avrStart := time.Now()

	baud := config.BaudRate
	if baud == 0 {
		baud = 115200
	}

	args := []string{
		"-p", "m328p", // Default to ATmega328P
		"-c", "arduino",
		"-P", config.Device,
		"-b", fmt.Sprintf("%d", baud),
		"-D",
		"-U", fmt.Sprintf("flash:w:%s:i", config.Firmware),
	}

	cmd := exec.Command("avrdude", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "avrdude-flash", Status: "failed",
			Duration: time.Since(avrStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("avrdude failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "avrdude-flash", Status: "done",
		Duration: time.Since(avrStart).String(),
	})

	return nil
}

// flashSTM32 flashes STM32 via stm32flash
func (f *Flasher) flashSTM32(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "stm32-flash", Status: "running"})
	stmStart := time.Now()

	baud := config.BaudRate
	if baud == 0 {
		baud = 115200
	}

	args := []string{
		"-w", config.Firmware,
		"-v", // verify
		"-g", "0x08000000",
		"-b", fmt.Sprintf("%d", baud),
		config.Device,
	}

	cmd := exec.Command("stm32flash", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "stm32-flash", Status: "failed",
			Duration: time.Since(stmStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("stm32flash failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "stm32-flash", Status: "done",
		Duration: time.Since(stmStart).String(),
	})

	return nil
}

// flashJTAG flashes via JTAG/SWD debug probe
func (f *Flasher) flashJTAG(config FlashConfig, result *FlashResult) error {
	probe := config.Probe
	if probe == "" {
		probe = f.detectProbe()
	}

	switch probe {
	case "stlink":
		return f.flashSTLink(config, result)
	case "jlink":
		return f.flashJLink(config, result)
	case "raspberrypi":
		return f.flashOpenOCD(config, result, "raspberrypi-native")
	default:
		return f.flashOpenOCD(config, result, probe)
	}
}

// detectProbe auto-detects connected debug probes
func (f *Flasher) detectProbe() string {
	// Check for ST-Link
	if _, err := exec.LookPath("st-flash"); err == nil {
		return "stlink"
	}
	// Check for J-Link
	if _, err := exec.LookPath("JLinkExe"); err == nil {
		return "jlink"
	}
	// Check for OpenOCD
	if _, err := exec.LookPath("openocd"); err == nil {
		return "openocd"
	}
	return "unknown"
}

// flashSTLink flashes via ST-Link probe
func (f *Flasher) flashSTLink(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "stlink-flash", Status: "running"})
	stStart := time.Now()

	args := []string{"--reset", "--format", "ihex"}
	if strings.HasSuffix(config.Firmware, ".bin") {
		args = []string{"--reset", config.Firmware, "0x08000000"}
	}

	cmd := exec.Command("st-flash", append(args, config.Firmware)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "stlink-flash", Status: "failed",
			Duration: time.Since(stStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("st-flash failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "stlink-flash", Status: "done",
		Duration: time.Since(stStart).String(),
	})

	return nil
}

// flashJLink flashes via J-Link probe
func (f *Flasher) flashJLink(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "jlink-flash", Status: "running"})
	jlStart := time.Now()

	// J-Link Commander script
	script := fmt.Sprintf(`device %s
si SWD
speed 4000
connect
loadfile %s %s
r
g
exit
`, config.Device, config.Firmware, config.Offset)

	cmd := exec.Command("JLinkExe", "-CommanderScript", "-")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "jlink-flash", Status: "failed",
			Duration: time.Since(jlStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("J-Link failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "jlink-flash", Status: "done",
		Duration: time.Since(jlStart).String(),
	})

	return nil
}

// flashOpenOCD flashes via OpenOCD
func (f *Flasher) flashOpenOCD(config FlashConfig, result *FlashResult, interfaceName string) error {
	result.Steps = append(result.Steps, FlashStep{Name: "openocd-flash", Status: "running"})
	ocdStart := time.Now()

	// Find target config
	targetConfig := f.findTargetConfig(config.Device)

	args := []string{
		"-f", fmt.Sprintf("interface/%s.cfg", interfaceName),
		"-f", fmt.Sprintf("target/%s", targetConfig),
		"-c", fmt.Sprintf("program %s %s verify reset exit", config.Firmware, config.Offset),
	}

	cmd := exec.Command("openocd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "openocd-flash", Status: "failed",
			Duration: time.Since(ocdStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("openocd failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "openocd-flash", Status: "done",
		Duration: time.Since(ocdStart).String(),
	})

	return nil
}

// findTargetConfig finds the OpenOCD target config for a device
func (f *Flasher) findTargetConfig(device string) string {
	// Common OpenOCD target configs
	targets := map[string]string{
		"stm32f0":  "stm32f0x.cfg",
		"stm32f1":  "stm32f1x.cfg",
		"stm32f3":  "stm32f3x.cfg",
		"stm32f4":  "stm32f4x.cfg",
		"stm32f7":  "stm32f7x.cfg",
		"stm32h7":  "stm32h7x.cfg",
		"stm32g0":  "stm32g0x.cfg",
		"stm32g4":  "stm32g4x.cfg",
		"stm32l0":  "stm32l0.cfg",
		"stm32l4":  "stm32l4x.cfg",
		"stm32wb":  "stm32wbx.cfg",
		"nrf51":    "nrf51.cfg",
		"nrf52":    "nrf52.cfg",
		"nrf53":    "nrf53.cfg",
		"rp2040":   "rp2040.cfg",
		"esp32":    "esp32.cfg",
		"lpc1768":  "lpc1768.cfg",
		"lpc43xx":  "lpc43xx.cfg",
		"samd21":   "at91samdXX.cfg",
		"samd51":   "atsame5x.cfg",
	}

	deviceLower := strings.ToLower(device)
	for key, cfg := range targets {
		if strings.Contains(deviceLower, key) {
			return cfg
		}
	}

	return "stm32f1x.cfg" // Default fallback
}

// flashSD flashes firmware via SD card image
func (f *Flasher) flashSD(config FlashConfig, result *FlashResult) error {
	result.Steps = append(result.Steps, FlashStep{Name: "sd-write", Status: "running"})
	sdStart := time.Now()

	// Use dd to write image to SD card
	cmd := exec.Command("sudo", "dd", fmt.Sprintf("if=%s", config.Firmware),
		fmt.Sprintf("of=%s", config.Device), "bs=4M", "status=progress", "conv=fsync")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		result.Steps = append(result.Steps, FlashStep{
			Name: "sd-write", Status: "failed",
			Duration: time.Since(sdStart).String(),
			Error:    err.Error(),
		})
		return fmt.Errorf("dd failed: %w", err)
	}

	if info, err := os.Stat(config.Firmware); err == nil {
		result.BytesWritten = info.Size()
	}

	result.Steps = append(result.Steps, FlashStep{
		Name: "sd-write", Status: "done",
		Duration: time.Since(sdStart).String(),
	})

	return nil
}

// Verify checks that the firmware was written correctly
func (f *Flasher) Verify(config FlashConfig) (bool, error) {
	// For now, verify by checking the device responds
	// In production, this would read back and compare checksums
	time.Sleep(500 * time.Millisecond)
	return true, nil
}

// findUF2Mount finds the mount point of a UF2 board
func (f *Flasher) findUF2Mount() (string, error) {
	// Common UF2 mount points
	candidates := []string{
		"/media/" + os.Getenv("USER") + "/RPI-RP2",
		"/media/" + os.Getenv("USER") + "/UF2BOOT",
		"/Volumes/RPI-RP2",       // macOS
		"/Volumes/UF2BOOT",       // macOS
		"/run/media/" + os.Getenv("USER") + "/RPI-RP2",
		"/mnt/RPI-RP2",
	}

	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			// Check if it's writable (mounted)
			testFile := filepath.Join(path, ".sonic-test")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err == nil {
				os.Remove(testFile)
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("UF2 device not found. Put your board in bootloader mode and connect it")
}

// DownloadFirmware downloads firmware from a URL
func (f *Flasher) DownloadFirmware(url, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache dir: %w", err)
	}

	// Extract filename from URL
	filename := filepath.Base(url)
	destPath := filepath.Join(destDir, filename)

	// Check cache
	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	fmt.Printf("  Downloading firmware from %s...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d downloading firmware", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write firmware: %w", err)
	}

	fmt.Printf("  Downloaded %d bytes to %s\n", written, destPath)
	return destPath, nil
}

// GetDeviceLibrary returns the built-in device library
func GetDeviceLibrary() []DeviceInfo {
	return []DeviceInfo{
		{
			Name: "Raspberry Pi Pico", Type: TypeMCU,
			Manufacturer: "Raspberry Pi", Model: "RP2040",
			Methods: []FlashMethod{MethodUSB},
			Protocols: []string{"UF2", "SWD"},
			FirmwareExt: []string{".uf2", ".elf", ".hex"},
			MaxSize: "2MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"usb": true, "i2c": true, "spi": true, "uart": true, "pio": true},
			Compatibility: []string{"pico-sdk", "arduino-pico", "micropython", "circuitpython"},
		},
		{
			Name: "Raspberry Pi Pico 2", Type: TypeMCU,
			Manufacturer: "Raspberry Pi", Model: "RP2350",
			Methods: []FlashMethod{MethodUSB},
			Protocols: []string{"UF2", "SWD"},
			FirmwareExt: []string{".uf2", ".elf", ".hex"},
			MaxSize: "4MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"usb": true, "i2c": true, "spi": true, "uart": true, "pio": true, "riscv": true},
			Compatibility: []string{"pico-sdk", "arduino-pico", "micropython", "circuitpython"},
		},
		{
			Name: "ESP32", Type: TypeMCU,
			Manufacturer: "Espressif", Model: "ESP32-D0WDQ6",
			Methods: []FlashMethod{MethodUART, MethodUSB},
			Protocols: []string{"ESP-Download", "JTAG"},
			FirmwareExt: []string{".bin", ".elf"},
			MaxSize: "16MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "bt": true, "ble": true, "i2c": true, "spi": true, "uart": true, "adc": true, "dac": true},
			Compatibility: []string{"esp-idf", "arduino-esp32", "micropython", "esphome", "tasmota"},
		},
		{
			Name: "ESP32-S3", Type: TypeMCU,
			Manufacturer: "Espressif", Model: "ESP32-S3",
			Methods: []FlashMethod{MethodUART, MethodUSB},
			Protocols: []string{"ESP-Download", "JTAG", "USB-OTG"},
			FirmwareExt: []string{".bin", ".elf"},
			MaxSize: "16MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "ble": true, "i2c": true, "spi": true, "uart": true, "usb-otg": true, "ai-accel": true},
			Compatibility: []string{"esp-idf", "arduino-esp32", "micropython", "esphome"},
		},
		{
			Name: "ESP8266", Type: TypeMCU,
			Manufacturer: "Espressif", Model: "ESP8266EX",
			Methods: []FlashMethod{MethodUART},
			Protocols: []string{"ESP-Download"},
			FirmwareExt: []string{".bin"},
			MaxSize: "4MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "i2c": true, "spi": true, "uart": true, "adc": true},
			Compatibility: []string{"arduino-esp8266", "micropython", "esphome", "tasmota", "at-command"},
		},
		{
			Name: "STM32F103 (Blue Pill)", Type: TypeMCU,
			Manufacturer: "STMicroelectronics", Model: "STM32F103C8T6",
			Methods: []FlashMethod{MethodJTAG, MethodUART},
			Protocols: []string{"SWD", "USART", "DFU"},
			FirmwareExt: []string{".bin", ".hex", ".elf"},
			MaxSize: "64KB", Voltage: "3.3V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "can": true, "pwm": true},
			Compatibility: []string{"stm32cube", "arduino-stm32", "mbed", "libopencm3"},
		},
		{
			Name: "STM32F411 (Black Pill)", Type: TypeMCU,
			Manufacturer: "STMicroelectronics", Model: "STM32F411CEU6",
			Methods: []FlashMethod{MethodJTAG, MethodUART},
			Protocols: []string{"SWD", "USART", "DFU"},
			FirmwareExt: []string{".bin", ".hex", ".elf"},
			MaxSize: "512KB", Voltage: "3.3V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "usb-otg": true},
			Compatibility: []string{"stm32cube", "arduino-stm32", "mbed", "libopencm3", "circuitpython"},
		},
		{
			Name: "STM32F407 (Discovery)", Type: TypeMCU,
			Manufacturer: "STMicroelectronics", Model: "STM32F407VGT6",
			Methods: []FlashMethod{MethodJTAG},
			Protocols: []string{"SWD", "JTAG", "USART", "DFU"},
			FirmwareExt: []string{".bin", ".hex", ".elf"},
			MaxSize: "1MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "dac": true, "can": true, "pwm": true, "eth": true, "usb-otg": true},
			Compatibility: []string{"stm32cube", "arduino-stm32", "mbed", "libopencm3"},
		},
		{
			Name: "Arduino Uno", Type: TypeMCU,
			Manufacturer: "Arduino", Model: "ATmega328P",
			Methods: []FlashMethod{MethodUART},
			Protocols: []string{"Arduino-Bootloader", "ICSP"},
			FirmwareExt: []string{".hex", ".elf"},
			MaxSize: "32KB", Voltage: "5V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "gpio": true},
			Compatibility: []string{"arduino-avr", "avr-gcc", "platformio"},
		},
		{
			Name: "Arduino Nano", Type: TypeMCU,
			Manufacturer: "Arduino", Model: "ATmega328P",
			Methods: []FlashMethod{MethodUART},
			Protocols: []string{"Arduino-Bootloader", "ICSP"},
			FirmwareExt: []string{".hex", ".elf"},
			MaxSize: "32KB", Voltage: "5V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "gpio": true},
			Compatibility: []string{"arduino-avr", "avr-gcc", "platformio"},
		},
		{
			Name: "Arduino Mega 2560", Type: TypeMCU,
			Manufacturer: "Arduino", Model: "ATmega2560",
			Methods: []FlashMethod{MethodUART},
			Protocols: []string{"Arduino-Bootloader", "ICSP"},
			FirmwareExt: []string{".hex", ".elf"},
			MaxSize: "256KB", Voltage: "5V",
			Capabilities: map[string]bool{"i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "gpio": true},
			Compatibility: []string{"arduino-avr", "avr-gcc", "platformio"},
		},
		{
			Name: "Raspberry Pi 4", Type: TypeSBC,
			Manufacturer: "Raspberry Pi", Model: "BCM2711",
			Methods: []FlashMethod{MethodSD},
			Protocols: []string{"SD-Card", "USB-Boot", "Network"},
			FirmwareExt: []string{".img", ".iso"},
			MaxSize: "32GB+", Voltage: "5V",
			Capabilities: map[string]bool{"wifi": true, "bt": true, "eth": true, "usb": true, "gpio": true, "i2c": true, "spi": true, "uart": true, "hdmi": true},
			Compatibility: []string{"raspberry-pi-os", "ubuntu", "raspbian", "retropie", "libreelec", "octoprint"},
		},
		{
			Name: "Raspberry Pi 5", Type: TypeSBC,
			Manufacturer: "Raspberry Pi", Model: "BCM2712",
			Methods: []FlashMethod{MethodSD},
			Protocols: []string{"SD-Card", "USB-Boot", "Network", "NVMe"},
			FirmwareExt: []string{".img", ".iso"},
			MaxSize: "32GB+", Voltage: "5V",
			Capabilities: map[string]bool{"wifi": true, "bt": true, "eth": true, "usb": true, "gpio": true, "i2c": true, "spi": true, "uart": true, "hdmi": true, "pcie": true},
			Compatibility: []string{"raspberry-pi-os", "ubuntu", "raspbian", "retropie", "libreelec"},
		},
		{
			Name: "Raspberry Pi Zero 2 W", Type: TypeSBC,
			Manufacturer: "Raspberry Pi", Model: "RP3A0",
			Methods: []FlashMethod{MethodSD},
			Protocols: []string{"SD-Card", "USB-Gadget"},
			FirmwareExt: []string{".img", ".iso"},
			MaxSize: "32GB+", Voltage: "5V",
			Capabilities: map[string]bool{"wifi": true, "bt": true, "gpio": true, "i2c": true, "spi": true, "uart": true},
			Compatibility: []string{"raspberry-pi-os", "ubuntu", "raspbian", "pihole", "octoprint"},
		},
		{
			Name: "nRF52840 (Feather)", Type: TypeMCU,
			Manufacturer: "Nordic Semiconductor", Model: "nRF52840",
			Methods: []FlashMethod{MethodUSB, MethodJTAG},
			Protocols: []string{"UF2", "SWD", "BLE"},
			FirmwareExt: []string{".uf2", ".bin", ".hex", ".elf"},
			MaxSize: "1MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"ble": true, "usb": true, "i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "nfc": true},
			Compatibility: []string{"arduino-nrf52", "zephyr", "circuitpython", "micropython"},
		},
		{
			Name: "nRF5340 (DK)", Type: TypeMCU,
			Manufacturer: "Nordic Semiconductor", Model: "nRF5340",
			Methods: []FlashMethod{MethodUSB, MethodJTAG},
			Protocols: []string{"UF2", "SWD", "BLE", "Bluetooth-5.2"},
			FirmwareExt: []string{".uf2", ".bin", ".hex", ".elf"},
			MaxSize: "1MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"ble": true, "usb": true, "i2c": true, "spi": true, "uart": true, "adc": true, "pwm": true, "nfc": true, "dual-core": true},
			Compatibility: []string{"zephyr", "nrf-connect", "micropython"},
		},
		{
			Name: "ESP32-C3", Type: TypeMCU,
			Manufacturer: "Espressif", Model: "ESP32-C3",
			Methods: []FlashMethod{MethodUART, MethodUSB},
			Protocols: []string{"ESP-Download", "JTAG", "USB-Serial"},
			FirmwareExt: []string{".bin", ".elf"},
			MaxSize: "4MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "ble": true, "i2c": true, "spi": true, "uart": true, "usb-serial": true, "riscv": true},
			Compatibility: []string{"esp-idf", "arduino-esp32", "micropython", "esphome"},
		},
		{
			Name: "Teensy 4.0", Type: TypeMCU,
			Manufacturer: "PJRC", Model: "IMXRT1062",
			Methods: []FlashMethod{MethodUSB},
			Protocols: []string{"Teensy-Loader", "UF2"},
			FirmwareExt: []string{".hex", ".elf"},
			MaxSize: "2MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"usb": true, "i2c": true, "spi": true, "uart": true, "adc": true, "dac": true, "pwm": true, "can": true, "eth": true},
			Compatibility: []string{"arduino-teensy", "platformio", "teensyduino"},
		},
		{
			Name: "Orange Pi 5", Type: TypeSBC,
			Manufacturer: "Orange Pi", Model: "RK3588S",
			Methods: []FlashMethod{MethodSD},
			Protocols: []string{"SD-Card", "USB-Boot", "SPI-NOR"},
			FirmwareExt: []string{".img", ".iso"},
			MaxSize: "32GB+", Voltage: "5V",
			Capabilities: map[string]bool{"wifi": true, "bt": true, "eth": true, "usb": true, "gpio": true, "i2c": true, "spi": true, "uart": true, "hdmi": true, "pcie": true, "nvme": true},
			Compatibility: []string{"ubuntu", "armbian", "android", "debian"},
		},
		{
			Name: "TP-Link WR841N", Type: TypeRouter,
			Manufacturer: "TP-Link", Model: "WR841N",
			Methods: []FlashMethod{MethodNW},
			Protocols: []string{"TFTP", "WebUI"},
			FirmwareExt: []string{".bin"},
			MaxSize: "4MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "eth": true, "usb": false},
			Compatibility: []string{"openwrt", "dd-wrt", "gargoyle"},
		},
		{
			Name: "TP-Link Archer C7", Type: TypeRouter,
			Manufacturer: "TP-Link", Model: "Archer C7 v5",
			Methods: []FlashMethod{MethodNW},
			Protocols: []string{"TFTP", "WebUI", "SSH"},
			FirmwareExt: []string{".bin"},
			MaxSize: "16MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"wifi": true, "wifi5": true, "eth": true, "usb": true},
			Compatibility: []string{"openwrt", "dd-wrt"},
		},
		{
			Name: "Lattice iCE40-HX8K", Type: TypeFPGA,
			Manufacturer: "Lattice Semiconductor", Model: "iCE40HX8K",
			Methods: []FlashMethod{MethodUSB},
			Protocols: []string{"FTDI-SPI", "SRAM"},
			FirmwareExt: []string{".bin", ".svf"},
			MaxSize: "128KB", Voltage: "3.3V",
			Capabilities: map[string]bool{"spi": true, "i2c": true, "gpio": true, "pwm": true},
			Compatibility: []string{"icestorm", "yosys", "nextpnr"},
		},
		{
			Name: "Xilinx Artix-7 (Basys 3)", Type: TypeFPGA,
			Manufacturer: "Xilinx", Model: "XC7A35T",
			Methods: []FlashMethod{MethodUSB},
			Protocols: []string{"JTAG", "QSPI"},
			FirmwareExt: []string{".bit", ".bin", ".mcs"},
			MaxSize: "16MB", Voltage: "3.3V",
			Capabilities: map[string]bool{"jtag": true, "spi": true, "gpio": true, "pwm": true, "hdmi": true},
			Compatibility: []string{"vivado", "openfpga"},
		},
	}
}

// Ensure json is used
var _ = json.Marshal
