package iso

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Distro defines an OS distribution available for download
type Distro struct {
Name    string `json:"name"`
Version string `json:"version"`
URL     string `json:"url"`
Mirror  string `json:"mirror,omitempty"`
SHA256  string `json:"sha256"`
Size    string `json:"size"`
Arch    string `json:"arch"`
}

// DownloadStatus tracks download progress
type DownloadStatus struct {
Distro     Distro  `json:"distro"`
FilePath   string  `json:"file_path"`
TotalBytes int64   `json:"total_bytes"`
Downloaded int64   `json:"downloaded"`
Progress   float64 `json:"progress"`
SHA256OK   bool    `json:"sha256_ok"`
Complete   bool    `json:"complete"`
Error      string  `json:"error,omitempty"`
}

var (
Ubuntu2404 = Distro{
Name:    "ubuntu", Version: "24.04",
URL:    "https://releases.ubuntu.com/24.04/ubuntu-24.04.2-desktop-amd64.iso",
Mirror: "https://mirror.aarnet.edu.au/pub/ubuntu/releases/24.04/ubuntu-24.04.2-desktop-amd64.iso",
Size:   "~5.7GB", Arch: "amd64",
}
LinuxMint22 = Distro{
Name:    "linuxmint", Version: "22",
URL:    "https://mirrors.edge.kernel.org/linuxmint/stable/22/linuxmint-22-cinnamon-64bit.iso",
Mirror: "https://mirror.aarnet.edu.au/pub/linuxmint/stable/22/linuxmint-22-cinnamon-64bit.iso",
Size:   "~2.8GB", Arch: "amd64",
}
// Classic Modern Mint is a theme/config overlay on Linux Mint 21.3 Cinnamon
// It uses the Linux Mint ISO + post-install theming via .she bundle
ClassicModern = Distro{
Name:    "classicmodern", Version: "1.0",
URL:  LinuxMint22.URL,  // Uses Linux Mint 22 ISO as base
Mirror: LinuxMint22.Mirror,
Size: LinuxMint22.Size, Arch: "amd64",
}
)

func GetCacheDir() string {
home, _ := os.UserHomeDir()
return filepath.Join(home, ".sonic", "iso-cache")
}

func GetDistro(name string) (Distro, error) {
switch strings.ToLower(name) {
case "ubuntu", "ubuntu-24.04", "ubuntu2404":
return Ubuntu2404, nil
case "mint", "linuxmint", "linux-mint", "mint22":
return LinuxMint22, nil
case "classicmodern", "classic-modern", "classic-modern-mint":
d := ClassicModern
d.Name = "classicmodern"
d.Version = "1.0 (Linux Mint " + LinuxMint22.Version + " base)"
return d, nil
default:
return Distro{}, fmt.Errorf("unknown distro: %s (available: ubuntu, mint, classicmodern)", name)
}
}

func ListDistros() []Distro {
return []Distro{Ubuntu2404, LinuxMint22, ClassicModern}
}

// Download downloads an ISO to the cache directory
func Download(distro Distro, progressCh chan<- DownloadStatus) (string, error) {
cacheDir := GetCacheDir()
os.MkdirAll(cacheDir, 0755)

filename := fmt.Sprintf("%s-%s-%s.iso", distro.Name, distro.Version, distro.Arch)
filePath := filepath.Join(cacheDir, filename)

// Check cache
if _, err := os.Stat(filePath); err == nil {
if distro.SHA256 == "" || verifySHA256(filePath, distro.SHA256) {
if progressCh != nil {
progressCh <- DownloadStatus{Distro: distro, FilePath: filePath, Progress: 100, SHA256OK: true, Complete: true}
}
return filePath, nil
}
os.Remove(filePath)
}

// Try URLs
urls := []string{distro.URL}
if distro.Mirror != "" {
urls = append(urls, distro.Mirror)
}

var lastErr error
for _, url := range urls {
if url == "" {
continue
}
fmt.Printf("  Downloading %s from %s...\n", distro.Name, url)
err := downloadFile(url, filePath, distro, progressCh)
if err == nil {
lastErr = nil
break
}
fmt.Printf("  Mirror failed: %v\n", err)
lastErr = err
}
if lastErr != nil {
return "", fmt.Errorf("all mirrors failed: %w", lastErr)
}

// Verify
if distro.SHA256 != "" && !verifySHA256(filePath, distro.SHA256) {
os.Remove(filePath)
return "", fmt.Errorf("SHA256 verification failed")
}

if progressCh != nil {
progressCh <- DownloadStatus{Distro: distro, FilePath: filePath, Progress: 100, SHA256OK: true, Complete: true}
}
return filePath, nil
}

// WriteISOToDisk writes an ISO directly to a block device
func WriteISOToDisk(isoPath, devicePath string) error {
fmt.Printf("  Writing %s to %s (this may take several minutes)...\n", filepath.Base(isoPath), devicePath)
cmd := execWithSudo("dd", "if="+isoPath, "of="+devicePath, "bs=4M", "status=progress", "conv=fsync")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("dd failed: %w", err)
}
fmt.Printf("  ✅ ISO written to %s\n", devicePath)
return nil
}

func downloadFile(url, filePath string, distro Distro, progressCh chan<- DownloadStatus) error {
resp, err := http.Get(url)
if err != nil {
return err
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
return fmt.Errorf("HTTP %d", resp.StatusCode)
}

out, err := os.Create(filePath)
if err != nil {
return err
}
defer out.Close()

total := resp.ContentLength
downloaded := int64(0)
buf := make([]byte, 32*1024)
lastReport := time.Now()

for {
n, readErr := resp.Body.Read(buf)
if n > 0 {
out.Write(buf[:n])
downloaded += int64(n)
if progressCh != nil && time.Since(lastReport) > 500*time.Millisecond {
progress := 0.0
if total > 0 {
progress = float64(downloaded) / float64(total) * 100.0
}
progressCh <- DownloadStatus{
Distro: distro, FilePath: filePath,
TotalBytes: total, Downloaded: downloaded, Progress: progress,
}
lastReport = time.Now()
}
}
if readErr == io.EOF {
break
}
if readErr != nil {
return readErr
}
}
return nil
}


func execWithSudo(name string, args ...string) *exec.Cmd {
if os.Geteuid() == 0 {
return exec.Command(name, args...)
}
sudoArgs := append([]string{name}, args...)
return exec.Command("sudo", sudoArgs...)
}

func verifySHA256(filePath, expected string) bool {
if expected == "" {
return true
}
f, err := os.Open(filePath)
if err != nil {
return false
}
defer f.Close()
h := sha256.New()
io.Copy(h, f)
return fmt.Sprintf("%x", h.Sum(nil)) == expected
}
