---
title: "Sonic-Screwdriver — Minimal OS Installation Brief"
status: draft
last_updated: 2026-05-20T22:32:15+10:00
category: readme
tags: [cli, sonic, sonicscrewdriver]
description: "| Property              | Value                                                |"
---
# Sonic-Screwdriver — Minimal OS Installation Brief

## 📋 Document Control

| Property              | Value                                                |
| --------------------- | ---------------------------------------------------- |
| **Brief Version**     | 1.0.0                                                |
| **Target OS**         | Linux Mint (Classic Modern) + Ubuntu (uHome Console) |
| **Installer**         | Sonic-Screwdriver (ventoy-based)                     |
| **Bundle Philosophy** | Minimal base + service modules = no bloat            |
| **Last Updated**      | 2026-04-21                                           |

***

## 🎯 Core Philosophy

**Sonic installs only what's needed — nothing more. One USB key, two personalities, infinite deployments.**

````
Sonic Installer
├── Linux Mint (Classic Modern) → Desktop-first, Classic Mac aesthetic
├── Ubuntu (uHome Console) → 10-foot UI, controller-first
└── Shared → No bundled media, no preinstalled apps (except Firefox)
````

***

## 📦 1. Minimal Base Requirements

### 1.1 Shared Components (Both OS)

| Component           | Version       | Purpose                        |
| ------------------- | ------------- | ------------------------------ |
| **Linux Kernel**    | 6.5+          | Hardware support               |
| **Systemd**         | 249+          | Service management             |
| **Network Manager** | 1.42+         | Network orchestration          |
| **PipeWire**        | 0.3+          | Audio (minimal)                |
| **Xorg/Wayland**    | Latest        | Display server                 |
| **Python 3**        | 3.10+         | Sonic dependencies             |
| **Go**              | 1.21+         | Sonic runtime                  |
| **Docker/Podman**   | Latest        | Container services             |
| **OpenSSH**         | Latest        | Remote access                  |
| **Firefox**         | ESR or Latest | Web browser (only bundled app) |

### 1.2 Excluded (No Bundled Media/Apps)

````
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
````

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

### 2.2 Mint-Specific Packages

```bash
# Core desktop (minimal)
cinnamon
cinnamon-session
cinnamon-control-center
nemo (file manager)
xapps-common

# No mint-meta-core (excludes bloat)
# No mint-meta-codecs
# No libreoffice-common
# No thunderbird
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

## 🔵 3. Ubuntu — uHome Console Edition

### 3.1 Base System

| Component             | Specification                             |
| --------------------- | ----------------------------------------- |
| **Base**              | Ubuntu 22.04 LTS (Minimal)                |
| **Desktop**           | uHome Console (custom) + GNOME (fallback) |
| **ISO Size Target**   | < 1.8 GB                                  |
| **RAM Usage (idle)**  | < 500 MB (console mode)                   |
| **Disk Usage (base)** | < 5 GB                                    |

### 3.2 Ubuntu-Specific Packages

```bash
# Minimal Xorg/Wayland
xorg
wayland-protocols

# Console UI (custom)
uhome-console (custom package)
uhome-console-controller

# GNOME (minimal for admin fallback)
gnome-shell --no-install-recommends
gnome-terminal
gnome-control-center --no-install-recommends

# No ubuntu-desktop-minimal (still too much)
# No snapd preloads
# No gnome-games
# No gnome-music
# No totem
```

### 3.3 uHome Console Theme (Pre-installed)

| Component           | Location                                         |
| ------------------- | ------------------------------------------------ |
| Console UI          | `/opt/uhome-console/`                            |
| Controller Service  | `/usr/lib/systemd/user/uhome-controller.service` |
| Tailwind Config     | `/etc/uhome/tailwind.config.js`                  |
| Admin Desktop Theme | `/usr/share/themes/Classic-Modern-Admin/`        |

### 3.4 Ubuntu Service Profile

```yaml
# /etc/sonic/profiles/ubuntu-uhome.yaml
profile:
  name: "ubuntu-uhome-console"
  role: "Smart Home Console"
  
  enabled_services:
    - sonic-home        # Local monitor
    - sonic-express     # Remote bridge
    - uhome-controller  # Controller input service
    - device-hub        # IoT device management
    - automation-engine # Scene automation
    
  disabled_services:
    - vault-master      # Optional (cloud sync)
    - code-vault        # Not needed on console
    
  console_features:
    - default_mode: "console"  # Boots to 10-foot UI
    - controller_paired: false # First-time setup
    - fallback_desktop: "gnome"
    - desktop_switch_key: "Ctrl+Alt+F7"
```

***

## 🔧 4. Sonic Service Modules (Both OS)

### 4.1 Required Modules

| Module            | Purpose               | Install Method                           |
| ----------------- | --------------------- | ---------------------------------------- |
| **sonic-core**    | Base orchestration    | Pre-installed                            |
| **sonic-home**    | Always-on monitor     | Pre-installed                            |
| **sonic-network** | Family network config | Pre-installed                            |
| **sonic-vault**   | Document spine        | Pre-installed (Mint) / Optional (Ubuntu) |
| **sonic-express** | Remote bridge         | Optional (Ubuntu primary)                |
| **sonic-boot**    | USB/ISO management    | Pre-installed                            |

### 4.2 Optional Modules (User Install)

```bash
# Install additional modules as needed
sonic module install --name=media-server
sonic module install --name=print-server
sonic module install --name=backup-daemon
```

***

## 🔐 5. User & Security Features

### 5.1 First-Time Setup (Sonic-Driven)

```bash
# Boot from USB → Sonic installer launches
# User answers minimal questions:

1. Select OS: [Mint Classic Modern] or [Ubuntu uHome Console]
2. Username: _______
3. Password: _______
4. Hostname: _______ (default: classic-modern or uhome-console)
5. Encryption: [ ] Encrypt home folder
6. Sonic Profile: [Default] [Custom]
```

### 5.2 Default User Configuration

```yaml
# /etc/sonic/user-defaults.yaml
users:
  default:
    name: "sonic"
    uid: 1000
    groups: ["sudo", "sonic-admin", "docker", "netdev"]
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

### 5.3 Firewall Configuration

```bash
# Sonic configures minimal firewall (UFW)
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp    # SSH
ufw allow 8042/tcp  # Sonic API
ufw allow 8043/tcp  # Web UI
# No media ports (optional install)
```

### 5.4 SSH Hardening

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

## 💾 6. USB Boot Disk Creation

### 6.1 Ventoy Configuration

```bash
# Sonic creates bootable USB with Ventoy + custom theme
sonic boot usb --create \
  --device=/dev/sdb \
  --format \
  --ventoy-theme=classic-modern \
  --iso=linux-mint-classic-modern.iso \
  --iso=ubuntu-uhome-console.iso
```

### 6.2 Boot Menu Options

````
┌─────────────────────────────────────────────────────────────┐
│                 CLASSIC MODERN BOOT SELECTOR                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Linux Mint — Classic Modern Edition                     │
│     Minimal desktop with Classic Mac aesthetic              │
│                                                             │
│  2. Ubuntu — uHome Console                                  │
│     10-foot smart home interface, controller-ready          │
│                                                             │
│  3. Sonic Recovery Tools                                    │
│     System repair, vault recovery, network diagnostics      │
│                                                             │
│  4. Memory Test                                             │
│                                                             │
│  5. Boot from Local Disk                                    │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  F1: Help  |  F2: Language  |  F3: Display  |  F4: About   │
└─────────────────────────────────────────────────────────────┘
````

### 6.3 Ventoy Custom Skin

```yaml
# /boot/ventoy/ventoy.json
{
  "theme": "classic-modern",
  "background": "/boot/ventoy/linen-bg.png",
  "font": "ChicagoFLF",
  "menu": {
    "color_normal": "#111111/#F2F2F2",
    "color_highlight": "#FFFFFF/#3A7BD5"
  },
  "timeout": 10,
  "default": 0
}
```

***

## 📊 7. Installation Size Comparison

| Component           | Mint Classic Modern | Ubuntu uHome Console                       |
| ------------------- | ------------------- | ------------------------------------------ |
| Base OS             | 1.2 GB              | 1.0 GB                                     |
| Desktop Environment | 400 MB (Cinnamon)   | 200 MB (Console) + 300 MB (GNOME fallback) |
| Sonic Services      | 150 MB              | 150 MB                                     |
| Firefox             | 100 MB              | 100 MB                                     |
| **Total ISO Size**  | **~1.85 GB**        | **~1.75 GB**                               |
| **Installed Size**  | **~6 GB**           | **~5 GB**                                  |

***

## 🚀 8. Installation Commands

### 8.1 Create Boot USB

```bash
# On any machine with sonic installed
sonic boot usb --create \
  --device=/dev/sdb \
  --os=mint,ubuntu \
  --edition=classic-modern,uhome-console
```

### 8.2 Network Install (PXE)

```bash
# Serve installer over network
sonic network pxe --serve \
  --os=mint \
  --edition=classic-modern \
  --ip-range=192.168.1.100-200
```

### 8.3 Automated Install (Preseed/Kickstart)

```bash
# Create answer file
sonic install generate --preseed \
  --os=mint \
  --user=sonic \
  --encrypt-home=true \
  --output=preseed.cfg

# Boot with preseed
# Append to kernel line: auto url=file:///preseed.cfg
```

***

## ✅ 9. Post-Install Verification

### 9.1 Sonic Health Check

```bash
# After installation, run:
sonic doctor --full

# Expected output:
# ┌─────────────────────────────────────────────────────────────┐
# │ SONIC INSTALLATION VERIFICATION                             │
# ├─────────────────────────────────────────────────────────────┤
# │ OS: Linux Mint 21.3 (Classic Modern Edition)               │
# │ Kernel: 6.5.0-generic                                       │
# │                                                             │
# │ ✓ Base system minimal (no bloat detected)                  │
# │ ✓ Firefox installed                                         │
# │ ✓ Sonic services running (sonic-home, network-bridge)      │
# │ ✓ User 'sonic' configured                                   │
# │ ✓ Firewall active                                           │
# │ ✓ SSH enabled (key auth)                                    │
# │                                                             │
# │ Installation type: MINIMAL                                  │
# │ Status: READY FOR DEPLOYMENT                                │
# └─────────────────────────────────────────────────────────────┘
```

### 9.2 Validate No Bundled Media

```bash
# Check for excluded packages
dpkg -l | grep -E "libreoffice|thunderbird|rhythmbox|gimp" | wc -l
# Expected: 0
```

***

## 📦 10. Package Sources

### 10.1 Official Repositories

```bash
# /etc/apt/sources.list.d/sonic.sources
Types: deb
URIs: http://archive.ubuntu.com/ubuntu/
Suites: jammy jammy-updates jammy-security
Components: main restricted universe
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

# Sonic custom repository (minimal packages)
Types: deb
URIs: https://repo.sonic.sh/stable
Suites: main
Components: sonic-minimal
Signed-By: /usr/share/keyrings/sonic-archive-keyring.gpg
```

### 10.2 Excluded Repositories

```bash
# Not added by default (user can enable)
# - multiverse (non-free)
# - backports
# - partner
# - ppa:libreoffice/ppa
```

***

## 🎯 11. Success Criteria

| Requirement               | Mint    | Ubuntu  | Status |
| ------------------------- | ------- | ------- | ------ |
| No bundled media/apps     | ✅      | ✅      | Design |
| Firefox only app          | ✅      | ✅      | Design |
| < 2 GB ISO                | 1.85 GB | 1.75 GB | Target |
| < 600 MB RAM idle         | 550 MB  | 450 MB  | Target |
| Sonic service integration | ✅      | ✅      | Design |
| USB boot (Ventoy)         | ✅      | ✅      | Design |
| User/security features    | ✅      | ✅      | Design |
| Classic Modern theme      | ✅      | N/A     | Design |
| uHome Console theme       | N/A     | ✅      | Design |

***

**Sonic now delivers minimal, purpose-built OS images — no bloat, just the essentials for uHome and uDos infrastructure.**