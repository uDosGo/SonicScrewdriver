---
title: "Sonic-Screwdriver v2 — Minimal OS Installation Brief"
status: draft
last_updated: 2026-06-13T20:43:00+10:00
category: readme
tags: [cli, sonic, sonicscrewdriver, v2]
description: "Sonic v2 — Python CLI + C bootloader for minimal Linux Mint ISO deployment"
---

# Sonic-Screwdriver v2 — Minimal OS Installation Brief

## 📋 Document Control

| Property              | Value                                                |
| --------------------- | ---------------------------------------------------- |
| **Brief Version**     | 2.0.0                                                |
| **Target OS**         | Linux Mint (Classic Modern)                          |
| **Installer**         | Sonic-Screwdriver v2 (Python CLI + C bootloader)     |
| **Bundle Philosophy** | Minimal base + service modules = no bloat            |
| **Last Updated**      | 2026-06-13                                           |

***

## 🎯 Core Philosophy

**Sonic v2 installs only what's needed — nothing more. One USB key, minimal Mint ISO, zero bloat.**

```
Sonic v2 Installer
├── Linux Mint (Classic Modern) → Desktop-first, Classic Mac aesthetic
├── SonicScrewloader → Custom bootloader (Teletext UI, UEFI + BIOS)
└── Recovery tools → Disk repair, password reset, data recovery, memtest
```

**Key changes from v1:**
- **Python CLI** (Click) replaces legacy Go CLI — lighter, faster to develop
- **C bootloader** replaces Ventoy dependency — Teletext-themed, UEFI + BIOS native
- **Device Library** — YAML-based hardware database for scan/identify/lookup
- **Recovery tools** bundled on ISO — no separate recovery media needed
- **ISO size target** reduced to < 1.5 GB (was 1.85 GB)

***

## 📦 1. Minimal Base Requirements

### 1.1 Shared Components

| Component           | Version       | Purpose                        |
| ------------------- | ------------- | ------------------------------ |
| **Linux Kernel**    | 6.5+          | Hardware support               |
| **Systemd**         | 249+          | Service management             |
| **Network Manager** | 1.42+         | Network orchestration          |
| **PipeWire**        | 0.3+          | Audio (minimal)                |
| **Xorg/Wayland**    | Latest        | Display server                 |
| **Python 3**        | 3.10+         | Sonic v2 CLI dependencies      |
| **OpenSSH**         | Latest        | Remote access                  |
| **Firefox**         | ESR or Latest | Web browser (only bundled app) |

### 1.2 Excluded (No Bundled Media/Apps)

```
✗ LibreOffice / Office suites
✗ GIMP / Image editors
✗ Thunderbird / Email clients
✗ Rhythmbox / Music players
✗ Celluloid / Video players
✗ Games
✗ Printing services (CUPS) - optional install
✗ Bluetooth stack - optional install
✗ Scanner drivers - optional install
✗ Language packs (except English)
✗ Docker/Podman - optional install
✗ Go runtime - not needed (v2 is Python)
```

***

## 🟢 2. Linux Mint — Classic Modern Edition

### 2.1 Base System

| Component             | Specification                    |
| --------------------- | -------------------------------- |
| **Base**              | Linux Mint 21.3 (Minimal)        |
| **Desktop**           | Cinnamon (Classic Modern themed) |
| **ISO Size Target**   | < 1.5 GB                         |
| **RAM Usage (idle)**  | < 600 MB                         |
| **Disk Usage (base)** | < 6 GB                           |

### 2.2 Build Process (Sonic v2)

```bash
# Build customized Mint ISO
sonic mint build --output=classic-modern-mint.iso

# Apply chroot customizations
sonic mint customize --iso=classic-modern-mint.iso

# Install Sonic services
sonic mint install --iso=classic-modern-mint.iso
```

Or use the shell script directly:

```bash
cd SonicScrewdriver/mint/
sudo ./build-iso.sh --output=classic-modern-mint.iso
```

### 2.3 Classic Modern Theme (Pre-installed)

| Component     | Location                                           |
| ------------- | -------------------------------------------------- |
| GTK Theme     | `/usr/share/themes/Classic-Modern-Mint/`           |
| Icon Theme    | `/usr/share/icons/Classic-Modern-Icons/`           |
| Cursor Theme  | `/usr/share/icons/Classic-Modern-Cursors/`         |
| Plymouth Boot | `/usr/share/plymouth/themes/classic-modern/`       |
| LightDM Theme | `/usr/share/lightdm-webkit/themes/classic-modern/` |

### 2.4 Mint Service Profile

```yaml
# /etc/sonic/profiles/mint-classic-modern.yaml
profile:
  name: "mint-classic-modern"
  role: "Desktop + Server"
  
  enabled_services:
    - sonic-home        # Local monitor
    - vault-master      # Document spine
    - network-bridge    # Family network
    - code-vault        # Source management
    
  disabled_services:
    - media-server      # No bundled media
    - print-server      # Optional install
    
  user_features:
    - default_user: "sonic"
    - auto_login: false
    - sudo_group: "sonic-admin"
    - home_encryption: optional
```

***

## 🔧 3. Sonic v2 Service Modules

### 3.1 Required Modules

| Module            | Purpose               | Install Method                           |
| ----------------- | --------------------- | ---------------------------------------- |
| **sonic-core**    | Base orchestration    | Pre-installed                            |
| **sonic-home**    | Always-on monitor     | Pre-installed                            |
| **sonic-network** | Family network config | Pre-installed                            |
| **sonic-vault**   | Document spine        | Pre-installed                            |
| **sonic-boot**    | USB/ISO management    | Pre-installed                            |

### 3.2 Optional Modules (User Install)

```bash
# Install additional modules as needed
sonic module install --name=media-server
sonic module install --name=print-server
sonic module install --name=backup-daemon
```

***

## 🔐 4. User & Security Features

### 4.1 First-Time Setup (Sonic-Driven)

```bash
# Boot from USB → SonicScrewloader launches
# User selects from Teletext-themed menu:

1. Install Linux Mint Classic Modern
2. Boot from Local Disk
3. Recovery Tools
4. Memory Test
```

### 4.2 Default User Configuration

```yaml
# /etc/sonic/user-defaults.yaml
users:
  default:
    name: "sonic"
    uid: 1000
    groups: ["sudo", "sonic-admin", "netdev"]
    shell: "/bin/bash"
    
  security:
    password_expiry: 90
    sudo_no_password: false
    ssh_enabled: true
    ssh_key_only: false
    
  home:
    encryption: false  # Optional during install
    vault_auto_mount: true
    code_vault_path: "~/code-vault"
```

### 4.3 Firewall Configuration

```bash
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp    # SSH
ufw allow 8042/tcp  # Sonic API
```

### 4.4 SSH Hardening

```bash
# /etc/ssh/sshd_config.d/sonic.conf
PermitRootLogin no
PasswordAuthentication no (optional)
PubkeyAuthentication yes
AllowUsers sonic
MaxAuthTries 3
ClientAliveInterval 300
```

***

## 💾 5. USB Boot Disk Creation

### 5.1 Sonic v2 USB Creation

```bash
# Create triple-partition USB (EFI, data, ISO)
sonic usb create --device=/dev/sdb

# Install SonicScrewloader to USB
sonic bootloader install --device=/dev/sdb
```

### 5.2 Boot Menu (Teletext Theme)

```
┌─────────────────────────────────────────────────────────────┐
│                 SONIC SCREWDRIVER v2                         │
│                 Boot Selection Menu                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Install Linux Mint — Classic Modern Edition             │
│     Minimal desktop with Classic Mac aesthetic              │
│                                                             │
│  2. Boot from Local Disk                                    │
│                                                             │
│  3. Sonic Recovery Tools                                    │
│     ├── Disk Repair                                         │
│     ├── Password Reset                                      │
│     ├── Data Recovery                                       │
│     └── Memory Test                                         │
│                                                             │
│  4. Reboot                                                  │
│                                                             │
│                                                             │
│  Use ↑/↓ to navigate, Enter to select                       │
└─────────────────────────────────────────────────────────────┘
```

### 5.3 Ventoy Alternative (Legacy)

For users who prefer Ventoy, the ISO can still be used with any Ventoy USB:

```bash
# Just copy the ISO to a Ventoy USB
cp classic-modern-mint.iso /media/ventoy/
```

***

## 📊 6. Installation Size Comparison

| Component           | v1 (Go + Ventoy) | v2 (Python + C bootloader) |
| ------------------- | ----------------- | -------------------------- |
| Base OS             | 1.2 GB            | 1.2 GB                     |
| Desktop Environment | 400 MB (Cinnamon) | 400 MB (Cinnamon)          |
| Sonic Services      | 150 MB            | 50 MB (Python, no Go)      |
| Firefox             | 100 MB            | 100 MB                     |
| **Total ISO Size**  | **~1.85 GB**      | **~1.75 GB**               |
| **Installed Size**  | **~6 GB**         | **~5.5 GB**                |

***

## ✅ 7. Post-Install Verification

### 7.1 Sonic Health Check

```bash
# After installation, run:
sonic doctor --full

# Expected output:
# ┌─────────────────────────────────────────────────────────────┐
# │ SONIC v2 INSTALLATION VERIFICATION                          │
# ├─────────────────────────────────────────────────────────────┤
# │ OS: Linux Mint 21.3 (Classic Modern Edition)               │
# │ Kernel: 6.5.0-generic                                       │
# │                                                             │
# │ ✓ Base system minimal (no bloat detected)                  │
# │ ✓ Firefox installed                                         │
# │ ✓ Sonic v2 services running                                 │
# │ ✓ User 'sonic' configured                                   │
# │ ✓ Firewall active                                           │
# │ ✓ SSH enabled (key auth)                                    │
# │                                                             │
# │ Installation type: MINIMAL                                  │
# │ Status: READY FOR DEPLOYMENT                                │
# └─────────────────────────────────────────────────────────────┘
```

### 7.2 Validate No Bundled Media

```bash
dpkg -l | grep -E "libreoffice|thunderbird|rhythmbox|gimp" | wc -l
# Expected: 0
```

***

## 🎯 8. Success Criteria

| Requirement               | v1 (Go) | v2 (Python) | Status |
| ------------------------- | ------- | ----------- | ------ |
| No bundled media/apps     | ✅      | ✅          | Design |
| Firefox only app          | ✅      | ✅          | Design |
| < 2 GB ISO                | 1.85 GB | 1.75 GB     | Target |
| < 600 MB RAM idle         | 550 MB  | 550 MB      | Target |
| Sonic service integration | ✅      | ✅          | Design |
| USB boot (native)         | Ventoy  | SonicScrewloader | Design |
| Recovery tools on ISO     | ❌      | ✅          | Design |
| Device library            | ❌      | ✅          | Design |

***

**Sonic v2 delivers minimal, purpose-built OS images — Python CLI, C bootloader, zero bloat.**
