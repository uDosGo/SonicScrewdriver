#!/bin/bash
# =============================================================================
# SonicScrewdriver — Install Sonic Tools in Linux Mint ISO
# =============================================================================
# Installs the Sonic CLI tools and dependencies into the chroot environment.
# Called by chroot-customize.sh during the ISO build process.
# =============================================================================

set -euo pipefail

SONIC_VERSION="2.0.0"
SONIC_REPO_URL="https://packages.sonic.sh/apt"

echo "=============================================="
echo "  Installing SonicScrewdriver v${SONIC_VERSION}"
echo "=============================================="

# ---------------------------------------------------------------------------
# Add Sonic APT repository
# ---------------------------------------------------------------------------
echo "Adding Sonic APT repository..."
cat > /etc/apt/sources.list.d/sonic.list << EOF
deb ${SONIC_REPO_URL} stable main
EOF

# Add GPG key
curl -fsSL "${SONIC_REPO_URL}/sonic.gpg" | gpg --dearmor -o /etc/apt/trusted.gpg.d/sonic.gpg 2>/dev/null || true

# ---------------------------------------------------------------------------
# Install system dependencies
# ---------------------------------------------------------------------------
echo "Installing system dependencies..."
apt-get update
apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    git \
    curl \
    wget \
    parted \
    exfatprogs \
    efibootmgr \
    grub-efi-amd64-bin \
    grub-pc-bin \
    xorriso \
    squashfs-tools \
    openssh-client \
    gnupg \
    pcscd \
    scdaemon \
    yubikey-personalization \
    || echo "Some packages failed to install (non-critical)"

# ---------------------------------------------------------------------------
# Install Sonic CLI via pip
# ---------------------------------------------------------------------------
echo "Installing Sonic CLI..."
pip3 install sonic-screwdriver=="${SONIC_VERSION}" || {
    echo "Installing from local source..."
    if [[ -d /tmp/sonic-cli ]]; then
        pip3 install /tmp/sonic-cli
    fi
}

# ---------------------------------------------------------------------------
# Install MeshCore dependencies
# ---------------------------------------------------------------------------
echo "Installing MeshCore dependencies..."
pip3 install pynacl zeroconf pyserial esptool 2>/dev/null || true

# ---------------------------------------------------------------------------
# Create Sonic configuration directory
# ---------------------------------------------------------------------------
mkdir -p /etc/sonic
cat > /etc/sonic/config.yaml << EOF
# SonicScrewdriver v${SONIC_VERSION} — System Configuration
version: ${SONIC_VERSION}
hostname: sonic-mint
mesh:
  port: 8765
  discovery: true
security:
  default_type: auto
usb:
  esp_size_mb: 512
  linux_size_mb: 32768
EOF

# ---------------------------------------------------------------------------
# Create MOTD
# ---------------------------------------------------------------------------
cat > /etc/motd << 'EOF'
╔══════════════════════════════════════════════════════════════╗
║              SonicScrewdriver v2 — Universal USB            ║
║              Bootloader & System Toolkit                    ║
║                                                            ║
║  Type 'sonic --help' to get started                        ║
║  Type 'sonic usb create' to build a Sonic USB drive        ║
║  Type 'sonic mesh init' to join the mesh network           ║
╚══════════════════════════════════════════════════════════════╝
EOF

# ---------------------------------------------------------------------------
# Create desktop shortcut
# ---------------------------------------------------------------------------
mkdir -p /usr/share/applications
cat > /usr/share/applications/sonic.desktop << EOF
[Desktop Entry]
Name=SonicScrewdriver
Comment=Universal USB Bootloader & System Toolkit
Exec=sonic
Terminal=true
Type=Application
Icon=utilities-terminal
Categories=System;Utility;
EOF

echo "=============================================="
echo "  SonicScrewdriver v${SONIC_VERSION} installed!"
echo "=============================================="

exit 0
