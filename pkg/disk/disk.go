package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Device represents a detected block device
type Device struct {
Name       string `json:"name"`
Path       string `json:"path"`
Size       string `json:"size"`
Model      string `json:"model"`
Type       string `json:"type"`
MountPoint string `json:"mountpoint"`
FSType     string `json:"fstype"`
Label      string `json:"label"`
Removable  bool   `json:"removable"`
}

// Partition represents a partition to create
type Partition struct {
Number int    `json:"number"`
Size   string `json:"size"`
Type   string `json:"type"`
FS     string `json:"fs"`
Label  string `json:"label"`
Boot   bool   `json:"boot"`
}

// Layout defines the partition layout for an OS install
type Layout struct {
Name       string      `json:"name"`
LabelType  string      `json:"label_type"`
Partitions []Partition `json:"partitions"`
BootType   string      `json:"boot_type"`
}

type lsblkDevice struct {
Name       string         `json:"name"`
Size       uint64         `json:"size"`
Model      *string        `json:"model"`
Type       string         `json:"type"`
MountPoint *string        `json:"mountpoint"`
FSType     *string        `json:"fstype"`
Label      *string        `json:"label"`
Rm         bool           `json:"rm"`
Children   []lsblkDevice  `json:"children,omitempty"`
}

type lsblkOutput struct {
BlockDevices []lsblkDevice `json:"blockdevices"`
}

// DetectDevices returns all block devices
func DetectDevices(removableOnly bool) ([]Device, error) {
cmd := exec.Command("lsblk", "-o", "NAME,SIZE,MODEL,TYPE,MOUNTPOINT,FSTYPE,LABEL,RM", "-J", "-b")
output, err := cmd.Output()
if err != nil {
return nil, fmt.Errorf("lsblk failed: %w", err)
}

var lsblkData lsblkOutput
if err := json.Unmarshal(output, &lsblkData); err != nil {
return nil, fmt.Errorf("failed to parse lsblk output: %w", err)
}

var devices []Device
for _, d := range lsblkData.BlockDevices {
if removableOnly && !d.Rm {
continue
}
if d.Type != "disk" {
continue
}
model := ""
if d.Model != nil {
model = *d.Model
}
dev := Device{
Name:      d.Name,
Path:      "/dev/" + d.Name,
Size:      formatBytes(d.Size),
Model:     model,
Type:      d.Type,
Removable: d.Rm,
}
devices = append(devices, dev)
}
return devices, nil
}

// WipeDevice removes all partitions and creates a fresh GPT table
func WipeDevice(devicePath string) error {
fmt.Printf("  Wiping partition table on %s...\n", devicePath)

// 1. Unmount ALL partitions on the device (thorough)
fmt.Println("  Unmounting all partitions...")
unmountAllPartitions(devicePath)

// 2. Wipe filesystem signatures (with force flag, using sudo if needed)
fmt.Println("  Wiping filesystem signatures...")
cmd := execWithSudo("wipefs", "-a", "-f", devicePath)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
// Try with individual partitions
fmt.Printf("  wipefs on %s failed, trying per-partition...\n", devicePath)
if err := wipeAllPartitions(devicePath); err != nil {
return fmt.Errorf("wipefs failed: %w", err)
}
}

cmd = execWithSudo("parted", "-s", devicePath, "mklabel", "gpt")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("parted mklabel failed: %w", err)
}

fmt.Printf("  ✅ GPT partition table created\n")
return nil
}

// CreatePartitions creates partitions according to the layout
func CreatePartitions(devicePath string, layout Layout) error {
fmt.Printf("  Creating partition layout: %s\n", layout.Name)

for _, p := range layout.Partitions {
var fsType string
switch p.Type {
case "efi":
fsType = "fat32"
case "linux":
fsType = "ext4"
case "swap":
fsType = "linux-swap"
default:
fsType = p.Type
}

cmd := execWithSudo("parted", "-s", devicePath,
"mkpart", p.Label, fsType, "0%", p.Size)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("mkpart %s failed: %w", p.Label, err)
}

partNum := fmt.Sprintf("%d", p.Number)
if p.Boot {
execWithSudo("parted", "-s", devicePath, "set", partNum, "boot", "on").Run()
}
if p.Type == "efi" {
execWithSudo("parted", "-s", devicePath, "set", partNum, "esp", "on").Run()
}

fmt.Printf("  ✅ Partition %d: %s (%s)\n", p.Number, p.Label, fsType)
}

execWithSudo("partprobe", devicePath).Run()
return nil
}

// FormatPartition formats a partition
func FormatPartition(devicePath string, partNum int, fsType, label string) error {
partPath := getPartitionPath(devicePath, partNum)
fmt.Printf("  Formatting %s as %s...\n", partPath, fsType)

var cmd *exec.Cmd
switch fsType {
case "fat32", "vfat":
cmd = execWithSudo("mkfs.fat", "-F32", "-n", label, partPath)
case "ext4":
cmd = execWithSudo("mkfs.ext4", "-F", "-L", label, partPath)
case "swap":
cmd = execWithSudo("mkswap", "-L", label, partPath)
default:
return fmt.Errorf("unsupported filesystem: %s", fsType)
}

cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("mkfs failed: %w", err)
}
fmt.Printf("  ✅ Formatted %s\n", partPath)
return nil
}

// FormatAllPartitions formats all partitions in a layout
func FormatAllPartitions(devicePath string, layout Layout) error {
for _, p := range layout.Partitions {
if err := FormatPartition(devicePath, p.Number, p.FS, p.Label); err != nil {
return err
}
}
return nil
}

// MountPartition mounts a partition
func MountPartition(devicePath string, partNum int, mountPoint string) error {
partPath := getPartitionPath(devicePath, partNum)
if err := os.MkdirAll(mountPoint, 0755); err != nil {
return err
}
cmd := execWithSudo("mount", partPath, mountPoint)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("mount failed: %w", err)
}
fmt.Printf("  ✅ Mounted %s -> %s\n", partPath, mountPoint)
return nil
}

// UnmountPath unmounts a path
func UnmountPath(mountPoint string) error {
exec.Command("umount", mountPoint).Run()
exec.Command("umount", "-l", mountPoint).Run()
return nil
}

// UbuntuLayout returns the standard Ubuntu partition layout
func UbuntuLayout() Layout {
return Layout{
Name:      "Ubuntu",
LabelType: "gpt",
BootType:  "uefi",
Partitions: []Partition{
{Number: 1, Size: "512MiB", Type: "efi", FS: "fat32", Label: "EFI", Boot: true},
{Number: 2, Size: "100%", Type: "linux", FS: "ext4", Label: "ubuntu"},
},
}
}

// LinuxMintLayout returns the standard Linux Mint partition layout
func LinuxMintLayout() Layout {
return Layout{
Name:      "Linux Mint",
LabelType: "gpt",
BootType:  "uefi",
Partitions: []Partition{
{Number: 1, Size: "512MiB", Type: "efi", FS: "fat32", Label: "EFI", Boot: true},
{Number: 2, Size: "100%", Type: "linux", FS: "ext4", Label: "mint"},
},
}
}

// ClassicModernLayout returns the Classic Modern Mint layout with swap
func ClassicModernLayout() Layout {
return Layout{
Name:      "Classic Modern Mint",
LabelType: "gpt",
BootType:  "uefi",
Partitions: []Partition{
{Number: 1, Size: "1GiB", Type: "efi", FS: "fat32", Label: "EFI", Boot: true},
{Number: 2, Size: "8GiB", Type: "swap", FS: "swap", Label: "swap"},
{Number: 3, Size: "100%", Type: "linux", FS: "ext4", Label: "classicmodern"},
},
}
}

// unmountAllPartitions unmounts every partition on a device
func unmountAllPartitions(devicePath string) {
// Get all mount points for this device
cmd := exec.Command("sh", "-c", 
fmt.Sprintf("lsblk -ln -o MOUNTPOINT %s | grep -v '^$' || true", devicePath))
out, _ := cmd.Output()
for _, mp := range strings.Split(strings.TrimSpace(string(out)), "\n") {
mp = strings.TrimSpace(mp)
if mp != "" {
fmt.Printf("  Unmounting %s...\n", mp)
exec.Command("umount", mp).Run()
exec.Command("umount", "-l", mp).Run() // lazy unmount as fallback
}
}
// Also try partition-specific paths
for i := 1; i <= 10; i++ {
partPath := getPartitionPath(devicePath, i)
exec.Command("umount", partPath).Run()
exec.Command("umount", "-l", partPath).Run()
}
}

// wipeAllPartitions wipes each partition individually
func wipeAllPartitions(devicePath string) error {
for i := 1; i <= 10; i++ {
partPath := getPartitionPath(devicePath, i)
if _, err := os.Stat(partPath); os.IsNotExist(err) {
continue
}
cmd := execWithSudo("wipefs", "-a", "-f", partPath)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
return fmt.Errorf("wipefs %s failed: %w", partPath, err)
}
fmt.Printf("  Wiped %s\n", partPath)
}
return nil
}

// execWithSudo runs a command with sudo if not already root
func execWithSudo(name string, args ...string) *exec.Cmd {
if os.Geteuid() == 0 {
return exec.Command(name, args...)
}
// Prepend sudo
sudoArgs := append([]string{name}, args...)
return exec.Command("sudo", sudoArgs...)
}

func formatBytes(bytes uint64) string {
const unit = 1024
if bytes < unit {
return fmt.Sprintf("%dB", bytes)
}
div, exp := uint64(unit), 0
for n := bytes / unit; n >= unit; n /= unit {
div *= unit
exp++
}
return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getPartitionPath(devicePath string, partNum int) string {
if strings.Contains(devicePath, "nvme") {
return fmt.Sprintf("%sp%d", devicePath, partNum)
}
return fmt.Sprintf("%s%d", devicePath, partNum)
}

func unmountPartitions(devicePath string) {
exec.Command("sh", "-c",
fmt.Sprintf("mount | grep %s | awk '{print $3}' | xargs -r umount 2>/dev/null", devicePath)).Run()
}
