# Sonic USB Drive Creation Guide

This guide explains how to create a Sonic USB drive with the triple-partition layout.

## Overview

A Sonic USB drive has three partitions:

1. **ESP (FAT32, 512MB)** — SonicScrewloader bootloader
2. **Linux (ext4, 32GB)** — Customized Linux Mint installation
3. **exFAT (remaining)** — Cross-platform data partition

## Prerequisites

- A USB drive (128GB recommended)
- Linux or macOS system
- `sonic` CLI installed

## Quick Start

```bash
# Create a Sonic USB drive
sudo sonic usb create /dev/sdX

# With custom options
sudo sonic usb create /dev/sdX \
    --macos-installer /Applications/Install\\ macOS\\ Sonoma.app \
    --linux-iso ./sonic-mint-22-cinnamon-64bit.iso \
    --recovery-dir ./recovery
```

## Step-by-Step

### 1. Identify Your USB Drive

```bash
# Linux
lsblk

# macOS
diskutil list
```

Look for your USB drive (e.g., `/dev/sdb` on Linux, `/dev/disk2` on macOS).

### 2. Create the Sonic USB Drive

```bash
sudo sonic usb create /dev/sdX
```

Replace `/dev/sdX` with your actual device path.

### 3. Verify the Installation

```bash
sonic usb info /dev/sdX
sonic bootloader status /dev/sdX
```

### 4. Boot from the USB

1. Restart your computer
2. Enter boot menu (F12/F2/Option key)
3. Select the Sonic USB drive
4. You'll see the SonicScrewloader teletext menu

## Building the Linux Mint ISO

To create a customized Linux Mint ISO with Sonic tools pre-installed:

```bash
# Download Linux Mint 22
wget https://mirrors.edge.kernel.org/linuxmint/stable/22/linuxmint-22-cinnamon-64bit.iso

# Build customized ISO
cd mint
sudo ./build-iso.sh --input ../linuxmint-22-cinnamon-64bit.iso
```

The customized ISO will be at `mint/output/sonic-mint-22-cinnamon-64bit.iso`.

## macOS Installer

To include the macOS Sonoma installer:

1. Download from the App Store
2. The installer app goes to `/Applications/Install macOS Sonoma.app`
3. Pass `--macos-installer` to `sonic usb create`

## Recovery Tools

Recovery scripts are automatically copied to the exFAT partition.
See `recovery/README.md` for details.

## Testing

### Test on Mac (Intel)

1. Insert Sonic USB
2. Restart holding `Option` key
3. Select "EFI Boot"
4. You should see the SonicScrewloader menu

### Test on Mac (Apple Silicon)

1. Insert Sonic USB
2. Restart holding power button
3. Select boot volume
4. Note: Apple Silicon requires signed bootloaders

### Test on PC (UEFI)

1. Insert Sonic USB
2. Enter BIOS/UEFI setup (F2/Del)
3. Enable "Legacy Boot" or "CSM" if needed
4. Set USB as first boot device
5. Save and restart

### Test on PC (Legacy BIOS)

1. Insert Sonic USB
2. Enter boot menu (F12)
3. Select USB drive
4. The BIOS bootloader will load

## Troubleshooting

### "Device or resource busy"

Unmount the drive first:
```bash
# Linux
sudo umount /dev/sdX*

# macOS
sudo diskutil unmountDisk /dev/diskX
```

### "Permission denied"

Use `sudo` for all USB operations.

### Bootloader not showing

- Ensure UEFI boot is enabled in BIOS
- Try the other USB port (USB 2.0 vs 3.0)
- Some systems need Secure Boot disabled

## Reference USB Build

For reproducible builds:

```bash
# Build the reference 128GB USB
sonic usb create /dev/sdX \
    --label "SONIC-REF" \
    --macos-installer /Applications/Install\\ macOS\\ Sonoma.app \
    --linux-iso ./sonic-mint-22-cinnamon-64bit.iso \
    --recovery-dir ./recovery

# Verify
sonic usb info /dev/sdX
sonic bootloader status /dev/sdX
```
