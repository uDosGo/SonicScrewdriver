# SonicScrewdriver v2.0.0

## Universal USB Bootloader & System Toolkit

SonicScrewdriver v2 is a **Python CLI + C bootloader** for creating bootable USB drives, managing device firmware, enrolling security keys, and deploying Linux Mint systems. The legacy Go v1 CLI is preserved in `cmd/sonic/` for backward compatibility.

## 🎯 What It Does (v2 Python CLI)

```
sonic usb create <device>              — Create triple-partition USB (ESP + ext4 + exFAT)
sonic security enroll <type>           — Enroll FIDO2/U2F, GPG, or SSH keys
sonic mint <check|build|apply>         — Linux Mint ISO customization & deployment
sonic mesh init|join <network>         — Peer-to-peer mesh networking
sonic chasis launch <game>             — Game library management
sonic bootloader install <device>      — Install SonicScrewloader to USB
sonic device scan|identify|lookup      — Device library management
sonic device add|remove|repurpose      — CRUD + router repurposing
```

## 🚀 Quick Start

```bash
# Install Python CLI
cd cli && pip install -e .

# View help
sonic --help
```

## 🏗️ Project Structure

```
SonicScrewdriver/
├── bootloader/             # C/asm bootloader (Teletext, UEFI, BIOS)
│   ├── src/                # teletext.c, framebuffer.c, detect.c, menu.c, chainload.c
│   ├── include/            # teletext.h, detect.h, menu.h
│   ├── config/menus/       # YAML menu configs (mac.yaml, pc.yaml, bios.yaml)
│   └── Makefile            # UEFI (gnu-efi) + BIOS targets
├── cli/                    # Python CLI (v2)
│   ├── sonic/              # CLI modules (usb, security, mint, mesh, chasis, bootloader, device)
│   │   └── data/devices/   # YAML device entries (PCs, routers, ESP32)
│   └── setup.py            # v2.0.0
├── mint/                   # Linux Mint ISO build system
│   ├── build-iso.sh        # Full ISO build pipeline
│   ├── chroot-customize.sh # Hostname, locale, user, desktop config
│   ├── install-sonic.sh    # Sonic CLI + dependencies
│   └── overlay/            # Plymouth theme, skel, backgrounds
├── recovery/               # Recovery tools
│   └── scripts/            # disk-repair, password-reset, data-recovery, memtest
├── cmd/sonic/              # Legacy Go CLI (v1, preserved for backward compat)
├── pkg/                    # Legacy Go packages (v1)
├── docs/                   # Documentation
└── version                 # v2.0.0
```

## 📖 Documentation

- **[docs/USB-CREATION.md](docs/USB-CREATION.md)** — Full USB creation guide
- **[docs/legacy/](docs/legacy/)** — Archived documentation from earlier aspirational scope

## Related Repositories

- **uServer** — Backend services, secret store, API central
- **DevStudio** — Development environment configuration and tooling

