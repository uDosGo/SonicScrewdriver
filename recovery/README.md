# SonicScrewdriver — Recovery Tools

This directory contains recovery and diagnostic tools for the Sonic USB drive.

## Directory Structure

```
recovery/
├── README.md
├── scripts/
│   ├── disk-repair.sh      — Disk repair & diagnostics
│   ├── password-reset.sh   — Linux password reset
│   ├── data-recovery.sh    — Data recovery (ddrescue, rsync, photorec, testdisk)
│   └── memtest.sh          — Memory diagnostics
└── tools/
    └── (place additional tools here)
```

## Usage

These scripts are copied to the exFAT data partition during `sonic usb create`.
They can be run directly from the USB drive on any Linux system.

### Quick Start

```bash
# Run a disk repair
sudo ./scripts/disk-repair.sh

# Reset a forgotten Linux password
sudo ./scripts/password-reset.sh

# Recover data from a failing drive
sudo ./scripts/data-recovery.sh

# Test system memory
sudo ./scripts/memtest.sh
```

### Adding Custom Tools

Place additional recovery tools (static binaries, scripts) in the `tools/` directory.
They will be automatically included when creating a Sonic USB drive.

## Included Tools

| Tool | Description |
|------|-------------|
| disk-repair.sh | Check/repair filesystems, bad blocks, wipe signatures |
| password-reset.sh | Mount Linux root and reset user passwords |
| data-recovery.sh | Clone drives, recover files, carve data |
| memtest.sh | Run memory diagnostics |

## Requirements

- Most scripts require root/sudo access
- Linux environment (or macOS with some limitations)
- For data recovery: ddrescue, testdisk, photorec (install via apt)
