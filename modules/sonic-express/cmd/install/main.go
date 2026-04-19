package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	usb := flag.Bool("usb", false, "auto-detect install bundle from USB")
	channel := flag.String("channel", "", "install channel")
	from := flag.String("from", "", "update server URL")
	dryRun := flag.Bool("dry-run", false, "print install plan only")
	flag.Parse()

	target := ""
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	if target == "" && !*usb && *from == "" {
		fmt.Fprintln(os.Stderr, "usage: install [bundle.she] [--usb] [--channel <name> --from <url>]")
		os.Exit(2)
	}

	mode := "bundle"
	switch {
	case *usb:
		mode = "usb"
		var err error
		target, err = detectUSBBundle()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to detect USB bundle: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("detected USB bundle: %s\n", target)
	case *from != "":
		mode = "channel"
	}

	if *dryRun {
		fmt.Printf("DRY RUN: Would install %s (mode=%s)\n", target, mode)
		return
	}

	fmt.Printf("installing: %s (mode=%s)\n", target, mode)
	// TODO: Implement actual installation logic
}

func detectUSBBundle() (string, error) {
	// Common USB mount points to check
	usbMountPoints := []string{
		"/Volumes",           // macOS
		"/media",             // Linux
		"/mnt",               // Linux/Windows
		"/run/media",         // Linux (user mounts)
	}

	// Look for .she files in USB mount points
	for _, mountPoint := range usbMountPoints {
		if _, err := os.Stat(mountPoint); os.IsNotExist(err) {
			continue
		}

		bundlePath, err := findSheFile(mountPoint)
		if err == nil {
			return bundlePath, nil
		}
	}

	// Also check common USB device names
	commonUSBNames := []string{
		"/Volumes/SONIC",
		"/Volumes/UDOS",
		"/Volumes/GAMES",
		"/media/usb0",
		"/media/usb1",
	}

	for _, usbPath := range commonUSBNames {
		if _, err := os.Stat(usbPath); os.IsNotExist(err) {
			continue
		}

		bundlePath, err := findSheFile(usbPath)
		if err == nil {
			return bundlePath, nil
		}
	}

	return "", fmt.Errorf("no .she bundle found on USB devices")
}

func findSheFile(root string) (string, error) {
	var sheFiles []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip directories we can't access
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".she") {
			sheFiles = append(sheFiles, path)
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to search %s: %w", root, err)
	}

	if len(sheFiles) == 0 {
		return "", fmt.Errorf("no .she files found in %s", root)
	}

	// Return the first .she file found
	return sheFiles[0], nil
}
