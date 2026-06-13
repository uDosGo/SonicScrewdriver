# SonicScrewdriver v2.0.0 — Roadmap

> **Last Updated:** 2026-06-13
> **Status:** v2.0.0 — Round 1 complete (bootloader, CLI, device library, Mint ISO)

## ✅ Current State (v2.0.0)

SonicScrewdriver v2 is a **Universal USB Bootloader & System Toolkit** — Python CLI + C bootloader for USB creation, security enrollment, device management, and Linux Mint deployment.

### Implemented Commands (v2 Python CLI)

| Command | Status | Notes |
|---------|--------|-------|
| `sonic usb create` | ✅ | Triple-partition USB (ESP + ext4 + exFAT) |
| `sonic security enroll` | ✅ | FIDO2/U2F, GPG, SSH key enrollment |
| `sonic mint` | ✅ | Linux Mint ISO customization & deployment |
| `sonic mesh init/join` | ✅ | Peer-to-peer mesh networking |
| `sonic chasis launch` | ✅ | Game library management |
| `sonic bootloader install` | ✅ | Install SonicScrewloader to USB |
| `sonic device scan` | ✅ | Connected device detection (USB, PCI, Bluetooth) |
| `sonic device identify` | ✅ | Hardware capability identification |
| `sonic device lookup` | ✅ | Query local device database |
| `sonic device add/remove` | ✅ | CRUD for device entries |
| `sonic device repurpose` | ✅ | Router → Sonic Beacon / OpenWrt |

### Bootloader (C/asm)

| Component | Status | Notes |
|-----------|--------|-------|
| Teletext renderer | ✅ | 80x25 grid, 16 colors, block graphics, ANSI output |
| Framebuffer abstraction | ✅ | UEFI GOP + BIOS VGA text mode |
| OS detection | ✅ | SMBIOS/ACPI (Mac vs PC vs BIOS) |
| Menu engine | ✅ | YAML config → C structs at build time |
| Chainloading | ✅ | GRUB, rEFInd, EFI stub |
| Makefile | ✅ | UEFI (gnu-efi) + BIOS targets |
| YAML menu configs | ✅ | mac.yaml, pc.yaml, bios.yaml, teletext.yaml |

### Linux Mint ISO Build

| Component | Status | Notes |
|-----------|--------|-------|
| build-iso.sh | ✅ | Full pipeline: extract → unsquashfs → overlay → chroot → repack → xorriso |
| chroot-customize.sh | ✅ | Hostname, locale, timezone, keyboard, user, desktop |
| install-sonic.sh | ✅ | Sonic CLI + dependencies installation |
| Plymouth theme | ✅ | Sonic bolt logo, progress bar |
| Recovery tools | ✅ | disk-repair, password-reset, data-recovery, memtest |

### Legacy Go v1 (cmd/sonic/)

| Command | Status | Notes |
|---------|--------|-------|
| `sonic container` | ✅ | Preserved for backward compat |
| `sonic vault` | ✅ | Preserved for backward compat |
| `sonic gui` | ✅ | Preserved for backward compat |
| `sonic catalogue` | ✅ | Preserved for backward compat |
| `sonic knowledge` | ✅ | Preserved for backward compat |
| `sonic library` | ✅ | Preserved for backward compat |
| `sonic ventoy` | ✅ | Preserved for backward compat |
| `sonic remote` | ✅ | Preserved for backward compat |
| `sonic reflash` | ✅ | Preserved for backward compat |
| `sonic driver` | ✅ | Preserved for backward compat |
| `sonic recovery` | ✅ | Preserved for backward compat |

---

## 🎯 Next Steps

### Round 2 (v2.1.0)

1. **Build 128GB reference USB** (requires hardware)
   - Test on Mac (Intel + Apple Silicon)
   - Test on PC (UEFI + Legacy BIOS)
   - Ship 10 beta USBs

2. **Python CLI hardening**
   - Add unit tests for all CLI modules
   - Add `--help` for each subcommand
   - Add error handling improvements

3. **Device Library A3-A4**
   - Device flashing support (ESP32, routers)
   - Global device registry integration

### Round 3 (v2.2.0)

4. **Bootloader finalisation**
   - UEFI + BIOS bootloader final
   - Cross-compile for arm64 (Raspberry Pi)
   - Integration with SonicScrewloader

5. **Device Library A5**
   - AI-powered device recommendations
   - Commercial USB run

---

## 📊 Known Gaps

- No unit tests for Python CLI modules
- Bootloader requires gnu-efi to build (not installed on macOS)
- 128GB reference USB not yet built (requires hardware)
- Legacy Go v1 CLI has no CI pipeline

