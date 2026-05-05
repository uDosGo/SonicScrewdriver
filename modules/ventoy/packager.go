package ventoy

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/OkAgentDigital/universal/library"
)

// Packager handles creating Ventoy-compatible installer bundles
type Packager struct {
	SourceDir string
	OutputDir string
	Config    *library.Manager
}

// NewPackager creates a new packager instance
func NewPackager(source, output string, config *library.Manager) *Packager {
	return &Packager{
		SourceDir: source,
		OutputDir: output,
		Config:    config,
	}
}

// CreateBundle creates a .she (Sonic Home Edition) bundle for USB deployment
func (p *Packager) CreateBundle(name string) (string, error) {
	if p.SourceDir == "" {
		return "", fmt.Errorf("source directory not specified")
	}

	if p.OutputDir == "" {
		p.OutputDir = "."
	}

	// Ensure output directory exists
	if err := os.MkdirAll(p.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove .she extension if already present
	cleanName := strings.TrimSuffix(name, ".she")
	bundlePath := filepath.Join(p.OutputDir, cleanName+".she")
	file, err := os.Create(bundlePath)
	if err != nil {
		return "", fmt.Errorf("failed to create bundle: %w", err)
	}
	defer file.Close()

	// Use gzip + tar for the bundle
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk the source directory and add files to the bundle
	err = filepath.Walk(p.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create header
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		// Update the name to be relative to source directory
		relPath, err := filepath.Rel(p.SourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to create bundle: %w", err)
	}

	return bundlePath, nil
}

// CreateUSBImage creates a bootable USB image using Ventoy
func (p *Packager) CreateUSBImage(bundlePath, usbDevice string) error {
	// This would integrate with Ventoy's tools
	// For now, we'll document the expected workflow
	
	fmt.Printf("USB packaging workflow:\n")
	fmt.Printf("1. Format USB device: %s\n", usbDevice)
	fmt.Printf("2. Install Ventoy: ventoy -i %s\n", usbDevice)
	fmt.Printf("3. Copy bundle: cp %s /Volumes/VENTOY/\n", bundlePath)
	fmt.Printf("4. Configure Ventoy: copy ventoy.json to USB\n")
	fmt.Printf("5. Eject USB: diskutil eject %s\n", usbDevice)
	
	return fmt.Errorf("USB packaging requires Ventoy tools - see documentation")
}

// ValidateBundle checks if a bundle is valid
func (p *Packager) ValidateBundle(bundlePath string) error {
	file, err := os.Open(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to open bundle: %w", err)
	}
	defer file.Close()

	// Check if it's a gzip file
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("invalid gzip format: %w", err)
	}
	defer gzipReader.Close()

	// Check if it's a tar archive
	tarReader := tar.NewReader(gzipReader)
	_, err = tarReader.Next()
	if err != nil {
		return fmt.Errorf("invalid tar format: %w", err)
	}

	return nil
}

// GetBundleInfo extracts information from a bundle
func GetBundleInfo(bundlePath string) (map[string]string, error) {
	// This would parse the bundle and extract metadata
	// For now, return basic info
	
	stat, err := os.Stat(bundlePath)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"path":      bundlePath,
		"size":      fmt.Sprintf("%d", stat.Size()),
		"type":      "Sonic Home Edition Bundle",
		"extension": strings.TrimPrefix(filepath.Ext(bundlePath), "."),
	}, nil
}