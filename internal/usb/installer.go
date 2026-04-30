package usb

import (
	"fmt"
	"os"

	"github.com/sonic-family/sonic-screwdriver/internal/disk"
	"github.com/sonic-family/sonic-screwdriver/internal/iso"
)

// InstallConfig holds the configuration for a USB install
type InstallConfig struct {
	Device     string      // /dev/sdX
	DistroName string      // ubuntu, mint, classicmodern
	Layout     disk.Layout // partition layout
	DryRun     bool        // don't actually write
	Force      bool        // skip confirmation
}

// InstallResult holds the result of an install operation
type InstallResult struct {
	Success  bool   `json:"success"`
	Device   string `json:"device"`
	Distro   string `json:"distro"`
	ISOPath  string `json:"iso_path"`
	Error    string `json:"error,omitempty"`
}

// PrepareDisk wipes, partitions, and formats a disk for OS installation
func PrepareDisk(config InstallConfig) error {
	fmt.Println("=== Preparing Disk ===")
	fmt.Printf("Device: %s | Layout: %s\n", config.Device, config.Layout.Name)

	if config.DryRun {
		fmt.Println("  [DRY RUN] Would wipe, partition, and format")
		return nil
	}

	// 1. Wipe existing partition table
	if err := disk.WipeDevice(config.Device); err != nil {
		return fmt.Errorf("wipe failed: %w", err)
	}

	// 2. Create partitions
	if err := disk.CreatePartitions(config.Device, config.Layout); err != nil {
		return fmt.Errorf("partition failed: %w", err)
	}

	// 3. Format partitions
	if err := disk.FormatAllPartitions(config.Device, config.Layout); err != nil {
		return fmt.Errorf("format failed: %w", err)
	}

	fmt.Println("✅ Disk prepared successfully")
	return nil
}

// InstallDistro downloads and writes a distro ISO to a USB device
func InstallDistro(config InstallConfig) (*InstallResult, error) {
	result := &InstallResult{Device: config.Device, Distro: config.DistroName}

	// Get distro info
	distro, err := iso.GetDistro(config.DistroName)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	fmt.Println("=== Installing", distro.Name, distro.Version, "===")
	fmt.Printf("Target: %s\n", config.Device)

	if config.DryRun {
		fmt.Println("  [DRY RUN] Would download ISO and write to device")
		result.Success = true
		return result, nil
	}

	// 1. Download ISO
	fmt.Println("\n--- Step 1: Download ISO ---")
	progressCh := make(chan iso.DownloadStatus, 10)
	go func() {
		for status := range progressCh {
			if status.Complete {
				break
			}
			fmt.Printf("\r  Downloading: %.1f%% (%d/%d bytes)", status.Progress, status.Downloaded, status.TotalBytes)
		}
	}()

	isoPath, err := iso.Download(distro, progressCh)
	close(progressCh)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}
	fmt.Printf("\n  ✅ Downloaded: %s\n", isoPath)
	result.ISOPath = isoPath

	// 2. Write ISO to device
	fmt.Println("\n--- Step 2: Write to Device ---")
	if err := iso.WriteISOToDisk(isoPath, config.Device); err != nil {
		result.Error = err.Error()
		return result, err
	}

	result.Success = true
	fmt.Println("\n✅ Installation complete!")
	fmt.Printf("   Distro: %s %s\n", distro.Name, distro.Version)
	fmt.Printf("   Device: %s\n", config.Device)
	fmt.Printf("   ISO:    %s\n", isoPath)

	return result, nil
}

// FullInstall prepares the disk AND installs the distro
func FullInstall(config InstallConfig) (*InstallResult, error) {
	// 1. Prepare disk
	if err := PrepareDisk(config); err != nil {
		return &InstallResult{Device: config.Device, Distro: config.DistroName, Error: err.Error()}, err
	}

	// 2. Install distro
	return InstallDistro(config)
}

// ListUSBDevices lists available removable devices
func ListUSBDevices() ([]disk.Device, error) {
	devices, err := disk.DetectDevices(true)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		fmt.Println("No removable USB devices found")
		return nil, nil
	}
	fmt.Println("\nAvailable USB devices:")
	for _, d := range devices {
		fmt.Printf("  %s  %s  %s  %s\n", d.Path, d.Size, d.Model, d.Name)
	}
	return devices, nil
}

// GetLayout returns the partition layout for a distro
func GetLayout(distroName string) (disk.Layout, error) {
	switch distroName {
	case "ubuntu", "ubuntu-24.04":
		return disk.UbuntuLayout(), nil
	case "mint", "linuxmint", "linux-mint":
		return disk.LinuxMintLayout(), nil
	case "classicmodern", "classic-modern", "classic-modern-mint":
		return disk.ClassicModernLayout(), nil
	default:
		return disk.Layout{}, fmt.Errorf("unknown distro: %s", distroName)
	}
}

// Confirm asks for user confirmation before destructive operations
func Confirm(device string) bool {
	fmt.Printf("\n⚠️  WARNING: This will DESTROY ALL DATA on %s\n", device)
	fmt.Print("Are you sure? Type 'YES' to confirm: ")
	var response string
	fmt.Scanln(&response)
	return response == "YES"
}

// Ensure we reference os for compilation
var _ = os.Stderr
