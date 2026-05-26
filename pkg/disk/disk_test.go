package disk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"bytes", 500, "500B"},
		{"kilobytes", 1024, "1.0KB"},
		{"megabytes", 1048576, "1.0MB"},
		{"gigabytes", 1073741824, "1.0GB"},
		{"terabytes", 1099511627776, "1.0TB"},
		{"petabytes", 1125899906842624, "1.0PB"},
		{"exabytes", 1152921504606846976, "1.0EB"},
		{"zero", 0, "0B"},
		{"custom", 2048, "2.0KB"},
		{"large", 5368709120, "5.0GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPartitionPath(t *testing.T) {
	tests := []struct {
		name       string
		devicePath string
		partNum    int
		expected   string
	}{
		{"sda", "/dev/sda", 1, "/dev/sda1"},
		{"sda partition 2", "/dev/sda", 2, "/dev/sda2"},
		{"nvme", "/dev/nvme0n1", 1, "/dev/nvme0n1p1"},
		{"nvme partition 3", "/dev/nvme0n1", 3, "/dev/nvme0n1p3"},
		// mmcblk doesn't match nvme pattern, so it falls through to sda-style
		{"mmcblk", "/dev/mmcblk0", 1, "/dev/mmcblk01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPartitionPath(tt.devicePath, tt.partNum)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUbuntuLayout(t *testing.T) {
	layout := UbuntuLayout()
	assert.Equal(t, "Ubuntu", layout.Name)
	assert.Equal(t, "gpt", layout.LabelType)
	assert.Equal(t, "uefi", layout.BootType)
	assert.Len(t, layout.Partitions, 2)

	assert.Equal(t, 1, layout.Partitions[0].Number)
	assert.Equal(t, "efi", layout.Partitions[0].Type)
	assert.Equal(t, "fat32", layout.Partitions[0].FS)
	assert.True(t, layout.Partitions[0].Boot)

	assert.Equal(t, 2, layout.Partitions[1].Number)
	assert.Equal(t, "linux", layout.Partitions[1].Type)
	assert.Equal(t, "ext4", layout.Partitions[1].FS)
}

func TestLinuxMintLayout(t *testing.T) {
	layout := LinuxMintLayout()
	assert.Equal(t, "Linux Mint", layout.Name)
	assert.Equal(t, "gpt", layout.LabelType)
	assert.Len(t, layout.Partitions, 2)
}

func TestClassicModernLayout(t *testing.T) {
	layout := ClassicModernLayout()
	assert.Equal(t, "Classic Modern Mint", layout.Name)
	assert.Equal(t, "gpt", layout.LabelType)
	assert.Len(t, layout.Partitions, 3)

	// Should have swap partition
	assert.Equal(t, 2, layout.Partitions[1].Number)
	assert.Equal(t, "swap", layout.Partitions[1].Type)
	assert.Equal(t, "swap", layout.Partitions[1].FS)
}

func TestDeviceStruct(t *testing.T) {
	dev := Device{
		Name:      "sda",
		Path:      "/dev/sda",
		Size:      "1.0TB",
		Model:     "Test Disk",
		Type:      "disk",
		Removable: true,
	}
	assert.Equal(t, "sda", dev.Name)
	assert.Equal(t, "/dev/sda", dev.Path)
	assert.Equal(t, "1.0TB", dev.Size)
	assert.Equal(t, "Test Disk", dev.Model)
	assert.True(t, dev.Removable)
}

func TestPartitionStruct(t *testing.T) {
	p := Partition{
		Number: 1,
		Size:   "512MiB",
		Type:   "efi",
		FS:     "fat32",
		Label:  "EFI",
		Boot:   true,
	}
	assert.Equal(t, 1, p.Number)
	assert.Equal(t, "512MiB", p.Size)
	assert.Equal(t, "efi", p.Type)
	assert.Equal(t, "fat32", p.FS)
	assert.Equal(t, "EFI", p.Label)
	assert.True(t, p.Boot)
}
